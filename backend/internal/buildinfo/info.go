package buildinfo

import "runtime"

var (
	version = "dev"
)

type Info struct {
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
}

func Read() Info {
	return Info{
		Version:   version,
		GoVersion: runtime.Version(),
	}
}
