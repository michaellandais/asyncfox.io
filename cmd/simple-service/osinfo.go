package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type OSInfo struct {
	Family       string
	Architecture string
	ID           string
	Name         string
	Codename     string
	Version      string
	Build        string
}

func readCommandOutput(cmd string, arg ...string) (result string, err error) {
	command := exec.Command(cmd, arg...)
	var bytes []byte
	bytes, err = command.CombinedOutput()
	if err == nil {
		result = strings.TrimSpace(string(bytes))
	}

	return
}

func GetOSInfo() (*OSInfo, error) {
	switch runtime.GOOS {
	case "linux":
		return getOSInfoLinux()
	default:
		return getOSInfoUnknown()
	}
}

func readTextFile(path string) (result string, err error) {
	var bytes []byte
	bytes, err = os.ReadFile(path)
	if err == nil {
		result = string(bytes)
	}
	return
}

func populateFromRuntime(info *OSInfo) {
	info.Architecture = runtime.GOARCH
	info.Family = runtime.GOOS
}

func parseEtcOSRelease(info *OSInfo, contents string) {
	keyvalues := parseKeyValues(contents)

	if v, ok := keyvalues["ID"]; ok && info.ID == "" {
		info.ID = v
	}
	if v, ok := keyvalues["VERSION_ID"]; ok && info.Version == "" {
		info.Version = v
	}
	if v, ok := keyvalues["NAME"]; ok && info.Name == "" {
		info.Name = v
	}
	if v, ok := keyvalues["VERSION_CODENAME"]; ok && info.Codename == "" {
		info.Codename = v
	}
}

func parseEtcLSBRelease(info *OSInfo, contents string) {
	keyvalues := parseKeyValues(contents)

	if v, ok := keyvalues["DISTRIB_ID"]; ok && info.ID == "" {
		info.ID = v
	}
	if v, ok := keyvalues["DISTRIB_RELEASE"]; ok && info.Version == "" {
		info.Version = v
	}
	if v, ok := keyvalues["DISTRIB_CODENAME"]; ok && info.Codename == "" {
		info.Codename = v
	}
	if v, ok := keyvalues["DISTRIB_DESCRIPTION"]; ok && info.Name == "" {
		info.Name = v
	}
}

func parseKeyValues(contents string) (kvmap map[string]string) {
	kvmap = make(map[string]string)
	re := regexp.MustCompile(`\b(.+)="?([^"\n]*)"?`)
	for _, found := range re.FindAllStringSubmatch(contents, -1) {
		kvmap[found[1]] = found[2]
	}
	return
}

func getOSInfoLinux() (info *OSInfo, err error) {
	info = new(OSInfo)
	populateFromRuntime(info)

	var contents string
	if contents, err = readTextFile("/etc/os-release"); err == nil {
		parseEtcOSRelease(info, contents)
	}

	lastError := err

	if contents, err = readTextFile("/etc/lsb-release"); err == nil {
		parseEtcLSBRelease(info, contents)
	}

	// Only propagate an error if both files failed to load
	if lastError == nil {
		err = nil
	}

	return
}

func getOSInfoUnknown() (info *OSInfo, err error) {
	info = new(OSInfo)
	populateFromRuntime(info)
	info.ID = "unknown"
	info.Name = "unknown"
	info.Version = "unknown"

	// Try to fill with contents of `uname`.
	var contents string
	contents, err = readCommandOutput("/usr/bin/uname")
	if err == nil {
		info.Name = contents
	}

	err = fmt.Errorf("%v: Unhandled OS", runtime.GOOS)

	return
}
