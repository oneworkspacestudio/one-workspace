# An installed app did not start

The installer lists file installation and startup checks separately. A red
**Start Chat**, **Start Calendar** or **Start Contacts** row means that app did
not become ready. Other running apps remain available; it does not mean your
files were deleted or that every app needs reinstalling.

Read the explanation beside the failed row. Choose **Save startup report** to
retain the app versions, check results and times. The report contains no account
credentials or raw app logs. The displayed log paths identify more detailed
diagnostics inside your Workspace Private folder. Do not post private logs
publicly without reviewing their contents.

Choose **Retry apps that need attention** to retry failed apps while keeping
healthy apps running. You can also retry from Workspace's App status page.

Calendar and Contacts share one background service. One service failure can
affect both apps. Windows versions of that engine require Microsoft Visual C++
x64 libraries. New installers supply verified app-local release libraries;
older installers relied on libraries already present on the PC. Update using the
latest installer and the same Workspace folder to retain your installation.
Use its protected repair flow rather than deleting your Private folder.

Chat uses a separate browser address and an existing MyAccounts client. Its
startup now allows three minutes for a cold start. A **chat-sign-in** failure
identifies account connection preparation; **process-exited** identifies a child
process that closed; **app-readiness** identifies a failed readiness check.
**windows-memory**, **access-denied**, **file-missing** and **windows-runtime**
provide specific recovery guidance. Raw errors are retained in the protected
app log, not shown in the browser.

The screenshot alone cannot establish a remote machine's exact error. A saved
startup report and the relevant private log distinguish runtime dependencies,
resource pressure, security software and identity connection failures.
