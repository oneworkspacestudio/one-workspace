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

Publication preflight discovered remote main had advanced through public73 and
prepared74, including My3D/navigation/remote administration/Mail18. Initial
unpublished72 candidate is abandoned, never overwrite public72. Merged remote
changes; only release-status needed textual conflict resolution (kept both).
Next fix release75 uses verified public73 bundle plus verified Mail18 from the
current lockfile, retains all26 apps and remote controls. Local user is still
I:/Workspace71; the remote history's C:/tmp deployment is another machine.

First75 candidate b72cb84 exact native repair passed in125.72s. Final review
caught the channel fast path lacking the API path's publication-time sort;
fixed sorting a copy and added a delayed-publication ordering regression test.
Candidate b72 is not released. Final75 gets a new immutable source commit,
source75-final/release75-final and a fresh final-fixture NTFS installation.
Old task fixture was stopped by native test cleanup and retained until final
bounded cleanup. Do not treat first candidate receipt as final acceptance.
