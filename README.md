# dist-midi
A way for my friends to play MIDI on my machine.

## Why?
I use [Farrago](https://rogueamoeba.com/farrago/) and [Loopback](https://rogueamoeba.com/loopback/) to play audio into Discord
when my friends and I play DnD. This has proven to be the most stable solution for playing diverse music and sound effects while on Discord. However, I want all of my friends to be able to trigger audio. Farrago exposes a MIDI interface, so I developed this tool to receive input MIDI from my friends' computers, publish MIDI messages through Google Cloud Platform pubsub, and consume MIDI messages on my machine, thus allowing my friends to control Farrago.
