package uk.co.m8e.workspace;

import android.content.*;

/** Restores an already enabled connection after reboot when Android permits it. */
public final class PushBootReceiver extends BroadcastReceiver {
  @Override
  public void onReceive(Context c, Intent i) {
    if (Intent.ACTION_BOOT_COMPLETED.equals(i.getAction())
        || Intent.ACTION_MY_PACKAGE_REPLACED.equals(i.getAction())) WorkspacePushService.start(c);
  }
}
