package tunebooklib

import (
	"fmt"
)

var AppName = "tunebook-lib"

var Version = "no-version"
var Commit = "no-commit"

func VersionTemplate(appName string) string {
	return fmt.Sprintf(
		"%s: %s (%s)\n",
		appName, Version, Commit)
}
