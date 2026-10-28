# My Shield in your Workspace

Open **My Shield** on your One Workspace page. It uses your existing OneAccount
session and starts with the other Workspace services. The tile and host security
evidence are available to the Workspace owner; ordinary members cannot read them.
You do not create a separate My Shield account or launch a separate Start shortcut.

This is a Monitor-only preview. AI remains disabled because neither candidate
completed qualification. Source gaps, including unavailable Windows event channels,
are shown explicitly. No driver, firewall permission, Defender exclusion or
privileged response action is installed by this integration.

The manager provides `/shield/`, registers its OIDC callback and private client,
and stores the app's data under `Private/shield`. The listener uses the Workspace
base port plus 32, independently of Mail and the earlier standalone preview.
Workspace sign-out revokes the captured app session and clears its cookies.
Normal Workspace backups and protected updates include the private app store.

My Shield's source repository remains private. Verified executable and licence
payloads are distributed inside the Workspace bundle. Private release metadata
and inaccessible private download links are not published in the public channel.

Earlier standalone previews are separate evaluations. Preserve their private data
and encrypted rollback material when transitioning; never silently merge local
accounts into OneAccount. The deployment receipt states whether old history was
imported or retained separately. Managed sign-in has no local owner-creation form.

The [native acceptance receipt](evidence/shield/native-integration-acceptance.json)
records the exact app binary tested with fresh MyAccounts identities: owner SSO,
member catalogue/API denial, disabled local setup, native logout revocation,
shared Apps navigation and phone layout. The fixture used protected NTFS storage;
no live owner identity, Drive library or account database was used as a test
fixture. This does not claim full model, MFA step-up, Windows VM or reboot acceptance.

The [protected upgrade receipt](evidence/shield/protected-upgrade-acceptance.json)
confirms a normal Workspace-only update from preview 94 adds Shield automatically,
preserves the existing owner and Accounts executable, respects other app exclusions
and completes verified recovery and startup checks in a separate fresh fixture.
