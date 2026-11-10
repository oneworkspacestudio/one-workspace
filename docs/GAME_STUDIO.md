# Game Studio

Game Studio turns a written idea into a small browser game. Choose a starter or
describe an idea to DeepSeek. Ask for changes in the same conversation and each
successful build opens in a new browser tab. If the browser blocks the tab,
choose **Play latest build** or **Play** in the arcade.

Four starters are included: treasure catcher, maze, runner, and matching cards.
DeepSeek chooses the template, title, character, collectible, theme, and
supported difficulty settings. Maze routes and memory board size are fixed.
Games support keyboard and touch. The app bundles the game engine and art
choices. Conversation is text. The starter catalogue is internal to Game Studio.

## Set up for your child

1. Select **Game Studio** under **Creative & Development** in One Workspace
   setup. Install **My Git** and **My Hosting** if you plan to publish games.
2. Have your child sign in with their own MyAccounts account. Copy the account
   code shown by Game Studio into **Grown-up settings** from the Workspace
   owner's session. Clear the code and save to revoke access.
3. Choose **Local DeepSeek via Ollama** and enter an installed DeepSeek R1 tag
   such as `deepseek-r1:8b`. Ollama must listen on loopback port 11434 and
   provide its OpenAI-compatible chat endpoint.
   Alternatively select **DeepSeek Cloud**, enter your API key, and opt in to
   sending game ideas to DeepSeek. Review its
   [API terms](https://cdn.deepseek.com/policies/en-US/deepseek-open-platform-terms-of-service.html)
   and [privacy policy](https://cdn.deepseek.com/policies/en-US/deepseek-privacy-policy.html)
   together before enabling cloud mode. Leave personal information out of ideas.

Generation is off by default. The service accepts one generation request at a
time and up to 30 turns per UTC day. It keeps up to 50 games and 30 builds of
each. Saved games contain the rendered game and validated settings. Conversation
text stays in the active browser page and recent turns are sent to the selected
model with each request. Keep a grown-up involved in cloud use and publishing.

## Publish a game

Create a dedicated, initialized **private** My Git repository and a token for a
Git identity with write access to it. Disable repository automation. Enter the
repository as `owner/repository` and the token in **Grown-up settings**. Game
Studio encrypts the saved keys with the Workspace platform secret store.

Sign into My Hosting with the same MyAccounts owner. Choose **Publish** on a
game, review the sharing notice, then choose **Save and deploy**. Game Studio
commits the selected build as an HTML file in My Git. It deploys a static ZIP
containing the same HTML through My Hosting and opens the published address in
a separate tab. Public availability follows My Hosting's publication settings.
If a deployment fails, the Git copy remains and the owner can retry. Use
My Hosting to inspect, roll back, or remove a publication.

Private play links require a Game Studio session with access to the game. A
publicly hosted link can be shared independently of MyAccounts. Preview and
publication use trusted template code and a restricted browser policy. The
model supplies validated data and has no shell, repository, or Hosting access.

## Managed runtime

Game Studio is a native Workspace process under `/gamestudio`, with private
loopback port 9045 (manager port plus 55 on rebased installs). It uses the
standard MyAccounts callback, launcher, logout, backup, and per-app recovery
contracts. Builds and encrypted settings live in the app's own private-data
directory. The game engine and browser UI are embedded, so setup requires no
separate JavaScript runtime. The ordinary organization rebuild packages the
executable, compatibility metadata, guide, and notices under the
`gamestudio` package namespace.
