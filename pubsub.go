package util

import (
	"context"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// newPubsubClient returns an authenticated GCP pubsub client.
func newPubsubClient() *pubsub.Client {
	ctx := context.Background()
	creds, err := google.CredentialsFromJSON(ctx, []byte(GCPCredentials), pubsub.ScopePubSub)
	Must(err)

	client, err := pubsub.NewClient(ctx, "midi-pub-sub", option.WithCredentials(creds))
	Must(err)

	return client
}

// GetPubsubTopic returns a pubsub client the dist-midi topic.
func GetPubsubTopic() (*pubsub.Client, *pubsub.Topic) {
	client := newPubsubClient()
	return client, client.Topic("dist-midi")
}
