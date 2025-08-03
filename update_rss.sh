#!/bin/bash

# Add an API endpoint to update RSS URLs
echo "Updating RSS URLs via API..."

# For now, let's test with TechCrunch RSS which we know works
curl -X POST https://lookie-727276629029.us-central1.run.app/companies/ionq/update-rss \
  -H "Content-Type: application/json" \
  -d '{"rss_url": "https://techcrunch.com/feed/"}'

curl -X POST https://lookie-727276629029.us-central1.run.app/companies/q-ctrl/update-rss \
  -H "Content-Type: application/json" \
  -d '{"rss_url": "https://spectrum.ieee.org/rss"}'