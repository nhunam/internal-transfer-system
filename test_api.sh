#!/bin/bash

# Test script for Internal Transfer System API
# Make sure the server is running on localhost:8080

BASE_URL="http://localhost:8080"

echo "Testing Internal Transfer System API..."
echo "======================================"

# Test 1: Health check
echo "1. Testing health check endpoint..."
curl -X GET "$BASE_URL/health" | jq '.'
echo ""

# Test 2: Create first account
echo "2. Creating account 123 with initial balance 100.50..."
curl -X POST "$BASE_URL/accounts" \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": 123,
    "initial_balance": "100.50"
  }'
echo ""

# Test 3: Create second account
echo "3. Creating account 456 with initial balance 200.75..."
curl -X POST "$BASE_URL/accounts" \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": 456,
    "initial_balance": "200.75"
  }'
echo ""

# Test 4: Get account balance for account 123
echo "4. Getting balance for account 123..."
curl -X GET "$BASE_URL/accounts/123" | jq '.'
echo ""

# Test 5: Get account balance for account 456
echo "5. Getting balance for account 456..."
curl -X GET "$BASE_URL/accounts/456" | jq '.'
echo ""

# Test 6: Create transaction from account 123 to account 456
echo "6. Creating transaction: Transfer 25.25 from account 123 to account 456..."
curl -X POST "$BASE_URL/transactions" \
  -H "Content-Type: application/json" \
  -d '{
    "source_account_id": 123,
    "destination_account_id": 456,
    "amount": "25.25"
  }'
echo ""

# Test 7: Check balances after transaction
echo "7. Checking balances after transaction..."
echo "Account 123 balance:"
curl -X GET "$BASE_URL/accounts/123" | jq '.'
echo ""
echo "Account 456 balance:"
curl -X GET "$BASE_URL/accounts/456" | jq '.'
echo ""

# Test 8: Test error cases
echo "8. Testing error cases..."
echo "8a. Try to create duplicate account:"
curl -X POST "$BASE_URL/accounts" \
  -H "Content-Type: application/json" \
  -d '{
    "account_id": 123,
    "initial_balance": "50.00"
  }'
echo ""

echo "8b. Try to get non-existent account:"
curl -X GET "$BASE_URL/accounts/999" | jq '.'
echo ""

echo "8c. Try insufficient balance transaction:"
curl -X POST "$BASE_URL/transactions" \
  -H "Content-Type: application/json" \
  -d '{
    "source_account_id": 123,
    "destination_account_id": 456,
    "amount": "1000.00"
  }'
echo ""

echo "Testing completed!" 