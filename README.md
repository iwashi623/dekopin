# Dekopin

Dekopin is a command-line tool for managing Google Cloud Run deployments, revisions, and traffic routing. It provides an efficient workflow for deploying Cloud Run services with revision tagging and traffic management.

## Features

- Deploy new Cloud Run revisions with or without traffic
- Create and manage revision tags
- Switch traffic between revisions
- Support for multiple deployment environments (local, GitHub Actions, Cloud Build)
- YAML configuration format
- Built-in timeout handling (default 30 seconds)
- Consistent revision naming using commit hashes

## Installation

```bash
go install github.com/iwashi623/dekopin/cmd/dekopin@latest
```

## Configuration

Dekopin uses a YAML configuration file named `dekopin.yml` by default. The structure of the configuration file is as follows:

```yaml
project: your-gcp-project-id
region: gcp-region
service: your-cloud-run-service-name
runner: github-actions  # or: cloud-build, local
```

## Usage

### Basic Commands

```bash
# Deploy a new revision with traffic
dekopin deploy --image gcr.io/project/image:tag

# Create a new revision without traffic
dekopin create-revision --image gcr.io/project/image:tag

# Assign a tag to a revision
dekopin create-tag --tag v1-0-0 --revision service-abcdef

# Remove a tag from a revision
dekopin remove-tag --tag v1-0-0

# Switch traffic to a specific revision
dekopin sr-deploy --revision service-abcdef

# Switch traffic to a tag
dekopin st-deploy --tag v1-0-0
```

### Global Flags

```
--project    GCP project ID
--region     GCP region
--service    Cloud Run service name
--runner     Runner type (github-actions, cloud-build, local)
--file, -f   Path to configuration file (default: dekopin.yml)
```

## CI/CD Integration

### GitHub Actions

Dekopin automatically detects GitHub Actions environments and can use environment variables for tag names or commit hashes. The first 7 characters of the commit hash are used for revision naming.

Example workflow:

```yaml
name: Deploy to Cloud Run

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'
          
      - name: Install Dekopin
        run: go install github.com/iwashi623/dekopin/cmd/dekopin@latest
        
      - name: Deploy to Cloud Run
        run: dekopin deploy --image gcr.io/project/image:${{ github.sha }}
```

### Google Cloud Build

Dekopin also supports Cloud Build integration using build environment variables. If the commit hash is 7 characters or fewer, it will be used as is; if longer, only the first 7 characters are used.

Example `cloudbuild.yaml`:

```yaml
steps:
  - name: 'golang'
    entrypoint: 'go'
    args: ['install', 'github.com/iwashi623/dekopin/cmd/dekopin@latest']
  
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/$PROJECT_ID/image:$COMMIT_SHA', '.']
  
  - name: 'gcr.io/cloud-builders/docker'
    args: ['push', 'gcr.io/$PROJECT_ID/image:$COMMIT_SHA']
  
  - name: 'golang'
    entrypoint: 'dekopin'
    args: ['deploy', '--image', 'gcr.io/$PROJECT_ID/image:$COMMIT_SHA']
```

## Validation

Dekopin includes validation for various input values:

- Tags must consist of lowercase alphanumeric characters and hyphens (e.g., `release-v1`, `v1-0-0`)
- Commands have appropriate required flags
- Input values are validated before execution

## License

This project is licensed under the MIT License - see the LICENSE file for details.
