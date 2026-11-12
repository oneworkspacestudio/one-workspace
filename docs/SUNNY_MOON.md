# Sunny & Moon

Sunny & Moon gives children a gentle daily routine for toilet trips and bedtime.
Select it under **General** during Workspace setup or in **Updates**, then open
its tile from the **Apps** menu. Sign in with the MyAccounts account that will
use the routine. Each signed-in account sees its own settings and check-ins.

- **I went to the toilet** empties the two toilet bars, including when logging a
  school trip on arriving home. **Log an earlier visit** records its actual time.
- The body bars follow time. They are visual reminders. Children can go whenever
  they feel the need.
- The default routine is bedtime at **9 pm**, a toilet reminder at **8 pm** and
  wake-up at **7 am**. Adjust all three in **Settings**. The sleep figure shows
  the time available until wake-up. Sunny gets sleepy as bedtime approaches and
  starts each morning happy, regardless of the previous evening.
- **I went before bed** marks the evening reminder complete and empties both
  toilet bars. Reminders appear in the app while it is open.

The service is mounted at `/sunnymoon/`. Its MyAccounts callback is
`/sunnymoon/auth/callback`, and its supervised listener is loopback-only on
Workspace's port +57. Data stays under the Workspace private-data folder
`sunnymoon/` and participates in standard Workspace backup and app recovery.
Only the latest toilet trip and evening check-in are kept for each account.
Expired values clear on the next routine request. The app needs no external
services. DM Sans is bundled under the SIL Open Font License, with the licence
and original variable font in the distribution.
