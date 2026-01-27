# CI/CD and Build Automation

This document describes the build automation and CI/CD setup for GLess.

## Build Tools

### 1. Makefile (Linux/macOS)

The Makefile provides cross-platform build targets for Unix-based systems:

**Available targets:**
- `make build` - Build for current platform
- `make build-all` - Build for all platforms (Windows, Linux, macOS)
- `make build-linux` - Linux builds (amd64, arm64, arm)
- `make build-macos` - macOS builds (amd64, arm64)
- `make build-windows` - Windows builds (amd64, arm64)
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make help` - Show help

### 2. PowerShell Build Script (Windows)

The `build.ps1` script provides the same functionality on Windows:

**Usage:**
```powershell
.\build.ps1 -All                    # Build all platforms
.\build.ps1 -Windows                # Windows only
.\build.ps1 -Linux -MacOS           # Specific platforms
.\build.ps1 -All -Version 1.0.0    # Set version
```

## GitHub Workflows

### CI Workflow (`ci.yml`)

Runs on every push and pull request to `main` and `develop` branches.

**Jobs:**

1. **Test**
   - Runs on: Ubuntu, macOS, Windows
   - Go versions: 1.21, 1.22
   - Executes: Tests with race detection and coverage
   - Uploads coverage to Codecov (Ubuntu + Go 1.22 only)

2. **Lint**
   - Runs golangci-lint for code quality
   - Configuration: `.golangci.yml`

3. **Build**
   - Builds for all platform combinations:
     - Linux: amd64, arm64
     - macOS: amd64, arm64
     - Windows: amd64, arm64
   - Uploads build artifacts (7-day retention)

4. **Format Check**
   - Ensures code is properly formatted with `gofmt`

5. **Mod Tidy Check**
   - Verifies `go.mod` and `go.sum` are up to date

### Release Workflow (`release.yml`)

Triggers on git tag push (format: `v*`)

**Process:**
1. Checks out code with full history
2. Sets up Go 1.21
3. Builds all platforms using `make build-all`
4. Generates SHA256 checksums
5. Creates GitHub release with:
   - All binaries
   - Checksums file
   - Auto-generated release notes

**Supported platforms:**
- Windows: amd64, arm64
- Linux: amd64, arm64, arm
- macOS: amd64, arm64 (Apple Silicon)

## Creating a Release

To create a new release:

1. **Tag the release:**
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

2. **Workflow automatically:**
   - Builds all platform binaries
   - Generates checksums
   - Creates GitHub release
   - Uploads all artifacts

3. **Users can download:**
   - Pre-built binaries from releases page
   - Verify checksums for security

## Build Artifacts

All builds output to `dist/` directory with naming convention:
- `gless-{os}-{arch}[.exe]`

Examples:
- `gless-linux-amd64`
- `gless-darwin-arm64`
- `gless-windows-amd64.exe`

## Linting Configuration

Located in `.golangci.yml`:
- Enabled linters: gofmt, govet, errcheck, staticcheck, unused, gosimple, ineffassign, typecheck
- Timeout: 5 minutes
- Simplified formatting enabled

## Platforms Supported

| Platform | Architecture | Build Tool | CI | Release |
|----------|-------------|------------|:--:|:-------:|
| Linux | amd64 | ✅ | ✅ | ✅ |
| Linux | arm64 | ✅ | ✅ | ✅ |
| Linux | arm | ✅ | ❌ | ✅ |
| macOS | amd64 (Intel) | ✅ | ✅ | ✅ |
| macOS | arm64 (M1/M2) | ✅ | ✅ | ✅ |
| Windows | amd64 | ✅ | ✅ | ✅ |
| Windows | arm64 | ✅ | ✅ | ✅ |
