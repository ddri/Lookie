package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

func main() {
	ctx := context.Background()
	projectID := "lookie-quantum-intelligence"

	// Initialize Firebase App
	config := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	// Get Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	// Update IonQ with TechCrunch RSS (temporary for testing)
	_, err = client.Collection("companies").Doc("ionq").Update(ctx, []firestore.Update{
		{Path: "rss_url", Value: "https://techcrunch.com/feed/"},
	})
	if err != nil {
		log.Fatalf("Failed to update IonQ RSS URL: %v", err)
	}

	log.Println("Successfully updated IonQ RSS URL to TechCrunch feed for testing")

	// Update Q-CTRL with IEEE RSS (temporary for testing)
	_, err = client.Collection("companies").Doc("q-ctrl").Update(ctx, []firestore.Update{
		{Path: "rss_url", Value: "https://spectrum.ieee.org/rss"},
	})
	if err != nil {
		log.Fatalf("Failed to update Q-CTRL RSS URL: %v", err)
	}

	log.Println("Successfully updated Q-CTRL RSS URL to IEEE Spectrum feed for testing")
}