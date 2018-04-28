package ham

import "os"

// FileOrDirectoryExists returns true if the file or directory exists.
func FileOrDirectoryExists(fn string) bool {
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return false
	}
	return true
}
