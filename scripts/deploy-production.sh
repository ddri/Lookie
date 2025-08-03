#!/bin/bash

# Deploy Lookie to Google Cloud Run
# Usage: ./scripts/deploy-production.sh

set -e

# Configuration
PROJECT_ID="lookie-quantum-intelligence"
REGION="us-central1"
SERVICE_NAME="lookie"
IMAGE_NAME="us-central1-docker.pkg.dev/${PROJECT_ID}/lookie/${SERVICE_NAME}"

echo "🚀 Deploying Lookie to Google Cloud Run..."

# Check if gcloud is authenticated
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | head -n1 > /dev/null; then
    echo "❌ Not authenticated with gcloud. Please run: gcloud auth login"
    exit 1
fi

# Set the project
echo "📋 Setting project to ${PROJECT_ID}..."
gcloud config set project ${PROJECT_ID}

# Enable required APIs if not already enabled
echo "🔧 Enabling required APIs..."
gcloud services enable cloudbuild.googleapis.com \
    run.googleapis.com \
    artifactregistry.googleapis.com \
    --quiet

# Build and submit to Cloud Build
echo "🏗️  Building container image..."
gcloud builds submit --tag ${IMAGE_NAME}:latest .

# Deploy to Cloud Run
echo "🚀 Deploying to Cloud Run..."
gcloud run deploy ${SERVICE_NAME} \
    --image ${IMAGE_NAME}:latest \
    --region ${REGION} \
    --platform managed \
    --allow-unauthenticated \
    --port 8080 \
    --memory 512Mi \
    --cpu 1 \
    --min-instances 0 \
    --max-instances 10 \
    --timeout 300s \
    --set-env-vars "GOOGLE_CLOUD_PROJECT=${PROJECT_ID}" \
    --quiet

# Get the service URL
SERVICE_URL=$(gcloud run services describe ${SERVICE_NAME} \
    --region ${REGION} \
    --format 'value(status.url)')

echo ""
echo "✅ Deployment complete!"
echo ""
echo "🌐 Service URL: ${SERVICE_URL}"
echo "🏥 Health Check: ${SERVICE_URL}/health"
echo "🏢 Companies API: ${SERVICE_URL}/companies"
echo "📄 Manual Scrape: curl -X POST ${SERVICE_URL}/scrape"
echo ""
echo "📊 Monitor your deployment:"
echo "   Cloud Run Console: https://console.cloud.google.com/run/detail/${REGION}/${SERVICE_NAME}"
echo "   Logs: gcloud logs tail --follow --format='value(textPayload)' --filter='resource.labels.service_name=${SERVICE_NAME}'"
echo ""

# Test the health endpoint
echo "🧪 Testing health endpoint..."
if curl -f "${SERVICE_URL}/health" > /dev/null 2>&1; then
    echo "✅ Health check passed!"
else
    echo "⚠️  Health check failed - service may still be starting up"
    echo "   Check logs: gcloud logs tail --filter='resource.labels.service_name=${SERVICE_NAME}'"
fi

echo ""
echo "🎉 Lookie is now running in production!"