# CSGolang
 Allows users to see enemies through walls in the video game Counter-Strike: Global Offense (CSGO), and adds other advantages.
 
 Reads and writes to the memory of CS:GO via [WinAPI](https://en.wikipedia.org/wiki/Windows_API), and draws an overlay using https://github.com/faiface/pixel
 
 > todo: add gif here 
## Features
- Configured via GET/POST over a [Svelte UI](https://github.com/sveltejs/svelte)
- Bunnyhop - Automates technique increasing movement up to 20%.
- ESP/X-Ray - View enemies name and skeleton through walls.
  - Colored based on team or a health gradient.
- Updates offsets/netvars via GET https://github.com/frk1/hazedumper/blob/master/csgo.json


 ## Usage
 1. Download and extract the [latest release](https://github.com/f0nkey/F0nkHack/releases).
 2. Open CS:GO in **fullscreen windowed** mode at 1920x1080.
 3. Run CGOlang.exe.
 4. Navigate to http://localhost:8085 in your web browser (or Steam browser) to change configuration.
 
