# Dekopin

Dekopin is a command-line tool for managing Google Cloud Run deployments, revisions, and traffic routing. It provides an efficient workflow for deploying Cloud Run services with revision tagging and traffic management.

## Features

- Deploy new Cloud Run revisions with or without traffic
- Create and manage revision tags
- Switch traffic between revisions
- Support for multiple deployment environments (local, GitHub Actions, Cloud Build)
- YAML configuration format
- Built-in timeout handling (default 120 seconds)
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

### Global Flags

The following flags can be used with any subcommand:

```
--project    GCP project ID
--region     GCP region
--service    Cloud Run service name
--runner     Runner type (github-actions, cloud-build, local)
--file, -f   Path to configuration file (default: dekopin.yml)
```

### Tag Naming Rules

Tags in Dekopin must follow these rules:

- Must consist of only lowercase alphanumeric characters and hyphens (`a-z`, `0-9`, `-`)
- Cannot contain uppercase letters, periods, underscores, spaces, or special characters
- Empty tags are handled differently depending on the runner type:
  - For GitHub Actions and Cloud Build: Automatically generates a tag based on the reference
  - For local runner: Empty tags are not allowed and will result in an error

Examples of valid tags:
- `production`
- `staging`
- `release-v1`
- `v1-0-0`
- `feature-123`

Examples of invalid tags:
- `Production` (contains uppercase)
- `staging.1` (contains period)
- `test_tag` (contains underscore)
- `tag with spaces` (contains spaces)

### Subcommands

#### deploy

Deploy a new revision with traffic directed to it.

```bash
dekopin deploy --image [IMAGE_URL]
```

Options:
- `--image, -i` (required): Container image URL (e.g., gcr.io/project/image:tag)
- `--tag, -t`: Tag name for the new revision (must follow tag naming rules)
- `--create-tag`: Create a revision tag after deployment
- `--remove-tags`: Remove all revision tags before deployment

Examples:
```bash
# Deploy with a specific image
dekopin deploy --image gcr.io/project/image:latest

# Deploy and create a tag
dekopin deploy --image gcr.io/project/image:latest --create-tag --tag release-v1
```

#### create-revision

Create a new revision without directing traffic to it.
If the runner is GitHub Actions or Cloud Build, the revision name will use the first 7 characters of the commit hash.

```bash
dekopin create-revision --image [IMAGE_URL]
```

Options:
- `--image, -i` (required): Container image URL

Example:
```bash
# Create a new revision with no traffic
dekopin create-revision --image gcr.io/project/image:v2
```

#### create-tag

Assign a tag to an existing revision. The tag must follow the naming rules specified above.

```bash
dekopin create-tag --tag [TAG_NAME] --revision [REVISION_NAME]
```

Options:
- `--tag, -t`: Tag name to create (must follow tag naming rules)
- `--revision`: Revision name to tag (default is latest)
- `--update-traffic`: Update traffic to the tagged revision after deployment
- `--remove-tags`: Remove all existing revision tags before creating the new tag

Examples:
```bash
# Tag the latest revision
dekopin create-tag --tag production

# Tag a specific revision
dekopin create-tag --tag staging --revision service-abcdef

# Tag a revision and direct traffic to it
dekopin create-tag --tag production --update-traffic

# Create a tag, making it the only tag by removing other tags
dekopin create-tag --tag production --remove-tags
```

#### remove-tag

Remove a tag from a revision.

```bash
dekopin remove-tag --tag [TAG_NAME]
```

Options:
- `--tag, -t` (required): Tag name to remove

Example:
```bash
# Remove a tag
dekopin remove-tag --tag old-release
```

#### sr-deploy (Switch Revision Deploy)

Switch traffic to a specific revision.

```bash
dekopin sr-deploy --revision [REVISION_NAME]
```

Options:
- `--revision`: Revision name to direct traffic to

Example:
```bash
# Direct all traffic to a specific revision
dekopin sr-deploy --revision service-abcdef
```

#### st-deploy (Switch Tag Deploy)

Switch traffic to a revision with a specific tag. The tag must already exist and follow the tag naming rules.

```bash
dekopin st-deploy --tag [TAG_NAME]
```

Options:
- `--tag, -t` (required): Tag name to direct traffic to
- `--remove-tags`: Remove all revision tags except the deployment target revision tag

Examples:
```bash
# Direct all traffic to the tagged revision
dekopin st-deploy --tag production

# Switch to tagged revision and clean up other tags
dekopin st-deploy --tag production --remove-tags
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

## Troubleshooting

### Common Errors

- **Timeout Errors**: By default, Dekopin has a 120-second timeout. For long-running operations, consider increasing this value in your code.
- **Tag Format Errors**: If you receive errors about invalid tag formats, ensure your tags follow the naming rules (lowercase alphanumeric and hyphens only).

## License

This project is licensed under the MIT License - see the LICENSE file for details.
