#!/bin/bash
# Test execution script for the generated unit tests

set -e

echo "=========================================="
echo "Running Unit Tests for Git Diff Changes"
echo "=========================================="
echo ""

echo "📦 Testing internal/model package..."
go test -v ./internal/model/... 2>&1 | grep -E "(PASS|FAIL|RUN|---)"
echo ""

echo "📦 Testing internal/delivery/http package..."
go test -v ./internal/delivery/http/... 2>&1 | grep -E "(PASS|FAIL|RUN|---)"
echo ""

echo "📦 Testing internal/service package..."
go test -v ./internal/service/... 2>&1 | grep -E "(PASS|FAIL|RUN|---)"
echo ""

echo "📦 Testing cmd package..."
go test -v ./cmd/... 2>&1 | grep -E "(PASS|FAIL|RUN|---)"
echo ""

echo "=========================================="
echo "Running tests with coverage..."
echo "=========================================="
echo ""

go test -cover ./internal/model/... 2>&1 | grep -v "telemetry"
go test -cover ./internal/delivery/http/... 2>&1 | grep -v "telemetry"
go test -cover ./internal/service/... 2>&1 | grep -v "telemetry"
go test -cover ./cmd/... 2>&1 | grep -v "telemetry"

echo ""
echo "✅ All tests completed!"