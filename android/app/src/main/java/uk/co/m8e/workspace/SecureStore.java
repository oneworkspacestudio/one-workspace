package uk.co.m8e.workspace;

import android.content.Context;
import android.security.keystore.KeyGenParameterSpec;
import android.security.keystore.KeyProperties;
import android.util.Base64;
import java.security.KeyStore;
import javax.crypto.Cipher;
import javax.crypto.KeyGenerator;
import javax.crypto.SecretKey;
import javax.crypto.spec.GCMParameterSpec;
import org.json.JSONObject;

/** Native device credentials and notification history never enter web storage or backups. */
final class SecureStore {
  private static final String ALIAS = "oneworkspace.device.v1";

  private static SecretKey key() throws Exception {
    KeyStore ks = KeyStore.getInstance("AndroidKeyStore");
    ks.load(null);
    if (!ks.containsAlias(ALIAS)) {
      KeyGenerator g = KeyGenerator.getInstance(KeyProperties.KEY_ALGORITHM_AES, "AndroidKeyStore");
      g.init(
          new KeyGenParameterSpec.Builder(
                  ALIAS, KeyProperties.PURPOSE_ENCRYPT | KeyProperties.PURPOSE_DECRYPT)
              .setBlockModes(KeyProperties.BLOCK_MODE_GCM)
              .setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_NONE)
              .build());
      g.generateKey();
    }
    return (SecretKey) ks.getKey(ALIAS, null);
  }

  static synchronized void put(Context c, String name, JSONObject value) throws Exception {
    String payload = encrypt(name, value);
    if (!c.getSharedPreferences("secure", 0).edit().putString(name, payload).commit())
      throw new java.io.IOException("Device storage unavailable");
  }

  private static String encrypt(String name, JSONObject value) throws Exception {
    Cipher x = Cipher.getInstance("AES/GCM/NoPadding");
    x.init(Cipher.ENCRYPT_MODE, key());
    x.updateAAD(name.getBytes(java.nio.charset.StandardCharsets.UTF_8));
    String payload =
        Base64.encodeToString(x.getIV(), Base64.NO_WRAP)
            + ":"
            + Base64.encodeToString(
                x.doFinal(value.toString().getBytes(java.nio.charset.StandardCharsets.UTF_8)),
                Base64.NO_WRAP);
    return payload;
  }

  /** Atomically replace configuration and discard every old account's cached data/reply. */
  static synchronized void replaceServers(Context c, JSONObject servers) throws Exception {
    String payload = encrypt("servers", servers);
    if (!c.getSharedPreferences("secure", 0).edit().clear().putString("servers", payload).commit())
      throw new java.io.IOException("Device storage unavailable");
  }

  static synchronized JSONObject get(Context c, String name) {
    try {
      String p = c.getSharedPreferences("secure", 0).getString(name, null);
      if (p == null) return null;
      String[] v = p.split(":", 2);
      Cipher x = Cipher.getInstance("AES/GCM/NoPadding");
      x.init(
          Cipher.DECRYPT_MODE,
          key(),
          new GCMParameterSpec(128, Base64.decode(v[0], Base64.NO_WRAP)));
      x.updateAAD(name.getBytes(java.nio.charset.StandardCharsets.UTF_8));
      return new JSONObject(
          new String(
              x.doFinal(Base64.decode(v[1], Base64.NO_WRAP)),
              java.nio.charset.StandardCharsets.UTF_8));
    } catch (Exception e) {
      return null;
    }
  }

  static synchronized void remove(Context c, String name) {
    c.getSharedPreferences("secure", 0).edit().remove(name).commit();
  }

  static synchronized boolean isDevice(Context c, String id) {
    JSONObject d = get(c, "device");
    return d != null && id != null && id.equals(d.optString("id"));
  }

  static synchronized boolean putForDevice(Context c, String id, String name, JSONObject value)
      throws Exception {
    if (!isDevice(c, id)) return false;
    put(c, name, value);
    return true;
  }

  static synchronized void clear(Context c) {
    c.getSharedPreferences("secure", 0).edit().clear().commit();
  }
}
