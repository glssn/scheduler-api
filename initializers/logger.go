package initializers

import (
	"log"
	"os"
)

// Logger returns a logger instance that writes to the standard output stream.
func Logger() *log.Logger {
	// create a new logger instance
	return log.New(os.Stdout, "", log.LstdFlags)
}
