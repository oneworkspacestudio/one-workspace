package uk.co.m8e.workspace;

import android.content.*;
import android.os.Bundle;
import androidx.core.app.RemoteInput;

public final class NotificationActions extends BroadcastReceiver {
  @Override
  public void onReceive(Context c, Intent i) {
    String topic = i.getStringExtra("topic");
    if (!ChatConnection.validTopic(topic) || !SecureStore.isDevice(c, i.getStringExtra("device")))
      return;
    if ("reply".equals(i.getAction())) {
      Bundle r = RemoteInput.getResultsFromIntent(i);
      CharSequence value = r == null ? null : r.getCharSequence("reply");
      if (value != null && !value.toString().trim().isEmpty() && value.length() <= 4000)
        ChatSyncWorker.enqueue(c, topic, value.toString(), true);
    } else if ("read".equals(i.getAction()))
      ChatSyncWorker.markRead(c, topic, i.getIntExtra("seq", 0));
  }
}
