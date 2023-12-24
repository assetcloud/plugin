package version

import (
	"fmt"
	"runtime"
)

//var version control
var (
<<<<<<< HEAD
	Version   = "1.68.4"
=======
	Version   = "1.68.2"
>>>>>>> c7f90405aaee21be0ead6b9b99d81458fc46c947
	GitCommit string
	BuildTime string
	// GoVersion system go version
	GoVersion = runtime.Version()
	// Platform info
	Platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
)

//GetVersion 获取版本信息
func GetVersion() string {
	if GitCommit != "" {
		return Version + "-" + GitCommit
	}
	return Version
}
