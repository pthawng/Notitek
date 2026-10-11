// Package buildinfo exposes the version control revision embedded by the Go
// toolchain, so binaries need no -ldflags to report what they were built from.
package buildinfo

import (
	"runtime/debug"
	"sync"
)

// Info describes the running binary.
type Info struct {
	GoVersion string
	Revision  string
	Modified  bool
}

var (
	once   sync.Once
	cached Info
)

// Read returns the build information of the running binary.
func Read() Info {
	once.Do(func() {
		cached = Info{Revision: "unknown"}
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			return
		}
		cached.GoVersion = bi.GoVersion
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				cached.Revision = s.Value
			case "vcs.modified":
				cached.Modified = s.Value == "true"
			}
		}
	})
	return cached
}
