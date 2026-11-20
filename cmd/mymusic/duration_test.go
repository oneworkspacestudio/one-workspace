package main

import (
	"strings"
	"testing"
)

func TestLyricDurationAccountsForSectionsAndTempo(t *testing.T) {
	opening := "[Verse 1]\nWe walked along the quiet shore\nAnd watched the rain fall through the door\n\n[Chorus]\nFollow the light, follow the light"
	full := opening + "\n\n[Verse 2]\nThe morning bells began to ring\nWe heard the whole town start to sing\n\n[Chorus]\nFollow the light, follow the light\n\n[End]"
	short := lyricDuration(opening, 105)
	long := lyricDuration(full, 105)
	faster := lyricDuration(full, 140)
	if short.Lines != 3 || short.Sections != 2 || long.Lines != 6 || long.Sections != 4 {
		t.Fatalf("counted headings or ignored repeated sections: short %+v, full %+v", short, long)
	}
	if long.Seconds <= short.Seconds || faster.Seconds >= long.Seconds {
		t.Fatalf("duration must grow for more lyrics and shrink for a faster tempo: %+v %+v %+v", short, long, faster)
	}
	if withoutOpening := lyricDuration("[Guitar solo]\n\n"+full, 105); withoutOpening.Seconds <= long.Seconds || withoutOpening.Lines != long.Lines {
		t.Fatalf("an instrumental section should add time and no sung lines: %+v", withoutOpening)
	}
}

func TestCustomLyricsReachACEUnchanged(t *testing.T) {
	lyrics := "[Verse 1]\r\n  Walk with me\r\n\r\n[Verse 2]\r\n Sing it back"
	song := generation{Prompt: "Folk", Lyrics: lyrics, AutoDuration: true, VocalLanguage: "en", Describe: true}
	payload, err := song.acePayload()
	if err != nil {
		t.Fatal(err)
	}
	if payload["lyrics"] != lyrics || payload["audio_duration"] != lyricDuration(lyrics, 0).Seconds || payload["bpm"] != defaultLyricBPM || payload["vocal_language"] != "en" {
		t.Fatalf("custom lyric payload changed: %#v", payload)
	}
	for _, key := range []string{"thinking", "use_format", "use_cot_caption", "use_cot_language"} {
		if payload[key] != false {
			t.Errorf("%s must be disabled for supplied lyrics", key)
		}
	}
	if _, exists := payload["sample_query"]; exists {
		t.Fatal("custom lyrics entered description mode")
	}
	song.AutoDuration, song.Seconds, song.BPM = false, 180, 105
	manual, err := song.acePayload()
	if err != nil || manual["audio_duration"] != 180 || manual["bpm"] != 105 {
		t.Fatalf("manual length or tempo was ignored: %#v, %v", manual, err)
	}
}

func TestAutoDurationReportsOverLimit(t *testing.T) {
	song := generation{Prompt: "Folk", Lyrics: strings.Repeat("[Verse]\nWords to sing again and again\n", 40), AutoDuration: true}
	if _, err := song.acePayload(); err == nil || !strings.Contains(err.Error(), "Split them into shorter tracks") {
		t.Fatalf("oversize lyrics need a clear length warning, got %v", err)
	}
}
