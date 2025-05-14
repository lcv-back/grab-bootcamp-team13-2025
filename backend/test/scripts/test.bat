@echo off
echo Starting integration tests...
echo Testing health endpoint...
curl -f http://localhost:3001/health

echo.
echo Testing predict endpoint...
curl -X POST ^
     -H "Content-Type: application/json" ^
     -H "Authorization: Bearer %JWT_TOKEN%" ^
     -d @../test/data/predict_request.json ^
     https://isymptom.id.vn/api/symptoms/predict

echo.
echo Testing followup endpoint...
curl -X POST ^
     -H "Content-Type: application/json" ^
     -H "Authorization: Bearer %JWT_TOKEN%" ^
     -d @../test/data/followup_request.json ^
     https://isymptom.id.vn/api/symptoms/followup