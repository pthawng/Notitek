// Command scheduler runs the Notitek scheduler runtime.
package main

import (
	"os"

	"github.com/pthawng/notitek/internal/bootstrap"
)

func main() {
	os.Exit(bootstrap.Main(bootstrap.RuntimeScheduler))
}
