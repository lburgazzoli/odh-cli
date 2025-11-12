package version

// Version information set by ldflags at build time.
//
//nolint:gochecknoglobals // These variables are set by ldflags during build.
var (
	// Version is the semantic version of the CLI.
	Version = "dev"
	// Commit is the git commit hash.
	Commit = "unknown"
	// Date is the build date.
	Date = "unknown"
)

// GetVersion returns the current version string.
func GetVersion() string {
	return Version
}

// GetCommit returns the git commit hash.
func GetCommit() string {
	return Commit
}

// GetDate returns the build date.
func GetDate() string {
	return Date
}
