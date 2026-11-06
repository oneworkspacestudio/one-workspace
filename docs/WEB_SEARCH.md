# Web Search

## Remote browser control

Open Web Search and choose **Open remote browser**. The live view controls the
same dedicated browser and private profile used for research. Click or tap its
screen to focus it, then type, paste, drag or scroll. Navigation, back/forward,
reload, tabs, full screen and a private text field are also available. Text in
that field is cleared after sending. This is the research browser's remote view,
not a remote desktop into the host PC; no RDP/VNC server or Docker is needed.

Remote control works with research turned off, even without model weights or
Windows OCR. Taking control pauses the queue and stops the local pilot model.
An active request is interrupted and must be submitted again. Waiting requests
stay queued and retain their normal ten-minute limit. Trusted API clients can
continue queuing while research is enabled. A handoff keeps the current page,
including an unfinished form or sign-in prompt, alive.

Choose **Release control** when finished. The browser closes cleanly to preserve
its profile; research resumes if it was enabled. Reconnect keeps the current
page when the connection is interrupted. A crashed browser can be reopened.
Closing the view or signing out ends control. An unreachable tab's lease expires
after 90 seconds without a heartbeat. Only the exact, live MyAccounts owner
session can view frames or send input; API search credentials cannot use these
routes. Frames stop on logout, expiry or account changes and are not recorded.

Use this view to complete legitimate sign-in and challenge prompts yourself.
Public HTTPS browsing restrictions still apply. Clipboard/OS shortcuts remain
blocked on the host; pasting sends only text you supply in the controlling
browser. The pilot does not bypass authentication or CAPTCHA.

## Research API

Web Search gives an LLM a private browser for web research. A small local model
looks at screenshots and uses browser mouse and keyboard input. Requests wait in
one queue and use a persistent, dedicated profile. The service does not use DOM
selectors, injected page scripts, page-text scraping or search-engine APIs. Windows screenshot text recognition helps the model read exact numbers and dates from the same pixels.

Choose **Web Search** in optional Administration apps, open its Workspace tile,
install its local model, then turn it on. Browser and model-server executables are
portable package payloads. No Docker, WSL, Python, Ollama or separate system
service is required. The browser retains Chromium's renderer sandbox. This is a
self-contained process runtime, not a virtual machine or an OS container.

The local model is Microsoft Fara1.5-4B, quantized Q4_K_M. Its text weights and
vision projector take 3,383,083,424 bytes, plus browser/server binaries, working
memory and the profile. Model setup is an explicit owner action, verified against
the pinned SHA256 values in source. Ordinary startup never downloads dependencies.
Models live in private app storage and survive program replacement. CPU fallback
is available; Vulkan acceleration needs a compatible GPU driver. An English Windows OCR language must be available (English UK and US are verified on this PC); activation reports a clear error when it is missing. Reliability and
latency depend on hardware, site changes and blocking pages.

Create a named API credential in the console. It is revealed once and stored
hashed. Credentials are for trusted server apps: they share the owner's dedicated
browser profile, including any sessions deliberately established there. Never put
one in browser JavaScript, logs or a URL. This service is owner-operated; Workspace
members cannot open its browser console or issue grants.

```http
POST /websearch/api/client/search
Authorization: Bearer YOUR_SERVER_CREDENTIAL
Content-Type: application/json

{"query":"Find the official documentation for Python pathlib.Path.is_file"}
```

A successful submission returns HTTP202 and an unguessable request `id`. Poll
`GET /websearch/api/client/jobs/ID` with the same credential. A completed response
contains a retrieval status in `result.answer`, browser-observed `result.sources` (each with `url`, `title` and screenshot-recognized `text`), `result.steps` and
`result.model`. The small model pilots navigation; the calling LLM should formulate its answer from the returned source text. The pilot's generated factual summary is deliberately not used as evidence. `DELETE` to that same job URL cancels it. Each credential can read
and cancel only its own requests. Revocation cancels its active and queued work.
Other native and administrative routes require the normal MyAccounts owner
session and CSRF checks. Browser-origin calls cannot use the server exemption.

The queue allows 64 active/waiting requests globally and eight per credential or
owner. Waiting work expires after ten minutes; each run has five minutes and at
most 32 actions. Responses may be completed, blocked, failed, cancelled, expired
or interrupted. Restarted requests are interrupted rather than replayed against
a potentially changed browser. Queries/results expire after 24 hours; at most 1,000 recent completed requests and eight source pages per request are retained. They are
private SQLite data; screenshots and model reasoning are not persisted.

Sign-in pages, CAPTCHA and security challenges pause research. The service does
not bypass challenges. Downloads are disabled. Browser network traffic goes
through a public-HTTPS proxy that rejects private, loopback, metadata and reserved
addresses and pins validated DNS results. Arbitrary code and DOM tools do not
exist. The model's read-only policy also excludes purchases, account changes,
messages and uploads; model-based navigation is not a formal proof that every
site click is side-effect free. Do not establish sensitive account sessions for
untrusted API clients.

Managed contract: `websearch`, `/websearch`, canonical9043 / manager+53;
`WebSearch.exe`; `websearch-*` package ownership; private `websearch` state;
MyAccounts owner binding and exact logout; optional Administration category;
installer integration13 / minimum Workspace120. Windows Job Objects own every
browser/model child, including forced stops. Profile/settings/queue are included
in the managed app recovery scope; model libraries are preserved in place.

Implementation and acceptance status is recorded in `agents/history` and the
release evidence. Do not infer publication, owner activation, physical-phone or
reboot acceptance from source availability.
