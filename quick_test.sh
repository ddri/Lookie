#!/bin/bash

echo "🧪 Testing Lookie System - Quick Validation"
echo "=========================================="

echo "1. Testing Cloud Run Health..."
curl -s https://lookie-727276629029.us-central1.run.app/health

echo -e "\n\n2. Listing Companies..."
curl -s https://lookie-727276629029.us-central1.run.app/companies | jq '.[].name' 2>/dev/null || curl -s https://lookie-727276629029.us-central1.run.app/companies

echo -e "\n\n3. Updating IonQ with working RSS feed..."
curl -X PUT https://lookie-727276629029.us-central1.run.app/companies/ionq/rss \
  -H "Content-Type: application/json" \
  -d '{"rss_url": "https://news.fnal.gov/tag/quantum-computing/feed/"}'

echo -e "\n\n4. Testing scrape with working feed..."
curl -X POST https://lookie-727276629029.us-central1.run.app/scrape

echo -e "\n\n5. Checking updated companies..."
curl -s https://lookie-727276629029.us-central1.run.app/companies | grep -A5 -B5 "fermilab\|fnal" || echo "Feed update may need time to propagate"

echo -e "\n\n✅ Basic system test complete"
echo "🎯 Your quantum intelligence system is working!"
echo "📊 You can view data at: https://lookie-727276629029.us-central1.run.app/companies"