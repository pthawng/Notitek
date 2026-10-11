// Command worker runs the Notitek worker runtime.
package main

import (
	"os"

	"github.com/pthawng/notitek/internal/bootstrap"
)

func main() {
	os.Exit(bootstrap.Main(bootstrap.RuntimeWorker))
}
