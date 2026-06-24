package builtin

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

type stubCommand struct{}

func (s stubCommand) Execute(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	return nil
}

func setWorkingDir(t *testing.T, dir string) {
	t.Helper()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error: %v", err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("os.Chdir(%q) error: %v", dir, err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatalf("restore working directory error: %v", err)
		}
	})
}

func resolvedPath(t *testing.T, path string) string {
	t.Helper()

	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("filepath.EvalSymlinks(%q) error: %v", path, err)
	}

	return resolved
}

func TestRegistryRegisterAndLookup(t *testing.T) {
	registry := NewRegistry()
	command := stubCommand{}

	registry.Register("stub", command)

	got, ok := registry.Get("stub")
	if !ok {
		t.Fatalf("Get(%q) reported missing command", "stub")
	}
	if got != command {
		t.Fatalf("Get(%q) returned wrong command instance", "stub")
	}
	if !registry.IsBuiltin("stub") {
		t.Fatalf("IsBuiltin(%q) = false, want true", "stub")
	}
	if registry.IsBuiltin("missing") {
		t.Fatalf("IsBuiltin(%q) = true, want false", "missing")
	}
}

func TestNewRegistryRegistersDefaultBuiltins(t *testing.T) {
	registry := NewRegistry()

	for _, name := range []string{"cd", "echo", "kill", "ps", "pwd"} {
		if !registry.IsBuiltin(name) {
			t.Fatalf("IsBuiltin(%q) = false, want true", name)
		}
	}
}

func TestEchoExecutePrintsArguments(t *testing.T) {
	var stdout bytes.Buffer

	err := NewEcho().Execute([]string{"echo", "hello", "world"}, nil, &stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Echo.Execute() error: %v", err)
	}
	if stdout.String() != "hello world\n" {
		t.Fatalf("Echo.Execute() output = %q, want %q", stdout.String(), "hello world\n")
	}
}

func TestCDExecuteErrorsOnInvalidArgCount(t *testing.T) {
	cd := NewCD()

	if err := cd.Execute([]string{"cd"}, nil, nil, nil); !errors.Is(err, ErrMissingPath) {
		t.Fatalf("CD.Execute() missing path error = %v, want %v", err, ErrMissingPath)
	}

	if err := cd.Execute([]string{"cd", "one", "two"}, nil, nil, nil); !errors.Is(err, ErrTooManyValues) {
		t.Fatalf("CD.Execute() too many args error = %v, want %v", err, ErrTooManyValues)
	}
}

func TestCDExecuteChangesWorkingDirectory(t *testing.T) {
	targetDir := t.TempDir()
	setWorkingDir(t, filepath.Dir(targetDir))
	wantDir := resolvedPath(t, targetDir)

	if err := NewCD().Execute([]string{"cd", targetDir}, nil, nil, nil); err != nil {
		t.Fatalf("CD.Execute() error: %v", err)
	}

	gotDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error: %v", err)
	}
	if gotDir != wantDir {
		t.Fatalf("working directory = %q, want %q", gotDir, wantDir)
	}
}

func TestPWDExecutePrintsCurrentWorkingDirectory(t *testing.T) {
	targetDir := t.TempDir()
	setWorkingDir(t, targetDir)
	wantDir := resolvedPath(t, targetDir)

	var stdout bytes.Buffer

	err := NewPWD().Execute([]string{"pwd"}, nil, &stdout, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("PWD.Execute() error: %v", err)
	}
	if stdout.String() != wantDir+"\n" {
		t.Fatalf("PWD.Execute() output = %q, want %q", stdout.String(), wantDir+"\n")
	}
}

func TestPSExecuteRunsSystemPS(t *testing.T) {
	if _, err := exec.LookPath("ps"); err != nil {
		t.Skip("ps command is not available")
	}

	var stdout bytes.Buffer
	err := NewPS().Execute([]string{"ps"}, nil, &stdout, &bytes.Buffer{})

	if err != nil {
		t.Fatalf("PS.Execute() error: %v", err)
	}
	if stdout.String() == "" {
		t.Fatal("PS.Execute() produced empty output")
	}
}

func TestKillExecuteRejectsInvalidArguments(t *testing.T) {
	kill := NewKill()

	if err := kill.Execute([]string{"kill", "1", "2"}, nil, nil, nil); err == nil {
		t.Fatal("Kill.Execute() with too many args returned nil error")
	}

	if err := kill.Execute([]string{"kill", "not-a-pid"}, nil, nil, nil); err == nil {
		t.Fatal("Kill.Execute() with non-numeric pid returned nil error")
	}
}

func TestKillExecuteRequiresPID(t *testing.T) {
	kill := NewKill()

	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("Kill.Execute() panicked without pid: %v", recovered)
		}
	}()

	if err := kill.Execute([]string{"kill"}, nil, nil, nil); err == nil {
		t.Fatal("Kill.Execute() without pid returned nil error")
	}
}

func TestKillExecuteTerminatesProcess(t *testing.T) {
	sleepPath, err := exec.LookPath("sleep")
	if err != nil {
		t.Skip("sleep command is not available")
	}

	cmd := exec.Command(sleepPath, "30")
	if err := cmd.Start(); err != nil {
		t.Fatalf("sleep process start error: %v", err)
	}

	waitDone := make(chan error, 1)
	go func() {
		waitDone <- cmd.Wait()
	}()

	t.Cleanup(func() {
		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			return
		}

		_ = cmd.Process.Kill()
		select {
		case <-waitDone:
		case <-time.After(2 * time.Second):
		}
	})

	if err := NewKill().Execute([]string{"kill", strconv.Itoa(cmd.Process.Pid)}, nil, nil, nil); err != nil {
		t.Fatalf("Kill.Execute() error: %v", err)
	}

	select {
	case <-waitDone:
	case <-time.After(2 * time.Second):
		t.Fatal("Kill.Execute() did not terminate the process in time")
	}
}
