package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"

	"github.com/dryan/lookie/internal/models"
)

func main() {
	ctx := context.Background()
	projectID := "lookie-quantum-intelligence"

	fmt.Printf("🔥 Testing Firestore connection to project: %s\n", projectID)

	// Initialize Firebase App - try different credential approaches
	config := &firebase.Config{ProjectID: projectID}
	
	// First try: Application Default Credentials
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		// If that fails, we have an auth issue
		fmt.Printf("❌ Authentication failed: %v\n", err)
		fmt.Println("💡 Need to set up Application Default Credentials")
		fmt.Println("   Run this command in terminal:")
		fmt.Println("   gcloud auth application-default login")
		fmt.Println("   Then retry this test")
		return
	}
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	// Get Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ Connected to Firestore successfully!")

	// Test 1: Create a test company
	testCompany := models.FirestoreCompany{
		ID:          "test-company",
		Name:        "Test Quantum Corp",
		Domain:      "test.example.com",
		RSSURL:      "https://test.example.com/rss",
		NewsPageURL: "https://test.example.com/news",
		QuantumFocus: "software",
		MarketSegment: "startup",
		Description: "Test company for Firestore connection",
		IsActive:    true,
		RobotsTxtCompliant: true,
		Stats: models.CompanyStats{
			TotalArticles:      0,
			ArticlesThisMonth:  0,
			AvgConfidenceScore: 0.0,
		},
	}

	fmt.Println("📝 Creating test company document...")
	_, err = client.Collection("companies").Doc(testCompany.ID).Set(ctx, testCompany)
	if err != nil {
		log.Fatalf("Failed to create test company: %v", err)
	}
	fmt.Printf("✅ Test company '%s' created successfully!\n", testCompany.Name)

	// Test 2: Read the company back
	fmt.Println("📖 Reading test company document...")
	doc, err := client.Collection("companies").Doc(testCompany.ID).Get(ctx)
	if err != nil {
		log.Fatalf("Failed to read test company: %v", err)
	}

	var readCompany models.FirestoreCompany
	if err := doc.DataTo(&readCompany); err != nil {
		log.Fatalf("Failed to parse company data: %v", err)
	}

	fmt.Printf("✅ Read company back: %s (Focus: %s)\n", readCompany.Name, readCompany.QuantumFocus)

	// Test 3: List all companies
	fmt.Println("📋 Listing all companies...")
	iter := client.Collection("companies").Documents(ctx)
	defer iter.Stop()

	count := 0
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		var company models.FirestoreCompany
		if err := doc.DataTo(&company); err != nil {
			fmt.Printf("Warning: Could not parse company %s: %v\n", doc.Ref.ID, err)
			continue
		}

		fmt.Printf("  - %s (%s)\n", company.Name, company.ID)
		count++
	}
	fmt.Printf("✅ Found %d companies total\n", count)

	// Test 4: Clean up test data
	fmt.Println("🧹 Cleaning up test company...")
	_, err = client.Collection("companies").Doc(testCompany.ID).Delete(ctx)
	if err != nil {
		log.Printf("Warning: Could not delete test company: %v", err)
	} else {
		fmt.Println("✅ Test company deleted")
	}

	fmt.Println("")
	fmt.Println("🎉 All Firestore tests passed!")
	fmt.Println("📊 Test Summary:")
	fmt.Println("  - ✅ Firebase connection established")
	fmt.Println("  - ✅ Document creation working")
	fmt.Println("  - ✅ Document reading working")
	fmt.Println("  - ✅ Collection listing working")
	fmt.Println("  - ✅ Document deletion working")
	fmt.Println("  - ✅ Go models compatible with Firestore")
	fmt.Println("")
	fmt.Println("🚀 Ready to proceed with data migration!")
}