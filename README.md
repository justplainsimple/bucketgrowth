# BucketGrowth

A command-line tool for tracking and forecasting AWS S3 bucket size and object count growth.

## Overview

BucketGrowth analyzes CloudWatch metrics for an S3 bucket and provides:

- Current size and object count
- Monthly and yearly growth rates
- Size and object count projections for 1 and 5 years

This tool helps with capacity planning, cost forecasting, and understanding storage growth patterns.

## Installation

### Binary Installation

1. Download the latest version from the [Releases](https://github.com/justplainsimple/bucketgrowth/releases) page
2. Extract the archive to a directory in your `PATH`
3. Make the file executable if needed: `chmod +x bucketgrowth`

### Building from Source

Requirements:
- Go 1.13 or higher
- AWS SDK for Go

```bash
# Clone the repository
git clone https://github.com/justplainsimple/bucketgrowth.git
cd bucketgrowth

# Build the binary
make build

# The binary will be available at .build/bucketgrowth
```

## Usage

```bash
bucketgrowth [options] BUCKET_NAME
```

### Example

```bash
# Basic usage
bucketgrowth my-bucket

# With AWS profile
bucketgrowth --profile production my-bucket

# With JSON output
bucketgrowth --output json my-bucket
```

### Output Example

```
Bucket Growth
=============

Total Size: 1.2 GB
Total Objects: 1,354

Size Growth: 3.41%/mo, 41.30%/yr
Object Growth: 2.87%/mo, 34.44%/yr

Size Projection: 1.7 GB (1 yr), 4.3 GB (5 yr)
Object Count Projection: 1,833 (1 yr), 3,415 (5 yr)
```

## Options

| Option | Description |
|--------|-------------|
| `--profile` | AWS profile to use (also uses `AWS_PROFILE` environment variable) |
| `--region` | AWS region to use (also uses `AWS_DEFAULT_REGION` environment variable) |
| `--verbose`, `-v` | Enable verbose logging |
| `--output TYPE` | Change output format to `text` (default) or `json` |
| `--skip-banner` | Suppress the banner in text output |

## AWS Credentials

BucketGrowth uses the standard AWS SDK credential resolution:

1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. Shared credential file (`~/.aws/credentials`)
3. IAM role for Amazon EC2 or ECS task role

Your IAM user or role needs the following permissions:
- `cloudwatch:GetMetricStatistics`
- `s3:ListAllMyBuckets` (if listing buckets)

## How It Works

BucketGrowth retrieves CloudWatch metrics for the specified S3 bucket:
- `BucketSizeBytes` metric for bucket size
- `NumberOfObjects` metric for object count

Growth calculations:
- Monthly growth rate is calculated based on the average month-over-month change
- Yearly growth rate is calculated based on the change between the earliest and latest data points
- Projections use compound growth formulas applied to the monthly growth rate

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development

```bash
# Format code
make format

# Run tests
make test

# Run linter
make vet

# Run security scanner
make sast

# Generate test coverage report
make coverage
```

## License

[MIT](LICENSE)
