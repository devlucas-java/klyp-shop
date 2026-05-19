package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(t *testing.T, target **os.File, fn func()) string {
	old := *target
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		*target = old
		_ = r.Close()
	}()

	*target = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	fn()
	_ = w.Close()
	output := <-outC
	return output
}

func TestNewLoggerWithComponentPrefix(t *testing.T) {
	output := captureOutput(t, &os.Stdout, func() {
		log := NewLoggerWithComponent(DEBUG, "AuthService")
		log.Info("user logged in")
	})

	if !strings.Contains(output, "[AuthService]") {
		t.Fatalf("expected component prefix in output, got %q", output)
	}
	if !strings.Contains(output, "[INFO]") {
		t.Fatalf("expected INFO level in output, got %q", output)
	}
	if !strings.Contains(output, "user logged in") {
		t.Fatalf("expected message in output, got %q", output)
	}
}

func TestLoggerWithComponentKeepsLevel(t *testing.T) {
	output := captureOutput(t, &os.Stdout, func() {
		log := NewLogger(INFO).WithComponent("OrderService")
		log.Debug("should not appear")
		log.Warn("inventory low")
	})

	if strings.Contains(output, "should not appear") {
		t.Fatalf("did not expect debug output when level is INFO, got %q", output)
	}
	if !strings.Contains(output, "[OrderService]") {
		t.Fatalf("expected component prefix in output, got %q", output)
	}
	if !strings.Contains(output, "inventory low") {
		t.Fatalf("expected warn message in output, got %q", output)
	}
}

func TestLoggerLevelFiltering(t *testing.T) {
	output := captureOutput(t, &os.Stdout, func() {
		log := NewLogger(WARN)
		log.Debug("debug hidden")
		log.Info("info hidden")
		log.Warn("warn shown")
	})

	if strings.Contains(output, "debug hidden") {
		t.Fatalf("expected debug message to be filtered out, got %q", output)
	}
	if strings.Contains(output, "info hidden") {
		t.Fatalf("expected info message to be filtered out, got %q", output)
	}
	if !strings.Contains(output, "warn shown") {
		t.Fatalf("expected warn message to appear, got %q", output)
	}
}

func TestLoggerErrorWritesStderr(t *testing.T) {
	output := captureOutput(t, &os.Stderr, func() {
		log := NewLogger(ERROR)
		log.Error("fatal error")
	})

	if !strings.Contains(output, "[ERROR]") {
		t.Fatalf("expected error level on stderr, got %q", output)
	}
	if !strings.Contains(output, "fatal error") {
		t.Fatalf("expected error message on stderr, got %q", output)
	}
}
