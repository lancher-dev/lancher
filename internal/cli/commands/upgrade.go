package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/lancher-dev/lancher/internal/cli/shared"
	"github.com/lancher-dev/lancher/internal/version"
)

const (
	githubAPIURL    = "https://api.github.com/repos/lancher-dev/lancher/releases/latest"
	installScriptURL = "https://raw.githubusercontent.com/lancher-dev/lancher/main/bin/install.sh"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	HTMLURL string `json:"html_url"`
}

// RunUpgradeHelp displays help for upgrade command
func RunUpgradeHelp() error {
	fmt.Printf("%slancher upgrade%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset)
	fmt.Printf("Check for updates and upgrade to the latest version\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher upgrade [options]\n\n")

	fmt.Printf("%sOPTIONS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s-f%s, %s--force%s   %sForce upgrade even if already on latest version%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "", "")
	fmt.Printf("    %s-h%s, %s--help%s    %sShow this help message%s\n\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "", "")

	return nil
}

// RunUpgrade checks for updates and upgrades lancher
func RunUpgrade(args []string) error {
	var force bool

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-f", "--force":
			force = true
		default:
			if strings.HasPrefix(args[i], "-") {
				usage := "USAGE:\n    lancher upgrade [OPTIONS]"
				return shared.FormatUnknownCommandError(args[i], usage, "lancher upgrade ")
			}
		}
	}

	// Check if we're on a supported platform
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		return shared.FormatError("upgrade is only supported on Linux and macOS")
	}

	currentVersion := version.Get()
	if currentVersion == "" || currentVersion == "dev" {
		currentVersion = "development build"
	}

	fmt.Printf("%sCurrent version:%s %s\n", shared.ColorCyan, shared.ColorReset, currentVersion)

	// Check for latest release
	spinner := shared.NewSpinner("Checking for updates...")
	spinner.Start()

	latestRelease, err := getLatestRelease()
	if err != nil {
		spinner.Fail("Failed to check for updates")
		return shared.FormatError(fmt.Sprintf("failed to check for updates: %v", err))
	}

	latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")
	currentVersionClean := strings.TrimPrefix(currentVersion, "v")

	spinner.Stop()

	fmt.Printf("%sLatest version:%s %s\n", shared.ColorCyan, shared.ColorReset, latestVersion)

	// Compare versions
	if !force && currentVersionClean == latestVersion {
		fmt.Printf("%s✓ You are already on the latest version%s\n", shared.ColorGreen, shared.ColorReset)
		return nil
	}

	if !force && currentVersionClean > latestVersion {
		fmt.Printf("%s✓ You are on a newer version than the latest release%s\n", shared.ColorGreen, shared.ColorReset)
		return nil
	}

	// Ask for confirmation
	fmt.Println()
	var confirmed bool
	if force {
		conf, err := shared.PromptConfirmWithDefault(fmt.Sprintf("Reinstall version %s?", latestVersion), false)
		if err != nil {
			return shared.FormatError("failed to read confirmation")
		}
		confirmed = conf
	} else {
		conf, err := shared.PromptConfirmWithDefault(fmt.Sprintf("Upgrade to version %s?", latestVersion), false)
		if err != nil {
			return shared.FormatError("failed to read confirmation")
		}
		confirmed = conf
	}

	if !confirmed {
		fmt.Printf("%sSkipped upgrade%s\n", shared.ColorYellow, shared.ColorReset)
		return nil
	}

	// Download and execute install script
	upgradeSpinner := shared.NewSpinner("Downloading installer...")
	upgradeSpinner.Start()

	script, err := downloadInstallScript()
	if err != nil {
		upgradeSpinner.Fail("Failed to download installer")
		return shared.FormatError(fmt.Sprintf("failed to download installer: %v", err))
	}

	upgradeSpinner.Stop()

	// Execute the install script
	fmt.Printf("%sRunning installer...%s\n\n", shared.ColorYellow, shared.ColorReset)

	cmd := exec.Command("sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return shared.FormatError(fmt.Sprintf("installation failed: %v", err))
	}

	fmt.Printf("\n%s✓ Upgrade completed successfully%s\n", shared.ColorGreen, shared.ColorReset)
	fmt.Printf("%sRun 'lancher --version' to verify the new version%s\n", shared.ColorYellow, shared.ColorReset)

	return nil
}

// getLatestRelease fetches the latest release from GitHub API
func getLatestRelease() (*GitHubRelease, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", githubAPIURL, nil)
	if err != nil {
		return nil, err
	}

	// Set User-Agent to avoid rate limiting
	req.Header.Set("User-Agent", "lancher")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, err
	}

	return &release, nil
}

// downloadInstallScript downloads the install script from GitHub
func downloadInstallScript() (string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(installScriptURL)
	if err != nil {
		return "", fmt.Errorf("failed to download script: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download script: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
