# AI Fieldnotes

AI Fieldnotes is a private reading room for AI engineering leads. It opens with
three short articles on loop engineering, AI maturity and accessible interfaces.

## Research and read

1. In Workspace setup, select **AI Fieldnotes** under **Creative & Development**.
   Open it from **Apps** and sign in as the Workspace owner.
2. Open **Settings** and save your DeepSeek API key. DeepSeek Search finds public
   web sources. DeepSeek also writes briefings from the sources you choose.
   Search requests and model usage are billed by your provider.
3. Open **Research**, enter a phrase, or choose a topic under **On the radar**.
   Review the results before choosing up to five and selecting **Write briefing**.
4. Follow the numbered citations in the article to inspect the source list. A
   generated briefing is a synthesis. Consult the original publication when
   you need to confirm a standard, requirement or reported result.

### Use the Workspace browser

When **Web Search** is installed, you can create an API credential there and
save it in AI Fieldnotes Settings. Enable and configure Web Search first.
Select **Workspace Web Search** as the provider for a research query. The app
uses the browser's observed source text and sends the selected evidence to
DeepSeek to write the briefing. Review which browser profile you use and search
public material you are comfortable sharing with your model provider.

## Read offline on Android

Download [AI Fieldnotes for Android](https://github.com/oneworkspacestudio/one-workspace/releases/download/v0.3.0-preview.145/AIFieldnotes-Android-1.0.0.apk).
The app includes the same three starter articles.

Save articles in the web app, open **Reading list**, and choose **Download
reading pack**. Move the JSON file to your phone. Open the native reader and
choose **Import pack** to add the articles to its reading list. Import merges
briefings by ID. The offline reader includes the article text, citation titles,
and original URLs. The file contains the selected articles in plain text, so
store and transfer it with the same care as the underlying material. A pack and
the device reading list each support up to 100 briefings and 2 MB.

## Storage and updates

The supervised service listens on a private loopback port. The public app
route is `/fieldnotes/`. Workspace controls sign-in, owner authorization,
session expiry, health, and app updates. Library, saved articles, and recent
research history live in the dedicated private data directory. Provider
credentials use Workspace's platform secret storage and are never included in
reading packs. The native reader stores imports privately and has no network
permission. Current research runs on demand in the Workspace service. Android
receives reading material through explicit pack imports.
