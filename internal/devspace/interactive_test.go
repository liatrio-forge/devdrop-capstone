package devspace

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"charm.land/huh/v2"
)

// Confirm/Input fields on a bytes.Buffer (non-terminal) always take the
// non-interactive fallback path, since isInteractiveTerminal requires real
// os.File terminals for both in and out. This is the same environment every
// other command test in this package already runs under.

func TestConfirmSetupRunFallsBackToPlainConfirmWhenPiped(t *testing.T) {
	var out bytes.Buffer
	err := confirmSetupRun(strings.NewReader("wrong\n"), &out, "web", "npm install", "apps/web")
	if err == nil || !strings.Contains(err.Error(), "confirmation did not match") {
		t.Fatalf("confirmSetupRun error = %v, want confirmation mismatch", err)
	}
	if !strings.Contains(out.String(), "Type web to run `npm install` in apps/web: ") {
		t.Fatalf("expected the original typed-phrase prompt text, got %q", out.String())
	}
}

func TestConfirmSetupRunAcceptsMatchingPhraseWhenPiped(t *testing.T) {
	var out bytes.Buffer
	err := confirmSetupRun(strings.NewReader("web\n"), &out, "web", "npm install", "apps/web")
	if err != nil {
		t.Fatalf("confirmSetupRun error = %v, want nil for a matching phrase", err)
	}
}

func TestConfirmSetupApplyFallsBackToPlainConfirmWhenPiped(t *testing.T) {
	var out bytes.Buffer
	prompt := "Type run all to run install commands for every runnable project: "
	err := confirmSetupApply(strings.NewReader("nope\n"), &out, prompt, "run all")
	if err == nil || !strings.Contains(err.Error(), "confirmation did not match") {
		t.Fatalf("confirmSetupApply error = %v, want confirmation mismatch", err)
	}
}

func TestScanCommandPrintsPlainProgressLineWhenPiped(t *testing.T) {
	initCommandWorkspace(t)
	stdout, _, err := executeCommand(t, "test", "scan")
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}
	if !strings.Contains(stdout, "Scanning workspace...") {
		t.Fatalf("expected plain progress line for non-terminal output, got %q", stdout)
	}
}

func TestHydrateCommandPrintsPlainProgressLineWhenPiped(t *testing.T) {
	workspace := initCommandWorkspace(t)
	remote := hardeningBareRepo(t)
	if err := SaveManifest(workspace, Manifest{
		Version:       ManifestVersion,
		WorkspaceRoot: workspace,
		Projects:      []Project{hardeningProject("apps/lazy", ProjectTypeGit, remote)},
	}); err != nil {
		t.Fatal(err)
	}
	stdout, _, err := executeCommand(t, "test", "project", "hydrate", "lazy")
	if err != nil {
		t.Fatalf("hydrate error: %v", err)
	}
	if !strings.Contains(stdout, "Hydrating lazy...") {
		t.Fatalf("expected plain progress line for non-terminal output, got %q", stdout)
	}
}

// The following two tests drive huh's accessible mode directly (rather than
// through confirmSetupRun/confirmSetupApply, which require a real pty to
// reach their interactive branch) to prove the huh Confirm/Input
// configuration used in production -- same defaults, same validation rule --
// behaves correctly. This is the automated-coverage remediation for the
// planning audit's TTY-only-rendering regression-risk flag.

func TestHuhConfirmAccessibleDefaultsToNo(t *testing.T) {
	var out bytes.Buffer
	var confirmed bool
	form := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Run `npm install` in apps/web?").
			Affirmative("Yes").
			Negative("No").
			Value(&confirmed),
	)).WithShowHelp(false).WithAccessible(true).WithInput(strings.NewReader("\n")).WithOutput(&out)
	if err := form.Run(); err != nil {
		t.Fatalf("form.Run: %v", err)
	}
	if confirmed {
		t.Fatal("expected the Confirm field to default to No when the user presses enter with no input")
	}
}

func TestHuhInputAccessibleValidatesExactPhraseWithRetry(t *testing.T) {
	var out bytes.Buffer
	var answer string
	form := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title("Type run all: ").
			Value(&answer).
			Validate(func(s string) error {
				if s != "run all" {
					return errors.New("confirmation did not match")
				}
				return nil
			}),
	)).WithShowHelp(false).WithAccessible(true).WithInput(strings.NewReader("wrong\nrun all\n")).WithOutput(&out)
	if err := form.Run(); err != nil {
		t.Fatalf("form.Run: %v", err)
	}
	if answer != "run all" {
		t.Fatalf("answer = %q, want %q after the validator rejected the first attempt", answer, "run all")
	}
	if !strings.Contains(out.String(), "confirmation did not match") {
		t.Fatalf("expected the validator's rejection message to be shown before the retry, got %q", out.String())
	}
}
