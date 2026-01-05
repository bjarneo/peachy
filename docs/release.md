# Releasing Peachy

Releases are automated via GitHub Actions when you push a version tag.

## Creating a Release

1. Update the version in `main.go` if needed
2. Commit your changes
3. Tag the release:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The workflow automatically:
- Builds binaries for Linux and macOS (amd64 and arm64)
- Generates SHA-256 checksums
- Creates a GitHub release with all artifacts

## Binaries Produced

| Platform | Architecture | Filename |
|----------|--------------|----------|
| Linux    | amd64        | `peachy-linux-amd64` |
| Linux    | arm64        | `peachy-linux-arm64` |
| macOS    | amd64        | `peachy-darwin-amd64` |
| macOS    | arm64        | `peachy-darwin-arm64` |

## Version Format

Use semantic versioning with a `v` prefix: `v0.1.0`, `v1.0.0`, `v2.1.3`
