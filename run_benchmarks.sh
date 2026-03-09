#!/bin/bash

echo "=========================================="
echo "  Structured Logger - Performance Tests"
echo "=========================================="
echo ""

echo "Running benchmarks..."
echo ""

go test ./benchmarks -bench=. -benchmem | grep -E "Benchmark|ns/op"

echo ""
echo "=========================================="
echo "Benchmark complete!"
echo "See BENCHMARKS.md for detailed analysis"
echo "=========================================="
