# Dekopin

Dekopin is a command-line tool for managing Google Cloud Run deployments, revisions, and traffic routing. It provides a streamlined workflow for deploying services to Cloud Run with support for revision tagging and traffic management.

## Features

- Deploy new Cloud Run revisions with or without traffic
- Create and manage revision tags
- Switch traffic between revisions
- Support for multiple deployment environments (local, GitHub Actions, Cloud Build)
- YAML-based configuration

## Installation

```bash
go install github.com/iwashi623/dekopin/cmd/dekopin@latest
```

## Configuration

Dekopin uses a YAML configuration file (`dekopin.yml` by default) with the following structure:

```yaml
project: your-gcp-project-id
region: gcp-region
service: your-cloud-run-service-name
runner: github-actions  # Or: cloud-build, local
```

## Usage

### Basic Commands

```bash
# Deploy a new revision with traffic
dekopin deploy --image gcr.io/project/image:tag

# Create a new revision without traffic
dekopin create-revision --image gcr.io/project/image:tag

# Assign a tag to a revision
dekopin create-tag --tag v1.0.0 --revision service-abcdef

# Remove a tag from a revision
dekopin remove-tag --tag v1.0.0

# Switch traffic to a specific revision
dekopin sr-deploy --revision service-abcdef
```

### Global Flags

```
--project    GCP project ID
--region     GCP region
--service    Cloud Run service name
--runner     Runner type (github-actions, cloud-build, local)
--file, -f   Config file path (default: dekopin.yml)
```

## CI/CD Integration

### GitHub Actions

Dekopin automatically detects GitHub Actions environments and can use environment variables for tag names and commit hashes.

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

Dekopin also supports Cloud Build integration using build environment variables.

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

## License

This project is licensed under the MIT License - see the LICENSE file for details.
