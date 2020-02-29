package util

import "fmt"

// Description returns a description of how to use dist-midi.
func Description(server bool) string {
	s := "\nWelcome to dist-midi (%s)!\n"

	s += "By default, MIDI messages from the selected input MIDI device will be sent to the server.\n"
	s += "\tPress 'h' to turn on HELP mode, which prints info on the sound mapped to that MIDI note.\n"
	s += "\tWhile in HELP mode, MIDI messages are not sent to the server. Pressing 'h' again turns off HELP mode.\n"
	s += "\tpress 'a' to print all currently mapped sounds.\n"

	if server {
		return fmt.Sprintf(s, "server")
	}
	return fmt.Sprintf(s, "client")
}
