# Reopen or continue setup

Refreshing or closing the setup webpage does not cancel an installation. Reopen
the installer on the same Windows account to see its saved status. The installer
shows the last recorded task, when it was recorded and the selected app files.

- **Installing:** another installer is still working. Follow its progress; a
  second copy will not start the same installation again.
- **Installed:** open or launch Workspace, and check each app's startup status.
  Retry apps that need attention without reinstalling healthy apps.
- **Interrupted:** choose **Continue installation**. Setup checks the recorded
  package hashes, keeps verified files and finishes the remaining work. An
  interrupted package is kept alongside its verified replacement.
- **Needs attention:** read the error and last task. Setup never labels a lost
  connection or incomplete operation as successful.

If setup does not recognise the location, expand **Already have One Workspace?**,
enter the same Workspace parent folder and choose **Check installation / continue
setup**. Do not choose a different folder just to bypass an error. If your storage
uses separate locations, reopen the original installer session or use its saved
private-data location with the existing recovery tools.

Installers older than this recovery feature did not save enough information to
continue every unfinished operation. The new installer still recognises their
completed installations and last saved task. Where a safe checkpoint is missing,
**Review recovery options** opens the existing protected reinstall flow. Read its
account/data retention notice before continuing; an incomplete older setup may
need a new sign-in setup while retaining its old state in a recovery folder.

If setup stopped while changing domain records, review the domain connection
through recovery. The installer does not blindly repeat an uncertain DNS/tunnel
operation. Existing records and files are retained.

Saved account setup details are encrypted and access-controlled on the local
NTFS settings volume. They are removed after successful setup. Browser refresh
keeps only the local setup access capability, not account passwords. The saved
status and app logs remain available for troubleshooting. Do not share private
logs or setup links publicly.
