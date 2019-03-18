
GameMusicRepeater
====

Game Music player for windows.
GameMusicRepeater can be distributed with a bundled playlist.


Usage
===

1. Write config.json to config list of bgm

config.hjson
```hjson
// This file is hjson.
// So you can write comment like this.
{
  // list of units
  "Songs":[
    {
      // One unit must contains Name 
      "Name":"Opening",
      // IntroFile will be played once when playing is started. 
      "IntroFile":"ogg/1.ogg",
      // File will be played repeatedly.
      "File":"ogg/2.ogg"
    },
    {
      // You can omitt IntroFile, then File will be played immediately. 
      "Name":"Club",
      "File":"ogg/ggg.ogg"
    },
    {
      "Name":"Dungeon",
      // mp3 files will be played 
      "File":"ogg/test.mp3"
    }
  ]
}


```

2. Run the GameMusicRepeater.exe with config.hjson in same directory.


Note
---

We use this player to distribute SoundTracks of 矢澤にこは夢を見る, is a fungame we made.
This may be good to distribute your own game musics.