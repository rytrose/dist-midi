package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	util "github.com/rytrose/dist-midi"

	driver "gitlab.com/gomidi/rtmididrv"
)

func main() {
	// Create MIDI driver
	drv, err := driver.New()
	util.Must(err)

	// Close all open ports on close
	defer drv.Close()

	// Prompt for and open output port
	out, err := util.PromptForOutput(drv)
	util.Must(err)
	fmt.Println(fmt.Sprintf("Using output device (%d) %s", out.Number(), out.String()))
	util.Must(out.Open())

	// Connect to GCP pubsub
	client, _ := util.GetPubsubTopic()
	defer client.Close()

	sub := client.Subscription("server")

	// *Subscription.Receive blocks
	err = sub.Receive(context.Background(), func(c context.Context, m *pubsub.Message) {
		// ACK regardless
		m.Ack()
		err := out.Send(m.Data)
		if err != nil {
			fmt.Println(fmt.Sprintf("[WARN] Unable to send MIDI to output: %s", err))
		}
	})
	if err != context.Canceled {
		util.Must(err)
	}
}
