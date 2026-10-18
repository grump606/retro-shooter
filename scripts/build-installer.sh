#!/usr/bin/env bash
# Build the Windows installer from WSL (or Linux with Inno Setup under Wine).
#
#   scripts/build-installer.sh
#
# Output: dist/RetroShooter-Setup-<version>.exe
# Needs Inno Setup 6 on Windows (winget install JRSoftware.InnoSetup).
set -euo pipefail

cd "$(dirname "$0")/.."
VERSION=$(tr -d '[:space:]' < VERSION)
export CGO_ENABLED=0

echo "==> Icon and version resources ($VERSION)"
go run ./tools/mkicon -out build/windows
go run github.com/tc-hib/go-winres@v0.3.3 make \
  --in build/windows/winres.json --out rsrc --arch amd64 \
  --file-version "$VERSION.0" --product-version "$VERSION.0"

echo "==> Game executable"
mkdir -p dist
# -H windowsgui: no console window when launched from the desktop icon.
GOOS=windows GOARCH=amd64 go build -trimpath \
  -ldflags="-s -w -H windowsgui" -o dist/RetroShooter.exe .

echo "==> Installer"
ISCC=${ISCC:-}
if [[ -z "$ISCC" ]]; then
  for c in \
    "$(command -v iscc.exe 2>/dev/null || true)" \
    /mnt/c/Users/*/AppData/Local/Programs/"Inno Setup 6"/ISCC.exe \
    "/mnt/c/Program Files (x86)/Inno Setup 6/ISCC.exe" \
    "/mnt/c/Program Files/Inno Setup 6/ISCC.exe"; do
    if [[ -n "$c" && -f "$c" ]]; then ISCC=$c; break; fi
  done
fi
if [[ -z "$ISCC" ]]; then
  echo "Inno Setup not found. Install it on Windows with:" >&2
  echo "  winget install JRSoftware.InnoSetup --scope user" >&2
  echo "or set ISCC=/path/to/ISCC.exe" >&2
  exit 1
fi

# ISCC is a Windows program, so give it a Windows path to the script.
script=installer/retro-shooter.iss
if command -v wslpath >/dev/null; then script=$(wslpath -w "$script"); fi
"$ISCC" /Q "/DAppVersion=$VERSION" "$script"
# ISCC doesn't set the Unix execute bit; without it WSL can't launch the file.
chmod +x "dist/RetroShooter-Setup-$VERSION.exe"

echo "==> Done: dist/RetroShooter-Setup-$VERSION.exe"
