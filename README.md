![CSGolang](https://i.imgur.com/LNNcd3u.png "Logo CSGolang")

# CSGolang
 Allows users to see enemies through walls in the video game Counter-Strike: Global Offense (CS:GO), and adds other advantages.
 
 Reads and writes to the memory of CS:GO via [WinAPI](https://en.wikipedia.org/wiki/Windows_API), and draws an overlay using https://github.com/faiface/pixel
 
 > todo: add gif here 
## Features
- Configured via GET/POST over a [Svelte UI](https://github.com/sveltejs/svelte)
- Bunnyhop - Automates a technique increasing movement up to 20%.
- ESP/X-Ray - View enemies name and skeleton through walls.
  - Colored based on team or a health gradient.
- Updates offsets/netvars via GET https://github.com/frk1/hazedumper/blob/master/csgo.json

 ## Usage
 1. Download and extract the [latest release](https://github.com/f0nkey/F0nkHack/releases).
 2. Set -insecure flag on CS:GO (or risk suspension from game)
 3. Open CS:GO, set **fullscreen windowed** mode at 1920x1080.
 4. Run CSGOlang.exe.
 5. Navigate to http://localhost:8085 in your web browser (or Steam browser) to change configuration.
 
## Previews

![CSGolang](https://i.imgur.com/F1ypEnr.gif "CS UI Preview")
![CSGolang](https://thumbs.gfycat.com/NeighboringEasygoingAfricanbushviper-small.gif "CS Wall Preview")
