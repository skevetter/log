package log

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestWithFields(t *testing.T) {
	t.Run("TextFormat", func(t *testing.T) {
		var output bytes.Buffer
		logger := NewStreamLogger(&output, &output, logrus.InfoLevel)
		logger.SetFormat(TextFormat)

		logger.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Info("A group of walrus emerges from the ocean")

		logOutput := output.String()
		if !strings.Contains(logOutput, "animal=walrus") {
			t.Errorf("expected output to contain animal=walrus, got %s", logOutput)
		}
		if !strings.Contains(logOutput, "size=10") {
			t.Errorf("expected output to contain size=10, got %s", logOutput)
		}
	})

	t.Run("JSONFormat", func(t *testing.T) {
		var output bytes.Buffer
		logger := NewStreamLogger(&output, &output, logrus.InfoLevel)
		logger.SetFormat(JSONFormat)

		logger.WithFields(logrus.Fields{
			"animal": "walrus",
			"size":   10,
		}).Info("A group of walrus emerges from the ocean")

		logOutput := output.String()
		var logEntry map[string]any
		err := json.Unmarshal([]byte(logOutput), &logEntry)
		if err != nil {
			t.Fatalf("failed to unmarshal JSON output: %v", err)
		}

		fields, ok := logEntry["fields"].(map[string]any)
		if !ok {
			t.Fatalf("expected fields in JSON output, got %v", logEntry)
		}

		if fields["animal"] != "walrus" {
			t.Errorf("expected animal=walrus, got %v", fields["animal"])
		}
		if fields["size"] != 10.0 { // JSON numbers are float64
			t.Errorf("expected size=10, got %v", fields["size"])
		}
	})

	t.Run("GlobalWithFields", func(t *testing.T) {
		// This just checks if it compiles and runs without panic,
		// capturing output of global logger is harder without mocking os.Stdout
		l := WithFields(logrus.Fields{"foo": "bar"})
		if l == nil {
			t.Fatal("expected logger, got nil")
		}
	})

	t.Run("DiscardLogger", func(t *testing.T) {
		logger := NewDiscardLogger(logrus.InfoLevel)
		l := logger.WithFields(logrus.Fields{"test": "value"})
		if l == nil {
			t.Fatal("expected logger, got nil")
		}
		// Should not panic
		l.Info("test message")
	})

	t.Run("FileLogger", func(t *testing.T) {
		tmpFile := "/tmp/test.log"
		logger := NewFileLogger(tmpFile, logrus.InfoLevel)
		l := logger.WithFields(logrus.Fields{"file": "test"})
		if l == nil {
			t.Fatal("expected logger, got nil")
		}
		l.Info("test message")
	})
}
