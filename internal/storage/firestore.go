package storage

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"

	"github.com/dryan/lookie/internal/models"
)

// FirestoreDB wraps the Firestore client for Lookie operations
type FirestoreDB struct {
	client    *firestore.Client
	projectID string
	ctx       context.Context
}

// NewFirestoreDB creates a new Firestore database connection
func NewFirestoreDB(ctx context.Context, projectID string, credentialsFile string) (*FirestoreDB, error) {
	// Initialize Firebase App
	var app *firebase.App
	var err error
	
	if credentialsFile != "" {
		// Use service account file for local development
		opt := option.WithCredentialsFile(credentialsFile)
		config := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, config, opt)
	} else {
		// Auto-discover credentials in Cloud Run
		config := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, config)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	// Get Firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}

	log.Printf("Connected to Firestore database: %s", projectID)

	return &FirestoreDB{
		client:    client,
		projectID: projectID,
		ctx:       ctx,
	}, nil
}

// Client returns the underlying Firestore client
func (f *FirestoreDB) Client() *firestore.Client {
	return f.client
}

// Close closes the Firestore client
func (f *FirestoreDB) Close() error {
	return f.client.Close()
}

// GetAllCompanies retrieves all active companies from Firestore
func (f *FirestoreDB) GetAllCompanies() ([]models.FirestoreCompany, error) {
	iter := f.client.Collection("companies").Where("is_active", "==", true).Documents(f.ctx)
	defer iter.Stop()

	var companies []models.FirestoreCompany
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		var company models.FirestoreCompany
		if err := doc.DataTo(&company); err != nil {
			continue
		}

		companies = append(companies, company)
	}

	return companies, nil
}

// CheckArticleExists checks if an article with the given content hash already exists
func (f *FirestoreDB) CheckArticleExists(contentHash string) (bool, error) {
	iter := f.client.Collection("articles").Where("content_hash", "==", contentHash).Limit(1).Documents(f.ctx)
	defer iter.Stop()
	
	_, err := iter.Next()
	if err != nil {
		return false, nil // No document found
	}
	
	return true, nil // Document exists
}