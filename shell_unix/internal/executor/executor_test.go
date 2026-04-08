package executor

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/parser"
)

func TestExecuteRedirectWritesBuiltinOutputToFile(t *testing.T) {
	exec := New().withIO(nil, io.Discard, io.Discard)
	outputFile := filepath.Join(t.TempDir(), "out.txt")

	node := &parser.Redirect{
		Type:   parser.RedirectOutput,
		Target: outputFile,
		Node: &parser.Command{
			Name: "echo",
			Args: []string{"hello", "redirect"},
		},
	}

	if err := exec.Execute(node); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("os.ReadFile() error: %v", err)
	}

	if string(content) != "hello redirect\n" {
		t.Fatalf("redirect output = %q, want %q", string(content), "hello redirect\n")
	}
}

func TestExecuteInputRedirectFeedsExternalCommand(t *testing.T) {
	if _, err := exec.LookPath("cat"); err != nil {
		t.Skip("cat command is not available")
	}

	inputFile := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(inputFile, []byte("redirect input\n"), 0o644); err != nil {
		t.Fatalf("os.WriteFile() error: %v", err)
	}

	var stdout bytes.Buffer
	exec := New().withIO(nil, &stdout, io.Discard)

	node := &parser.Redirect{
		Type:   parser.RedirectInput,
		Target: inputFile,
		Node: &parser.Command{
			Name: "cat",
		},
	}

	if err := exec.Execute(node); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if stdout.String() != "redirect input\n" {
		t.Fatalf("redirect output = %q, want %q", stdout.String(), "redirect input\n")
	}
}

func TestExecutePipelineRunsBuiltinIntoExternalCommand(t *testing.T) {
	if _, err := exec.LookPath("cat"); err != nil {
		t.Skip("cat command is not available")
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exec := New().withIO(nil, &stdout, &stderr)

	node := &parser.Pipeline{
		Commands: []parser.Node{
			&parser.Command{
				Name: "echo",
				Args: []string{"hello", "pipeline"},
			},
			&parser.Command{
				Name: "cat",
			},
		},
	}

	if err := exec.Execute(node); err != nil {
		t.Fatalf("Execute() error: %v (stderr=%q)", err, stderr.String())
	}

	if stdout.String() != "hello pipeline\n" {
		t.Fatalf("pipeline output = %q, want %q", stdout.String(), "hello pipeline\n")
	}
}

func TestInterruptForegroundStopsExternalCommand(t *testing.T) {
	sleepPath, err := exec.LookPath("sleep")
	if err != nil {
		t.Skip("sleep command is not available")
	}

	exec := New().withIO(nil, io.Discard, io.Discard)
	done := make(chan error, 1)

	go func() {
		done <- exec.Execute(&parser.Command{
			Name: sleepPath,
			Args: []string{"30"},
		})
	}()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if exec.state.current != nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if exec.state.current == nil {
		t.Fatal("foreground process was not registered")
	}

	if err := exec.InterruptForeground(); err != nil {
		t.Fatalf("InterruptForeground() error: %v", err)
	}

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("external command exited without interruption error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("external command was not interrupted in time")
	}
}

func TestInterruptForegroundStopsExternalPipeline(t *testing.T) {
	sleepPath, err := exec.LookPath("sleep")
	if err != nil {
		t.Skip("sleep command is not available")
	}
	catPath, err := exec.LookPath("cat")
	if err != nil {
		t.Skip("cat command is not available")
	}

	exec := New().withIO(nil, io.Discard, io.Discard)
	done := make(chan error, 1)

	go func() {
		done <- exec.Execute(&parser.Pipeline{
			Commands: []parser.Node{
				&parser.Command{
					Name: sleepPath,
					Args: []string{"30"},
				},
				&parser.Command{
					Name: catPath,
				},
			},
		})
	}()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if exec.state.current != nil && len(exec.state.current.List()) > 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	if exec.state.current == nil {
		t.Fatal("foreground pipeline was not registered")
	}

	if err := exec.InterruptForeground(); err != nil {
		t.Fatalf("InterruptForeground() error: %v", err)
	}

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("pipeline exited without interruption error")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("pipeline was not interrupted in time")
	}
}
