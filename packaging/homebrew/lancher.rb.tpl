# Filled by envsubst in .github/workflows/publish-homebrew.yml and written to
# Formula/lancher.rb in lancher-dev/homebrew-lancher.
class Lancher < Formula
  desc "A minimal local project template manager"
  homepage "https://lancher.dev"
  version "${VERSION}"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/lancher-dev/lancher/releases/download/v${VERSION}/lancher-darwin-amd64"
      sha256 "${SHA_DARWIN_AMD64}"
    end

    on_arm do
      url "https://github.com/lancher-dev/lancher/releases/download/v${VERSION}/lancher-darwin-arm64"
      sha256 "${SHA_DARWIN_ARM64}"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/lancher-dev/lancher/releases/download/v${VERSION}/lancher-linux-amd64"
      sha256 "${SHA_LINUX_AMD64}"
    end

    on_arm do
      url "https://github.com/lancher-dev/lancher/releases/download/v${VERSION}/lancher-linux-arm64"
      sha256 "${SHA_LINUX_ARM64}"
    end
  end

  def install
    bin.install "lancher-darwin-amd64" => "lancher" if OS.mac? && Hardware::CPU.intel?
    bin.install "lancher-darwin-arm64" => "lancher" if OS.mac? && Hardware::CPU.arm?
    bin.install "lancher-linux-amd64"  => "lancher" if OS.linux? && Hardware::CPU.intel?
    bin.install "lancher-linux-arm64"  => "lancher" if OS.linux? && Hardware::CPU.arm?
    (bin/"lancher").chmod 0755
    generate_completions_from_executable(bin/"lancher", "completion")
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/lancher --version")
  end
end
