# One Authenticator and OneAccount MFA

Open One Authenticator from Administration or the public signed-out link. Install
its PWA from a trusted HTTPS Workspace origin. Enrol from MyAccounts account MFA
settings, verify the phone code and save ten recovery codes separately. The
phone vault has its own passphrase and needs no OneAccount session or internet.

The managed app is a single embedded Windows executable, loopback-only on base
port+37 with private runtime directory, ordinary supervision, navigation,
update group, package verification and recovery ownership. The public proxy
allowlist serves only static shell/assets and health/version with GET/HEAD;
account and vault APIs are absent. Service-worker scope/cache remain app-local.
No token, key or secret is injected into its public shell.

The paired MyAccounts MFA component uses the existing identity database, pending
server-side enrolment, separate DPAPI-protected seed encryption key, atomic code
consumption, ten recovery digests, step-up and security events. Updated manager109
revalidates provider token/session proof for protected requests and every five
seconds on proxied streams, with a five-second provider timeout. An outage fails
closed. Old portal cookies must reauthenticate after MFA activation. This does
not replace an individual native app's independently issued API-key policy.

Installer integration6 and minimum manager109 prevent old controllers accepting
incomplete runtime/security integration. Preserve issuer, subject, client IDs,
credentials, libraries, app selections and recovery copies on upgrade. Enabling
MFA is the account owner's deliberate security operation; installation alone does
not enrol their phone or change an existing password.

Physical Android Chrome and iPhone Safari, reboot/service and clean-host tests
remain open. Browser software storage is not hardware-backed; TOTP is not phishing
resistant. Security events appear in account settings; email/push notification is
not promised. Supervised recovery when every factor is lost remains manual.
See the component security model and evidence for tested scope:
https://github.com/oneworkspacestudio/one-authenticator
