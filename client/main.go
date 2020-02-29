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

	// Connect to GCP pubsub
	client, topic := util.GetPubsubTopic()
	defer topic.Stop()
	defer client.Close()

	// Read and publish MIDI
	rd := mid.NewReader()
	rd.Msg.Channel.NoteOn = func(p *mid.Position, channel, key, vel uint8) {
		data, err := json.Marshal(&util.MIDINote{
			IsOn:     true,
			Key:      key,
			Velocity: vel,
		})
		util.Must(err)
		publish(topic, data)
	}
	rd.Msg.Channel.NoteOff = func(p *mid.Position, channel, key, vel uint8) {
		data, err := json.Marshal(&util.MIDINote{
			IsOn:     false,
			Key:      key,
			Velocity: 0,
		})
		util.Must(err)
		publish(topic, data)
	}

	// Use WaitGroup to block
	var wg sync.WaitGroup
	wg.Add(1)

	// Listen for MIDI
	fmt.Println("Starting to publish MIDI...")
	go mid.ConnectIn(in, rd)

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
