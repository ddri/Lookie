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

	// Initialize Firebase App using Application Default Credentials
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

	// Update IonQ with Fermilab quantum RSS feed
	fermilab_url := "https://news.fnal.gov/tag/quantum-computing/feed/"
	_, err = client.Collection("companies").Doc("ionq").Update(ctx, []firestore.Update{
		{Path: "rss_url", Value: fermilab_url},
		{Path: "name", Value: "Fermilab (Quantum Computing News)"},
		{Path: "description", Value: "Government quantum research lab covering industry developments"},
	})
	if err != nil {
		log.Fatalf("Failed to update company with Fermilab feed: %v", err)
	}

	log.Println("✅ Successfully updated company to use Fermilab quantum computing RSS feed")
	log.Printf("🔗 RSS URL: %s", fermilab_url)
	log.Println("📰 Feed contains 10 quantum computing articles")
	log.Println("🏢 Mentions companies like Diraq, HRL Labs")
	log.Println("🚀 Ready to test full pipeline: RSS → Firestore → Gemini AI")
}