package uk.co.m8e.workspace;

import android.content.Context;
import org.json.JSONObject;

/** My Push credentials are scoped to the current phone and encrypted at rest. */
final class PushConfig {
  static synchronized void initialize(Context c, JSONObject p) {
    if (p == null || !p.optBoolean("configured") || !"workspace".equals(p.optString("provider")))
      return;
    JSONObject d = SecureStore.get(c, "device");
    if (d == null || !p.optString("device").equals("chat_" + d.optString("id"))) return;
    try {
      WorkspaceApi.origin(p.getString("url"));
      p.put("owner", d.getString("id"));
      SecureStore.putForDevice(c, d.getString("id"), "push-config", p);
    } catch (Exception ignored) {
    }
  }

  static String token(Context c) {
    return "";
  }

  static boolean ready(Context c) {
    JSONObject p = SecureStore.get(c, "push-config");
    return p != null && SecureStore.isDevice(c, p.optString("owner")) && p.optBoolean("configured");
  }
}
