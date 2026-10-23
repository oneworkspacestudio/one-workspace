# Remove shared floating Workspace button

User requests removal across every app, push, release and local installation.
Base source8eac4d9/public71; current machine's actual owner is C:/tmp/workspace
on67 (the earlier handoff's I:/Workspace owner path does not exist here). Preserve
unrelated pkg/apphost/host.go changes for My3D. No agents, Actions, owner functional
tests or physical print. Use isolated named C NTFS cache workspace-nav-20261023;
source/build/download artifacts I:/Projects/tooling/navigation-no-float-20261023.

Central navigation.js used a floating toggle whenever no visible explicit header
trigger existed; header icon menus without the marker and scrolled headers caused
it to cover app bottom tabs. Remove toggle creation entirely, not conditional CSS.
Existing app menus, shared header trigger/event and Ctrl/Cmd+K remain. Center the
modal and constrain its scroll area on short phones. No app payload rebuild needed.

Target Workspace71.1 from verified public71 bundle. Preserve all app contracts and
payloads; only manager and suite release notice change. Tests/release/local update
still pending; completion evidence will supersede this paragraph.
