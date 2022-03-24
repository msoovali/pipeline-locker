package logger

import (
	"log"
	"os"
	"testing"
)

func TestLogger_New(t *testing.T) {
	t.Run("initializesLogger_returnsInfoAndErrorLoggers", func(t *testing.T) {
		logger := New()

		if logger.Info.Writer() != os.Stdout {
			t.Errorf("Expected infologger to write into %v (stdout), got %v", os.Stdout, logger.Info.Writer())
		}
		if logger.Info.Prefix() != infoPrefix {
			t.Errorf("Expected infologger with prefix %s, got %s", infoPrefix, logger.Info.Prefix())
		}
		if logger.Info.Flags() != log.Ldate|log.Ltime {
			t.Errorf("Expected infologger with flags %d, got %d", log.Ldate|log.Ltime, logger.Info.Flags())
		}
		if logger.Error.Writer() != os.Stderr {
			t.Errorf("Expected errorlogger to write into %v (stderr), got %v", os.Stderr, logger.Info.Writer())
		}
		if logger.Error.Prefix() != errorPrefix {
			t.Errorf("Expected errorlogger with prefix %s, got %s", errorPrefix, logger.Info.Prefix())
		}
		if logger.Error.Flags() != log.Ldate|log.Ltime|log.Lshortfile {
			t.Errorf("Expected errorlogger with flags %d, got %d", log.Ldate|log.Ltime|log.Lshortfile, logger.Error.Flags())
		}
	})
}
