package uk.co.m8e.workspace;

import android.content.Context;
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
    HttpUrl u = HttpUrl.parse(raw);
    if (u == null || !u.isHttps() || !u.username().isEmpty() || !u.password().isEmpty())
      throw new IllegalArgumentException("Use a secure https:// Workspace address");
    return u.newBuilder()
        .encodedPath("/")
        .query(null)
        .fragment(null)
        .build()
        .toString()
        .replaceAll("/+$", "");
  }

  static String home(Context c) {
    JSONObject servers = SecureStore.get(c, "servers");
    if (servers != null) return servers.optString("workspace", "");
    return c.getSharedPreferences("workspace", 0).getString("home", "");
  }

  static String chatOverride(Context c) {
    JSONObject servers = SecureStore.get(c, "servers");
    return servers == null ? "" : servers.optString("chatOverride", "");
  }

  static JSONObject savedDiscovery(Context c) {
    JSONObject servers = SecureStore.get(c, "servers");
    return servers == null ? SecureStore.get(c, "discovery") : servers.optJSONObject("discovery");
  }

  static String chat(Context c) {
    String override = chatOverride(c);
    JSONObject d = savedDiscovery(c);
    String selected = override.isEmpty() ? (d == null ? "" : d.optString("chat", "")) : override;
    return selected.isEmpty() ? "" : serverAddress(selected, false);
  }

  /** A server base, never a sign-in callback, token-bearing URL or embedded credential. */
  static String serverAddress(String raw, boolean workspace) {
    HttpUrl u = HttpUrl.parse(raw.trim());
    if (u == null
        || !u.isHttps()
        || !u.username().isEmpty()
        || !u.password().isEmpty()
        || u.query() != null
        || u.fragment() != null
        || (workspace && !u.encodedPath().equals("/")))
      throw new IllegalArgumentException("Enter an HTTPS server address");
    return u.toString().replaceAll("/+$", "");
  }

  static void validateDiscovery(String home, JSONObject d) throws Exception {
    if (!serverAddress(d.getString("workspace"), true).equals(home))
      throw new IOException("Workspace returned a different server address");
    d.put("workspace", home);
    origin(d.getString("issuer"));
    String chat = d.optString("chat", "");
    if (!chat.isEmpty()) d.put("chat", serverAddress(chat, false));
    JSONArray origins = d.optJSONArray("origins");
    if (origins != null) for (int i = 0; i < origins.length(); i++) origin(origins.getString(i));
  }

  static String resolveChat(String requested, JSONObject d) throws Exception {
    String base = serverAddress(requested, false);
    String discovered = d.optString("chat", "");
    // Accept a hostname pasted from the address bar even when Chat lives at /chat.
    if (base.equals(origin(base)) && !discovered.isEmpty() && origin(discovered).equals(base))
      base = serverAddress(discovered, false);
    String[] candidates =
        base.equals(origin(base)) ? new String[] {base, base + "/chat"} : new String[] {base};
    for (String candidate : candidates) {
      try (Response r =
          HTTP.newCall(new Request.Builder().url(candidate + "/api/me").build()).execute()) {
        // This marker belongs to the Workspace Chat gateway, including signed-out responses.
        if ((r.code() == 200 || r.code() == 401) && r.header("X-Chat-Instance") != null) {
          // A gateway can answer API routes at / while its UI and cookies live at /chat.
          // Follow only its same-origin sign-in redirect to recover that canonical base.
          HttpUrl login = HttpUrl.get(redirect(candidate + "/"));
          String path = login.encodedPath();
          if (!origin(login.toString()).equals(origin(candidate)) || !path.endsWith("/auth/login"))
            throw new IOException("Chat returned an unexpected sign-in address");
          String canonical =
              origin(candidate) + path.substring(0, path.length() - "/auth/login".length());
          String issuer = origin(redirect(canonical + "/auth/login"));
          d.put("chatIssuer", issuer);
          return canonical;
        }
      } catch (IOException ignored) {
        // A hostname may serve a different page at / and mount Chat under /chat.
      }
    }
    throw new IOException("No Workspace Chat server at that address");
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
