# My Finance

Open My Finance from Workspace's General apps and continue with OneAccount. Use
Actual Budget to create a GBP household budget, add manual accounts, import bank
statements and plan the month. Ledger transfers are bookkeeping entries; the app
cannot send money from a bank.

The Windows portable runtime is bundled and starts with Workspace. No separate
Node, Docker or WSL installation is needed. Keep Finance's private data and your
browser budget when uninstalling unless you deliberately want to remove them.

Members use their own OneAccount. Share an entire budget through Actual's User
Access menu; every member with access sees all accounts and transactions in it.
Keep private spending in another budget. Bank setup and server backups belong to
the Workspace owner. Use separate browser/OS profiles on shared devices.

Bank setup is optional. Enable Banking requires the owner's eligible application,
exact HTTPS callback and personal consent on the provider/bank page. Its native
integration remains experimental. Lunch Flow through SimpleFIN requires the user's
own supported subscription. Manual statement import remains available without a
bank connection. No bank has been connected automatically by this installation.

Finance can encrypt a native budget export or save credential-free server
snapshots, optionally to My Drive's One Workspace Backups/My Finance folder. Daily
snapshots are off until enabled with explicit Windows DPAPI recovery-key consent.
Keep your backup recovery and budget E2EE passphrases separately. Finance snapshots
exclude bank credentials; Workspace's full recovery is a separate protected copy.

This is a preview. Provider mock tests, local real OneAccount sharing, ledger and
clean backup/restore tests pass. Official provider sandbox/live banking, public
bank callbacks, complete keyboard-only editing and reboot/power-loss qualification
remain outstanding. Headless scheduled bank sync is deferred; use native Sync with
the budget open. Do not infer production qualification from installation.

- [Finance setup and UK banking guide](https://github.com/oneworkspacestudio/my-finance/blob/main/docs/OPEN_BANKING_SETUP_UK.md)
- [Backups and clean restoration](https://github.com/oneworkspacestudio/my-finance/blob/main/docs/MY_DRIVE_BACKUPS.md)
- [Full acceptance matrix and evidence](https://github.com/oneworkspacestudio/my-finance/blob/main/docs/TEST_REPORT.md)
- [Portable release and matching source](https://github.com/oneworkspacestudio/my-finance/releases/tag/v0.1.0-preview.1)
