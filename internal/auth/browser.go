package auth

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var ErrNoBrowser = errors.New("no supported browser found (install Chrome, Edge, Brave or Chromium)")

type macBrowser struct {
	bundle     string
	executable string
}

var macBrowsers = []macBrowser{
	{bundle: "Google Chrome.app", executable: "Google Chrome"},
	{bundle: "Microsoft Edge.app", executable: "Microsoft Edge"},
	{bundle: "Brave Browser.app", executable: "Brave Browser"},
	{bundle: "Chromium.app", executable: "Chromium"},
}

func FindBrowser() (string, error) {
	for _, candidate := range browserCandidates() {
		if path, ok := resolveExecutable(candidate); ok {
			return path, nil
		}
	}
	return "", ErrNoBrowser
}

func resolveExecutable(candidate string) (string, bool) {
	if filepath.IsAbs(candidate) {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, true
		}
		return "", false
	}

	path, err := exec.LookPath(candidate)
	return path, err == nil
}

func browserCandidates() []string {
	switch runtime.GOOS {
	case "darwin":
		return macBrowserCandidates()
	case "windows":
		return windowsBrowserCandidates()
	default:
		return []string{
			"google-chrome",
			"google-chrome-stable",
			"microsoft-edge",
			"brave-browser",
			"chromium",
			"chromium-browser",
		}
	}
}

func macBrowserCandidates() []string {
	applicationDirs := []string{"/Applications"}
	if home, err := os.UserHomeDir(); err == nil {
		applicationDirs = append(applicationDirs, filepath.Join(home, "Applications"))
	}

	var candidates []string
	for _, dir := range applicationDirs {
		for _, browser := range macBrowsers {
			candidates = append(candidates, filepath.Join(dir, browser.bundle, "Contents", "MacOS", browser.executable))
		}
	}
	return candidates
}

func windowsBrowserCandidates() []string {
	var candidates []string
	for _, envVar := range []string{"ProgramFiles", "ProgramFiles(x86)", "LocalAppData"} {
		dir := os.Getenv(envVar)
		if dir == "" {
			continue
		}
		candidates = append(candidates,
			filepath.Join(dir, "Google", "Chrome", "Application", "chrome.exe"),
			filepath.Join(dir, "Microsoft", "Edge", "Application", "msedge.exe"),
			filepath.Join(dir, "BraveSoftware", "Brave-Browser", "Application", "brave.exe"),
			filepath.Join(dir, "Chromium", "Application", "chrome.exe"),
		)
	}
	return candidates
}
