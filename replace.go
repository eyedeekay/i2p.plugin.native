package shellservice

import "strings"

var OSES = []string{"mac", "linux", "windows"}
var ARCHES = []string{"386", "amd64", "arm", "arm64"}

func Replace(r string) string {
	for _, os := range OSES {
		r = strings.Replace(r, os, "$OS", -1)
	}
	for _, arch := range ARCHES {
		r = strings.Replace(r, arch, "$ARCH", -1)
	}
	return r
}
