package version

import "fmt"

var (
	name       = "mirage"
	version    = "unknown"
	commit     = "unknown"
	commitLong = "unknown"
	branch     = "unknown"
	date       = "unknown"
)

// Name returns main application name
func Name() string {
	return name
}

// Version returns current application version
func Version() string {
	return version
}

// FullName returns current app version, commit and build time
func FullName() string {
	return fmt.Sprintf(
		`Name: %v
Version: %v
branch: %v
commit: %v
commit long: %v
build date: %v`, Name(), Version(), branch, commit, commitLong, date)
}
