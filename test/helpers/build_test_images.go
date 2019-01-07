package helpers

import (
	"fmt"
	"github.com/solo-io/solo-kit/test/helpers"
	"hash/crc32"
	"os"
)

var glooComponents = []string{
	"gloo",
	"discovery",
	"gateway",
	"ingress",
}
var versionTag = ""

func Version() string {
	if versionTag != "" {
		return versionTag
	}
	tag := os.Getenv("VERSION")
	// if no tag set, default to a hash of the user's hostname
	if tag == "" {
		if host, err := os.Hostname(); err == nil {
			tag = hash(host)
		} else {
			tag = helpers.RandString(4)
		}
	}

	versionTag = "testing-" + tag
	return versionTag
}

func hash(h string) string {
	crc32q := crc32.MakeTable(0xD5828281)
	return fmt.Sprintf("%08x", crc32.Checksum([]byte(h), crc32q))
}

// builds and pushes all docker containers needed for test
func BuildPushContainers(push, verbose bool) (string, error) {
	if os.Getenv("SKIP_BUILD") == "1" {
		return "", nil
	}
	version := Version()
	os.Setenv("VERSION", version)

	// make the gloo containers
	if err := RunCommand(verbose, "make", "docker"); err != nil {
		return "", err
	}

	if push {
		if err := RunCommand(verbose, "make", "docker-push"); err != nil {
			return "", err
		}
	}

	return version, nil
}