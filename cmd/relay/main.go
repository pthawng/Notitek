// Command relay runs the Notitek relay runtime.
package main

import (
	"os"

	"github.com/pthawng/notitek/internal/bootstrap"
)

func main() {
	os.Exit(bootstrap.Main(bootstrap.RuntimeRelay))
}
