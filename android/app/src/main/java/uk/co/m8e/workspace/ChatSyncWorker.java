package uk.co.m8e.workspace;

import android.content.Context;
import androidx.annotation.NonNull;
import androidx.core.app.NotificationManagerCompat;
import androidx.work.*;
import java.util.*;
import java.util.concurrent.TimeUnit;
import org.json.*;

/** WorkManager owns finite background work, retries and reboot rescheduling. */
public final class ChatSyncWorker extends Worker {
  private volatile ChatConnection active;
  private static final Constraints NETWORK =
      new Constraints.Builder().setRequiredNetworkType(NetworkType.CONNECTED).build();

  public ChatSyncWorker(Context c, WorkerParameters p) {
    super(c, p);
  }

  static void schedule(Context c) {
    WorkManager.getInstance(c)
        .enqueueUniquePeriodicWork(
            "workspace-chat-poll",
            ExistingPeriodicWorkPolicy.KEEP,
            new PeriodicWorkRequest.Builder(ChatSyncWorker.class, 15, TimeUnit.MINUTES)
                .setConstraints(NETWORK)
                .addTag("workspace-chat")
                .build());
    enqueue(c, null, null, false);
  }

  static void enqueue(Context c, String topic, String reply, boolean urgent) {
    JSONObject device = SecureStore.get(c, "device");
    if (device == null) return;
    String owner = device.optString("id");
    String id = UUID.randomUUID().toString();
    Data.Builder d = new Data.Builder().putString("id", id).putString("device", owner);
    if (topic != null) d.putString("topic", topic);
    // WorkManager input is not encrypted. Put reply text only in Keystore storage.
    if (reply != null)
      try {
        if (!SecureStore.putForDevice(
            c, owner, "reply-" + id, new JSONObject().put("text", reply).put("device", owner)))
          return;
        d.putBoolean("reply", true);
      } catch (Exception e) {
        ChatNotifications.actionFailed(c, topic);
        return;
      }
    OneTimeWorkRequest.Builder b =
        new OneTimeWorkRequest.Builder(ChatSyncWorker.class)
            .setInputData(d.build())
            .setConstraints(NETWORK)
            .addTag("workspace-chat");
    if (urgent) b.setExpedited(OutOfQuotaPolicy.RUN_AS_NON_EXPEDITED_WORK_REQUEST);
    WorkManager.getInstance(c)
        .enqueueUniqueWork(
            reply == null ? "workspace-chat-sync" : "workspace-chat-reply-" + id,
            reply == null ? ExistingWorkPolicy.APPEND_OR_REPLACE : ExistingWorkPolicy.KEEP,
            b.build());
  }

  static void markRead(Context c, String topic, int seq) {
    JSONObject device = SecureStore.get(c, "device");
    if (device == null) return;
    Data d =
        new Data.Builder()
            .putString("device", device.optString("id"))
            .putString("topic", topic)
            .putInt("read", seq)
            .build();
    WorkManager.getInstance(c)
        .enqueue(
            new OneTimeWorkRequest.Builder(ChatSyncWorker.class)
                .setInputData(d)
                .setConstraints(NETWORK)
                .addTag("workspace-chat")
                .build());
  }

  @NonNull
  @Override
  public Result doWork() {
    Context c = getApplicationContext();
    JSONObject device = SecureStore.get(c, "device");
    if (device == null) return Result.success();
    String topic = getInputData().getString("topic"), id = getInputData().getString("id");
    boolean reply = getInputData().getBoolean("reply", false), publishing = false;
    String owner = device.optString("id"), expected = getInputData().getString("device");
    if (expected != null && !expected.equals(owner)) {
      if (reply) SecureStore.remove(c, "reply-" + id);
      return Result.success();
    }
    try {
      JSONObject bootstrap = WorkspaceApi.bootstrap(device);
      if (!SecureStore.isDevice(c, owner)) return Result.success();
      PushConfig.initialize(c, bootstrap.optJSONObject("push"));
      try (ChatConnection chat = new ChatConnection(c, bootstrap, device)) {
        active = chat;
        if (isStopped() || !SecureStore.isDevice(c, owner)) return Result.success();
        if (reply) {
          JSONObject r = SecureStore.get(c, "reply-" + id);
          if (r == null || !owner.equals(r.optString("device"))) return Result.success();
          publishing = true;
          chat.reply(topic, r.getString("text"), id);
          SecureStore.remove(c, "reply-" + id);
          if (SecureStore.isDevice(c, owner))
            NotificationManagerCompat.from(c).cancel(ChatNotifications.id(topic));
          return Result.success();
        }
        int read = getInputData().getInt("read", 0);
        if (read > 0) {
          chat.read(topic, read);
          NotificationManagerCompat.from(c).cancel(ChatNotifications.id(topic));
          return Result.success();
        }
        JSONArray subs = chat.subscriptions();
        if (subs == null) return Result.success();
        JSONObject history = SecureStore.get(c, "notification-history");
        if (history == null) history = new JSONObject();
        HashMap<String, String> names = new HashMap<>();
        for (int i = 0; i < subs.length(); i++) {
          JSONObject s = subs.getJSONObject(i);
          JSONObject p = s.optJSONObject("public");
          names.put(s.optString("topic"), p == null ? "Chat" : p.optString("fn", "Chat"));
        }
        int processed = 0;
        for (int i = 0; i < subs.length() && processed < 20 && !isStopped(); i++) {
          if (!SecureStore.isDevice(c, owner)) return Result.success();
          JSONObject s = subs.getJSONObject(i);
          String t = s.optString("topic");
          if (!ChatConnection.validTopic(t)) continue;
          int seq = s.optInt("seq"), seen = s.optInt("read");
          if (seq <= seen) {
            history.remove(t);
            NotificationManagerCompat.from(c).cancel(ChatNotifications.id(t));
            continue;
          }
          JSONObject previous = history.optJSONObject(t);
          int notified = previous == null ? seen : Math.max(seen, previous.optInt("last"));
          if (seq <= notified) continue;
          JSONObject acs = s.optJSONObject("acs");
          if (acs != null) {
            String mode = acs.optString("mode", "");
            if (mode.isEmpty()) {
              String want = acs.optString("want", ""), given = acs.optString("given", "");
              for (char permission : want.toCharArray())
                if (given.indexOf(permission) >= 0) mode += permission;
            }
            if (!mode.contains("R") || !mode.contains("P")) continue;
          }
          JSONArray incoming = chat.messages(t, Math.max(notified + 1, seq - 19)),
              messages = new JSONArray();
          for (int n = 0; n < incoming.length(); n++) {
            JSONObject m = incoming.getJSONObject(n);
            seq = Math.max(seq, m.optInt("seq"));
            if (chat.uid.equals(m.optString("from"))) continue;
            Object content = m.opt("content");
            String text =
                content instanceof String
                    ? (String) content
                    : (content instanceof JSONObject
                        ? ((JSONObject) content).optString("txt", "Attachment")
                        : "Attachment");
            if (text.length() > 4000) text = text.substring(0, 4000);
            long time = System.currentTimeMillis();
            try {
              time = java.time.Instant.parse(m.getString("ts")).toEpochMilli();
            } catch (Exception ignored) {
            }
            messages.put(
                new JSONObject()
                    .put("seq", m.optInt("seq"))
                    .put("from", m.optString("from"))
                    .put(
                        "sender",
                        names.getOrDefault(m.optString("from"), names.getOrDefault(t, "Chat")))
                    .put("text", text)
                    .put("time", time));
          }
          history.put(t, new JSONObject().put("last", seq));
          if (messages.length() > 0) {
            if (!SecureStore.putForDevice(c, owner, "notification-history", history))
              return Result.success();
            ChatNotifications.show(
                c, owner, t, names.getOrDefault(t, "Chat"), messages, seq - seen, seq);
          }
          processed++;
        }
        SecureStore.putForDevice(c, owner, "notification-history", history);
        if (SecureStore.isDevice(c, owner)) NotificationManagerCompat.from(c).cancel(42);
        return Result.success();
      }
    } catch (WorkspaceApi.SignInRequired e) {
      synchronized (SecureStore.class) {
        if (SecureStore.isDevice(c, owner)) {
          SecureStore.remove(c, "device");
          SecureStore.remove(c, "push-config");
          c.stopService(new android.content.Intent(c, WorkspacePushService.class));
          ChatNotifications.clear(c);
        }
      }
      return Result.success();
    } catch (Exception e) {
      // Retrying after a publish may duplicate a reply whose acknowledgement was lost.
      if (reply && publishing) {
        SecureStore.remove(c, "reply-" + id);
        if (SecureStore.isDevice(c, owner)) ChatNotifications.actionFailed(c, topic);
        return Result.failure();
      }
      if (getRunAttemptCount() < 3) return Result.retry();
      if (reply) {
        SecureStore.remove(c, "reply-" + id);
        ChatNotifications.actionFailed(c, topic);
      }
      return Result.failure();
    } finally {
      active = null;
    }
  }

  @Override
  public void onStopped() {
    ChatConnection chat = active;
    if (chat != null) chat.close();
    super.onStopped();
  }
}
