# My 3D

My 3D is a native Windows app with an embedded browser CAD editor. Add it in setup or Updates, open it from your app collection and sign in with MyAccounts. Its data stays in the protected `my3d` app folder; removing the program retains projects.

Build models from shapes, turn cylinders into holes, group/union/cut, align and edit dimensions. Dimensioned rectangle, circle and closed-polyline profiles extrude on XY/XZ/YZ. Export STL/3MF for a slicer or STEP for CAD. A1 mini and A1 build-volume overlays warn about oversized models. Downloads do not transmit to a printer.

The desktop layout follows the supplied blue/teal phone-stand reference. Phones use a bottom Tools sheet and numeric editing; tablets preserve the workplane/sidebar workflow. Pointer handles are easier on desktop.

This first preview is a new Replicad/OpenCascade application, not a TweakCAD fork. The official TweakCAD repository was unavailable during audit. General sketch constraints/arcs, several interchange formats and advanced manufacturing checks are still unavailable. See the [app capability matrix](https://github.com/oneworkspacestudio/my-3d/blob/main/docs/ACCEPTANCE.md) before depending on those workflows.

The app runs at base port + 28, route `/my3d`, health `/healthz`, with the standard apphost OIDC and private storage contract. Server and WASM assets are embedded in My3D.exe. No end-user Node, Python or Docker is required. My Drive writes new editable `.my3d.json` copies using the signed-in user's existing Drive session. It does not share files publicly.
