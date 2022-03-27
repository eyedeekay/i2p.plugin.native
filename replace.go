package shellservice

import "strings"

func Replace(r string) string {
	oses := []string{"mac", "linux", "windows"}
	for _, os := range oses {
		r = strings.Replace(r, os, "$OS", -1)
	}
	arches := []string{"386", "amd64", "arm", "arm64"}
	for _, arch := range arches {
		r = strings.Replace(r, arch, "$ARCH", -1)
	}
	return r
}
