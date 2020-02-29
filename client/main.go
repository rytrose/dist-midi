package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
	util "github.com/rytrose/dist-midi"

	"gitlab.com/gomidi/midi/mid"
	driver "gitlab.com/gomidi/rtmididrv"
)

const helpKey = 'h'
const allKey = 'a'

func main() {
	// Get SoundMap
	soundMap := util.GetSoundMap()

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

	// Set up keyboard reading
	kr := util.NewKeyboardReader()
	kr.Register(helpKey, func(prev bool, current bool) {
		if current {
			fmt.Println("[HELP] Now printing info of played MIDI-mapped sounds.")
		} else {
			fmt.Println("[SENDING] Now sending MIDI to play mapped sounds.")
		}
	})
	kr.Register(allKey, func(prev bool, current bool) {
		fmt.Print(soundMap.String())
	})
	kr.Read()
	defer kr.Close()

	// Connect to GCP pubsub
	client, topic := util.GetPubsubTopic()
	defer topic.Stop()
	defer client.Close()

	// Read and publish MIDI
	rd := mid.NewReader(mid.NoLogger())
	rd.Msg.Channel.NoteOn = func(p *mid.Position, channel, key, vel uint8) {
		if kr.GetState(helpKey) {
			sound, ok := soundMap.GetEntry(int(key))
			if ok {
				fmt.Println(fmt.Sprintf("[HELP] MIDI Note: %d, Title: %s, Hold To Play: %t, Allow Pausing: %t, Loop: %t",
					key, sound.Title, sound.HoldToPlay, sound.AllowPausing, sound.Loop))
			} else {
				fmt.Println(fmt.Sprintf("[HELP] No MIDI sound mapped to MIDI note %d.", key))
			}
		} else {
			data, err := json.Marshal(&util.MIDINote{
				IsOn:     true,
				Key:      key,
				Velocity: vel,
			})
			util.Must(err)
			publish(topic, data)
		}
	}
	rd.Msg.Channel.NoteOff = func(p *mid.Position, channel, key, vel uint8) {
		if !kr.GetState(helpKey) {
			data, err := json.Marshal(&util.MIDINote{
				IsOn:     false,
				Key:      key,
				Velocity: 0,
			})
			util.Must(err)
			publish(topic, data)
		}
	}

	// Use WaitGroup to block
	var wg sync.WaitGroup
	wg.Add(1)

	// Listen for MIDI
	go mid.ConnectIn(in, rd)

	// Print description and state
	fmt.Println(util.Description(false))
	fmt.Println("[SENDING] Now sending MIDI to play mapped sounds.")

	wg.Wait()
}

func publish(topic *pubsub.Topic, data []byte) {
	res := topic.Publish(context.Background(), &pubsub.Message{
		Data: data,
	})
	_, err := res.Get(context.Background())
	if err != nil {
		fmt.Println(fmt.Sprintf("[WARN] Unable to publish MIDI message: %s", err))
	}
}
