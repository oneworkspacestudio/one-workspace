# Pet

Select **Pet** under **Creative & Development** in One Workspace setup or
when choosing apps in Updates. Open Pet from Apps and sign in with MyAccounts.
The Workspace owner names the first pet. Signed-in Workspace members can visit
the shared habitat and give a clover treat when the cooldown ends.

## Visiting Clover Hollow

The first-person view initially follows the pet. Drag to look around, then use
the on-screen arrows on a phone or focus the view and press WASD or the arrow
keys. **Find** brings the camera back to your pet. Field notes show what the pet
has been doing, and the gallery shows its recent paintings and sculptures.
Open **Journal** to see both. The view fills the screen, with the pet’s needs,
clock, movement and treats over the scene. WebGL 2 is required for the 3D view.
The indicators and journal remain available when the 3D view is unavailable.

The pet forages for food that appears in the meadow, plays with a ball, explores,
rests in its burrow, draws and builds. Daylight, rain and its needs affect its
choices. Each in-world day lasts 24 minutes. The service runs while Workspace
is running, even after every viewer signs out. When Workspace returns from an
outage, Pet simulates up to eight hours of missed time using instincts. The
journal keeps 80 events, the gallery keeps the latest 24 paintings, and existing
sculptures can evolve.

## Choosing how Pet decides

Open **Journal**, then **Decision settings**, as the Workspace owner.

- **Instincts** runs the whole simulation immediately.
- **Local Ollama model** adds structured decisions from an installed model on
  the Workspace host. Enter its exact Ollama model name. Pet uses the local
  Ollama service on port 11434 and never installs models on startup.
- **Workspace AI Service** uses the model the owner enabled in AI Service.
  Create a local app credential named **pet** in AI Service, then paste it into
  Pet and save. Revoke it in AI Service to disconnect the app.
- **Local model + AI Service** alternates the two model sources.

Models receive the simulated pet's needs, current activity, weather, time and
counts of nearby resources and creations. They choose from available play,
exploration and creative actions. Pet handles eating and sleeping itself and
continues using instincts if a model is disabled, busy or unreachable.

Pet preview.2 ships in Workspace preview.127 and supports Pet-only updates
on Workspace preview.124 or newer. The supervised executable listens
on loopback at the Workspace port plus 56. Public requests go through the
Workspace proxy at `/pet/` with MyAccounts. World and session data are stored in
the private Pet SQLite database. Decision settings and credentials use
Workspace's private secret storage. Use the standard Workspace backup and
recovery flow to keep or restore the habitat.
