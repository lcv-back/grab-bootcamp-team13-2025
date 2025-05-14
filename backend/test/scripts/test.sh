#!/bin/sh
echo "Starting integration tests..."
echo "Testing health endpoint..."
curl -f http://localhost:3001/health

echo "\nTesting predict endpoint..."
curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer ${JWT_TOKEN}" \
     -d @./test/data/predict_request.json \
     https://isymptom.id.vn/api/symptoms/predict

echo "\nTesting followup endpoint..."
curl -X POST \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer ${JWT_TOKEN}" \
     -d @./test/data/followup_request.json \
     https://isymptom.id.vn/api/symptoms/followup