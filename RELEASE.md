# Release Instructions

This document describes how to create and manage releases for tscgit.

## Automated Releases

Releases are automated using GitHub Actions and GoReleaser. Simply create and push a tag to trigger a release.

### Creating a Release

1. **Update the changelog**:
   ```bash
   # Edit CHANGELOG.md to move items from [Unreleased] to a new version section
   ```

2. **Create and push a tag**:
   ```bash
   # Create a new version tag (follow semantic versioning)
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. **Monitor the release**:
   - Go to the [Actions tab](https://github.com/rohit746/tscgit/actions) to monitor the release workflow
   - Check the [Releases page](https://github.com/rohit746/tscgit/releases) for the new release

### What Gets Built

The release pipeline automatically creates:

- **Cross-platform binaries**: Windows (amd64, arm64), macOS (amd64, arm64), Linux (amd64, arm64)
- **Archives**: `.tar.gz` for Unix systems, `.zip` for Windows
- **Checksums**: SHA256 checksums for all artifacts
- **Package manager support**: Homebrew formula, Debian packages, RPM packages
- **Release notes**: Auto-generated from git commits and changelog

### Manual Release (Emergency)

If the automated pipeline fails, you can trigger a manual release:

1. **Using GitHub UI**:
   - Go to Actions tab → Release workflow → "Run workflow"
   - Enter the tag name (e.g., `v1.0.0`)

2. **Local GoReleaser** (requires setup):
   ```bash
   # Install GoReleaser
   go install github.com/goreleaser/goreleaser@latest
   
   # Create a test release (dry run)
   goreleaser release --snapshot --clean
   
   # Create actual release (requires GITHUB_TOKEN)
   export GITHUB_TOKEN="your_token_here"
   goreleaser release --clean
   ```

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **Major version** (X.0.0): Breaking changes, incompatible API changes
- **Minor version** (X.Y.0): New features, backwards compatible
- **Patch version** (X.Y.Z): Bug fixes, backwards compatible

### Examples
- `v1.0.0`: First stable release
- `v1.1.0`: New lesson types added
- `v1.1.1`: Bug fix in verification logic
- `v2.0.0`: Complete rewrite of lesson format

## Testing Releases

Before creating a release:

1. **Run tests**:
   ```bash
   go test ./...
   ```

2. **Test cross-compilation**:
   ```bash
   GOOS=windows GOARCH=amd64 go build ./cmd/tscgit
   GOOS=darwin GOARCH=arm64 go build ./cmd/tscgit
   GOOS=linux GOARCH=amd64 go build ./cmd/tscgit
   ```

3. **Test installation scripts** (on different platforms):
   ```bash
   # Windows PowerShell
   .\install.ps1
   
   # Linux/macOS
   ./install.sh
   ```

4. **Test the CLI**:
   ```bash
   tscgit version
   tscgit lessons
   tscgit verify init-basics
   tscgit run 0
   ```

## Post-Release

After a successful release:

1. **Verify installation methods work**:
   - Test `go install github.com/rohit746/tscgit/cmd/tscgit@latest`
   - Test installation scripts on each platform
   - Check package manager installations if applicable

2. **Update documentation** if needed:
   - README.md installation instructions
   - Any version-specific documentation

3. **Announce the release**:
   - Social media, forums, or educational platforms
   - Include key new features and installation instructions

## Troubleshooting

### Release Failed

- Check the [Actions logs](https://github.com/rohit746/tscgit/actions) for error details
- Common issues:
  - Missing or invalid GitHub token
  - Build failures on specific platforms
  - GoReleaser configuration errors

### Installation Scripts Not Working

- Test scripts on each platform
- Check GitHub API rate limits
- Verify download URLs in scripts match actual release assets

### Binary Not Working

- Check that all platforms build successfully
- Verify version information is embedded correctly
- Test on clean systems without Go installed