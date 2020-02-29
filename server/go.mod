module github.com/rytrose/dist-midi/server

go 1.13

replace github.com/rytrose/dist-midi => ../

require (
	cloud.google.com/go/pubsub v1.2.0
	github.com/rytrose/dist-midi v0.0.0-00010101000000-000000000000
	gitlab.com/gomidi/midi v1.14.1
	gitlab.com/gomidi/rtmididrv v0.4.2
)
