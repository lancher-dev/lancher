#!/bin/sh
# Builds the signed APT repository into OUT_DIR from the .deb assets of the
# latest non-prerelease GitHub releases.
#
# Usage: build-apt.sh OUT_DIR
#
# Environment:
#   APT_GPG_PRIVATE_KEY  armored private key used to sign the repository (required)
#   GH_TOKEN             token for the GitHub CLI
#   APT_REPO             GitHub repository with the releases (default: lancher-dev/lancher)
#   APT_KEEP             number of releases to include (default: 3)
#
# Requires: gh, dpkg-scanpackages, apt-ftparchive, gpg, gzip.
set -eu

OUT_DIR=${1:?usage: build-apt.sh OUT_DIR}
REPO=${APT_REPO:-lancher-dev/lancher}
KEEP=${APT_KEEP:-3}
ARCHS="amd64 arm64"

if [ -z "${APT_GPG_PRIVATE_KEY:-}" ]; then
    echo "APT_GPG_PRIVATE_KEY is not set" >&2
    exit 1
fi

mkdir -p "$OUT_DIR/pool/main"
cd "$OUT_DIR"

# Releases are listed newest first; those without .deb assets are skipped.
found=0
for tag in $(gh release list --repo "$REPO" --exclude-drafts --exclude-pre-releases \
    --limit 50 --json tagName --jq '.[].tagName'); do
    [ "$found" -ge "$KEEP" ] && break
    if gh release download "$tag" --repo "$REPO" --pattern '*.deb' \
        --dir pool/main --skip-existing 2>/dev/null; then
        echo "Included $tag"
        found=$((found + 1))
    fi
done

if [ "$found" -eq 0 ]; then
    echo "No release with .deb assets found in $REPO" >&2
    exit 1
fi

for arch in $ARCHS; do
    dir="dists/stable/main/binary-$arch"
    mkdir -p "$dir"
    dpkg-scanpackages --arch "$arch" pool/main /dev/null > "$dir/Packages"
    gzip -k -f "$dir/Packages"
done

conf=$(mktemp)
release=$(mktemp)
cat > "$conf" <<EOF
APT::FTPArchive::Release::Origin "lancher-dev";
APT::FTPArchive::Release::Label "lancher";
APT::FTPArchive::Release::Suite "stable";
APT::FTPArchive::Release::Codename "stable";
APT::FTPArchive::Release::Architectures "$ARCHS";
APT::FTPArchive::Release::Components "main";
APT::FTPArchive::Release::Description "lancher official APT repository";
EOF
(cd dists/stable && apt-ftparchive -c "$conf" release . > "$release")
mv "$release" dists/stable/Release
rm -f "$conf"

GNUPGHOME=$(mktemp -d)
export GNUPGHOME
trap 'rm -rf "$GNUPGHOME"' EXIT
printf '%s\n' "$APT_GPG_PRIVATE_KEY" | gpg --batch --quiet --import

gpg --batch --yes --armor --clearsign --output dists/stable/InRelease dists/stable/Release
gpg --batch --yes --armor --detach-sign --output dists/stable/Release.gpg dists/stable/Release

echo "APT repository written to $OUT_DIR"
