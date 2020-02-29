package util

// MIDINote is a simple representation of a MIDI NoteOn/NoteOff
type MIDINote struct {
	IsOn     bool  `json:"isOn"`
	Key      uint8 `json:"key"`
	Velocity uint8 `json:"velocity"`
}
