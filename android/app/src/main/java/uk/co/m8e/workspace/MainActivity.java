package uk.co.m8e.workspace;

import android.Manifest;
import android.app.*;
import android.content.*;
import android.content.pm.PackageManager;
import android.graphics.Color;
import android.net.Uri;
import android.os.*;
import android.view.*;
import android.webkit.*;
import android.widget.*;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.view.ViewCompat;
import androidx.core.view.WindowInsetsCompat;
import androidx.webkit.*;
import androidx.work.WorkManager;
import com.google.android.material.appbar.MaterialToolbar;
import com.google.android.material.button.MaterialButton;
import java.util.*;
import java.util.concurrent.*;
import org.json.*;

/** A local launcher remains usable when the Workspace server is restarting. */
public class MainActivity extends AppCompatActivity {
  static volatile boolean foreground;
  static volatile String visibleTopic = "";

  static boolean viewing(String topic) {
    return foreground && topic != null && !topic.isEmpty() && topic.equals(visibleTopic);
  }

  private final ExecutorService io = Executors.newSingleThreadExecutor();
  private LinearLayout root, body;
  private WebView web;
  private MaterialToolbar toolbar;
  private JSONObject catalog, discovery;
  private final Set<String> trusted = new HashSet<>();
  private volatile long generation;
  private String chat = "", title = "", pendingTopic = "";
  private ValueCallback<Uri[]> chooser;
  private PermissionRequest media;
  private String[] mediaAllowed = new String[0];
  private GeolocationPermissions.Callback locationCallback;
  private String locationOrigin;
  private int BLUE = Color.rgb(52, 93, 229),
      BG = Color.rgb(245, 247, 252),
      INK = Color.rgb(23, 35, 59);

  private int dp(float v) {
    return Math.round(v * getResources().getDisplayMetrics().density);
  }

  @Override
  public void onCreate(Bundle state) {
    super.onCreate(state);
    boolean dark =
        (getResources().getConfiguration().uiMode
                & android.content.res.Configuration.UI_MODE_NIGHT_MASK)
            == android.content.res.Configuration.UI_MODE_NIGHT_YES;
    if (dark) {
      BG = 0xff12161f;
      INK = 0xffe6eaf5;
    }
    getWindow().setSoftInputMode(WindowManager.LayoutParams.SOFT_INPUT_ADJUST_RESIZE);
    ChatNotifications.channels(this);
    root = new LinearLayout(this);
    root.setOrientation(1);
    root.setBackgroundColor(BG);
    setContentView(root);
    androidx.core.view.WindowInsetsControllerCompat systemBars =
        new androidx.core.view.WindowInsetsControllerCompat(getWindow(), root);
    systemBars.setAppearanceLightStatusBars(!dark);
    systemBars.setAppearanceLightNavigationBars(!dark);
    getWindow().setStatusBarColor(BG);
    getWindow().setNavigationBarColor(BG);
    ViewCompat.setOnApplyWindowInsetsListener(
        root,
        (v, insets) -> {
          androidx.core.graphics.Insets bars =
              insets.getInsets(
                  WindowInsetsCompat.Type.systemBars() | WindowInsetsCompat.Type.displayCutout());
          v.setPadding(bars.left, bars.top, bars.right, bars.bottom);
          return insets;
        });
    toolbar = new MaterialToolbar(this);
    toolbar.setTitleTextColor(INK);
    toolbar.setNavigationIconTint(INK);
    toolbar.setBackgroundColor(BG);
    root.addView(toolbar, new LinearLayout.LayoutParams(-1, dp(64)));
    body = new LinearLayout(this);
    body.setOrientation(1);
    root.addView(body, new LinearLayout.LayoutParams(-1, 0, 1));
    catalog = SecureStore.get(this, "catalog");
    discovery = WorkspaceApi.savedDiscovery(this);
    restoreTrust();
    if (WorkspaceApi.home(this).isEmpty()) setup();
    else {
      home();
      refresh();
    }
    if (this instanceof DocumentActivity) {
      String target = getIntent().getStringExtra("documentUrl");
      if (target != null && allowed(target)) open(target, "Document");
      else finish();
      return;
    }
    getOnBackPressedDispatcher()
        .addCallback(
            this,
            new androidx.activity.OnBackPressedCallback(true) {
              @Override
              public void handleOnBackPressed() {
                navigateBack();
              }
            });
    handleIntent(getIntent());
  }

  private TextView text(String value, int size) {
    TextView t = new TextView(this);
    t.setText(value);
    t.setTextSize(size);
    t.setTextColor(INK);
    t.setPadding(dp(20), dp(12), dp(20), dp(12));
    return t;
  }

  private MaterialButton button(String label, Runnable action) {
    MaterialButton b = new MaterialButton(this);
    b.setText(label);
    b.setAllCaps(false);
    b.setMinHeight(dp(48));
    b.setOnClickListener(v -> action.run());
    return b;
  }

  private void clear(String label, boolean back) {
    visibleTopic = "";
    if (web != null) {
      web.stopLoading();
      web.destroy();
      web = null;
    }
    body.removeAllViews();
    toolbar.setVisibility(View.VISIBLE);
    toolbar.setTitle(label);
    toolbar.getMenu().clear();
    if (back) toolbar.setNavigationIcon(R.drawable.ic_back);
    else toolbar.setNavigationIcon((android.graphics.drawable.Drawable) null);
    toolbar.setNavigationOnClickListener(
        v -> {
          if (this instanceof DocumentActivity) finish();
          else home();
        });
  }

  private LinearLayout formPage() {
    ScrollView scroll = new ScrollView(this);
    scroll.setFillViewport(true);
    LinearLayout panel = new LinearLayout(this);
    panel.setOrientation(LinearLayout.VERTICAL);
    panel.setPadding(dp(20), dp(12), dp(20), dp(20));
    scroll.addView(panel);
    body.addView(scroll, new LinearLayout.LayoutParams(-1, -1));
    return panel;
  }

  private void setup() {
    setup(WorkspaceApi.home(this), WorkspaceApi.chatOverride(this), "");
  }

  private EditText serverField(LinearLayout panel, String label, String hint, String value) {
    com.google.android.material.textfield.TextInputLayout field =
        new com.google.android.material.textfield.TextInputLayout(this);
    field.setHint(label);
    field.setBoxBackgroundMode(
        com.google.android.material.textfield.TextInputLayout.BOX_BACKGROUND_OUTLINE);
    com.google.android.material.textfield.TextInputEditText input =
        new com.google.android.material.textfield.TextInputEditText(field.getContext());
    input.setSingleLine();
    input.setInputType(
        android.text.InputType.TYPE_CLASS_TEXT | android.text.InputType.TYPE_TEXT_VARIATION_URI);
    input.setHint(hint);
    input.setText(value);
    input.setImportantForAutofill(View.IMPORTANT_FOR_AUTOFILL_NO);
    field.addView(input, new LinearLayout.LayoutParams(-1, -2));
    LinearLayout.LayoutParams layout = new LinearLayout.LayoutParams(-1, -2);
    layout.setMargins(0, dp(12), 0, dp(8));
    panel.addView(field, layout);
    return input;
  }

  private void setup(String savedHome, String savedChat, String error) {
    boolean configured = !WorkspaceApi.home(this).isEmpty();
    clear(configured ? "Servers" : "One Workspace", configured);
    LinearLayout panel = formPage();
    panel.setPadding(dp(20), dp(20), dp(20), dp(20));
    ImageView logo = new ImageView(this);
    logo.setImageResource(R.drawable.workspace_brand);
    panel.addView(logo, new LinearLayout.LayoutParams(dp(72), dp(72)));
    panel.addView(text("Your Workspace, wherever you are", 27));
    panel.addView(
        text(
            "Connect to your own remote servers once. We’ll remember them next time you open the"
                + " app.",
            16));
    panel.addView(
        text("All your apps open as web apps here. No separate phone apps to install.", 14));
    EditText address =
        serverField(panel, "Workspace address", "https://workspace.example.com", savedHome);
    Switch automatic = new Switch(this);
    automatic.setText("Use Chat from this Workspace");
    automatic.setChecked(savedChat.isEmpty());
    automatic.setMinHeight(dp(48));
    panel.addView(automatic);
    LinearLayout chatPanel = new LinearLayout(this);
    chatPanel.setOrientation(LinearLayout.VERTICAL);
    EditText chatAddress =
        serverField(chatPanel, "Chat address", "https://chat.example.com/chat", savedChat);
    chatPanel.addView(
        text(
            "Paste the address you normally use to open your Chat server, including /chat if"
                + " shown.",
            14));
    chatPanel.setVisibility(automatic.isChecked() ? View.GONE : View.VISIBLE);
    panel.addView(chatPanel);
    automatic.setOnCheckedChangeListener(
        (b, on) -> chatPanel.setVisibility(on ? View.GONE : View.VISIBLE));
    if (!error.isEmpty()) panel.addView(text(error, 16));
    MaterialButton connect =
        button(
            configured ? "Save servers" : "Connect to Workspace",
            () -> {
              String h;
              try {
                h = WorkspaceApi.serverAddress(address.getText().toString(), true);
              } catch (Exception e) {
                address.setError("Enter your https:// Workspace address without a page path");
                return;
              }
              String selectedChat = "";
              if (!automatic.isChecked()) {
                try {
                  selectedChat =
                      WorkspaceApi.serverAddress(chatAddress.getText().toString(), false);
                } catch (Exception e) {
                  chatAddress.setError("Enter your https:// Chat address");
                  return;
                }
              }
              connect(h, selectedChat);
            });
    panel.addView(connect);
    panel.addView(
        text(
            "Sign in with MyAccounts when asked. We keep your session while it’s valid. Change"
                + " these addresses later in Settings → Servers.",
            14));
    if (configured) panel.addView(button("Cancel", this::settings));
  }

  private void connect(String h, String selectedChat) {
    final long operation = ++generation;
    status("Checking your servers securely…");
    body.addView(
        button(
            "Cancel",
            () -> {
              generation++;
              setup(h, selectedChat, "");
            }));
    io.execute(
        () -> {
          String stage = "Workspace";
          try {
            JSONObject d = WorkspaceApi.discovery(h);
            WorkspaceApi.validateDiscovery(h, d);
            stage = "Chat";
            String override =
                selectedChat.isEmpty() ? "" : WorkspaceApi.resolveChat(selectedChat, d);
            runOnUiThread(
                () -> {
                  if (generation != operation) return;
                  String effectiveChat = override.isEmpty() ? d.optString("chat", "") : override;
                  boolean changed =
                      !h.equals(WorkspaceApi.home(this)) || !effectiveChat.equals(chat);
                  if (changed && !WorkspaceApi.home(this).isEmpty()) {
                    new AlertDialog.Builder(this)
                        .setTitle("Switch servers?")
                        .setMessage(
                            "Connect to "
                                + h
                                + (effectiveChat.isEmpty() ? "" : "\nChat: " + effectiveChat)
                                + "\n\n"
                                + "This phone’s old sign-in, saved app list and notification access"
                                + " will be removed. Your server data stays where it is.")
                        .setPositiveButton(
                            "Switch servers",
                            (dialog, which) -> saveConnection(h, override, d, true))
                        .setNegativeButton("Cancel", (dialog, which) -> setup(h, selectedChat, ""))
                        .setOnCancelListener(dialog -> setup(h, selectedChat, ""))
                        .show();
                  } else saveConnection(h, override, d, changed);
                });
          } catch (Exception e) {
            String failed = stage;
            runOnUiThread(
                () -> {
                  if (generation != operation) return;
                  setup(
                      h,
                      selectedChat,
                      "Couldn’t connect to "
                          + failed
                          + ". Check the address, Internet connection and server HTTPS certificate,"
                          + " then try again. Your saved servers haven’t changed.");
                });
          }
        });
  }

  private void saveConnection(String h, String override, JSONObject d, boolean changed) {
    status("Saving your servers…");
    Runnable save =
        () -> {
          if (isFinishing() || isDestroyed()) return;
          try {
            JSONObject servers =
                new JSONObject()
                    .put("workspace", h)
                    .put("chatOverride", override)
                    .put("discovery", d);
            synchronized (SecureStore.class) {
              if (changed) SecureStore.replaceServers(this, servers);
              else SecureStore.put(this, "servers", servers);
              generation++;
            }
            getSharedPreferences("workspace", 0).edit().remove("home").apply();
            if (changed) catalog = null;
            discovery = d;
            pendingTopic = "";
            restoreTrust();
            if (catalog == null) open(h + "/auth/login", "Sign in");
            else {
              home();
              refresh();
            }
          } catch (Exception e) {
            setup(
                h, override, "Couldn’t save the server settings on this phone. Please try again.");
          }
        };
    if (changed) {
      disconnectDevice();
      // These stores belong to this Android app, never the user's external browser.
      WebStorage.getInstance().deleteAllData();
      CookieManager.getInstance()
          .removeAllCookies(
              removed -> {
                CookieManager.getInstance().flush();
                save.run();
              });
    } else save.run();
  }

  private void restoreTrust() {
    trusted.clear();
    String h = WorkspaceApi.home(this);
    if (!h.isEmpty()) trusted.add(h);
    chat = "";
    if (discovery == null) return;
    try {
      if (!WorkspaceApi.origin(discovery.getString("workspace")).equals(h)) return;
      trusted.add(WorkspaceApi.origin(discovery.getString("issuer")));
      String override = WorkspaceApi.chatOverride(this);
      chat = override.isEmpty() ? discovery.optString("chat", "") : override;
      if (!chat.isEmpty()) chat = WorkspaceApi.serverAddress(chat, false);
      String chatIssuer = discovery.optString("chatIssuer", "");
      if (!chatIssuer.isEmpty()) trusted.add(WorkspaceApi.origin(chatIssuer));
      JSONArray a = discovery.optJSONArray("origins");
      if (a != null)
        for (int i = 0; i < a.length(); i++) trusted.add(WorkspaceApi.origin(a.getString(i)));
      if (!chat.isEmpty()) trusted.add(WorkspaceApi.origin(chat));
    } catch (Exception e) {
      chat = "";
    }
  }

  private void status(String msg) {
    clear("One Workspace", false);
    body.addView(text(msg, 17));
  }

  private void home() {
    clear("One Workspace", false);
    body.setPadding(0, 0, 0, 0);
    toolbar.getMenu().add(0, 1, 1, "Refresh apps").setShowAsAction(MenuItem.SHOW_AS_ACTION_NEVER);
    android.graphics.drawable.Drawable settingsIcon =
        androidx.core.content.ContextCompat.getDrawable(this, R.drawable.ic_settings).mutate();
    androidx.core.graphics.drawable.DrawableCompat.setTint(settingsIcon, INK);
    toolbar
        .getMenu()
        .add(0, 2, 0, "Settings")
        .setIcon(settingsIcon)
        .setShowAsAction(MenuItem.SHOW_AS_ACTION_ALWAYS);
    toolbar.setOnMenuItemClickListener(
        i -> {
          if (i.getItemId() == 2) settings();
          else if (i.getItemId() == 1) refresh();
          else return false;
          return true;
        });
    if (catalog == null) {
      body.addView(
          text("Connect once to load your apps. Your launcher will remain available offline.", 17));
      body.addView(
          button("Sign in", () -> open(WorkspaceApi.home(this) + "/auth/login", "Sign in")));
    } else {
      body.addView(text("Hello, " + catalog.optString("name", "there"), 24));
      body.addView(text("Your web apps — no extra phone apps needed", 14));
      EditText search = new EditText(this);
      search.setSingleLine();
      search.setHint("Search your apps");
      search.setPadding(dp(20), 0, dp(20), 0);
      body.addView(search, new LinearLayout.LayoutParams(-1, dp(56)));
      ScrollView scroll = new ScrollView(this);
      GridLayout grid = new GridLayout(this);
      grid.setColumnCount(3);
      grid.setPadding(dp(10), dp(14), dp(10), dp(14));
      scroll.addView(grid);
      body.addView(scroll, new LinearLayout.LayoutParams(-1, 0, 1));
      Runnable render =
          () -> {
            grid.removeAllViews();
            JSONArray apps = launcherApps();
            String q = search.getText().toString().toLowerCase(Locale.ROOT);
            for (String category : new String[] {"general", "administration", "development"}) {
              boolean headingAdded = false;
              for (int n = 0; n < apps.length(); n++) {
                JSONObject app = apps.optJSONObject(n);
                if (app == null || !app.optString("Name").toLowerCase(Locale.ROOT).contains(q))
                  continue;
                if (!category.equals(app.optString("CategoryID", appCategory(app.optString("ID")))))
                  continue;
                if (!headingAdded) {
                  TextView heading =
                      text(
                          category.equals("general")
                              ? "General"
                              : category.equals("administration")
                                  ? "Administration"
                                  : "Creative & Development",
                          21);
                  heading.setAccessibilityHeading(true);
                  GridLayout.LayoutParams sectionLayout = new GridLayout.LayoutParams();
                  sectionLayout.columnSpec = GridLayout.spec(0, 3);
                  sectionLayout.width = GridLayout.LayoutParams.MATCH_PARENT;
                  grid.addView(heading, sectionLayout);
                  headingAdded = true;
                }
                LinearLayout tile = new LinearLayout(this);
                tile.setOrientation(1);
                tile.setGravity(Gravity.CENTER);
                tile.setPadding(dp(4), dp(12), dp(4), dp(12));
                tile.setBackgroundResource(android.R.drawable.list_selector_background);
                tile.setMinimumHeight(dp(112));
                tile.setContentDescription("Open " + app.optString("Name"));
                ImageView icon = new ImageView(this);
                String id = app.optString("ID", "");
                String iconName =
                    app.optString("Icon", "")
                        .replaceAll("^.*/", "")
                        .replaceAll("-[a-f0-9]{16}\\.svg$", "")
                        .replaceAll("\\.svg$", "");
                int resource =
                    getResources().getIdentifier("brand_" + iconName, "drawable", getPackageName());
                if (resource == 0)
                  resource =
                      getResources()
                          .getIdentifier(
                              "app_" + id.replace('-', '_'), "drawable", getPackageName());
                icon.setImageResource(resource == 0 ? R.drawable.workspace_brand : resource);
                tile.addView(icon, new LinearLayout.LayoutParams(dp(46), dp(46)));
                TextView name = text(app.optString("Name"), 13);
                name.setGravity(Gravity.CENTER);
                name.setMaxLines(2);
                name.setPadding(0, dp(8), 0, 0);
                tile.addView(name);
                GridLayout.LayoutParams lp = new GridLayout.LayoutParams();
                lp.width = 0;
                lp.columnSpec = GridLayout.spec(GridLayout.UNDEFINED, 1f);
                lp.setMargins(dp(4), dp(4), dp(4), dp(4));
                grid.addView(tile, lp);
                tile.setOnClickListener(
                    v -> {
                      if ("chat".equals(id)) {
                        openChat();
                        return;
                      }
                      String u = app.optString("URL", "");
                      try {
                        String target =
                            new java.net.URI(WorkspaceApi.home(this) + "/").resolve(u).toString();
                        WorkspaceApi.origin(target);
                        open(target, app.optString("Name"));
                      } catch (Exception ignored) {
                        toast("This app link is unavailable");
                      }
                    });
              }
            }
          };
      render.run();
      search.addTextChangedListener(
          new android.text.TextWatcher() {
            public void beforeTextChanged(CharSequence s, int st, int count, int after) {}

            public void onTextChanged(CharSequence s, int st, int before, int count) {
              render.run();
            }

            public void afterTextChanged(android.text.Editable e) {}
          });
    }
  }

  private JSONArray launcherApps() {
    JSONArray apps = new JSONArray();
    JSONArray saved = catalog == null ? null : catalog.optJSONArray("apps");
    boolean hasChat = false;
    if (saved != null) {
      for (int i = 0; i < saved.length(); i++) {
        JSONObject app = saved.optJSONObject(i);
        if (app == null) continue;
        hasChat |= "chat".equals(app.optString("ID"));
        apps.put(app);
      }
    }
    // A separate Chat server need not be installed on the chosen Workspace host.
    if (!hasChat && !chat.isEmpty()) {
      try {
        apps.put(
            new JSONObject()
                .put("ID", "chat")
                .put("Name", "Chat")
                .put("Icon", "my_chat.svg")
                .put("URL", chat + "/"));
      } catch (JSONException ignored) {
      }
    }
    return apps;
  }

  private static String appCategory(String id) {
    if (Arrays.asList("accounts", "backups", "monitor", "shield", "updates", "app-status")
        .contains(id)) return "administration";
    return Arrays.asList("gitea", "hosting", "myai", "my3d", "studio", "create", "slicer", "push").contains(id) ? "development" : "general";
  }

  private void refresh() {
    final String target = WorkspaceApi.home(this);
    final long operation = generation;
    if (target.isEmpty()) return;
    io.execute(
        () -> {
          try {
            JSONObject c = WorkspaceApi.json(target + "/_workspace/apps", "GET", null, null, true);
            synchronized (SecureStore.class) {
              if (generation != operation || !target.equals(WorkspaceApi.home(this))) return;
              SecureStore.put(this, "catalog", c);
            }
            runOnUiThread(
                () -> {
                  if (generation != operation) return;
                  catalog = c;
                  if (web == null
                      || ("Sign in".equals(title)
                          && web.getUrl() != null
                          && web.getUrl().equals(target + "/"))) home();
                });
          } catch (WorkspaceApi.SignInRequired e) {
            runOnUiThread(
                () -> {
                  if (generation == operation) toast("Sign in to refresh your apps");
                });
          } catch (Exception e) {
            runOnUiThread(
                () -> {
                  if (generation == operation) toast("Offline — showing your saved apps");
                });
          }
        });
  }

  private boolean allowed(String url) {
    try {
      return trusted.contains(WorkspaceApi.origin(url));
    } catch (Exception e) {
      return false;
    }
  }

  private boolean isChat(String url) {
    return !chat.isEmpty() && url != null && url.startsWith(chat + "/");
  }

  private boolean isChatUI(String url) {
    if (!isChat(url)) return false;
    String path = Uri.parse(url).getPath(), base = Uri.parse(chat).getPath();
    return path.equals(base + "/") || path.equals(base + "/index.html");
  }

  private void openChatPage() {
    String target = chat + "/" + (pendingTopic.isEmpty() ? "" : "#" + pendingTopic);
    pendingTopic = "";
    open(target, "Chat");
  }

  private void openChat() {
    if (chat.isEmpty()) {
      toast("Chat is not installed on this Workspace");
      return;
    }
    if (SecureStore.get(this, "device") == null) {
      openChatPage();
      return;
    }
    final long operation = generation;
    final String target = chat;
    final JSONObject device = SecureStore.get(this, "device");
    io.execute(
        () -> {
          try {
            WorkspaceApi.restoreChatSession(this);
          } catch (WorkspaceApi.SignInRequired e) {
            synchronized (SecureStore.class) {
              if (device != null && SecureStore.isDevice(this, device.optString("id")))
                SecureStore.remove(this, "device");
            }
          } catch (Exception ignored) {
          }
          runOnUiThread(
              () -> {
                if (generation == operation && target.equals(WorkspaceApi.chat(this)))
                  openChatPage();
              });
        });
  }

  private boolean selectedChatLink(String url) {
    if (WorkspaceApi.chatOverride(this).isEmpty() || isChat(url)) return false;
    try {
      Uri uri = Uri.parse(url);
      String original = discovery == null ? "" : discovery.optString("chat", "");
      String destination =
          uri.buildUpon().clearQuery().fragment(null).build().toString().replaceAll("/+$", "");
      if (!original.isEmpty() && destination.equals(WorkspaceApi.serverAddress(original, false)))
        return true;
      String path = uri.getPath();
      return WorkspaceApi.origin(url).equals(WorkspaceApi.home(this))
          && ("/chat".equals(path) || "/chat/".equals(path));
    } catch (Exception ignored) {
      return false;
    }
  }

  private void open(String url, String label) {
    if (selectedChatLink(url)) {
      openChat();
      return;
    }
    if (!allowed(url)) {
      external(url);
      return;
    }
    clear(label, true);
    title = label;
    web = new WebView(this);
    WebView w = web;
    WebSettings s = w.getSettings();
    s.setJavaScriptEnabled(true);
    s.setDomStorageEnabled(true);
    s.setAllowFileAccess(false);
    s.setAllowContentAccess(false);
    s.setMixedContentMode(WebSettings.MIXED_CONTENT_NEVER_ALLOW);
    s.setMediaPlaybackRequiresUserGesture(true);
    s.setSupportMultipleWindows(true);
    s.setJavaScriptCanOpenWindowsAutomatically(false);
    s.setSafeBrowsingEnabled(true);
    CookieManager.getInstance().setAcceptCookie(true);
    CookieManager.getInstance().setAcceptThirdPartyCookies(w, false);
    if (WebViewFeature.isFeatureSupported(WebViewFeature.WEB_MESSAGE_LISTENER))
      WebViewCompat.addWebMessageListener(
          w,
          "WorkspaceNative",
          chat.isEmpty()
              ? Collections.singleton(WorkspaceApi.home(this))
              : Collections.singleton(WorkspaceApi.origin(chat)),
          (view, message, source, main, reply) -> {
            if (!main || !isChat(view.getUrl())) return;
            try {
              JSONObject m = new JSONObject(message.getData());
              if ("signedOut".equals(m.optString("type"))) disconnectDevice();
            } catch (Exception ignored) {
            }
          });
    w.setWebViewClient(
        new WebViewClient() {
          boolean chatReady;

          @Override
          public void onPageStarted(WebView v, String u, android.graphics.Bitmap icon) {
            if (allowed(u)) {
              Uri uri = Uri.parse(u);
              if ("1".equals(uri.getQueryParameter("signedOut"))
                  || uri.getPath().endsWith("/auth/logout")) disconnectDevice();
            }
          }

          @Override
          public void doUpdateVisitedHistory(WebView v, String u, boolean reload) {
            visibleTopic = "";
            if (isChat(u)) {
              String hash = Uri.parse(u).getFragment();
              if (hash != null) {
                String candidate = hash.replaceFirst("^/", "").split("/", 2)[0];
                if (ChatConnection.validTopic(candidate)) visibleTopic = candidate;
              }
            }
          }

          @Override
          public boolean shouldOverrideUrlLoading(WebView v, WebResourceRequest r) {
            if (!r.isForMainFrame()) return false;
            String u = r.getUrl().toString();
            if (selectedChatLink(u)) {
              openChat();
              return true;
            }
            if (allowed(u)) return false;
            external(u);
            return true;
          }

          @Override
          public void onPageFinished(WebView v, String u) {
            CookieManager.getInstance().flush();
            if (allowed(u)) refresh();
            if (isChatUI(u)) {
              toolbar.setVisibility(View.GONE);
              if (!chatReady) {
                v.clearHistory();
                chatReady = true;
              }
              linkDevice(u);
            } else toolbar.setVisibility(View.VISIBLE);
          }

          @Override
          public void onReceivedSslError(
              WebView v, android.webkit.SslErrorHandler handler, android.net.http.SslError e) {
            handler.cancel();
            toast("Secure connection failed. Check your Workspace certificate.");
          }

          @Override
          public void onReceivedError(WebView v, WebResourceRequest r, WebResourceError e) {
            if (r.isForMainFrame())
              toast("Workspace is reconnecting. Your saved apps are available from Back to apps.");
          }
        });
    w.setWebChromeClient(
        new WebChromeClient() {
          @Override
          public boolean onCreateWindow(
              WebView parent, boolean dialog, boolean gesture, Message result) {
            if (!gesture) return false;
            WebView pending = new WebView(MainActivity.this);
            pending.getSettings().setAllowFileAccess(false);
            pending.getSettings().setAllowContentAccess(false);
            pending.getSettings().setMixedContentMode(WebSettings.MIXED_CONTENT_NEVER_ALLOW);
            pending.setWebViewClient(
                new WebViewClient() {
                  boolean done;

                  @Override
                  public boolean shouldOverrideUrlLoading(WebView v, WebResourceRequest r) {
                    String target = r.getUrl().toString();
                    if (target.equals("about:blank")) return false;
                    if (done) return true;
                    done = true;
                    if (allowed(target))
                      startActivity(
                          new Intent(MainActivity.this, DocumentActivity.class)
                              .putExtra("documentUrl", target));
                    else external(target);
                    v.stopLoading();
                    v.destroy();
                    return true;
                  }
                });
            ((WebView.WebViewTransport) result.obj).setWebView(pending);
            result.sendToTarget();
            new Handler(Looper.getMainLooper())
                .postDelayed(
                    () -> {
                      try {
                        pending.destroy();
                      } catch (Exception ignored) {
                      }
                    },
                    60000);
            return true;
          }

          @Override
          public boolean onShowFileChooser(
              WebView v, ValueCallback<Uri[]> callback, FileChooserParams p) {
            if (chooser != null) chooser.onReceiveValue(null);
            chooser = callback;
            try {
              startActivityForResult(p.createIntent(), 100);
              return true;
            } catch (Exception e) {
              chooser = null;
              return false;
            }
          }

          @Override
          public void onGeolocationPermissionsShowPrompt(
              String origin, GeolocationPermissions.Callback callback) {
            if (!allowed(origin) || locationCallback != null) {
              callback.invoke(origin, false, false);
              return;
            }
            locationCallback = callback;
            locationOrigin = origin;
            new AlertDialog.Builder(MainActivity.this)
                .setTitle("Allow location?")
                .setMessage(
                    "Allow "
                        + Uri.parse(origin).getHost()
                        + " to use your location while this app is open?")
                .setPositiveButton(
                    "Continue",
                    (d, which) ->
                        requestPermissions(
                            new String[] {
                              Manifest.permission.ACCESS_COARSE_LOCATION,
                              Manifest.permission.ACCESS_FINE_LOCATION
                            },
                            103))
                .setNegativeButton(
                    "Cancel",
                    (d, which) -> {
                      callback.invoke(origin, false, false);
                      locationCallback = null;
                    })
                .setOnCancelListener(
                    d -> {
                      callback.invoke(origin, false, false);
                      locationCallback = null;
                    })
                .show();
          }

          @Override
          public void onPermissionRequest(PermissionRequest p) {
            runOnUiThread(() -> requestMedia(p));
          }
        });
    w.setDownloadListener((u, agent, cd, mime, size) -> external(u));
    body.addView(w, new LinearLayout.LayoutParams(-1, -1));
    toolbar
        .getMenu()
        .add("Reload")
        .setOnMenuItemClickListener(
            i -> {
              w.reload();
              return true;
            });
    toolbar
        .getMenu()
        .add("Open in browser")
        .setOnMenuItemClickListener(
            i -> {
              external(w.getUrl());
              return true;
            });
    w.loadUrl(url);
  }

  private void linkDevice(String url) {
    final String target = chat;
    final long operation = generation;
    io.execute(
        () -> {
          try {
            JSONObject me = WorkspaceApi.json(target + "/api/me", "GET", null, null, true);
            if (generation != operation
                || !target.equals(chat)
                || !target.equals(WorkspaceApi.chat(this))) return;
            JSONObject old = SecureStore.get(this, "device");
            if (old != null
                && old.optString("accountKey").equals(me.getString("accountKey"))
                && old.optString("chat").equals(target)) {
              try {
                WorkspaceApi.bootstrap(this);
              } catch (WorkspaceApi.SignInRequired expired) {
                synchronized (SecureStore.class) {
                  if (generation != operation || !SecureStore.isDevice(this, old.optString("id")))
                    return;
                  SecureStore.remove(this, "device");
                }
                old = null;
              }
            }
            if (generation != operation || !target.equals(WorkspaceApi.chat(this))) return;
            if (old == null
                || !old.optString("accountKey").equals(me.getString("accountKey"))
                || !old.optString("chat").equals(target)) {
              if (old != null) disconnectDevice();
              long next = generation;
              JSONObject d =
                  WorkspaceApi.json(target + "/api/native/link", "POST", "{}", null, true);
              d.put("chat", target);
              d.put("accountKey", me.getString("accountKey"));
              boolean accept;
              synchronized (SecureStore.class) {
                accept =
                    generation == next
                        && target.equals(chat)
                        && target.equals(WorkspaceApi.chat(this));
                if (accept) SecureStore.put(this, "device", d);
              }
              if (!accept) {
                try {
                  WorkspaceApi.json(
                      target + "/api/native/revoke", "POST", "{}", d.getString("token"), false);
                } catch (Exception ignored) {
                }
                return;
              }
            }
            ChatSyncWorker.schedule(this);
            PushConfig.initialize(this, WorkspaceApi.bootstrap(this).optJSONObject("push"));
            runOnUiThread(
                () -> {
                  WorkspacePushService.start(this);
                  if (Build.VERSION.SDK_INT >= 33
                      && checkSelfPermission(Manifest.permission.POST_NOTIFICATIONS)
                          != PackageManager.PERMISSION_GRANTED)
                    requestPermissions(new String[] {Manifest.permission.POST_NOTIFICATIONS}, 101);
                });
          } catch (Exception e) {
            runOnUiThread(
                () ->
                    toast(
                        "Chat is open. Background notifications need the matching Chat server"
                            + " update."));
          }
        });
  }

  private void requestMedia(PermissionRequest p) {
    if (!allowed(p.getOrigin().toString()) || media != null) {
      p.deny();
      return;
    }
    ArrayList<String> permissions = new ArrayList<>(), resources = new ArrayList<>();
    for (String resource : p.getResources()) {
      if (resource.equals(PermissionRequest.RESOURCE_AUDIO_CAPTURE)) {
        permissions.add(Manifest.permission.RECORD_AUDIO);
        resources.add(resource);
      } else if (resource.equals(PermissionRequest.RESOURCE_VIDEO_CAPTURE)) {
        permissions.add(Manifest.permission.CAMERA);
        resources.add(resource);
      }
    }
    if (permissions.isEmpty()) {
      p.deny();
      return;
    }
    media = p;
    mediaAllowed = resources.toArray(new String[0]);
    new AlertDialog.Builder(this)
        .setTitle("Allow camera or microphone?")
        .setMessage(
            "Allow " + p.getOrigin().getHost() + " to use the requested camera/microphone access?")
        .setPositiveButton(
            "Continue", (d, which) -> requestPermissions(permissions.toArray(new String[0]), 102))
        .setNegativeButton(
            "Cancel",
            (d, which) -> {
              p.deny();
              media = null;
            })
        .setOnCancelListener(
            d -> {
              p.deny();
              media = null;
            })
        .show();
  }

  private void settings() {
    clear("Settings", true);
    LinearLayout panel = formPage();
    panel.addView(text("Your servers", 23));
    panel.addView(
        text(
            "Workspace: "
                + WorkspaceApi.home(this)
                + "\nChat: "
                + (chat.isEmpty() ? "Not configured" : chat),
            16));
    panel.addView(button("Servers", this::setup));
    panel.addView(text("Chat notifications", 23));
    panel.addView(
        text(
            SecureStore.get(this, "device") == null
                ? "Open Chat and sign in once to link this device."
                : (PushConfig.ready(this)
                    ? "My Push and background sync configured"
                    : "Background sync enabled. Instant notifications need push configured on your"
                        + " Workspace."),
            16));
    panel.addView(
        button(
            "Notification settings",
            () -> {
              if (Build.VERSION.SDK_INT >= 33
                  && checkSelfPermission(Manifest.permission.POST_NOTIFICATIONS)
                      != PackageManager.PERMISSION_GRANTED)
                requestPermissions(new String[] {Manifest.permission.POST_NOTIFICATIONS}, 101);
              else
                startActivity(
                    new Intent(android.provider.Settings.ACTION_APP_NOTIFICATION_SETTINGS)
                        .putExtra(android.provider.Settings.EXTRA_APP_PACKAGE, getPackageName()));
            }));
    if (!chat.isEmpty() && SecureStore.get(this, "device") != null)
      panel.addView(button("Manage linked phones", () -> open(chat + "/devices", "Linked phones")));
    Switch preview = new Switch(this);
    preview.setText("Show message previews");
    preview.setPadding(dp(20), dp(12), dp(20), dp(12));
    preview.setChecked(getSharedPreferences("workspace", 0).getBoolean("previews", false));
    preview.setOnCheckedChangeListener(
        (b, on) -> getSharedPreferences("workspace", 0).edit().putBoolean("previews", on).apply());
    panel.addView(preview);
    Switch background = new Switch(this);
    background.setText("Keep private Chat connection running");
    background.setPadding(dp(20), dp(12), dp(20), dp(12));
    background.setChecked(getSharedPreferences("workspace", 0).getBoolean("backgroundPush", true));
    background.setOnCheckedChangeListener(
        (b, on) -> {
          getSharedPreferences("workspace", 0).edit().putBoolean("backgroundPush", on).apply();
          if (on) WorkspacePushService.start(this);
          else stopService(new Intent(this, WorkspacePushService.class));
        });
    panel.addView(background);
    panel.addView(
        text(
            "My Push uses your own server. Android shows a quiet ongoing connection notification;"
                + " no Firebase account is needed.",
            14));
    panel.addView(
        button(
            "Android battery settings",
            () ->
                startActivity(
                    new Intent(
                        android.provider.Settings.ACTION_IGNORE_BATTERY_OPTIMIZATION_SETTINGS))));
    panel.addView(
        text(
            "Android may delay periodic sync. Force-stopping the app pauses notifications until you"
                + " open it again.",
            14));
    panel.addView(
        button(
            "Disconnect this device",
            () ->
                new AlertDialog.Builder(this)
                    .setMessage(
                        "Disconnect this phone and remove its Chat notification access? Server data"
                            + " stays in your Workspace.")
                    .setPositiveButton(
                        "Disconnect",
                        (d, w) -> {
                          disconnectDevice();
                          CookieManager.getInstance().removeAllCookies(null);
                          WebStorage.getInstance().deleteAllData();
                          SecureStore.clear(this);
                          catalog = null;
                          discovery = null;
                          chat = "";
                          trusted.clear();
                          getSharedPreferences("workspace", 0).edit().clear().apply();
                          setup();
                        })
                    .setNegativeButton("Cancel", null)
                    .show()));
  }

  private void disconnectDevice() {
    stopService(new Intent(this, WorkspacePushService.class));
    JSONObject d;
    synchronized (SecureStore.class) {
      generation++;
      d = SecureStore.get(this, "device");
      SecureStore.remove(this, "device");
      SecureStore.remove(this, "notification-history");
      SecureStore.remove(this, "push-config");
      SecureStore.remove(this, "push-seen");
      ChatNotifications.clear(this);
    }
    if (d != null)
      io.execute(
          () -> {
            try {
              WorkspaceApi.json(
                  d.getString("chat") + "/api/native/revoke",
                  "POST",
                  "{}",
                  d.getString("token"),
                  false);
            } catch (Exception ignored) {
            }
          });
    WorkManager.getInstance(this).cancelAllWorkByTag("workspace-chat");
  }

  private void external(String u) {
    try {
      Uri uri = Uri.parse(u);
      if (!Arrays.asList("https", "mailto", "tel").contains(uri.getScheme())) return;
      startActivity(new Intent(Intent.ACTION_VIEW, uri));
    } catch (Exception e) {
      toast("No app can open this link");
    }
  }

  private void toast(String s) {
    Toast.makeText(this, s, Toast.LENGTH_LONG).show();
  }

  @Override
  public void onBackPressed() {
    navigateBack();
  }

  private void navigateBack() {
    if (web != null && web.canGoBack()) web.goBack();
    else if (this instanceof DocumentActivity) finish();
    else if (web != null || !toolbar.getTitle().equals("One Workspace")) home();
    else finish();
  }

  @Override
  protected void onNewIntent(Intent intent) {
    super.onNewIntent(intent);
    setIntent(intent);
    handleIntent(intent);
  }

  private void handleIntent(Intent intent) {
    if (intent != null && intent.getBooleanExtra("openChat", false)) {
      String topic = intent.getStringExtra("topic");
      pendingTopic =
          SecureStore.isDevice(this, intent.getStringExtra("device"))
                  && ChatConnection.validTopic(topic)
              ? topic
              : "";
      openChat();
    }
  }

  @Override
  protected void onResume() {
    super.onResume();
    foreground = true;
    WorkspacePushService.start(this);
    if (SecureStore.get(this, "device") != null) ChatSyncWorker.enqueue(this, null, null, false);
  }

  @Override
  protected void onPause() {
    foreground = false;
    CookieManager.getInstance().flush();
    super.onPause();
  }

  @Override
  protected void onDestroy() {
    generation++;
    if (web != null) web.destroy();
    if (chooser != null) chooser.onReceiveValue(null);
    io.shutdown();
    super.onDestroy();
  }

  @Override
  public void onRequestPermissionsResult(int code, String[] p, int[] grants) {
    super.onRequestPermissionsResult(code, p, grants);
    if (code == 101) WorkspacePushService.start(this);
    if (code == 103 && locationCallback != null) {
      boolean allowed = false;
      for (int grant : grants) allowed |= grant == PackageManager.PERMISSION_GRANTED;
      locationCallback.invoke(locationOrigin, allowed, false);
      locationCallback = null;
    }
    if (code == 102 && media != null) {
      boolean all = true;
      for (int g : grants) all &= g == PackageManager.PERMISSION_GRANTED;
      if (all) media.grant(mediaAllowed);
      else media.deny();
      media = null;
    }
  }

  @Override
  protected void onActivityResult(int code, int result, Intent data) {
    super.onActivityResult(code, result, data);
    if (code == 100 && chooser != null) {
      chooser.onReceiveValue(WebChromeClient.FileChooserParams.parseResult(result, data));
      chooser = null;
    }
  }
}
