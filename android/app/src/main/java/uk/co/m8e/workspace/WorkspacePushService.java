package uk.co.m8e.workspace;

import android.app.*;
import android.content.*;
import android.content.pm.ServiceInfo;
import android.os.*;
import androidx.core.app.NotificationCompat;
import androidx.core.content.ContextCompat;
import java.util.concurrent.*;
import okhttp3.*;
import org.json.*;

/** Visible self hosted remote-messaging connection, with bounded reconnect backoff. */
public final class WorkspacePushService extends Service {
  private static final String CHANNEL = "workspace_connection";
  private static final int NOTICE = 73;
  private final ExecutorService connection = Executors.newSingleThreadExecutor();
  private final ExecutorService events = Executors.newSingleThreadExecutor();
  private volatile boolean running;
  private volatile WebSocket socket;
  private volatile CountDownLatch waiting;

  static void start(Context c) {
    if (SecureStore.get(c, "device") == null
        || !ChatNotifications.permitted(c)
        || !c.getSharedPreferences("workspace", 0).getBoolean("backgroundPush", true)) return;
    try {
      ContextCompat.startForegroundService(c, new Intent(c, WorkspacePushService.class));
    } catch (RuntimeException restricted) {
    }
  }

  @Override
  public void onCreate() {
    super.onCreate();
    NotificationChannel ch =
        new NotificationChannel(
            CHANNEL, "Private Chat connection", NotificationManager.IMPORTANCE_LOW);
    ch.setShowBadge(false);
    ch.setDescription("Keeps your self hosted Chat notification connection available");
    getSystemService(NotificationManager.class).createNotificationChannel(ch);
  }

  private Notification notice(String text) {
    PendingIntent open =
        PendingIntent.getActivity(
            this,
            73,
            new Intent(this, MainActivity.class),
            PendingIntent.FLAG_UPDATE_CURRENT | PendingIntent.FLAG_IMMUTABLE);
    return new NotificationCompat.Builder(this, CHANNEL)
        .setSmallIcon(R.drawable.ic_chat)
        .setContentTitle("One Workspace")
        .setContentText(text)
        .setContentIntent(open)
        .setOngoing(true)
        .setSilent(true)
        .setOnlyAlertOnce(true)
        .setCategory(NotificationCompat.CATEGORY_SERVICE)
        .build();
  }

  private void status(String text) {
    getSystemService(NotificationManager.class).notify(NOTICE, notice(text));
  }

  @Override
  public int onStartCommand(Intent intent, int flags, int id) {
    if (SecureStore.get(this, "device") == null
        || !ChatNotifications.permitted(this)
        || !getSharedPreferences("workspace", 0).getBoolean("backgroundPush", true)) {
      stopSelf();
      return START_NOT_STICKY;
    }
    if (Build.VERSION.SDK_INT >= 34)
      startForeground(
          NOTICE,
          notice("Connecting to private Chat notifications"),
          ServiceInfo.FOREGROUND_SERVICE_TYPE_REMOTE_MESSAGING);
    else startForeground(NOTICE, notice("Connecting to private Chat notifications"));
    if (!running) {
      running = true;
      connection.execute(this::listen);
    }
    return START_STICKY;
  }

  private void listen() {
    int delay = 2;
    while (running) {
      JSONObject d = SecureStore.get(this, "device");
      if (d == null) break;
      String owner = d.optString("id");
      CountDownLatch closed = new CountDownLatch(1);
      waiting = closed;
      try {
        if (!ChatNotifications.permitted(this)) break;
        JSONObject b = WorkspaceApi.bootstrap(d);
        if (!running || !SecureStore.isDevice(this, owner)) break;
        JSONObject options = b.optJSONObject("push");
        if (options == null
            || !options.optBoolean("configured")
            || !"workspace".equals(options.optString("provider")))
          throw new java.io.IOException("My Push is reconnecting");
        PushConfig.initialize(this, options);
        JSONObject p = SecureStore.get(this, "push-config");
        if (p == null
            || !owner.equals(p.optString("owner"))
            || p.optLong("expires") <= System.currentTimeMillis() / 1000)
          throw new java.io.IOException("Private push not configured");
        String url = p.getString("url").replaceAll("/+$", "");
        WorkspaceApi.origin(url);
        android.net.Uri u = android.net.Uri.parse(url);
        String ws =
            "wss://" + u.getEncodedAuthority() + u.getEncodedPath() + "/connection/websocket";
        socket =
            WorkspaceApi.HTTP.newWebSocket(
                new Request.Builder().url(ws).build(),
                new WebSocketListener() {
                  @Override
                  public void onOpen(WebSocket s, Response r) {
                    try {
                      s.send(
                          new JSONObject()
                              .put("id", 1)
                              .put(
                                  "connect",
                                  new JSONObject()
                                      .put("token", p.getString("token"))
                                      .put("name", "one-workspace-android"))
                              .toString());
                    } catch (Exception e) {
                      s.cancel();
                      closed.countDown();
                    }
                  }

                  @Override
                  public void onMessage(WebSocket s, String text) {
                    if (text.length() > 65536) {
                      s.cancel();
                      closed.countDown();
                      return;
                    }
                    for (String line : text.split("\n"))
                      try {
                        JSONObject x = new JSONObject(line);
                        if (x.length() == 0) {
                          s.send("{}");
                          continue;
                        }
                        JSONObject push = x.optJSONObject("push");
                        if (x.has("error") || push != null && push.has("disconnect")) {
                          s.close(1000, "Reconnect");
                          closed.countDown();
                          continue;
                        }
                        if (x.has("connect")) {
                          status("Private Chat notifications connected");
                          events.execute(() -> pending(owner, p));
                          ChatSyncWorker.enqueue(WorkspacePushService.this, null, null, true);
                        }
                        if (push != null) {
                          JSONObject pub = push.optJSONObject("pub");
                          if (pub != null) {
                            JSONObject event = pub.optJSONObject("data");
                            if (event != null) events.execute(() -> receive(owner, p, event));
                          }
                        }
                      } catch (Exception ignored) {
                      }
                  }

                  @Override
                  public void onClosed(WebSocket s, int code, String reason) {
                    closed.countDown();
                  }

                  @Override
                  public void onFailure(WebSocket s, Throwable t, Response r) {
                    closed.countDown();
                  }
                });
        long began = System.currentTimeMillis();
        closed.await(25, TimeUnit.MINUTES);
        delay = System.currentTimeMillis() - began > 60000 ? 2 : Math.min(120, delay * 2);
      } catch (WorkspaceApi.SignInRequired e) {
        status("Open Chat to renew your phone connection");
        break;
      } catch (Exception e) {
        delay = Math.min(120, delay * 2);
      } finally {
        WebSocket s = socket;
        socket = null;
        if (s != null) s.cancel();
      }
      if (running) {
        status("Chat connection interrupted · retrying automatically");
        try {
          Thread.sleep(
              (delay + java.util.concurrent.ThreadLocalRandom.current().nextInt(3)) * 1000L);
        } catch (InterruptedException e) {
          Thread.currentThread().interrupt();
          break;
        }
      }
    }
    running = false;
    stopSelf();
  }

  private void pending(String owner, JSONObject p) {
    try {
      JSONObject q =
          WorkspaceApi.json(
              p.getString("url") + "/api/device/pending", "GET", null, p.getString("token"), false);
      JSONArray list = q.optJSONArray("events");
      if (list != null)
        for (int i = 0; i < list.length(); i++) receive(owner, p, list.getJSONObject(i));
    } catch (Exception ignored) {
    }
  }

  private void receive(String owner, JSONObject p, JSONObject event) {
    if (!running || !SecureStore.isDevice(this, owner)) return;
    try {
      String id = event.getString("id"), what = event.optString("what");
      if (!id.matches("[a-f0-9]{48}") || !what.equals("sync") && !what.equals("check")) return;
      synchronized (SecureStore.class) {
        if (!running || !SecureStore.isDevice(this, owner)) return;
        JSONObject seen = SecureStore.get(this, "push-seen");
        if (seen == null || !owner.equals(seen.optString("owner")))
          seen = new JSONObject().put("owner", owner).put("ids", new JSONArray());
        JSONArray ids = seen.getJSONArray("ids");
        boolean found = false;
        for (int i = 0; i < ids.length(); i++) if (id.equals(ids.optString(i))) found = true;
        if (!found) {
          if ("sync".equals(what) && !MainActivity.foreground) ChatNotifications.pending(this);
          ChatSyncWorker.enqueue(this, null, null, true);
          JSONArray next = new JSONArray();
          for (int i = Math.max(0, ids.length() - 99); i < ids.length(); i++) next.put(ids.get(i));
          next.put(id);
          seen.put("ids", next);
          SecureStore.putForDevice(this, owner, "push-seen", seen);
        }
      }
      if (SecureStore.isDevice(this, owner))
        WorkspaceApi.json(
            p.getString("url") + "/api/device/ack",
            "POST",
            new JSONObject().put("id", id).toString(),
            p.getString("token"),
            false);
    } catch (Exception ignored) {
    }
  }

  @Override
  public void onDestroy() {
    running = false;
    WebSocket s = socket;
    if (s != null) s.cancel();
    CountDownLatch w = waiting;
    if (w != null) w.countDown();
    connection.shutdownNow();
    events.shutdownNow();
    stopForeground(STOP_FOREGROUND_REMOVE);
    super.onDestroy();
  }

  @Override
  public IBinder onBind(Intent intent) {
    return null;
  }
}
