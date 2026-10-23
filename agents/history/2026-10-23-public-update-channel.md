# Public update channel / HTTP 403

User requests investigation, installer fix, push, release and actual local reinstall.
Read AGENTS, RTK and update-expiry/installer-handoff history. No agents/Actions.
Original unauthenticated API response confirms shared-IP hourly quota: HTTP403,
limit60, remaining0, used60. Never publish the user's IP or maintenance capability.

Official public clients use a publisher-promoted raw HTTPS index for all app
repositories and public asset URLs; authenticated/custom sources retain API.
No silent API fallback after a broken index. Existing SHA256/size/repository/tag
verification remains mandatory. Update page shows channel publication date.
Publisher script reads GCM/token in memory, excludes drafts, atomically writes
the channel. AGENTS release gate and channel guide require promotion after every
app/suite release. Missing repo lists produce an error, not up-to-date status.

Artifacts I:/Projects/tooling/release-check-403-20261023; only disposable C root
C:/Projects/.cache/release-check-403-20261023. Authoritative user I:/Workspace
preserved. Keep independently updated Drive4/Mail17/Slicer4/Clips6. Keep Umbrel
saved/autostart off and Jellyfin/Radarr/Prowlarr disabled. All native fixture
processes must be stopped before bounded cleanup. Delivery pending verification.
