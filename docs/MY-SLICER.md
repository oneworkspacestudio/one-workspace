# My Slicer in One Workspace

Workspace67 includes My Slicer 0.1.0-preview.4 from the public
[organization repository](https://github.com/oneworkspacestudio/my-slicer).
New default apps are added on update; explicit advanced-install exclusions remain.
My Slicer depends on MyAccounts and Drive, uses the normal `/slicer/` route,
OIDC client/session, shared logout, app switcher and central backup/update handling.
Private data is under `Private/slicer`; models and completed native archives also
copy as ordinary files to My Slicer folders in Drive while its browser is open.

The main executable is included in the complete installer. Engine setup separately
downloads the checksum-pinned official Bambu Studio runtime (474 MB download,
approximately570 MB installed), resolves its complete profiles and performs real
A1 and A1 mini slice checks before activation. It does not bundle the separate
proprietary networking plugin. No Docker, router change or new public port is needed.

Open My Slicer and choose **Add model** on Prepare. No project form is required.
STL and Bambu 3MF models appear on the plate before slicing. The 3MF view shows
saved parts and positions while preserving the original file and print settings.
A loading message and STL outline show progress; **Retry model preview** recovers
a failed preview without another upload.
Recent plates are kept automatically; **Save project** is optional. Slice the
plate, inspect its real layers on Preview and add a printer from Printers. The
guided setup uses a normal screen with Details, Trust device and Check connection
steps. Install the verified slicing engine once when prompted. Only confirmed 0.4 mm A1/A1 mini configurations are enabled.
The browser supports STL transforms and preserves imported native3MF projects for
selected-plate slicing. Layers and estimates come from native output. Printer
start needs an explicit physical check and confirmation; unknown outcomes never
trigger an automatic second print. Certificate trust and access codes are per device.

See the app's [setup guide](https://github.com/oneworkspacestudio/my-slicer/blob/main/docs/SETUP.md)
and [acceptance scope](https://github.com/oneworkspacestudio/my-slicer/blob/main/docs/ACCEPTANCE.md).
Real software slicing and local protocol fixtures are distinct from actual printer
acceptance. Physical A1/A1 mini/AMS printing, representative desktop modifier/
painting comparisons, clean second-machine/reboot/SCM and independent audit remain
explicit checks until recorded. This is an unsigned preview, built locally without
GitHub Actions. Existing ecosystem component limitations remain.

The header **Apps** button opens Workspace navigation. Its floating shortcut is
hidden while that control is available, keeping phone slicer controls clear.
