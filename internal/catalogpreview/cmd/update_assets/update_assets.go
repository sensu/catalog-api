package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/mgutz/ansi"
)

// Bucket
const packageURL = "https://s3.us-west-2.amazonaws.com/sensu-ci-web-builds"
const packagePathPrefix = "enterprise/catalog"
const packageFilename = "build.tgz"

// Path
var (
	baseDir       = filepath.Join(".")
	assetsDir     = filepath.Join(baseDir, "assets")
	assetsRefPath = filepath.Join(baseDir, "assetsref")
)

// x_x
var (
	bold = ansi.ColorFunc("+hb")
	red  = ansi.ColorFunc("red+b")
	blue = ansi.ColorFunc("cyan+h")
)

//
// This script attempts to update the preview web application's assets by
// fetching the build from S3 or a specified path.
//
func main() {
	infoln("read the given ref from assets path")
	ref, err := readRef(assetsRefPath)
	if err != nil {
		fatalln(err.Error())
		return
	}

	infoln("attempt to fetch build from S3 bucke")
	buildDir, err := fetchBuildFromBucket(ref)
	if err != nil {
		fatalln("unable to fetch build from bucket:", err.Error())
		return
	}
	defer os.RemoveAll(buildDir)

	infoln("expanding assets; stats & source maps are omitted")
	mustRunCmd("sh", "-c", "find "+filepath.Join(buildDir, "build", "catalog")+" -type f -name 'stats.json' -exec rm {} +")
	mustRunCmd("sh", "-c", "find "+filepath.Join(buildDir, "build", "catalog")+" -type f -name '*.map'      -exec rm {} +")

	infoln("copying new files")
	mustRunCmd("mv", assetsDir, filepath.Join(assetsDir+".archived"))
	defer mustRunCmd("rm", "-rf", filepath.Join(assetsDir+".archived"))
	mustRunCmd("mv", filepath.Join(buildDir, "build", "catalog"), assetsDir)
}

func readRef(path string) (string, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	vs := string(f)
	if len(vs) < 40 {
		return vs, fmt.Errorf("given ref file doesn't look right, must contain valid git SHA")
	}
	return vs[0:40], err
}

func fetchBuildFromBucket(ref string) (string, error) {
	// ensure that git is installed
	if _, err := exec.LookPath("git"); err != nil {
		fatalln(
			"'git' was not found in your PATH, unable to bundle the web UI.",
			"See https://git-scm.com/downloads for installation instructions.",
		)
	}

	// attempt to pull the build from the S3 bucket
	fpath, err := fetchPackage(ref)
	if err != nil {
		return "", err
	}
	defer os.Remove(fpath)

	// extract the build
	return extractPackage(fpath)
}

func fetchPackage(ref string) (string, error) {
	tmpFile, err := ioutil.TempFile("", "*."+packageFilename)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	url := packageURL + "/" + path.Join(packagePathPrefix, ref, packageFilename)
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	debugln(fmt.Sprintf("%d: %s", res.StatusCode, url))

	if res.StatusCode != 200 {
		return "", errors.New("received non 200 status")
	}

	_, err = io.Copy(tmpFile, res.Body)
	if err != nil {
		return "", err
	}

	if err = tmpFile.Sync(); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func extractPackage(path string) (string, error) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	mustRunCmd("tar", "-zxpf", path, "-C", tmpDir, "--strip-components=1")
	return tmpDir, nil
}

func mustRunCmd(pro string, args ...string) string {
	cmd := exec.Command(pro, args...)
	cmdStr := strings.Join(append([]string{pro}, args...), " ")

	debugln(fmt.Sprintf("running '%s'", cmd.String()))

	buf, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = io.Copy(os.Stderr, bytes.NewReader(buf))
		fatalln("failed to run", cmdStr, err.Error())
	}
	return string(buf[:])
}

func debugln(str ...string) {
	fmt.Println(bold("debug"), strings.Join(str, " "))
}

func infoln(str ...string) {
	fmt.Println(blue("info"), "", strings.Join(str, " "))
}

func errln(str ...string) {
	fmt.Println(red("error"), strings.Join(str, " "))
}

func fatalln(str ...string) {
	errln(str...)
	os.Exit(1)
}
