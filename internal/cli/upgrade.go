package cli

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/minio/selfupdate"
	"github.com/spf13/cobra"
)

// Function variables for testability — matches provider.CommandRunner pattern.
var (
	getExecutable  = os.Executable
	lookupEnv      = os.LookupEnv
	runCommand     = defaultRunCommand
	fetchLatestTag = defaultFetchLatestTag
	applyUpdate    = defaultApplyUpdate
	checkWritable  = canWrite
)

// defaultRunCommand executes a system command, printing output to stdout/stderr.
func defaultRunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// defaultFetchLatestTag queries GitHub Releases API for the latest tag.
func defaultFetchLatestTag(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/cluion/paq/releases/latest", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch latest release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return release.TagName, nil
}

// defaultApplyUpdate downloads the latest release binary and applies the self-update.
func defaultApplyUpdate(ctx context.Context, tag string) error {
	url := buildDownloadURL(tag)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download release: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	binary, err := extractBinary(resp.Body, tag)
	if err != nil {
		return fmt.Errorf("failed to extract binary: %w", err)
	}

	if err := selfupdate.Apply(binary, selfupdate.Options{}); err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			return fmt.Errorf("update failed and rollback also failed: %w", rerr)
		}
		return fmt.Errorf("update failed, rolled back: %w", err)
	}

	return nil
}

// buildDownloadURL constructs the GitHub Release download URL for the current platform.
func buildDownloadURL(tag string) string {
	version := strings.TrimPrefix(tag, "v")
	ext := "tar.gz"
	if runtime.GOOS == "windows" {
		ext = "zip"
	}
	return fmt.Sprintf(
		"https://github.com/cluion/paq/releases/download/%s/paq_%s_%s_%s.%s",
		tag, version, runtime.GOOS, runtime.GOARCH, ext,
	)
}

// extractBinary extracts the paq binary from a tar.gz or zip archive.
// Returns a reader that must be closed by the caller.
func extractBinary(r io.Reader, tag string) (io.ReadSeeker, error) {
	if runtime.GOOS == "windows" {
		return extractBinaryFromZip(r)
	}
	return extractBinaryFromTarGz(r)
}

func extractBinaryFromTarGz(r io.Reader) (io.ReadSeeker, error) {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress: %w", err)
	}
	defer func() { _ = gzr.Close() }()

	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("paq binary not found in archive")
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read archive: %w", err)
		}

		if filepath.Base(hdr.Name) == "paq" && !hdr.FileInfo().IsDir() {
			data, err := io.ReadAll(tr)
			if err != nil {
				return nil, fmt.Errorf("failed to read binary: %w", err)
			}
			return strings.NewReader(string(data)), nil
		}
	}
}

func extractBinaryFromZip(r io.Reader) (io.ReadSeeker, error) {
	// zip requires ReadAt + Seeker, so buffer the entire response.
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read zip data: %w", err)
	}

	readerAt := bytesReaderAt(data)
	zr, err := zip.NewReader(readerAt, int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to open zip: %w", err)
	}

	for _, f := range zr.File {
		if filepath.Base(f.Name) == "paq.exe" && !f.FileInfo().IsDir() {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open file in zip: %w", err)
			}
			defer func() { _ = rc.Close() }()

			content, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("failed to read binary: %w", err)
			}
			return strings.NewReader(string(content)), nil
		}
	}

	return nil, fmt.Errorf("paq.exe not found in archive")
}

// bytesReaderAt wraps a []byte to implement io.ReaderAt.
type bytesReaderAtImpl struct {
	data []byte
}

func (r *bytesReaderAtImpl) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(r.data)) {
		return 0, io.EOF
	}
	n := copy(p, r.data[off:])
	return n, nil
}

func bytesReaderAt(data []byte) *bytesReaderAtImpl {
	return &bytesReaderAtImpl{data: data}
}

// detectInstallMethod determines how paq was installed based on binary path.
func detectInstallMethod() string {
	exePath, err := getExecutable()
	if err != nil {
		return "unknown"
	}

	// Resolve symlinks (important for Homebrew which uses cellar symlinks).
	resolved, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		resolved = exePath
	}
	resolved = filepath.Clean(resolved)

	// Check GOPATH/bin.
	gopath, ok := lookupEnv("GOPATH")
	if !ok {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}
	gopathBin := filepath.Clean(filepath.Join(gopath, "bin"))
	if strings.HasPrefix(resolved, gopathBin+string(filepath.Separator)) || resolved == filepath.Join(gopathBin, "paq") {
		return "gopath"
	}

	// Check Homebrew.
	lower := strings.ToLower(resolved)
	if strings.Contains(lower, "homebrew") || strings.Contains(lower, "cellar") {
		return "homebrew"
	}

	// Check /usr/local/bin (Unix) or C:\usr\local\bin (Windows via MSYS2/Cygwin).
	localBin := filepath.Clean("/usr/local/bin")
	if strings.HasPrefix(resolved, localBin+string(filepath.Separator)) || resolved == filepath.Join(localBin, "paq") {
		return "local"
	}

	return "unknown"
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade paq to the latest version",
	RunE:  runUpgrade,
}

var checkOnly bool

func init() {
	upgradeCmd.Flags().BoolVar(&checkOnly, "check", false, "only check for new version, do not upgrade")
	rootCmd.AddCommand(upgradeCmd)
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	v, _, _ := buildInfo()
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Current version: %s\n", v)

	if checkOnly {
		return runCheck(cmd, v)
	}
	return runUpgradeExec(cmd)
}

func runCheck(cmd *cobra.Command, currentVersion string) error {
	ctx := context.Background()
	tag, err := fetchLatestTag(ctx)
	if err != nil {
		return fmt.Errorf("failed to check latest version: %w", err)
	}

	latest := strings.TrimPrefix(tag, "v")
	current := strings.TrimPrefix(currentVersion, "v")

	w := cmd.OutOrStdout()

	if current == "dev" {
		_, _ = fmt.Fprintf(w, "Latest release: %s (current: dev build)\n", tag)
		return nil
	}

	if current == latest {
		_, _ = fmt.Fprintf(w, "Already up to date (%s)\n", currentVersion)
	} else {
		_, _ = fmt.Fprintf(w, "New version available: %s → %s\n", currentVersion, tag)
	}
	return nil
}

func runUpgradeExec(cmd *cobra.Command) error {
	method := detectInstallMethod()
	w := cmd.OutOrStdout()

	// Check if already up to date (skip for dev builds).
	currentV, _, _ := buildInfo()
	current := strings.TrimPrefix(currentV, "v")
	if current != "dev" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tag, err := fetchLatestTag(ctx)
		if err == nil {
			latest := strings.TrimPrefix(tag, "v")
			if current == latest {
				_, _ = fmt.Fprintf(w, "Already up to date (%s)\n", currentV)
				return nil
			}
			_, _ = fmt.Fprintf(w, "New version available: %s → %s\n", currentV, tag)
		}
	}

	switch method {
	case "gopath":
		_, _ = fmt.Fprintln(w, "Detected install method: go install")
		_, _ = fmt.Fprintln(w, "Running: go install github.com/cluion/paq/cmd/paq@latest")
		if err := runCommand("go", "install", "github.com/cluion/paq/cmd/paq@latest"); err != nil {
			return fmt.Errorf("upgrade failed: %w", err)
		}
		_, _ = fmt.Fprintln(w, "Upgrade complete!")
	case "homebrew":
		_, _ = fmt.Fprintln(w, "Detected install method: Homebrew")
		_, _ = fmt.Fprintln(w, "Running: brew upgrade cluion/tap/paq")
		if err := runCommand("brew", "upgrade", "cluion/tap/paq"); err != nil {
			return fmt.Errorf("upgrade failed: %w", err)
		}
		_, _ = fmt.Fprintln(w, "Upgrade complete!")
	default:
		selfUpdate(w, method)
		return nil
	}
	return nil
}

func selfUpdate(w io.Writer, method string) {
	if method == "local" {
		_, _ = fmt.Fprintln(w, "Detected install method: local binary")
	} else {
		_, _ = fmt.Fprintln(w, "Install method unknown, using direct binary update")
	}

	// Check write permission before attempting update.
	exePath, err := getExecutable()
	if err == nil {
		if !checkWritable(exePath) {
			_, _ = fmt.Fprintln(w, "Permission denied: cannot write to current binary.")
			_, _ = fmt.Fprintf(w, "Run with elevated privileges:\n")
			_, _ = fmt.Fprintf(w, "  sudo paq upgrade\n")
			_, _ = fmt.Fprintln(w)
			_, _ = fmt.Fprintln(w, "Or move the binary to a user-writable directory:")
			_, _ = fmt.Fprintln(w, "  mkdir -p ~/.local/bin")
			_, _ = fmt.Fprintf(w, "  sudo mv %s ~/.local/bin/paq\n", exePath)
			_, _ = fmt.Fprintln(w, "  echo 'export PATH=\"$HOME/.local/bin:$PATH\"' >> ~/.zshrc")
			return
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tag, err := fetchLatestTag(ctx)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error: failed to fetch latest version: %v\n", err)
		return
	}

	_, _ = fmt.Fprintln(w, "Downloading latest release...")

	if err := applyUpdate(ctx, tag); err != nil {
		_, _ = fmt.Fprintf(w, "Error: %v\n", err)
		return
	}

	_, _ = fmt.Fprintf(w, "Successfully updated to %s!\n", tag)
}

// canWrite checks if the current user can write to the given file.
func canWrite(path string) bool {
	// Try opening the file for writing (without actually modifying it).
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	_ = f.Close()
	return true
}
