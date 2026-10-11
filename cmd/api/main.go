// Command api runs the Notitek api runtime.
package main

import (
	"os"

	"github.com/pthawng/notitek/internal/bootstrap"
)

func main() {
	os.Exit(bootstrap.Main(bootstrap.RuntimeAPI))
}
