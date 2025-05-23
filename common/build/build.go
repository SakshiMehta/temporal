package build

import (
	"runtime/debug"
	"time"
)

type (
	Info struct {
		Available bool
		GoVersion string
		GoArch    string
		GoOs      string
		// CgoEnabled indicates whether cgo was enabled when this build was created.
		CgoEnabled bool

		// GitRevision is the git revision associated with this build.
		GitRevision string
		// GitTime is the git revision time.
		GitTime     time.Time
		GitModified bool
	}
)

var (
	InfoData Info
)

func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	InfoData.Available = true
	InfoData.GoVersion = buildInfo.GoVersion

	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "GOARCH":
			InfoData.GoArch = setting.Value
		case "GOOS":
			InfoData.GoOs = setting.Value
		case "CGO_ENABLED":
			InfoData.CgoEnabled = setting.Value == "1"
		case "vcs.revision":
			InfoData.GitRevision = setting.Value
		case "vcs.time":
			if gitTime, err := time.Parse(time.RFC3339, setting.Value); err == nil {
				InfoData.GitTime = gitTime
			}
		case "vcs.modified":
			InfoData.GitModified = setting.Value == "true"
		}
	}
}
