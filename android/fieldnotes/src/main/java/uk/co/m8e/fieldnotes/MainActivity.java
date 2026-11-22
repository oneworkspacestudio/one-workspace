package uk.co.m8e.fieldnotes;

import android.app.Activity;
import android.content.Intent;
import android.content.SharedPreferences;
import android.content.pm.ApplicationInfo;
import android.graphics.Color;
import android.net.Uri;
import android.os.Build;
import android.os.Bundle;
import android.view.View;
import android.view.WindowInsets;
import android.webkit.JavascriptInterface;
import android.webkit.WebResourceRequest;
import android.webkit.WebResourceResponse;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import androidx.webkit.WebViewAssetLoader;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.util.Collections;
import org.json.JSONObject;

public final class MainActivity extends Activity {
    private static final String HOME = "https://appassets.androidplatform.net/assets/reader/index.html";
    private static final int PICK_PACK = 115;
    private static final int MAX_PACK = 2_000_000;
    private WebView reader;

    @Override
    public void onCreate(Bundle state) {
        super.onCreate(state);
        reader = new WebView(this);
        reader.setBackgroundColor(Color.rgb(250, 248, 241));
        getWindow().getDecorView().setSystemUiVisibility(
            View.SYSTEM_UI_FLAG_LIGHT_STATUS_BAR | View.SYSTEM_UI_FLAG_LIGHT_NAVIGATION_BAR);
        reader.setOnApplyWindowInsetsListener((view, insets) -> {
            if (Build.VERSION.SDK_INT >= 30) {
                android.graphics.Insets safe = insets.getInsets(WindowInsets.Type.systemBars() | WindowInsets.Type.displayCutout());
                view.setPadding(safe.left, safe.top, safe.right, safe.bottom);
            }
            return insets;
        });

        WebSettings settings = reader.getSettings();
        settings.setJavaScriptEnabled(true);
        settings.setAllowFileAccess(false);
        settings.setAllowContentAccess(false);
        settings.setBlockNetworkLoads(true);
        settings.setMixedContentMode(WebSettings.MIXED_CONTENT_NEVER_ALLOW);
        settings.setSupportMultipleWindows(false);
        reader.addJavascriptInterface(new Device(), "FieldnotesDevice");
        WebView.setWebContentsDebuggingEnabled((getApplicationInfo().flags & ApplicationInfo.FLAG_DEBUGGABLE) != 0);
        WebViewAssetLoader assets = new WebViewAssetLoader.Builder()
            .addPathHandler("/assets/reader/", path -> new WebViewAssetLoader.AssetsPathHandler(this)
                .handle("reader/" + path)).build();
        reader.setWebViewClient(new WebViewClient() {
            @Override
            public boolean shouldOverrideUrlLoading(WebView view, WebResourceRequest request) {
                return !HOME.equals(request.getUrl().toString().split("#", 2)[0]);
            }

            @Override
            public WebResourceResponse shouldInterceptRequest(WebView view, WebResourceRequest request) {
                WebResourceResponse asset = assets.shouldInterceptRequest(request.getUrl());
                if (asset != null) return asset;
                return new WebResourceResponse("text/plain", "UTF-8", 403, "Offline reader",
                    Collections.emptyMap(), new ByteArrayInputStream(new byte[0]));
            }
        });
        setContentView(reader);
        reader.requestApplyInsets();
        reader.loadUrl(HOME);
    }

    private final class Device {
        private final SharedPreferences preferences = getSharedPreferences("reading-list", MODE_PRIVATE);

        @JavascriptInterface
        public String read() {
            return preferences.getString("pack", null);
        }

        @JavascriptInterface
        public boolean write(String value) {
            return value != null && value.getBytes(StandardCharsets.UTF_8).length <= MAX_PACK
                && preferences.edit().putString("pack", value).commit();
        }

        @JavascriptInterface
        public void pick() {
            runOnUiThread(() -> {
                Intent intent = new Intent(Intent.ACTION_OPEN_DOCUMENT);
                intent.addCategory(Intent.CATEGORY_OPENABLE);
                intent.setType("*/*");
                startActivityForResult(intent, PICK_PACK);
            });
        }
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        if (requestCode != PICK_PACK || resultCode != RESULT_OK || data == null || data.getData() == null) return;
        Uri uri = data.getData();
        new Thread(() -> {
            String content = loadPack(uri);
            runOnUiThread(() -> {
                if (isDestroyed()) return;
                reader.evaluateJavascript("window.fieldnotesImport?.(" +
                    (content == null ? "null" : JSONObject.quote(content)) + ")", null);
            });
        }, "Fieldnotes pack import").start();
    }

    private String loadPack(Uri uri) {
        try (InputStream stream = getContentResolver().openInputStream(uri);
             ByteArrayOutputStream buffer = new ByteArrayOutputStream()) {
            if (stream == null) return null;
            byte[] chunk = new byte[8192];
            for (int read; (read = stream.read(chunk)) != -1;) {
                if (buffer.size() + read > MAX_PACK) return null;
                buffer.write(chunk, 0, read);
            }
            return buffer.toString(StandardCharsets.UTF_8.name());
        } catch (IOException | SecurityException error) {
            return null;
        }
    }

    @Override
    public void onBackPressed() {
        reader.evaluateJavascript("window.fieldnotesBack?.() || false", result -> {
            if (!"true".equals(result)) finish();
        });
    }

    @Override
    protected void onDestroy() {
        reader.removeJavascriptInterface("FieldnotesDevice");
        reader.destroy();
        super.onDestroy();
    }
}
