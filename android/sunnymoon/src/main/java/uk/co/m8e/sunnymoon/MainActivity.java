package uk.co.m8e.sunnymoon;

import android.app.Activity;
import android.content.pm.ApplicationInfo;
import android.content.SharedPreferences;
import android.graphics.Color;
import android.os.Bundle;
import android.view.View;
import android.webkit.WebResourceRequest;
import android.webkit.JavascriptInterface;
import android.webkit.WebResourceResponse;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;
import androidx.webkit.WebViewAssetLoader;
import java.io.ByteArrayInputStream;
import java.util.Collections;

public final class MainActivity extends Activity {
    private static final String GAME_URL = "https://appassets.androidplatform.net/assets/game/index.html";
    private WebView game;

    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        game = new WebView(this);
        game.setBackgroundColor(Color.rgb(180, 217, 239));
        game.setVerticalScrollBarEnabled(false);
        game.setHorizontalScrollBarEnabled(false);
        game.setOverScrollMode(View.OVER_SCROLL_NEVER);
        WebSettings settings = game.getSettings();
        settings.setJavaScriptEnabled(true);
        settings.setDomStorageEnabled(true);
        settings.setAllowFileAccess(false);
        settings.setAllowContentAccess(false);
        settings.setBlockNetworkLoads(true);
        settings.setMixedContentMode(WebSettings.MIXED_CONTENT_NEVER_ALLOW);
        settings.setSupportMultipleWindows(false);
        game.addJavascriptInterface(new RoutineStore(), "SunnyRoutineStore");
        WebView.setWebContentsDebuggingEnabled((getApplicationInfo().flags & ApplicationInfo.FLAG_DEBUGGABLE) != 0);

        WebViewAssetLoader assets = new WebViewAssetLoader.Builder()
            .addPathHandler("/assets/", new WebViewAssetLoader.AssetsPathHandler(this)).build();
        game.setWebViewClient(new WebViewClient() {
            @Override
            public WebResourceResponse shouldInterceptRequest(WebView view, WebResourceRequest request) {
                WebResourceResponse asset = assets.shouldInterceptRequest(request.getUrl());
                if (asset != null) {
                    if (request.getUrl().getPath().endsWith(".mjs")) asset.setMimeType("text/javascript");
                    return asset;
                }
                return new WebResourceResponse("text/plain", "UTF-8", 403, "Offline app",
                    Collections.emptyMap(), new ByteArrayInputStream(new byte[0]));
            }

            @Override
            public boolean shouldOverrideUrlLoading(WebView view, WebResourceRequest request) {
                return !GAME_URL.equals(request.getUrl().toString());
            }
        });
        setContentView(game);
        game.loadUrl(GAME_URL);
    }

    private final class RoutineStore {
        private final SharedPreferences preferences = getSharedPreferences("routine", MODE_PRIVATE);

        @JavascriptInterface
        public String read() {
            return preferences.getString("routine.v1", null);
        }

        @JavascriptInterface
        public boolean write(String value) {
            return value != null && value.length() <= 4096
                && preferences.edit().putString("routine.v1", value).commit();
        }
    }

    @Override
    public void onWindowFocusChanged(boolean hasFocus) {
        super.onWindowFocusChanged(hasFocus);
        if (hasFocus) {
            getWindow().getDecorView().setSystemUiVisibility(View.SYSTEM_UI_FLAG_FULLSCREEN
                | View.SYSTEM_UI_FLAG_HIDE_NAVIGATION | View.SYSTEM_UI_FLAG_IMMERSIVE_STICKY);
        }
    }

    @Override
    public void onBackPressed() {
        game.evaluateJavascript("(() => { const dialog = document.querySelector('dialog[open]');"
            + "if (dialog) { dialog.close(); return true; } return false; })()", result -> {
                if (!"true".equals(result)) finish();
            });
    }

    @Override
    protected void onResume() {
        super.onResume();
        if (game != null) {
            game.onResume();
            game.evaluateJavascript("document.dispatchEvent(new Event('visibilitychange'))", null);
        }
    }

    @Override
    protected void onPause() {
        if (game != null) game.onPause();
        super.onPause();
    }

    @Override
    protected void onDestroy() {
        game.destroy();
        super.onDestroy();
    }
}
