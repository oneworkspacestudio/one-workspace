package uk.co.m8e.workspace;

/** Initialise native channels; a visible user action starts background messaging. */
public final class WorkspaceApplication extends android.app.Application {
  @Override
  public void onCreate() {
    super.onCreate();
    ChatNotifications.channels(this);
  }
}
