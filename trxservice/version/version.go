package version

import (
	"fmt"
	"runtime"
)

var (
	GitBranch string
	GitTag       string
	GitCommit    string
	GitTreeState string
	Version string
	BuildTime    string
	GoVersion    string
)


// Info contains versioning information.
type Info struct {
	Version      string `json:"version"`
	GitBranch    string `json:"gitBranch"`
	GitTag       string `json:"gitTag"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildTime    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	return info.Platform
}

func GetVersion() Info {
	return Info{
		Version:      Version,
		GitBranch:    GitBranch,
		GitTag:       GitTag,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildTime:    BuildTime,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func ShowVersion()  {
	v := GetVersion()
	fmt.Printf("Version: %s\nGitBranch: %s\nCommitId: %s\nBuild Date: %s\nGo Version: %s\nOS/Arch: %s\n", v.Version, v.GitBranch, v.GitCommit, v.BuildTime, v.GoVersion, v.Platform)
}
