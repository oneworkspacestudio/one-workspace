package uk.co.m8e.workspace;

import android.Manifest;
import android.app.*;
import android.content.*;
import android.content.pm.PackageManager;
import android.os.Build;
import androidx.core.app.*;
import androidx.core.app.Person;
import androidx.core.app.RemoteInput;
import androidx.core.content.ContextCompat;
import org.json.*;

final class ChatNotifications {
  static final String CHANNEL = "chat_messages", GROUP = "workspace_chat";

  static void channels(Context c) {
    NotificationChannel ch =
        new NotificationChannel(CHANNEL, "Chat messages", NotificationManager.IMPORTANCE_HIGH);
    ch.setDescription("Messages from your linked Workspace Chat");
    ch.setShowBadge(true);
    ch.enableVibration(true);
    ch.setLockscreenVisibility(Notification.VISIBILITY_PRIVATE);
    c.getSystemService(NotificationManager.class).createNotificationChannel(ch);
  }

  static boolean permitted(Context c) {
    return Build.VERSION.SDK_INT < 33
        || ContextCompat.checkSelfPermission(c, Manifest.permission.POST_NOTIFICATIONS)
            == PackageManager.PERMISSION_GRANTED;
  }

  static int id(String topic) {
    return 1000 + (topic.hashCode() & 0x3fffffff);
  }

  static PendingIntent open(Context c, String topic) {
    JSONObject d = SecureStore.get(c, "device");
    String owner = d == null ? "" : d.optString("id");
    Intent i =
        new Intent(c, MainActivity.class)
            .setData(android.net.Uri.parse("oneworkspace://open/" + owner + "/" + topic))
            .putExtra("device", owner)
            .putExtra("openChat", true)
            .putExtra("topic", topic)
            .addFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP | Intent.FLAG_ACTIVITY_SINGLE_TOP);
    return PendingIntent.getActivity(
        c, id(topic), i, PendingIntent.FLAG_UPDATE_CURRENT | PendingIntent.FLAG_IMMUTABLE);
  }

  static void pending(Context c) {
    if (!permitted(c)) return;
    Notification n =
        new NotificationCompat.Builder(c, CHANNEL)
            .setSmallIcon(R.drawable.ic_chat)
            .setContentTitle("Chat")
            .setContentText("New message")
            .setContentIntent(open(c, ""))
            .setAutoCancel(true)
            .setOnlyAlertOnce(true)
            .setVisibility(NotificationCompat.VISIBILITY_PRIVATE)
            .setCategory(NotificationCompat.CATEGORY_MESSAGE)
            .build();
    NotificationManagerCompat.from(c).notify(42, n);
  }

  static void show(
      Context c,
      String owner,
      String topic,
      String title,
      JSONArray messages,
      int unread,
      int latest)
      throws Exception {
    synchronized (SecureStore.class) {
      if (!SecureStore.isDevice(c, owner)) return;
      if (!permitted(c) || MainActivity.viewing(topic)) return;
      boolean previews = c.getSharedPreferences("workspace", 0).getBoolean("previews", false);
      Person self = new Person.Builder().setName("You").build();
      NotificationCompat.MessagingStyle style =
          new NotificationCompat.MessagingStyle(self)
              .setConversationTitle(title)
              .setGroupConversation(topic.startsWith("grp"));
      int seq = latest;
      Intent shortcutIntent =
          new Intent(c, MainActivity.class)
              .setAction(Intent.ACTION_VIEW)
              .putExtra("device", owner)
              .putExtra("openChat", true)
              .putExtra("topic", topic);
      androidx.core.content.pm.ShortcutManagerCompat.pushDynamicShortcut(
          c,
          new androidx.core.content.pm.ShortcutInfoCompat.Builder(c, "chat-" + topic)
              .setShortLabel(title.isEmpty() ? "Chat" : title)
              .setIntent(shortcutIntent)
              .setLongLived(true)
              .setIcon(
                  androidx.core.graphics.drawable.IconCompat.createWithResource(
                      c, R.drawable.brand_my_chat))
              .build());
      for (int i = 0; i < messages.length(); i++) {
        JSONObject m = messages.getJSONObject(i);
        seq = Math.max(seq, m.optInt("seq"));
        Person sender =
            new Person.Builder()
                .setName(m.optString("sender", "Chat"))
                .setKey(m.optString("from"))
                .build();
        style.addMessage(
            previews ? m.optString("text", "New message") : "New message",
            m.optLong("time", System.currentTimeMillis()),
            sender);
      }
      NotificationCompat.Builder b =
          new NotificationCompat.Builder(c, CHANNEL)
              .setSmallIcon(R.drawable.ic_chat)
              .setColor(0xff345de5)
              .setContentTitle(title)
              .setContentText("New message")
              .setStyle(style)
              .setCategory(NotificationCompat.CATEGORY_MESSAGE)
              .setContentIntent(open(c, topic))
              .setAutoCancel(true)
              .setGroup(GROUP)
              .setNumber(Math.max(1, unread))
              .setOnlyAlertOnce(false)
              .setVisibility(NotificationCompat.VISIBILITY_PRIVATE)
              .setShortcutId("chat-" + topic);
      b.setPublicVersion(
          new NotificationCompat.Builder(c, CHANNEL)
              .setSmallIcon(R.drawable.ic_chat)
              .setContentTitle("One Workspace Chat")
              .setContentText("New message")
              .build());
      Intent reply =
          new Intent(c, NotificationActions.class)
              .setAction("reply")
              .setData(android.net.Uri.parse("oneworkspace://reply/" + owner + "/" + topic))
              .putExtra("device", owner)
              .putExtra("topic", topic)
              .putExtra("seq", seq);
      PendingIntent rp =
          PendingIntent.getBroadcast(
              c, id(topic), reply, PendingIntent.FLAG_UPDATE_CURRENT | PendingIntent.FLAG_MUTABLE);
      b.addAction(
          new NotificationCompat.Action.Builder(R.drawable.ic_chat, "Reply", rp)
              .addRemoteInput(new RemoteInput.Builder("reply").setLabel("Reply").build())
              .setAllowGeneratedReplies(true)
              .setSemanticAction(NotificationCompat.Action.SEMANTIC_ACTION_REPLY)
              .setAuthenticationRequired(true)
              .build());
      Intent read =
          new Intent(c, NotificationActions.class)
              .setAction("read")
              .setData(android.net.Uri.parse("oneworkspace://read/" + owner + "/" + topic))
              .putExtra("device", owner)
              .putExtra("topic", topic)
              .putExtra("seq", seq);
      b.addAction(
          new NotificationCompat.Action.Builder(
                  0,
                  "Mark as read",
                  PendingIntent.getBroadcast(
                      c,
                      id(topic) + 1,
                      read,
                      PendingIntent.FLAG_UPDATE_CURRENT | PendingIntent.FLAG_IMMUTABLE))
              .setSemanticAction(NotificationCompat.Action.SEMANTIC_ACTION_MARK_AS_READ)
              .setAuthenticationRequired(true)
              .build());
      NotificationManagerCompat.from(c).cancel(42);
      NotificationManagerCompat.from(c).notify(id(topic), b.build());
    }
  }

  static void actionFailed(Context c, String topic) {
    if (!permitted(c)) return;
    NotificationManagerCompat.from(c)
        .notify(
            id(topic),
            new NotificationCompat.Builder(c, CHANNEL)
                .setSmallIcon(R.drawable.ic_chat)
                .setContentTitle("Check your reply in Chat")
                .setContentText(
                    "Delivery could not be confirmed. Open Chat before sending it again.")
                .setContentIntent(open(c, topic))
                .setAutoCancel(true)
                .setVisibility(NotificationCompat.VISIBILITY_PRIVATE)
                .build());
  }

  static void clear(Context c) {
    NotificationManagerCompat.from(c).cancelAll();
    java.util.List<String> ids = new java.util.ArrayList<>();
    for (androidx.core.content.pm.ShortcutInfoCompat s :
        androidx.core.content.pm.ShortcutManagerCompat.getDynamicShortcuts(c)) ids.add(s.getId());
    androidx.core.content.pm.ShortcutManagerCompat.disableShortcuts(
        c, ids, "Reconnect One Workspace");
    androidx.core.content.pm.ShortcutManagerCompat.removeAllDynamicShortcuts(c);
    androidx.core.content.pm.ShortcutManagerCompat.removeLongLivedShortcuts(c, ids);
  }
}
