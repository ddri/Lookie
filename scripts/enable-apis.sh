#!/bin/bash

# Enable required APIs for Lookie Firebase project
# Run this after billing is configured

set -e

PROJECT_ID="lookie-quantum-intelligence"
echo "🔧 Enabling APIs for project: $PROJECT_ID"

# Set project context
/Users/dryan/google-cloud-sdk/bin/gcloud config set project $PROJECT_ID

# Required APIs for Firebase + Lookie
APIS=(
    "firebase.googleapis.com"
    "firestore.googleapis.com" 
    "cloudfunctions.googleapis.com"
    "run.googleapis.com"
    "cloudbuild.googleapis.com"
    "cloudscheduler.googleapis.com"
    "aiplatform.googleapis.com"
    "storage.googleapis.com"
    "secretmanager.googleapis.com"
)

echo "Enabling APIs..."
for api in "${APIS[@]}"; do
    echo "  - Enabling $api..."
    /Users/dryan/google-cloud-sdk/bin/gcloud services enable $api --quiet
    echo "    ✅ $api enabled"
done

echo ""
echo "🎉 All APIs enabled successfully!"
echo ""
echo "Next: Add Firebase to your project at https://console.firebase.google.com/"