#!/bin/bash

# Script to run all escposimg examples
# This demonstrates that each example works independently

echo "Running all escposimg examples..."
echo "================================="

# Check if test image exists
if [ ! -f "test_pattern.png" ]; then
    echo "Error: test_pattern.png not found in examples directory"
    echo "Please copy an image file to examples/test_pattern.png"
    exit 1
fi

echo ""
echo "1. Running basic_usage.go..."
echo "----------------------------"
go run basic_usage.go > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ basic_usage.go completed successfully"
else
    echo "✗ basic_usage.go failed"
fi

echo ""
echo "2. Running dithering_comparison.go..."
echo "------------------------------------"
go run dithering_comparison.go > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ dithering_comparison.go completed successfully"
    echo "  Generated debug images for all dithering algorithms"
else
    echo "✗ dithering_comparison.go failed"
fi

echo ""
echo "3. Running output_methods.go..."
echo "------------------------------"
go run output_methods.go > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ output_methods.go completed successfully"
    echo "  Demonstrated stdout, file, and network output methods"
else
    echo "✗ output_methods.go failed"
fi

echo ""
echo "4. Running configuration_options.go..."
echo "-------------------------------------"
go run configuration_options.go > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ configuration_options.go completed successfully"
    echo "  Generated examples for different printer configurations"
else
    echo "✗ configuration_options.go failed"
fi

echo ""
echo "5. Running debug_features.go..."
echo "------------------------------"
go run debug_features.go > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✓ debug_features.go completed successfully"
    echo "  Demonstrated debugging and troubleshooting features"
else
    echo "✗ debug_features.go failed"
fi

echo ""
echo "All examples completed!"
echo "======================"
echo ""
echo "Generated files in this directory:"
ls -la *.png *.escpos 2>/dev/null | head -10
echo ""
echo "To run individual examples:"
echo "  go run basic_usage.go"
echo "  go run dithering_comparison.go"
echo "  go run output_methods.go"
echo "  go run configuration_options.go"
echo "  go run debug_features.go"
echo ""
echo "See README.md for detailed usage instructions."
