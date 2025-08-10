# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

BucketGrowth is a Go CLI application that analyzes AWS S3 bucket growth patterns using CloudWatch metrics. It calculates current size/object counts, growth rates, and provides 1-year and 5-year projections for capacity planning.

## Architecture

- **Main entry point**: `cmd/bucketgrowth/main.go` - CLI setup using urfave/cli/v2
- **Core logic**: `cmd/bucketgrowth/core.go` - Main application flow and output formatting
- **Business logic**: Root package `bucketgrowth` contains:
  - `cloudwatch.go` - AWS CloudWatch integration and data retrieval
  - `growth.go` - Growth rate calculations using monthly averages
  - `projection.go` - Future projections using compound growth formulas
- **Data structures**: `Metrics` struct contains all calculated values and projections

## Development Commands

```bash
# Build the application
make build

# Run tests
make test

# Format code
make format

# Run linter
make vet

# Security analysis
make sast

# Generate coverage report
make coverage
```

## Key Technical Details

- Uses AWS SDK for Go to retrieve CloudWatch metrics (`BucketSizeBytes` and `NumberOfObjects`)
- Growth calculations average 12 months of month-over-month changes
- Supports both text and JSON output formats
- AWS credentials resolved via standard SDK chain (env vars, profile, IAM roles)
- Binary built to `.build/bucketgrowth`

## Testing

Tests cover all core functionality including growth calculations and projections. Run individual test files with:
```bash
go test ./cloudwatch_test.go
go test ./growth_test.go  
go test ./projection_test.go
```