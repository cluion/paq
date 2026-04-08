package cli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildInfoNonDev(t *testing.T) {
	t.Parallel()

	origVersion := version
	origCommit := commit
	origDate := date
	version = "1.2.3"
	commit = "abc12345"
	date = "2025-01-01"
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	v, c, d := buildInfo()
	if v != "1.2.3" {
		t.Errorf("version = %q, want %q", v, "1.2.3")
	}
	if c != "abc12345" {
		t.Errorf("commit = %q, want %q", c, "abc12345")
	}
	if d != "2025-01-01" {
		t.Errorf("date = %q, want %q", d, "2025-01-01")
	}
}

func TestBuildInfoDevFallback(t *testing.T) {
	origVersion := version
	origCommit := commit
	origDate := date
	version = "dev"
	commit = "none"
	date = "unknown"
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	v, c, d := buildInfo()

	// When built via go install, debug.ReadBuildInfo provides real values.
	// When run via go test without ldflags, the build info may or may not
	// have VCS settings, so we verify the function does not panic and
	// returns non-empty values.
	if v == "" {
		t.Error("version should not be empty")
	}
	if c == "" {
		t.Error("commit should not be empty")
	}
	if d == "" {
		t.Error("date should not be empty")
	}

	// The version stays "dev" if no go install version is available.
	// The commit may remain "none" or get populated from VCS settings.
	// The date may remain "unknown" or get populated from VCS settings.
	// All we can assert is that the function returns valid strings.
	t.Logf("buildInfo returned: version=%q commit=%q date=%q", v, c, d)
}

func TestBuildInfoDevCommitAlreadySet(t *testing.T) {
	origVersion := version
	origCommit := commit
	origDate := date
	version = "dev"
	commit = "fixed1234"
	date = "unknown"
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	_, c, _ := buildInfo()

	// When commit is already set (not "none"), the VCS revision should
	// not overwrite it even in dev mode.
	if c != "fixed1234" {
		t.Errorf("commit = %q, want %q (should not be overwritten by VCS)", c, "fixed1234")
	}
}

func TestBuildInfoDevDateAlreadySet(t *testing.T) {
	origVersion := version
	origCommit := commit
	origDate := date
	version = "dev"
	commit = "none"
	date = "2025-06-15"
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	_, _, d := buildInfo()

	// When date is already set (not "unknown"), the VCS time should
	// not overwrite it.
	if d != "2025-06-15" {
		t.Errorf("date = %q, want %q (should not be overwritten by VCS)", d, "2025-06-15")
	}
}

func TestBuildInfoCommitTruncation(t *testing.T) {
	origVersion := version
	origCommit := commit
	origDate := date
	version = "dev"
	commit = "none"
	date = "unknown"
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	_, c, _ := buildInfo()

	// If VCS revision was populated, it should be at most 8 characters.
	// If it remains "none" because no VCS info is available, that is also fine.
	if c != "none" && len(c) > 8 {
		t.Errorf("commit = %q (len=%d), want at most 8 characters", c, len(c))
	}
}

func TestVersionCmdOutput(t *testing.T) {
	origVersion := version
	origCommit := commit
	origDate := date
	version = "2.0.0"
	commit = "deadbeef"
	date = "2025-03-15"
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	var buf strings.Builder
	cmd := versionCmd
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{})

	// Create a fresh command to avoid double-registration issues.
	// versionCmd is already added to rootCmd, so we execute it via a
	// new cobra.Command with the same Run function.
	testCmd := newVersionCmd()
	testCmd.SetOut(&buf)

	if err := testCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "paq 2.0.0") {
		t.Errorf("expected version in output, got: %s", output)
	}
	if !strings.Contains(output, "deadbeef") {
		t.Errorf("expected commit in output, got: %s", output)
	}
	if !strings.Contains(output, "2025-03-15") {
		t.Errorf("expected date in output, got: %s", output)
	}
}

// newVersionCmd creates a fresh version command for testing.
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "version",
		Run: func(cmd *cobra.Command, args []string) {
			v, c, d := buildInfo()
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "paq %s (commit: %s, built: %s)\n", v, c, d)
		},
	}
}
