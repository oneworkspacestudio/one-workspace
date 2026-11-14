# Pet

Select **Pet** under **Creative & Development** in One Workspace setup or
when choosing apps in Updates. Open Pet from Apps and sign in with MyAccounts.
The Workspace owner names the first pet. Signed-in Workspace members can visit
the shared habitat and give a clover treat when the cooldown ends.

## Visiting Clover Hollow

Each visit starts at the trail entrance. The pet continues wherever it was already headed. Drag to look around, then use
the on-screen arrows on a phone or focus the view and press WASD or the arrow
keys. **Find** brings the camera back to your pet. Field notes show what the pet
has been doing, and the gallery shows its recent paintings and structures.
Open **Journal** to see both. The view fills the screen, with the pet’s needs,
clock, movement and treats over the scene. WebGL 2 is required for the 3D view.
The indicators and journal remain available when the 3D view is unavailable.

The pet forages for food that appears in the meadow, plays with a ball, explores,
rests in its burrow, draws and builds. Daylight, rain and its needs affect its
choices. Each in-world day lasts 24 minutes. The service runs while Workspace
is running, even after every viewer signs out. When Workspace returns from an
outage, Pet simulates up to eight hours of missed time using instincts. The
journal records activity by habitat day and time. Older installations import
all entries still present in their previous 80-event journal. The gallery holds
the latest 24 paintings.

During clear daylight, Pet gathers stone, branches or reeds and builds in
three stages. Active projects and stock survive rest and restarts. Completing
projects teaches it to build arches, shelters and lookouts. It can rebuild a
site when the meadow is full.

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

Models see the pet’s simulated position, reachable landmarks, obstacles,
nearby food, material sources, current project, weather and recent activity.
They can pick a destination and a waypoint, compose a painting with geometric
strokes, or plan a structure from materials and learned designs. The pet checks
routes against the map and prioritises eating and sleeping. When a model is
unavailable, instincts use varied routes and compositions.

Pet preview.3 ships in Workspace preview.129 and supports Pet-only updates
on Workspace preview.124 or newer. The supervised executable listens
on loopback at the Workspace port plus 56. Public requests go through the
Workspace proxy at `/pet/` with MyAccounts. World and session data are stored in
the private Pet SQLite database. Decision settings and credentials use
Workspace's private secret storage. Use the standard Workspace backup and
recovery flow to keep or restore the habitat.
