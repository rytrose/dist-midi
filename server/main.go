package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	util "github.com/rytrose/dist-midi"

	"gitlab.com/gomidi/midi/mid"
	driver "gitlab.com/gomidi/rtmididrv"
)

func main() {
	// Create MIDI driver
	drv, err := driver.New()
	util.Must(err)

	// Close all open ports on close
	defer drv.Close()

	// Prompt for and open input port
	in, err := util.PromptForInput(drv)
	util.Must(err)
	fmt.Println(fmt.Sprintf("Using input device (%d) %s", in.Number(), in.String()))
	util.Must(in.Open())

	// Prompt for and open output port
	out, err := util.PromptForOutput(drv)
	util.Must(err)
	fmt.Println(fmt.Sprintf("Using output device (%d) %s", out.Number(), out.String()))
	util.Must(out.Open())

	// Create MIDI writer
	wr := mid.ConnectOut(out)

	// Read input MIDI and write to output
	rd := mid.NewReader()
	rd.Msg.Channel.NoteOn = func(p *mid.Position, channel, key, vel uint8) {
		err := wr.NoteOn(key, vel)
		if err != nil {
			fmt.Println(fmt.Sprintf("[WARN] Unable to send MIDI to output (local NoteOn): %s", err))
		}
	}
	rd.Msg.Channel.NoteOff = func(p *mid.Position, channel, key, vel uint8) {
		err := wr.NoteOff(key)
		if err != nil {
			fmt.Println(fmt.Sprintf("[WARN] Unable to send MIDI to output (local NoteOff): %s", err))
		}
	}

	// Start reading
	mid.ConnectIn(in, rd)

	// Connect to GCP pubsub
	client, _ := util.GetPubsubTopic()
	defer client.Close()
	sub := client.Subscription("server")

	// *Subscription.Receive blocks
	err = sub.Receive(context.Background(), func(c context.Context, m *pubsub.Message) {
		// ACK regardless
		m.Ack()

		note := &util.MIDINote{}
		err := json.Unmarshal(m.Data, note)
		util.Must(err)

		if note.IsOn {
			err := wr.NoteOn(note.Key, note.Velocity)
			if err != nil {
				fmt.Println(fmt.Sprintf("[WARN] Unable to send MIDI to output (remote NoteOn): %s", err))
			}
		} else {
			err := wr.NoteOff(note.Key)
			if err != nil {
				fmt.Println(fmt.Sprintf("[WARN] Unable to send MIDI to output (remote NoteOff): %s", err))
			}
		}
	})
	if err != context.Canceled {
		util.Must(err)
	}
}
