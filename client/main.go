package main

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
	util "github.com/rytrose/dist-midi"

	"gitlab.com/gomidi/midi"
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
	defer client.Close()

	// Read and publish MIDI
	rd := mid.NewReader()
	rd.Msg.Each = func(p *mid.Position, m midi.Message) {
		res := topic.Publish(context.Background(), &pubsub.Message{
			Data: m.Raw(),
		})
		_, err := res.Get(context.Background())
		if err != nil {
			fmt.Println(fmt.Sprintf("[WARN] Unable to publish MIDI message: %s", err))
		}
	}

	// Use WaitGroup to block
	var wg sync.WaitGroup
	wg.Add(1)

	// Listen for MIDI
	fmt.Println("Starting to publish MIDI...")
	go mid.ConnectIn(in, rd)

	wg.Wait()
}
