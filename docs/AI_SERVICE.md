# AI Service

AI Service is installed in the Administration category. It shares the My AI source and release, with a separate tile, executable, private data and MyAccounts session. Choose an installed local model in AI Service and turn it on. My AI then offers the active model and queues local chat and code requests through this service. Cloud providers retain their settings.

The host must already have Ollama installed and running on loopback port 11434. Neither executable dependencies nor models are downloaded at app startup. The owner can install a model through My AI. An alternative loopback Ollama port can be set with AISERVICE_OLLAMA_URL in the supervised process environment. No external inference endpoint is accepted.

The service is off on fresh setup. It keeps the chosen model resident while enabled and retries if Ollama restarts. Turning it off cancels active and waiting work and unloads only the selected model. It does not stop the host Ollama service. Settings survive restarts and protected updates.

Requests use one FIFO queue, one active inference, at most64 waiting and8 per app/user. Individual requests expire after20 minutes. Prompts and replies are held only in memory; waiting work is cancelled on restart. This is a host-wide service, so the owner chooses the shared model. Model context above32768 tokens is rejected to bound resource use.

My AI receives a distinct server credential from Workspace. Other apps may receive owner-issued credentials, stored hashed and revealed once. Server calls POST /_internal/chat on the app's private loopback port with Authorization: Bearer KEY, X-AI-App and X-AI-User. The app is trusted to delegate only its authenticated user's identity. Never expose these credentials to browser code. Users may alternatively POST /aiservice/api/chat through their MyAccounts app session, matching CSRF/origin/account headers. Credential revocation and exact session logout cancel waiting and streaming work within about one second. My AI forwards native request cancellation.

ID aiservice; route /aiservice; canonical port9042; rebased port manager+52; private state Private/aiservice; package namespace aiservice. Update group myai covers both apps while keeping explicit omissions. Installer integration12 / manager118 is required. Installed Ollama/model libraries are externally owned and are not included in the app recovery copy. Existing standalone My AI data and model settings are preserved; shared memory retention overrides individual keep-alive preferences while connected.

No physical-phone, reboot, clean-host Ollama install or model-answer-quality claim is made by the synthetic tests. Native queue acceptance uses a dedicated Ollama protocol fixture and synthetic accounts. Public release and owner installation are recorded separately in compact history.

Explicit owner model benchmarks and other clients that call host Ollama directly operate outside the shared queue. Pause AI Service before benchmarking other models; external Ollama clients may evict its resident model, which the service attempts to reload. The service does not replace or supervise the externally installed Ollama process.
