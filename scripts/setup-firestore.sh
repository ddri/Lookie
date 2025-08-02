#!/bin/bash

# Setup Firestore database for Lookie
# Run this after APIs are enabled

set -e

PROJECT_ID="lookie-quantum-intelligence"
REGION="us-central1"

echo "🗃️ Setting up Firestore database..."

# Set project context
/Users/dryan/google-cloud-sdk/bin/gcloud config set project $PROJECT_ID

# Create Firestore database
echo "Creating Firestore database in $REGION..."
/Users/dryan/google-cloud-sdk/bin/gcloud firestore databases create \
    --location=$REGION \
    --type=firestore-native \
    --quiet

echo "✅ Firestore database created successfully!"

# Create initial security rules
cat > firestore.rules << 'EOF'
rules_version = '2';

service cloud.firestore {
  match /databases/{database}/documents {
    // Development rules - replace with production rules later
    match /{document=**} {
      allow read, write: if request.auth != null;
    }
  }
}
EOF

echo "✅ Created basic Firestore security rules"

# Create indexes file
cat > firestore.indexes.json << 'EOF'
{
  "indexes": [
    {
      "collectionGroup": "articles",
      "queryScope": "COLLECTION",
      "fields": [
        {
          "fieldPath": "company_id",
          "order": "ASCENDING"
        },
        {
          "fieldPath": "published_at",
          "order": "DESCENDING" 
        }
      ]
    },
    {
      "collectionGroup": "articles",
      "queryScope": "COLLECTION", 
      "fields": [
        {
          "fieldPath": "is_processed",
          "order": "ASCENDING"
        },
        {
          "fieldPath": "scraped_at",
          "order": "ASCENDING"
        }
      ]
    },
    {
      "collectionGroup": "classifications",
      "queryScope": "COLLECTION_GROUP",
      "fields": [
        {
          "fieldPath": "category",
          "order": "ASCENDING"
        },
        {
          "fieldPath": "confidence_score", 
          "order": "DESCENDING"
        }
      ]
    }
  ],
  "fieldOverrides": []
}
EOF

echo "✅ Created Firestore indexes configuration"

echo ""
echo "🎉 Firestore setup complete!"
echo ""
echo "Database URL: https://console.firebase.google.com/project/$PROJECT_ID/firestore"