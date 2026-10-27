package uk.co.m8e.workspace;

import android.content.Context;
import java.io.*;
import java.util.*;
import java.util.concurrent.*;
import okhttp3.*;
import org.json.*;

/** Short-lived background Tinode connection. Server membership remains authoritative. */
final class ChatConnection implements AutoCloseable {
  private final Map<String, CompletableFuture<JSONObject>> awaiting = new ConcurrentHashMap<>();
  private final Map<String, JSONArray> data = new ConcurrentHashMap<>();
  private final BlockingQueue<Throwable> errors = new LinkedBlockingQueue<>();
  private WebSocket socket;
  private final CompletableFuture<Void> opened = new CompletableFuture<>();
  final String uid;

  ChatConnection(Context c, JSONObject b, JSONObject device) throws Exception {
    if (device == null) throw new WorkspaceApi.SignInRequired();
    uid = device.getString("uid");
    String chat = device.getString("chat");
    String url =
        "wss://"
            + android.net.Uri.parse(chat).getEncodedAuthority()
            + android.net.Uri.parse(chat).getEncodedPath()
            + "/api/native/socket";
    socket =
        WorkspaceApi.HTTP.newWebSocket(
            new Request.Builder()
                .url(url)
                .header("Authorization", "Bearer " + device.getString("token"))
                .header("Origin", WorkspaceApi.origin(chat))
                .build(),
            new WebSocketListener() {
              @Override
              public void onOpen(WebSocket s, Response r) {
                opened.complete(null);
              }

              @Override
              public void onFailure(WebSocket s, Throwable t, Response r) {
                opened.completeExceptionally(t);
                for (CompletableFuture<JSONObject> f : awaiting.values())
                  f.completeExceptionally(t);
              }

              @Override
              public void onMessage(WebSocket s, String value) {
                try {
                  if (value.length() > 1048576) throw new IOException("Message too large");
                  JSONObject p = new JSONObject(value);
                  if (p.has("data")) {
                    JSONObject m = p.getJSONObject("data");
                    JSONArray list = data.get(m.optString("topic"));
                    if (list != null)
                      synchronized (list) {
                        list.put(m);
                      }
                    return;
                  }
                  if (p.has("meta")) {
                    JSONObject m = p.getJSONObject("meta");
                    CompletableFuture<JSONObject> f = awaiting.get(m.optString("id"));
                    if (f != null) f.complete(m);
                  }
                  if (p.has("ctrl")) {
                    JSONObject ctrl = p.getJSONObject("ctrl");
                    CompletableFuture<JSONObject> f = awaiting.get(ctrl.optString("id"));
                    if (f != null) {
                      if (ctrl.optInt("code", 500) >= 400)
                        f.completeExceptionally(
                            new IOException("Chat permission or request failed"));
                      else f.complete(ctrl);
                    }
                  }
                } catch (Exception e) {
                  for (CompletableFuture<JSONObject> f : awaiting.values())
                    f.completeExceptionally(e);
                }
              }
            });
    try {
      opened.get(15, TimeUnit.SECONDS);
      JSONObject hi =
          new JSONObject()
              .put("ver", "0.25")
              .put("ua", "OneWorkspace/Android")
              .put("platf", "android")
              .put("bkg", true)
              .put("lang", Locale.getDefault().getLanguage());
      String token = PushConfig.token(c);
      if (!token.isEmpty()) hi.put("dev", token);
      call("hi", hi);
      call(
          "login",
          new JSONObject().put("scheme", "workspace").put("secret", b.getString("ticket")));
    } catch (Exception e) {
      socket.cancel();
      throw e;
    }
  }

  private JSONObject call(String op, JSONObject payload) throws Exception {
    String id = UUID.randomUUID().toString();
    payload.put("id", id);
    CompletableFuture<JSONObject> future = new CompletableFuture<>();
    awaiting.put(id, future);
    try {
      if (!socket.send(new JSONObject().put(op, payload).toString()))
        throw new IOException("Chat disconnected");
      return future.get(12, TimeUnit.SECONDS);
    } finally {
      awaiting.remove(id);
    }
  }

  JSONArray subscriptions() throws Exception {
    JSONObject r =
        call(
            "sub",
            new JSONObject().put("topic", "me").put("get", new JSONObject().put("what", "sub")));
    if (r.has("sub")) return r.getJSONArray("sub");
    return call("get", new JSONObject().put("topic", "me").put("what", "sub")).optJSONArray("sub");
  }

  JSONArray messages(String topic, int since) throws Exception {
    if (!validTopic(topic)) throw new IOException("Invalid conversation");
    JSONArray out = new JSONArray();
    data.put(topic, out);
    call("sub", new JSONObject().put("topic", topic));
    call(
        "get",
        new JSONObject()
            .put("topic", topic)
            .put("what", "data")
            .put("data", new JSONObject().put("since", since).put("limit", 20)));
    data.remove(topic);
    return out;
  }

  void reply(String topic, String text, String id) throws Exception {
    if (!validTopic(topic) || text.trim().isEmpty() || text.length() > 4000)
      throw new IOException("Invalid reply");
    call("sub", new JSONObject().put("topic", topic));
    call(
        "pub",
        new JSONObject()
            .put("topic", topic)
            .put("noecho", true)
            .put("head", new JSONObject().put("workspace-native-id", id))
            .put("content", text));
  }

  void read(String topic, int seq) throws Exception {
    if (validTopic(topic) && seq > 0)
      socket.send(
          new JSONObject()
              .put("note", new JSONObject().put("topic", topic).put("what", "read").put("seq", seq))
              .toString());
  }

  static boolean validTopic(String t) {
    return t != null && t.matches("(?:usr|grp)[A-Za-z0-9_-]{1,100}");
  }

  @Override
  public void close() {
    for (CompletableFuture<JSONObject> f : awaiting.values())
      f.completeExceptionally(new IOException("Background task closed"));
    if (socket != null && !socket.close(1000, "Background sync finished")) socket.cancel();
  }
}
