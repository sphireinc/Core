package core

import "runtime/debug"

// Version is intended to be set at build time, e.g.
// go build -ldflags "-X github.com/sphireinc/core/v1.Version=v1.3.9"
var Version = "dev"

// ResolveVersion prefers the injected Version. If not set, it attempts to
// read module version/build metadata as a fallback.
func ResolveVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return Version
	}

	// Prefer module version if present
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	// Fall back to short VCS revision if available
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" && s.Value != "" {
			rev := s.Value
			if len(rev) > 7 {
				rev = rev[:7]
			}
			return "dev-" + rev
		}
	}

	return Version
}
