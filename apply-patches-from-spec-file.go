package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("You need to provide one argument which points to the spec file.")
	}
	specFile := os.Args[1]

	specDir, err := filepath.Abs(filepath.Dir(specFile))
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadFile(specFile)
	if err != nil {
		log.Fatal(err)
	}
	spec := string(content)

	// Example patch tag line:
	// Patch17:	0001-Driver-Prefer-gcc-toolchains-with-libgcc_s.so-when-n.patch
	patchTagLines := regexp.MustCompilePOSIX(`^Patch[0-9]+:.*`).FindAllString(spec, -1)

	// Example patch line:
	// %patch17 -p1 -b .check-gcc_s
	// TODO(kwk): this only finds *numbered* patch lines which is the predominant case. -P is not considered!
	patchLines := regexp.MustCompilePOSIX(`^%patch[0-9]+.*`).FindAllString(spec, -1)

	// Create a lookup table from the Patch Tag Lines
	patchMap := map[uint64]string{}
	for _, patchTagLine := range patchTagLines {
		arr := strings.Split(patchTagLine, ":")
		match := regexp.MustCompilePOSIX(`[0-9]+`).FindString(arr[0])
		number, err := strconv.ParseUint(match, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		abs, err := filepath.Abs(filepath.Join(specDir, strings.TrimSpace(arr[1])))
		if err != nil {
			log.Fatal(err)
		}
		patchMap[number] = abs
	}

	// Create a list of patch commands to be executed
	commands := []string{}
	for _, patchLine := range patchLines {
		patchNum := regexp.MustCompilePOSIX(`[0-9]+`).FindString(patchLine)
		number, err := strconv.ParseUint(patchNum, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		patchFile := patchMap[number]
		// patchOptions := regexp.MustCompilePOSIX(`^%patch[0-9]+`).ReplaceAllString(patchLine, "")
		patchOptions := regexp.MustCompilePOSIX(`-p[:space:]*[0-9]+`).FindString(patchLine)
		cmd := fmt.Sprintf("/usr/bin/patch %s <%s\n", patchOptions, patchFile)
		commands = append(commands, cmd)
	}

	for _, cmd := range commands {
		fmt.Print("UPCOMING PATCH OPERATION: " + cmd)
		cmd := exec.Command("sh", "-c", cmd)
		stdoutStderr, err := cmd.CombinedOutput()
		fmt.Println(string(stdoutStderr))
		if err != nil {
			log.Fatal(err)
		}
	}
}
