# Microsoft Visual C++ runtime

Calendar and Contacts' shared Stalwart engine imports Microsoft C++ x64 runtime
libraries. Workspace supplies the release CRT files beside its private engine;
no administrator privilege or global runtime installation is required.

Copyright Microsoft Corporation. All rights reserved. Microsoft binaries are
proprietary redistributable components, not covered by Workspace's source license.
They remain unmodified and carry Microsoft's Authenticode signatures.

Source: Microsoft's official https://aka.ms/vc14/vc_redist.x64.exe, version
14.51.36247.0. The release provenance receipt records the downloaded installer
and each distributed DLL's SHA256. No files were copied from Windows System32.

Microsoft documentation and distribution terms:
- https://learn.microsoft.com/cpp/windows/redistributing-visual-cpp-files
- https://learn.microsoft.com/cpp/windows/determining-which-dlls-to-redistribute
- https://visualstudio.microsoft.com/license-terms/

App-local runtimes are serviced by Workspace package updates. Keep Workspace
updated to receive runtime security fixes; these private copies are not registered
as a system-wide Visual C++ Redistributable installation.
