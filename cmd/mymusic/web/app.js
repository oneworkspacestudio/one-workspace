"use strict";
const byId = id => document.getElementById(id);
const form = byId("compose");
const notice = byId("notice");
const player = byId("player");
let tracks = [];
let playingID;
let clips = [];
let context;
let playback = [];
let mixOperation = 0;
const buffers = new Map();
let session;
let managed = false;
const pause = ms => new Promise(resolve => setTimeout(resolve, ms));

function lock(message, reconnect) {
  session = undefined;
  playingID = undefined;
  player.pause();
  player.removeAttribute("src");
  player.load();
  stopMix();
  buffers.clear();
  clips = [];
  tracks = [];
  for (const id of ["work","transport","tabs","engine"]) byId(id).hidden = true;
  byId("entryMessage").textContent = message;
  byId("connect").hidden = !reconnect;
  byId("sign-in").hidden = false;
  if (byId("settings").open) byId("settings").close();
}

async function boot() {
  const response = await fetch("api/bootstrap");
  if (response.ok) {
    const config = await response.json();
    managed = config.id === "mymusic";
    if (!managed || !config.oidc) throw new Error("Vuno could not read its Workspace configuration.");
    const apps = byId("apps");
    apps.href = config.workspace || "/";
    apps.hidden = false;
    byId("connect").href = "auth/login";
    const identity = await fetch("api/me");
    if (identity.status === 401) {
      lock("Continue with MyAccounts to open your library.", true);
      return;
    }
    if (!identity.ok) throw new Error("MyAccounts could not check your session.");
    session = await identity.json();
    if (session.user.role !== "owner") {
      lock("Only the Workspace owner can open Vuno.", false);
      return;
    }
    byId("engineSettings").hidden = false;
    const checkSession = async () => {
      if (!session) return;
      if (Date.now() >= Date.parse(session.expires)) {
        lock("Your Vuno session expired. Reconnect with MyAccounts.", true);
        return;
      }
      try {
        const result = await fetch("api/me", {headers:{"X-Workspace-Expected-Account":session.user.id}});
        if (result.status === 401 || result.status === 409) {
          lock("Your Vuno session changed. Reconnect with MyAccounts.", true);
        }
      } catch {
        // A network interruption does not change the current sign-in.
      }
    };
    setInterval(checkSession, 30000);
    addEventListener("focus", checkSession);
  } else if (response.status !== 404) {
    throw new Error("Vuno could not start.");
  }
  for (const id of ["work","transport","engine"]) byId(id).hidden = false;
  await refresh();
}

function node(tag, text, className) {
  const item = document.createElement(tag);
  if (text) item.textContent = text;
  if (className) item.className = className;
  return item;
}
async function api(path, body) {
  const headers = new Headers(body?.headers);
  if (session) {
    headers.set("X-CSRF-Token", session.csrf);
    headers.set("X-Workspace-Expected-Account", session.user.id);
  }
  const response = await fetch("api/" + path, {...body, headers});
  if (managed && [401,403,409].includes(response.status)) {
    lock("Your Vuno session changed. Reconnect with MyAccounts.", true);
    throw new Error("Your MyAccounts session changed.");
  }
  const data = await response.json().catch(() => ({}));
  if (!response.ok) throw new Error(data.error || "Request failed");
  return data;
}
function post(path, data) {
  return api(path, {method:"POST", headers:{"Content-Type":"application/json"}, body:JSON.stringify(data)});
}
async function action(button, status, work) {
  button.disabled = true;
  status.textContent = "";
  try { await work(); }
  catch (error) { status.textContent = error.message; }
  finally { button.disabled = false; }
}
function navigate() {
  const page = ["home","create","library","studio"].find(id => location.hash === "#" + id) || "home";
  document.querySelectorAll(".page").forEach(section => { section.hidden = section.id !== page; });
  document.querySelectorAll("#tabs a").forEach(link => {
    if (link.hash === "#" + page) link.setAttribute("aria-current","page");
    else link.removeAttribute("aria-current");
  });
  notice.textContent = "";
  window.scrollTo(0, 0);
}

addEventListener("hashchange", navigate);
navigate();
byId("quick").onsubmit = event => {
  event.preventDefault();
  const idea = byId("idea").value.trim();
  if (!idea) return;
  form.elements.Prompt.value = idea;
  location.hash = "create";
  navigate();
  form.elements.Prompt.focus();
};
byId("engine").onclick = () => action(byId("engine"), notice, async () => {
  const settings = await api("settings");
  byId("connection").textContent = "ACE-Step " + (settings.ready ? "reachable" : "unavailable") + " at " + settings.engine + " · DeepSeek " + (settings.assistant ? "configured" : "off");
  byId("engineSettings").reset();
  byId("engineSettings").elements.Engine.value = settings.engine;
  byId("engineStatus").textContent = "";
  byId("settings").showModal();
});
byId("engineSettings").onsubmit = event => {
  event.preventDefault();
  const settings = event.currentTarget.elements;
  const data = {Engine:settings.Engine.value};
  if (settings.ACEKey.value || settings.clearACE.checked) data.ACEKey = settings.clearACE.checked ? "" : settings.ACEKey.value;
  if (settings.DeepSeekKey.value || settings.clearDeepSeek.checked) data.DeepSeekKey = settings.clearDeepSeek.checked ? "" : settings.DeepSeekKey.value;
  action(event.submitter, byId("engineStatus"), async () => {
    await api("settings", {method:"PUT",headers:{"Content-Type":"application/json"},body:JSON.stringify(data)});
    byId("settings").close();
    notice.textContent = "Engine settings saved.";
  });
};
form.onsubmit = event => {
  event.preventDefault();
  const data = Object.fromEntries(new FormData(form));
  data.Seconds = Number(data.Seconds);
  data.BPM = Number(data.BPM);
  data.Instrumental = data.Instrumental === "on";
  data.Describe = !data.Lyrics.trim();
  const status = byId("generationStatus");
  action(event.submitter, status, async () => {
    status.textContent = "ACE-Step is generating your song. Keep Vuno open.";
    let job = await post("generate", data);
    while (!job.done) {
      await pause(1500);
      job = await api("generate/" + job.id);
    }
    if (job.error) throw new Error(job.error);
    const track = job.track;
    await refresh();
    status.textContent = track.title + " is in your library.";
    play(track);
  });
};
for (const [id,mode,field] of [["lyrics","lyrics","Lyrics"],["refine","prompt","Prompt"]]) {
  byId(id).onclick = () => action(byId(id), byId("generationStatus"), async () => {
    byId("generationStatus").textContent = "Waiting for DeepSeek…";
    const result = await post("assist", {Prompt:form.elements.Prompt.value, Mode:mode});
    form.elements[field].value = result.text;
    if (mode === "lyrics") byId("custom").open = true;
    form.elements[field].focus();
    byId("generationStatus").textContent = "Draft ready.";
  });
}
const audioURL = track => "audio/" + track.id + "." + track.format;
function play(track) {
  if (!track) return;
  stopMix();
  playingID = track.id;
  player.src = audioURL(track);
  byId("playing").textContent = track.title;
  updatePlayer();
  player.play().catch(error => {
    if (error.name !== "AbortError") notice.textContent = error.message;
  });
}

function timeLabel(seconds) {
  const time = Math.floor(Number.isFinite(seconds) ? seconds : 0);
  return Math.floor(time / 60) + ":" + String(time % 60).padStart(2, "0");
}

function updatePlayer() {
  const playButton = byId("toggle");
  playButton.textContent = player.paused ? "▶" : "Ⅱ";
  playButton.setAttribute("aria-label", player.paused ? "Play" : "Pause");
  const seek = byId("seek");
  const hasDuration = Number.isFinite(player.duration) && player.duration > 0;
  seek.max = hasDuration ? player.duration : 1;
  seek.value = hasDuration ? player.currentTime : 0;
  seek.disabled = !hasDuration;
  byId("elapsed").textContent = timeLabel(player.currentTime);
  byId("total").textContent = timeLabel(player.duration);
}

for (const event of ["durationchange", "timeupdate", "play", "pause", "ended", "emptied"]) {
  player.addEventListener(event, updatePlayer);
}
byId("toggle").onclick = () => {
  if (!playingID) {
    play(tracks[0]);
    return;
  }
  if (player.paused) player.play().catch(error => { notice.textContent = error.message; });
  else player.pause();
};
function skipTrack(offset) {
  if (!tracks.length) return;
  const current = tracks.findIndex(track => track.id === playingID);
  play(tracks[current < 0 ? 0 : (current + offset + tracks.length) % tracks.length]);
}
byId("previous").onclick = () => skipTrack(-1);
byId("next").onclick = () => skipTrack(1);
byId("seek").oninput = event => { player.currentTime = Number(event.target.value); };
byId("volume").oninput = event => { player.volume = Number(event.target.value); };
updatePlayer();

function renderTracks(id, items, emptyMessage) {
  const list = byId(id);
  list.replaceChildren();
  if (!items.length) {
    list.append(node("p", emptyMessage, "empty"));
    return;
  }
  for (const track of items) {
    const row = node("div",null,"item");
    const details = node("div",null,"item-info");
    details.append(node("strong",track.title),node("small",track.format.toUpperCase()+" · "+new Date(track.created).toLocaleDateString()));
    const playButton = node("button","▶","cover");
    playButton.setAttribute("aria-label","Play "+track.title);
    playButton.style.setProperty("--hue", [...track.id].reduce((sum, char) => sum + char.charCodeAt(0) * 7, 0) % 360);
    playButton.onclick = () => play(track);
    const download = node("a","Download","download");
    download.href = audioURL(track);
    download.download = track.title + "." + track.format;
    row.append(playButton,details,download);
    list.append(row);
  }
}

function renderLibrary() {
  const query = byId("search").value.toLowerCase().trim();
  const visible = tracks.filter(track => track.title.toLowerCase().includes(query));
  const order = byId("sort").value;
  if (order === "old") visible.reverse();
  if (order === "title") visible.sort((a,b) => a.title.localeCompare(b.title));
  renderTracks("tracks", visible, query ? "No tracks match your search." : "Import audio or generate a song to start your library.");
}
byId("search").oninput = renderLibrary;
byId("sort").onchange = renderLibrary;
async function refresh() {
  tracks = (await api("library")).sort((a,b) => Date.parse(b.created) - Date.parse(a.created));
  renderLibrary();
  renderTracks("homeTracks", tracks.slice(0,3), "Your imported and generated tracks will appear here.");
  renderTracks("createTracks", tracks.slice(0,8), "Songs you generate will appear here. You can also import audio from Library.");
  for (const id of ["toggle", "previous", "next"]) byId(id).disabled = !tracks.length;
  const picker = byId("source");
  const previous = picker.value;
  picker.replaceChildren(...tracks.map(track => {
    const option = node("option",track.title);
    option.value = track.id;
    return option;
  }));
  if (tracks.some(track => track.id === previous)) picker.value = previous;
  byId("add").disabled = !tracks.length;
}
byId("import").onchange = () => action(byId("import"),notice,async () => {
  const input = byId("import");
  const file = input.files[0];
  if (!file) return;
  const dot = file.name.lastIndexOf(".");
  const title = file.name.slice(0,dot);
  try {
    await api("import",{method:"POST",headers:{"X-Audio-Name":encodeURIComponent(title)+"|"+file.name.slice(dot+1).toLowerCase()},body:file});
    await refresh();
    notice.textContent = "Audio imported.";
  } finally { input.value = ""; }
});
function getContext() {
  context ??= new AudioContext();
  return context;
}
async function decode(track) {
  if (!buffers.has(track.id)) {
    buffers.set(track.id,(async () => {
      const response = await fetch(audioURL(track));
      if (!response.ok) throw new Error("Could not open " + track.title);
      return getContext().decodeAudioData(await response.arrayBuffer());
    })().catch(error => { buffers.delete(track.id); throw error; }));
  }
  return buffers.get(track.id);
}
function drawWaveform(canvas, buffer) {
  const data = buffer.getChannelData(0);
  const ink = canvas.getContext("2d");
  ink.fillStyle = "#ff9b88";
  const count = canvas.width / 3;
  for (let i = 0; i < count; i++) {
    const start = Math.floor(i * data.length / count);
    const end = Math.floor((i + 1) * data.length / count);
    let peak = 0;
    for (let position = start; position < end; position += Math.max(1, Math.floor((end - start) / 64))) {
      peak = Math.max(peak, Math.abs(data[position]));
    }
    const height = Math.max(1, peak * canvas.height * .85);
    ink.fillRect(i * 3, (canvas.height - height) / 2 + 4, 2, height);
  }
}

function renderTimeline() {
  const length = Math.max(10,...clips.map(clip => clip.start + clip.buffer.duration));
  byId("ruler").replaceChildren(...[0,.25,.5,.75,1].map(point => node("span",(length * point).toFixed(1)+"s")));
  const lanes = byId("lanes");
  lanes.replaceChildren(...clips.map(clip => {
    const lane = node("div");
    const bar = node("div",null,"bar");
    bar.style.left = (clip.start / length * 100) + "%";
    bar.style.width = (clip.buffer.duration / length * 100) + "%";
    const wave = node("canvas");
    wave.width = 240;
    wave.height = 72;
    wave.setAttribute("aria-hidden", "true");
    drawWaveform(wave, clip.buffer);
    bar.append(wave, node("span", clip.track.title));
    lane.append(bar);
    return lane;
  }));
  if (!clips.length) {
    const empty = node("div",null,"studio-empty");
    empty.append(node("p","Add a library track to arrange it here."));
    lanes.append(empty);
  }
}
function renderClips() {
  const list = byId("clips");
  list.replaceChildren();
  for (const id of ["preview","stop","export"]) byId(id).disabled = !clips.length;
  for (const clip of clips) {
    const row = node("div",null,"clip");
    row.append(node("strong",clip.track.title));
    for (const [field,label,max,step] of [["start","Start (s)",600,.1],["gain","Level (%)",150,1]]) {
      const wrapper = node("label",label);
      const input = node("input");
      Object.assign(input,{type:"number", min:0, max, step, value:clip[field]});
      input.onchange = () => {
        clip[field] = Math.min(max,Math.max(0,Number(input.value)||0));
        input.value = clip[field];
        renderTimeline();
      };
      wrapper.append(input);
      row.append(wrapper);
    }
    const remove = node("button","Remove");
    remove.onclick = () => { clips = clips.filter(other => other !== clip); renderClips(); };
    row.append(remove);
    list.append(row);
  }
  renderTimeline();
}
byId("add").onclick = () => action(byId("add"),notice,async () => {
  if (clips.length >= 16) throw new Error("An arrangement can hold up to 16 clips.");
  const track = tracks.find(item => item.id === byId("source").value);
  if (track) clips.push({track,buffer:await decode(track),start:0,gain:100});
  renderClips();
});
function stopMix() {
  mixOperation++;
  for (const source of playback) source.stop();
  playback = [];
}
function schedule(target, items, when) {
  return items.map(clip => {
    const source = target.createBufferSource();
    const level = target.createGain();
    source.buffer = clip.buffer;
    level.gain.value = clip.gain / 100;
    source.connect(level).connect(target.destination);
    source.start(when + clip.start);
    source.onended = () => { source.disconnect(); level.disconnect(); };
    return source;
  });
}
byId("preview").onclick = () => action(byId("preview"),notice,async () => {
  if (!clips.length) throw new Error("Add a clip first.");
  stopMix();
  player.pause();
  const attempt = mixOperation;
  const audio = getContext();
  await audio.resume();
  if (attempt === mixOperation) playback = schedule(audio, clips, audio.currentTime + .05);
});
byId("stop").onclick = stopMix;
player.addEventListener("play",stopMix);
function wav(buffer) {
  const size = buffer.length * 4;
  const bytes = new DataView(new ArrayBuffer(44 + size));
  function tag(offset,text) { for (let i=0;i<text.length;i++) bytes.setUint8(offset+i,text.charCodeAt(i)); }
  tag(0,"RIFF"); bytes.setUint32(4,36+size,true); tag(8,"WAVE"); tag(12,"fmt ");
  bytes.setUint32(16,16,true); bytes.setUint16(20,1,true); bytes.setUint16(22,2,true);
  bytes.setUint32(24,buffer.sampleRate,true); bytes.setUint32(28,buffer.sampleRate*4,true);
  bytes.setUint16(32,4,true); bytes.setUint16(34,16,true); tag(36,"data"); bytes.setUint32(40,size,true);
  const left = buffer.getChannelData(0), right = buffer.getChannelData(1);
  for (let i=0;i<buffer.length;i++) {
    bytes.setInt16(44+i*4,Math.max(-1,Math.min(1,left[i]))*32767,true);
    bytes.setInt16(46+i*4,Math.max(-1,Math.min(1,right[i]))*32767,true);
  }
  return new Blob([bytes],{type:"audio/wav"});
}
byId("export").onclick = () => action(byId("export"),notice,async () => {
  if (!clips.length) throw new Error("Add a clip first.");
  const duration = Math.max(...clips.map(clip => clip.start+clip.buffer.duration));
  if (duration>600) throw new Error("WAV exports can run up to 10 minutes.");
  const offline = new OfflineAudioContext(2,Math.ceil(duration*44100),44100);
  schedule(offline,clips,0);
  notice.textContent = "Rendering WAV…";
  const url = URL.createObjectURL(wav(await offline.startRendering()));
  const link = node("a");
  link.href = url;
  link.download = "Vuno arrangement.wav";
  link.click();
  setTimeout(() => URL.revokeObjectURL(url),30000);
  notice.textContent = "WAV rendered.";
});
renderClips();
boot().catch(error => lock(error.message, managed));
