package cli

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestDetectInstallMethod(t *testing.T) {
	tests := []struct {
		name    string
		exePath string
		gopath  string
		want    string
	}{
		{
			name:    "GOPATH install",
			exePath: "/home/user/go/bin/paq",
			gopath:  "/home/user/go",
			want:    "gopath",
		},
		{
			name:    "Homebrew install",
			exePath: "/opt/homebrew/Cellar/paq/1.0.0/bin/paq",
			gopath:  "/home/user/go",
			want:    "homebrew",
		},
		{
			name:    "Homebrew linux",
			exePath: "/usr/local/Homebrew/Cellar/paq/1.0.0/bin/paq",
			gopath:  "/home/user/go",
			want:    "homebrew",
		},
		{
			name:    "Local install /usr/local/bin",
			exePath: "/usr/local/bin/paq",
			gopath:  "/home/user/go",
			want:    "local",
		},
		{
			name:    "Unknown path",
			exePath: "/custom/location/paq",
			gopath:  "/home/user/go",
			want:    "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getExecutable = func() (string, error) {
				return tt.exePath, nil
			}
			lookupEnv = func(key string) (string, bool) {
				if key == "GOPATH" {
					return tt.gopath, true
				}
				return "", false
			}
			defer func() {
				getExecutable = os.Executable
				lookupEnv = os.LookupEnv
			}()

			got := detectInstallMethod()
			if got != tt.want {
				t.Errorf("detectInstallMethod() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUpgradeCheckNewVersion(t *testing.T) {
	origVersion := version
	version = "1.0.0"
	defer func() { version = origVersion }()

	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v1.1.0", nil
	}
	defer func() { fetchLatestTag = defaultFetchLatestTag }()

	var buf bytes.Buffer
	cmd := newUpgradeCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--check"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "New version available: 1.0.0 → v1.1.0") {
		t.Errorf("expected new version message, got: %s", output)
	}
}

func TestUpgradeCheckUpToDate(t *testing.T) {
	origVersion := version
	version = "1.0.0"
	defer func() { version = origVersion }()

	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v1.0.0", nil
	}
	defer func() { fetchLatestTag = defaultFetchLatestTag }()

	var buf bytes.Buffer
	cmd := newUpgradeCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--check"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Already up to date") {
		t.Errorf("expected up to date message, got: %s", output)
	}
}

func TestUpgradeCheckNetworkError(t *testing.T) {
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "", fmt.Errorf("network error")
	}
	defer func() { fetchLatestTag = defaultFetchLatestTag }()

	var buf bytes.Buffer
	cmd := newUpgradeCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--check"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpgradeCheckDevVersion(t *testing.T) {
	origVersion := version
	version = "dev"
	defer func() { version = origVersion }()

	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v1.0.0", nil
	}
	defer func() { fetchLatestTag = defaultFetchLatestTag }()

	var buf bytes.Buffer
	cmd := newUpgradeCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"--check"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "dev build") {
		t.Errorf("expected dev build message, got: %s", output)
	}
}

func TestUpgradeGOPATH(t *testing.T) {
	var executedCmd string
	var executedArgs []string

	getExecutable = func() (string, error) {
		return "/home/user/go/bin/paq", nil
	}
	lookupEnv = func(key string) (string, bool) {
		if key == "GOPATH" {
			return "/home/user/go", true
		}
		return "", false
	}
	runCommand = func(name string, args ...string) error {
		executedCmd = name
		executedArgs = args
		return nil
	}
	defer func() {
		getExecutable = os.Executable
		lookupEnv = os.LookupEnv
		runCommand = defaultRunCommand
	}()

	var buf bytes.Buffer
	cmd := newUpgradeCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if executedCmd != "go" {
		t.Errorf("expected go command, got: %s", executedCmd)
	}
	if !strings.Contains(strings.Join(executedArgs, " "), "install") {
		t.Errorf("expected install arg, got: %v", executedArgs)
	}
}

func TestUpgradeHomebrew(t *testing.T) {
	var executedCmd string
	var executedArgs []string

	getExecutable = func() (string, error) {
		return "/opt/homebrew/bin/paq", nil
	}
	lookupEnv = func(key string) (string, bool) {
		return "", false
	}
	runCommand = func(name string, args ...string) error {
		executedCmd = name
		executedArgs = args
		return nil
	}
	defer func() {
		getExecutable = os.Executable
		lookupEnv = os.LookupEnv
		runCommand = defaultRunCommand
	}()

	var buf bytes.Buffer
	cmd := newUpgradeCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if executedCmd != "brew" {
		t.Errorf("expected brew command, got: %s", executedCmd)
	}
	if len(executedArgs) == 0 || !strings.Contains(strings.Join(executedArgs, " "), "upgrade") {
		t.Errorf("expected upgrade arg, got: %v", executedArgs)
	}
}

func TestUpgradeLocalSelfUpdate(t *testing.T) {
	getExecutable = func() (string, error) {
		return "/usr/local/bin/paq", nil
	}
	lookupEnv = func(key string) (string, bool) {
		return "", false
	}
	checkWritable = func(path string) bool {
		return true
	}
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v1.1.0", nil
	}
	applyUpdate = func(ctx context.Context, tag string) error {
		return nil
	}
	defer func() {
		getExecutable = os.Executable
		lookupEnv = os.LookupEnv
		checkWritable = canWrite
		fetchLatestTag = defaultFetchLatestTag
		applyUpdate = defaultApplyUpdate
	}()

	var buf bytes.Buffer
	cmd := newUpgradeCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "local binary") {
		t.Errorf("expected local binary message, got: %s", output)
	}
	if !strings.Contains(output, "Successfully updated to v1.1.0") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestUpgradeUnknownSelfUpdate(t *testing.T) {
	getExecutable = func() (string, error) {
		return "/some/random/path/paq", nil
	}
	lookupEnv = func(key string) (string, bool) {
		return "", false
	}
	checkWritable = func(path string) bool {
		return true
	}
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v1.1.0", nil
	}
	applyUpdate = func(ctx context.Context, tag string) error {
		return nil
	}
	defer func() {
		getExecutable = os.Executable
		lookupEnv = os.LookupEnv
		checkWritable = canWrite
		fetchLatestTag = defaultFetchLatestTag
		applyUpdate = defaultApplyUpdate
	}()

	var buf bytes.Buffer
	cmd := newUpgradeCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "direct binary update") {
		t.Errorf("expected direct binary update message, got: %s", output)
	}
	if !strings.Contains(output, "Successfully updated to v1.1.0") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestUpgradeSelfUpdateFetchError(t *testing.T) {
	getExecutable = func() (string, error) {
		return "/usr/local/bin/paq", nil
	}
	lookupEnv = func(key string) (string, bool) {
		return "", false
	}
	checkWritable = func(path string) bool {
		return true
	}
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "", fmt.Errorf("network error")
	}
	defer func() {
		getExecutable = os.Executable
		lookupEnv = os.LookupEnv
		checkWritable = canWrite
		fetchLatestTag = defaultFetchLatestTag
	}()

	var buf bytes.Buffer
	cmd := newUpgradeCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Error") {
		t.Errorf("expected error message, got: %s", output)
	}
}

func TestUpgradeSelfUpdateApplyError(t *testing.T) {
	getExecutable = func() (string, error) {
		return "/usr/local/bin/paq", nil
	}
	lookupEnv = func(key string) (string, bool) {
		return "", false
	}
	checkWritable = func(path string) bool {
		return true
	}
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v1.1.0", nil
	}
	applyUpdate = func(ctx context.Context, tag string) error {
		return fmt.Errorf("permission denied")
	}
	defer func() {
		getExecutable = os.Executable
		lookupEnv = os.LookupEnv
		checkWritable = canWrite
		fetchLatestTag = defaultFetchLatestTag
		applyUpdate = defaultApplyUpdate
	}()

	var buf bytes.Buffer
	cmd := newUpgradeCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Error") {
		t.Errorf("expected error message, got: %s", output)
	}
}

func TestUpgradeSelfUpdateAlreadyUpToDate(t *testing.T) {
	origVersion := version
	version = "1.1.0"
	defer func() { version = origVersion }()

	getExecutable = func() (string, error) {
		return "/usr/local/bin/paq", nil
	}
	lookupEnv = func(key string) (string, bool) {
		return "", false
	}
	checkWritable = func(path string) bool {
		return true
	}
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v1.1.0", nil
	}
	applyUpdate = func(ctx context.Context, tag string) error {
		t.Fatal("applyUpdate should not be called when already up to date")
		return nil
	}
	defer func() {
		getExecutable = os.Executable
		lookupEnv = os.LookupEnv
		checkWritable = canWrite
		fetchLatestTag = defaultFetchLatestTag
		applyUpdate = defaultApplyUpdate
	}()

	var buf bytes.Buffer
	cmd := newUpgradeCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Already up to date") {
		t.Errorf("expected up to date message, got: %s", output)
	}
	if strings.Contains(output, "Downloading") {
		t.Errorf("should not download when already up to date, got: %s", output)
	}
}

func TestBuildDownloadURL(t *testing.T) {
	url := buildDownloadURL("v1.0.0")

	// Should contain tag and current platform info.
	if !strings.Contains(url, "v1.0.0") {
		t.Errorf("URL should contain tag, got: %s", url)
	}
	if !strings.Contains(url, runtime.GOOS) || !strings.Contains(url, runtime.GOARCH) {
		t.Errorf("URL should contain GOOS/GOARCH, got: %s", url)
	}
	if !strings.HasPrefix(url, "https://github.com/cluion/paq/releases/download/") {
		t.Errorf("URL should start with GitHub release base, got: %s", url)
	}
}

// newUpgradeCmd creates a fresh upgrade command for testing.
// The global upgradeCmd is already registered with rootCmd,
// so we create a standalone command for isolated testing.
func newUpgradeCmd() *cobra.Command {
	var check bool
	cmd := &cobra.Command{
		Use:  "upgrade",
		RunE: func(cmd *cobra.Command, args []string) error {
			checkOnly = check
			return runUpgrade(cmd, args)
		},
	}
	cmd.Flags().BoolVar(&check, "check", false, "only check for new version")
	return cmd
}


func TestCanWrite(t *testing.T) {
	t.Parallel()

	t.Run("writable file", func(t *testing.T) {
		t.Parallel()

		f, err := os.CreateTemp(t.TempDir(), "writable")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		path := f.Name()
		_ = f.Close()

		if !canWrite(path) {
			t.Errorf("canWrite(%q) = false, want true", path)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		t.Parallel()

		if canWrite(filepath.Join(t.TempDir(), "no-such-file")) {
			t.Error("canWrite() for non-existent file = true, want false")
		}
	})

	t.Run("read-only file", func(t *testing.T) {
		t.Parallel()

		f, err := os.CreateTemp(t.TempDir(), "readonly")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		path := f.Name()
		_ = f.Close()

		_ = os.Chmod(path, 0o444)
		defer func() { _ = os.Chmod(path, 0o644) }()

		if canWrite(path) {
			t.Error("canWrite() for read-only file = true, want false")
		}
	})
}

func TestExtractBinaryFromTarGz(t *testing.T) {
	t.Parallel()

	t.Run("valid archive with paq binary", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		gw := gzip.NewWriter(&buf)
		tw := tar.NewWriter(gw)

		binaryContent := []byte("fake-paq-binary-content")
		_ = tw.WriteHeader(&tar.Header{
			Name: "paq",
			Mode: 0o755,
			Size: int64(len(binaryContent)),
		})
		_, _ = tw.Write(binaryContent)
		_ = tw.Close()
		_ = gw.Close()

		r, err := extractBinaryFromTarGz(&buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("failed to read result: %v", err)
		}
		if string(got) != string(binaryContent) {
			t.Errorf("got %q, want %q", got, binaryContent)
		}
	})

	t.Run("archive without paq binary", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		gw := gzip.NewWriter(&buf)
		tw := tar.NewWriter(gw)

		otherContent := []byte("not-paq")
		_ = tw.WriteHeader(&tar.Header{
			Name: "other-binary",
			Mode: 0o755,
			Size: int64(len(otherContent)),
		})
		_, _ = tw.Write(otherContent)
		_ = tw.Close()
		_ = gw.Close()

		_, err := extractBinaryFromTarGz(&buf)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "paq binary not found") {
			t.Errorf("error = %v, want paq binary not found", err)
		}
	})

	t.Run("paq in subdirectory", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		gw := gzip.NewWriter(&buf)
		tw := tar.NewWriter(gw)

		binaryContent := []byte("nested-paq")
		_ = tw.WriteHeader(&tar.Header{
			Name: "subdir/paq",
			Mode: 0o755,
			Size: int64(len(binaryContent)),
		})
		_, _ = tw.Write(binaryContent)
		_ = tw.Close()
		_ = gw.Close()

		r, err := extractBinaryFromTarGz(&buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("failed to read result: %v", err)
		}
		if string(got) != string(binaryContent) {
			t.Errorf("got %q, want %q", got, binaryContent)
		}
	})

	t.Run("directory entry named paq is skipped", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		gw := gzip.NewWriter(&buf)
		tw := tar.NewWriter(gw)

		_ = tw.WriteHeader(&tar.Header{
			Name:     "paq/",
			Mode:     0o755,
			Typeflag: tar.TypeDir,
		})
		// Add the real paq file after the directory entry.
		binaryContent := []byte("real-paq")
		_ = tw.WriteHeader(&tar.Header{
			Name: "paq",
			Mode: 0o755,
			Size: int64(len(binaryContent)),
		})
		_, _ = tw.Write(binaryContent)
		_ = tw.Close()
		_ = gw.Close()

		r, err := extractBinaryFromTarGz(&buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("failed to read result: %v", err)
		}
		if string(got) != string(binaryContent) {
			t.Errorf("got %q, want %q", got, binaryContent)
		}
	})

	t.Run("invalid gzip data", func(t *testing.T) {
		t.Parallel()

		_, err := extractBinaryFromTarGz(strings.NewReader("this is not gzip data"))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to decompress") {
			t.Errorf("error = %v, want failed to decompress", err)
		}
	})

	t.Run("empty archive", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		gw := gzip.NewWriter(&buf)
		_ = gw.Close()

		_, err := extractBinaryFromTarGz(&buf)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "paq binary not found") {
			t.Errorf("error = %v, want paq binary not found", err)
		}
	})
}

func TestExtractBinaryFromZip(t *testing.T) {
	t.Parallel()

	t.Run("valid zip with paq.exe", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		w := zip.NewWriter(&buf)

		binaryContent := []byte("fake-paq-exe-content")
		f, err := w.Create("paq.exe")
		if err != nil {
			t.Fatalf("failed to create zip entry: %v", err)
		}
		_, _ = f.Write(binaryContent)
		_ = w.Close()

		r, err := extractBinaryFromZip(&buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("failed to read result: %v", err)
		}
		if string(got) != string(binaryContent) {
			t.Errorf("got %q, want %q", got, binaryContent)
		}
	})

	t.Run("zip without paq.exe", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		w := zip.NewWriter(&buf)

		other, err := w.Create("other.exe")
		if err != nil {
			t.Fatalf("failed to create zip entry: %v", err)
		}
		_, _ = other.Write([]byte("other"))
		_ = w.Close()

		_, err = extractBinaryFromZip(&buf)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "paq.exe not found") {
			t.Errorf("error = %v, want paq.exe not found", err)
		}
	})

	t.Run("paq.exe in subdirectory", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		w := zip.NewWriter(&buf)

		binaryContent := []byte("nested-paq-exe")
		f, err := w.Create("subdir/paq.exe")
		if err != nil {
			t.Fatalf("failed to create zip entry: %v", err)
		}
		_, _ = f.Write(binaryContent)
		_ = w.Close()

		r, err := extractBinaryFromZip(&buf)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("failed to read result: %v", err)
		}
		if string(got) != string(binaryContent) {
			t.Errorf("got %q, want %q", got, binaryContent)
		}
	})

	t.Run("invalid zip data", func(t *testing.T) {
		t.Parallel()

		_, err := extractBinaryFromZip(strings.NewReader("this is not a zip file"))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "failed to open zip") {
			t.Errorf("error = %v, want failed to open zip", err)
		}
	})
}

func TestBytesReaderAt(t *testing.T) {
	t.Parallel()

	t.Run("normal read", func(t *testing.T) {
		t.Parallel()

		data := []byte("hello world")
		r := bytesReaderAt(data)

		p := make([]byte, 5)
		n, err := r.ReadAt(p, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != 5 {
			t.Errorf("ReadAt returned %d bytes, want 5", n)
		}
		if string(p) != "hello" {
			t.Errorf("got %q, want %q", p, "hello")
		}

		// Read from offset.
		n, err = r.ReadAt(p, 6)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != 5 {
			t.Errorf("ReadAt returned %d bytes, want 5", n)
		}
		if string(p) != "world" {
			t.Errorf("got %q, want %q", p, "world")
		}
	})

	t.Run("read at EOF", func(t *testing.T) {
		t.Parallel()

		data := []byte("hi")
		r := bytesReaderAt(data)

		p := make([]byte, 4)
		n, err := r.ReadAt(p, 3)
		if err != io.EOF {
			t.Errorf("error = %v, want io.EOF", err)
		}
		if n != 0 {
			t.Errorf("ReadAt returned %d bytes, want 0", n)
		}
	})

	t.Run("read partially beyond data", func(t *testing.T) {
		t.Parallel()

		data := []byte("hello")
		r := bytesReaderAt(data)

		p := make([]byte, 10)
		n, err := r.ReadAt(p, 2)
		// io.ReaderAt spec: returns EOF when no more data, but may return
		// the partial data first. Our implementation returns data without error
		// when off < len(data) even if the buffer is larger.
		if n != 3 {
			t.Errorf("ReadAt returned %d bytes, want 3", n)
		}
		if string(p[:n]) != "llo" {
			t.Errorf("got %q, want %q", p[:n], "llo")
		}
		_ = err
	})
}

func TestExtractBinary(t *testing.T) {
	t.Parallel()

	t.Run("dispatches based on GOOS", func(t *testing.T) {
		t.Parallel()

		if runtime.GOOS == "windows" {
			// On windows extractBinary should call extractBinaryFromZip.
			var buf bytes.Buffer
			w := zip.NewWriter(&buf)
			f, err := w.Create("paq.exe")
			if err != nil {
				t.Fatalf("failed to create zip entry: %v", err)
			}
			_, _ = f.Write([]byte("win-paq"))
			_ = w.Close()

			r, err := extractBinary(&buf, "v1.0.0")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("failed to read result: %v", err)
			}
			if string(got) != "win-paq" {
				t.Errorf("got %q, want %q", got, "win-paq")
			}
		} else {
			// On non-windows extractBinary should call extractBinaryFromTarGz.
			var buf bytes.Buffer
			gw := gzip.NewWriter(&buf)
			tw := tar.NewWriter(gw)

			content := []byte("unix-paq")
			_ = tw.WriteHeader(&tar.Header{
				Name: "paq",
				Mode: 0o755,
				Size: int64(len(content)),
			})
			_, _ = tw.Write(content)
			_ = tw.Close()
			_ = gw.Close()

			r, err := extractBinary(&buf, "v1.0.0")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("failed to read result: %v", err)
			}
			if string(got) != "unix-paq" {
				t.Errorf("got %q, want %q", got, "unix-paq")
			}
		}
	})
}

func TestSelfUpdatePermissionDenied(t *testing.T) {
	getExecutable = func() (string, error) {
		return "/usr/local/bin/paq", nil
	}
	checkWritable = func(path string) bool {
		return false
	}
	defer func() {
		getExecutable = os.Executable
		checkWritable = canWrite
	}()

	var buf bytes.Buffer
	selfUpdate(&buf, "local")

	output := buf.String()
	if !strings.Contains(output, "Permission denied") {
		t.Errorf("expected permission denied message, got: %s", output)
	}
	if !strings.Contains(output, "sudo paq upgrade") {
		t.Errorf("expected sudo hint, got: %s", output)
	}
	if strings.Contains(output, "Downloading") {
		t.Errorf("should not attempt download on permission denied, got: %s", output)
	}
}

func TestSelfUpdateUnknownMethod(t *testing.T) {
	getExecutable = func() (string, error) {
		return "/custom/path/paq", nil
	}
	checkWritable = func(path string) bool {
		return true
	}
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v2.0.0", nil
	}
	applyUpdate = func(ctx context.Context, tag string) error {
		return nil
	}
	defer func() {
		getExecutable = os.Executable
		checkWritable = canWrite
		fetchLatestTag = defaultFetchLatestTag
		applyUpdate = defaultApplyUpdate
	}()

	var buf bytes.Buffer
	selfUpdate(&buf, "unknown")

	output := buf.String()
	if !strings.Contains(output, "direct binary update") {
		t.Errorf("expected direct binary update message, got: %s", output)
	}
	if !strings.Contains(output, "Successfully updated to v2.0.0") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestSelfUpdateGetExecutableError(t *testing.T) {
	getExecutable = func() (string, error) {
		return "", fmt.Errorf("executable not found")
	}
	// checkWritable should not be called since getExecutable returns error.
	checkWritable = func(path string) bool {
		t.Error("checkWritable should not be called when getExecutable fails")
		return true
	}
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v2.0.0", nil
	}
	applyUpdate = func(ctx context.Context, tag string) error {
		return nil
	}
	defer func() {
		getExecutable = os.Executable
		checkWritable = canWrite
		fetchLatestTag = defaultFetchLatestTag
		applyUpdate = defaultApplyUpdate
	}()

	var buf bytes.Buffer
	selfUpdate(&buf, "local")

	output := buf.String()
	if !strings.Contains(output, "Successfully updated to v2.0.0") {
		t.Errorf("expected success message even when getExecutable fails, got: %s", output)
	}
}

func TestDefaultRunCommand(t *testing.T) {
	t.Parallel()

	t.Run("successful command", func(t *testing.T) {
		t.Parallel()

		err := defaultRunCommand("echo", "hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failing command", func(t *testing.T) {
		t.Parallel()

		err := defaultRunCommand("false")
		if err == nil {
			t.Fatal("expected error for failing command, got nil")
		}
	})

	t.Run("non-existent command", func(t *testing.T) {
		t.Parallel()

		err := defaultRunCommand("nonexistent-command-that-does-not-exist-12345")
		if err == nil {
			t.Fatal("expected error for non-existent command, got nil")
		}
	})
}

func TestSelfUpdateFetchTagError(t *testing.T) {
	getExecutable = func() (string, error) {
		return "/usr/local/bin/paq", nil
	}
	checkWritable = func(path string) bool {
		return true
	}
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "", fmt.Errorf("network timeout")
	}
	defer func() {
		getExecutable = os.Executable
		checkWritable = canWrite
		fetchLatestTag = defaultFetchLatestTag
	}()

	var buf bytes.Buffer
	selfUpdate(&buf, "local")

	output := buf.String()
	if !strings.Contains(output, "Error") {
		t.Errorf("expected error message, got: %s", output)
	}
	if !strings.Contains(output, "network timeout") {
		t.Errorf("expected network timeout detail, got: %s", output)
	}
}

func TestSelfUpdateApplyError(t *testing.T) {
	getExecutable = func() (string, error) {
		return "/usr/local/bin/paq", nil
	}
	checkWritable = func(path string) bool {
		return true
	}
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v3.0.0", nil
	}
	applyUpdate = func(ctx context.Context, tag string) error {
		return fmt.Errorf("disk full")
	}
	defer func() {
		getExecutable = os.Executable
		checkWritable = canWrite
		fetchLatestTag = defaultFetchLatestTag
		applyUpdate = defaultApplyUpdate
	}()

	var buf bytes.Buffer
	selfUpdate(&buf, "local")

	output := buf.String()
	if !strings.Contains(output, "Error") {
		t.Errorf("expected error message, got: %s", output)
	}
	if !strings.Contains(output, "disk full") {
		t.Errorf("expected disk full detail, got: %s", output)
	}
	if strings.Contains(output, "Successfully") {
		t.Errorf("should not show success on apply error, got: %s", output)
	}
}

func TestSelfUpdateLocalMethod(t *testing.T) {
	getExecutable = func() (string, error) {
		return "/usr/local/bin/paq", nil
	}
	checkWritable = func(path string) bool {
		return true
	}
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v2.5.0", nil
	}
	applyUpdate = func(ctx context.Context, tag string) error {
		return nil
	}
	defer func() {
		getExecutable = os.Executable
		checkWritable = canWrite
		fetchLatestTag = defaultFetchLatestTag
		applyUpdate = defaultApplyUpdate
	}()

	var buf bytes.Buffer
	selfUpdate(&buf, "local")

	output := buf.String()
	if !strings.Contains(output, "Detected install method: local binary") {
		t.Errorf("expected local binary message, got: %s", output)
	}
	if !strings.Contains(output, "Successfully updated to v2.5.0") {
		t.Errorf("expected success message, got: %s", output)
	}
}

func TestDefaultFetchLatestTag(t *testing.T) {
	t.Run("successful response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Accept") != "application/vnd.github+json" {
				t.Errorf("missing Accept header, got: %s", r.Header.Get("Accept"))
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprintf(w, `{"tag_name":"v2.3.4"}`)
		}))
		defer srv.Close()

		// Override http.DefaultClient to redirect to test server.
		origClient := http.DefaultClient
		http.DefaultClient = &http.Client{
			Transport: &redirectTransport{baseURL: srv.URL},
		}
		defer func() { http.DefaultClient = origClient }()

		tag, err := defaultFetchLatestTag(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tag != "v2.3.4" {
			t.Errorf("tag = %q, want %q", tag, "v2.3.4")
		}
	})

	t.Run("non-200 status", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer srv.Close()

		origClient := http.DefaultClient
		http.DefaultClient = &http.Client{
			Transport: &redirectTransport{baseURL: srv.URL},
		}
		defer func() { http.DefaultClient = origClient }()

		_, err := defaultFetchLatestTag(context.Background())
		if err == nil {
			t.Fatal("expected error for 404, got nil")
		}
		if !strings.Contains(err.Error(), "status 404") {
			t.Errorf("error = %v, want status 404", err)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not json`))
		}))
		defer srv.Close()

		origClient := http.DefaultClient
		http.DefaultClient = &http.Client{
			Transport: &redirectTransport{baseURL: srv.URL},
		}
		defer func() { http.DefaultClient = origClient }()

		_, err := defaultFetchLatestTag(context.Background())
		if err == nil {
			t.Fatal("expected error for invalid JSON, got nil")
		}
		if !strings.Contains(err.Error(), "failed to parse") {
			t.Errorf("error = %v, want failed to parse", err)
		}
	})

	t.Run("cancelled context", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(5 * time.Second)
		}))
		defer srv.Close()

		origClient := http.DefaultClient
		http.DefaultClient = &http.Client{
			Transport: &redirectTransport{baseURL: srv.URL},
		}
		defer func() { http.DefaultClient = origClient }()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := defaultFetchLatestTag(ctx)
		if err == nil {
			t.Fatal("expected error for cancelled context, got nil")
		}
	})
}

// redirectTransport is an http.RoundTripper that redirects all requests
// to a test server base URL, preserving the path.
type redirectTransport struct {
	baseURL string
}

func (t *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Redirect to test server while preserving the path.
	newURL := t.baseURL + req.URL.Path
	newReq, err := http.NewRequestWithContext(req.Context(), req.Method, newURL, req.Body)
	if err != nil {
		return nil, err
	}
	newReq.Header = req.Header
	return http.DefaultTransport.RoundTrip(newReq)
}

func TestDefaultApplyUpdate(t *testing.T) {
	t.Run("non-200 download status", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer srv.Close()

		origClient := http.DefaultClient
		http.DefaultClient = &http.Client{
			Transport: &redirectTransport{baseURL: srv.URL},
		}
		defer func() { http.DefaultClient = origClient }()

		err := defaultApplyUpdate(context.Background(), "v1.0.0")
		if err == nil {
			t.Fatal("expected error for non-200 download, got nil")
		}
		if !strings.Contains(err.Error(), "download failed") {
			t.Errorf("error = %v, want download failed", err)
		}
	})

	t.Run("invalid archive", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("not a valid archive"))
		}))
		defer srv.Close()

		origClient := http.DefaultClient
		http.DefaultClient = &http.Client{
			Transport: &redirectTransport{baseURL: srv.URL},
		}
		defer func() { http.DefaultClient = origClient }()

		err := defaultApplyUpdate(context.Background(), "v1.0.0")
		if err == nil {
			t.Fatal("expected error for invalid archive, got nil")
		}
		if !strings.Contains(err.Error(), "failed to extract binary") {
			t.Errorf("error = %v, want failed to extract binary", err)
		}
	})
}
