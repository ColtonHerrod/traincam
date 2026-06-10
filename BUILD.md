# Building the Train Cam Tracker

This project provides multiple ways to build the application:

## Local Development

```bash
# Install dependencies
cd frontend && npm install && npm run build
cd ..

# Run with Wails
wails dev
```

## Building Artifacts

### Using GitHub Actions (Recommended)

The workflow automatically builds for all platforms:

1. **Push a tag** to trigger a build:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **Manually trigger** the workflow:
   - Go to: Actions → Build and Release → Run workflow

3. **Download artifacts** from the workflow run:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)

### Using goreleaser

```bash
# Build locally
goreleaser build --clean

# Create release
goreleaser release --clean
```

### Using Wails Directly

```bash
# Build for all platforms
wails build

# Build for specific platform
wails build -platform linux/amd64
wails build -platform darwin/amd64
wails build -platform windows/amd64
wails build -platform linux/arm64
wails build -platform darwin/arm64
```

## Build Output

Artifacts are located in `build/bin/`:
- `traincam` - Linux binary
- `traincam-arm64` - Linux ARM64 binary
- `traincam` - macOS binary
- `traincam.exe` - Windows binary

All include the built frontend assets in the same directory.