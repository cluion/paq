package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

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
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v1.1.0", nil
	}
	applyUpdate = func(ctx context.Context, tag string) error {
		return nil
	}
	defer func() {
		getExecutable = os.Executable
		lookupEnv = os.LookupEnv
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
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v1.1.0", nil
	}
	applyUpdate = func(ctx context.Context, tag string) error {
		return nil
	}
	defer func() {
		getExecutable = os.Executable
		lookupEnv = os.LookupEnv
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
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "", fmt.Errorf("network error")
	}
	defer func() {
		getExecutable = os.Executable
		lookupEnv = os.LookupEnv
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
	fetchLatestTag = func(ctx context.Context) (string, error) {
		return "v1.1.0", nil
	}
	applyUpdate = func(ctx context.Context, tag string) error {
		return fmt.Errorf("permission denied")
	}
	defer func() {
		getExecutable = os.Executable
		lookupEnv = os.LookupEnv
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
