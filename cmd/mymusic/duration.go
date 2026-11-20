package main

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	minSongSeconds     = 5
	maxSongSeconds     = 240
	defaultSongSeconds = 30
	defaultLyricBPM    = 105
)

type durationEstimate struct {
	Seconds  int `json:"seconds"`
	BPM      int `json:"bpm"`
	Lines    int `json:"lines"`
	Sections int `json:"sections"`
}

func lyricDuration(lyrics string, bpm int) durationEstimate {
	if bpm <= 0 || bpm > 300 {
		bpm = defaultLyricBPM
	}
	estimate := durationEstimate{Seconds: defaultSongSeconds, BPM: bpm}
	var bars, phraseBars, instrumentalBars int
	finishSection := func() {
		if phraseBars > 0 {
			bars += phraseBars + 2
			estimate.Sections++
		} else {
			bars += instrumentalBars
		}
		phraseBars, instrumentalBars = 0, 0
	}
	for _, raw := range strings.Split(lyrics, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" && phraseBars > 0 {
			finishSection()
		}
		for strings.HasPrefix(line, "[") {
			tag, rest, found := strings.Cut(line[1:], "]")
			if !found {
				break
			}
			switch sectionTag(tag) {
			case "verse", "chorus", "pre-chorus", "post-chorus", "bridge", "refrain", "hook", "end":
				finishSection()
			case "intro", "outro", "ending":
				finishSection()
				instrumentalBars = 2
			case "interlude", "break", "instrumental", "solo", "guitar solo", "drum break":
				finishSection()
				instrumentalBars = 4
			}
			line = strings.TrimSpace(rest)
		}
		var letters, syllabicCharacters, words int
		for _, word := range strings.Fields(line) {
			var wordLetters int
			for _, r := range word {
				if unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana, unicode.Hangul) {
					syllabicCharacters++
				} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
					wordLetters++
				}
			}
			if wordLetters > 0 {
				letters += wordLetters
				words++
			}
		}
		units := max(words, (letters+2)/3) + syllabicCharacters
		if units > 0 {
			estimate.Lines++
			phraseBars += max(1, (units+5)/6)
		}
	}
	finishSection()
	if estimate.Lines > 0 {
		// Budget four beats for each bar, room at section transitions, and
		// additional phrasing time as verses and repeated choruses accumulate.
		bars += 4 + (bars+5)/6
		estimate.Seconds = max(12, (bars*240+bpm-1)/bpm)
	}
	return estimate
}

func sectionTag(tag string) string {
	tag = strings.ToLower(strings.TrimSpace(tag))
	tag, _, _ = strings.Cut(tag, ":")
	tag, _, _ = strings.Cut(tag, " – ")
	tag, _, _ = strings.Cut(tag, " - ")
	return strings.TrimRight(strings.TrimSpace(tag), " 0123456789")
}

func (song generation) acePayload() (map[string]any, error) {
	if strings.TrimSpace(song.Prompt) == "" || song.BPM < 0 || song.BPM > 300 {
		return nil, fmt.Errorf("enter a music brief and an optional tempo up to 300 BPM")
	}
	hasLyrics := !song.Instrumental && strings.TrimSpace(song.Lyrics) != ""
	if song.AutoDuration {
		song.Seconds = defaultSongSeconds
		if hasLyrics {
			estimate := lyricDuration(song.Lyrics, song.BPM)
			if estimate.Seconds > maxSongSeconds {
				return nil, fmt.Errorf("these lyrics need about %d seconds at %d BPM. Split them into shorter tracks or choose a manual length", estimate.Seconds, estimate.BPM)
			}
			song.Seconds, song.BPM = estimate.Seconds, estimate.BPM
		}
	}
	if song.Seconds < minSongSeconds || song.Seconds > maxSongSeconds {
		return nil, fmt.Errorf("choose a length between %d and %d seconds", minSongSeconds, maxSongSeconds)
	}
	language := strings.ToLower(strings.TrimSpace(song.VocalLanguage))
	if language == "" {
		language = "en"
	}
	if len(language) < 2 || len(language) > 3 || strings.Trim(language, "abcdefghijklmnopqrstuvwxyz") != "" {
		return nil, fmt.Errorf("choose a vocal language")
	}
	body := map[string]any{
		"prompt": song.Prompt, "lyrics": song.Lyrics, "audio_duration": song.Seconds,
		"batch_size": 1, "audio_format": "wav", "vocal_language": language,
		"thinking": false, "use_format": false, "use_cot_caption": false, "use_cot_language": false,
	}
	if song.BPM > 0 {
		body["bpm"] = song.BPM
	}
	if song.Instrumental {
		body["lyrics"] = "[Instrumental]"
	} else if song.Describe && !hasLyrics {
		body["sample_query"] = song.Prompt
		body["thinking"] = true
	}
	return body, nil
}
