package uk.co.m8e.workspace;

import android.content.Context;
import android.net.Uri;
import android.webkit.CookieManager;
import java.io.IOException;
import java.util.concurrent.TimeUnit;
import okhttp3.*;
import org.json.*;

final class WorkspaceApi {
  static final OkHttpClient HTTP =
      new OkHttpClient.Builder()
          .connectTimeout(12, TimeUnit.SECONDS)
          .readTimeout(20, TimeUnit.SECONDS)
          .followRedirects(false)
          .followSslRedirects(false)
          .build();

  static String origin(String raw) {
    Uri u = Uri.parse(raw);
    if (!"https".equals(u.getScheme()) || u.getHost() == null || u.getUserInfo() != null)
      throw new IllegalArgumentException("Use a secure https:// Workspace address");
    return "https://" + u.getEncodedAuthority();
  }

  static String home(Context c) {
    return c.getSharedPreferences("workspace", 0).getString("home", "");
  }

  static JSONObject discovery(String home) throws Exception {
    try (Response r =
        HTTP.newCall(new Request.Builder().url(home + "/.well-known/one-workspace.json").build())
            .execute()) {
      if (r.isSuccessful() && r.body() != null) return decode(r.body());
      if (r.code() != 404) throw new IOException("Workspace connection metadata unavailable");
    }
    // Older Workspace releases expose the same authenticated launch routes.
    // Resolve only redirects returned directly by the chosen HTTPS server.
    JSONObject d = new JSONObject().put("schema", 0).put("workspace", home);
    JSONArray origins = new JSONArray().put(home);
    String issuer = redirect(home + "/auth/login");
    d.put("issuer", origin(issuer));
    origins.put(origin(issuer));
    String chat = "";
    try {
      chat = redirect(home + "/chat/");
      origins.put(origin(chat));
    } catch (Exception ignored) {
    }
    try {
      origins.put(origin(redirect(home + "/vault/")));
    } catch (Exception ignored) {
    }
    return d.put("chat", chat.replaceAll("/+$", "")).put("origins", origins);
  }

  private static String redirect(String url) throws Exception {
    try (Response r = HTTP.newCall(new Request.Builder().url(url).build()).execute()) {
      String target = r.header("Location");
      if (r.code() < 300 || r.code() > 399 || target == null)
        throw new IOException("No app connection");
      target = new java.net.URI(url).resolve(target).toString();
      origin(target);
      return target;
    }
  }

  static JSONObject json(String url, String method, String body, String bearer, boolean cookies)
      throws Exception {
    origin(url);
    Request.Builder b = new Request.Builder().url(url).header("Accept", "application/json");
    if (bearer != null) b.header("Authorization", "Bearer " + bearer);
    if (cookies) {
      String ck = CookieManager.getInstance().getCookie(url);
      if (ck != null) b.header("Cookie", ck);
    }
    if (!method.equals("GET"))
      b.header("Origin", origin(url))
          .method(
              method,
              RequestBody.create(body == null ? "{}" : body, MediaType.get("application/json")));
    try (Response r = HTTP.newCall(b.build()).execute()) {
      if (r.code() == 401 || r.code() == 403) throw new SignInRequired();
      if (!r.isSuccessful() || r.body() == null)
        throw new IOException("Workspace is unavailable (" + r.code() + ")");
      if (r.body().contentLength() > 1048576) throw new IOException("Response too large");
      if (r.code() == 204) return new JSONObject();
      return decode(r.body());
    }
  }

  private static JSONObject decode(ResponseBody body) throws Exception {
    try (java.io.InputStream in = body.byteStream()) {
      java.io.ByteArrayOutputStream out = new java.io.ByteArrayOutputStream();
      byte[] buffer = new byte[8192];
      int n;
      while ((n = in.read(buffer)) != -1) {
        out.write(buffer, 0, n);
        if (out.size() > 1048576) throw new IOException("Response too large");
      }
      return new JSONObject(out.toString("UTF-8"));
    }
  }

  static JSONObject bootstrap(Context c) throws Exception {
    JSONObject d = SecureStore.get(c, "device");
    if (d == null) throw new SignInRequired();
    return bootstrap(d);
  }

  static JSONObject bootstrap(JSONObject d) throws Exception {
    return json(
        d.getString("chat") + "/api/native/bootstrap", "GET", null, d.getString("token"), false);
  }

  static void restoreChatSession(Context c) throws Exception {
    JSONObject d = SecureStore.get(c, "device");
    if (d == null) throw new SignInRequired();
    String chat = d.getString("chat");
    Request.Builder b =
        new Request.Builder()
            .url(chat + "/api/native/session")
            .header("Origin", origin(chat))
            .header("Authorization", "Bearer " + d.getString("token"))
            .post(RequestBody.create("{}", MediaType.get("application/json")));
    String cookie = CookieManager.getInstance().getCookie(chat + "/");
    if (cookie != null) b.header("Cookie", cookie);
    try (Response r = HTTP.newCall(b.build()).execute()) {
      if (r.code() == 401 || r.code() == 403) throw new SignInRequired();
      if (r.code() != 204) throw new IOException("Chat is reconnecting");
      for (String value : r.headers("Set-Cookie"))
        CookieManager.getInstance().setCookie(chat + "/", value);
      CookieManager.getInstance().flush();
    }
  }

  static final class SignInRequired extends IOException {}
}
