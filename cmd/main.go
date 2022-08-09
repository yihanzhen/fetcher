package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var module = flag.String("module", "", "module to list versions for")
var username = flag.String("username", "", "username used to retrieve data")
var password = flag.String("password", "", "password used to retrieve data")
var host = flag.String("host", "", "host of the remote repository")

func main() {
	flag.Parse()
	if err := saveCredentials(*host, *username, *password); err != nil {
		panic(err)
	}
	listVersionResult, err := listVersions(*module)
	if err != nil {
		panic(err)
	}
	fmt.Println(listVersionResult)
}

type listVersionsResult struct {
	Path     string
	Versions []string
}

type Info struct {
	Version string    // version string
	Time    time.Time // commit time
}

func listVersions(module string) (string, error) {
	cmd := exec.Command("go", "list", "-versions", "-m", "-json", module)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("listVersions: cmd.Output: %v, %v", err, string(output))
	}
	var result listVersionsResult
	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("listVersions: %v", err)
	}
	return strings.Join(result.Versions, "\n") + "\n", nil
}

func saveCredentials(host, username, password string) error {
	if host == "" {
		return fmt.Errorf("saveCredentials: host is required")
	}
	if username == "" {
		return fmt.Errorf("saveCredentials: username is required")
	}
	if password == "" {
		return fmt.Errorf("saveCredentials: password is required")
	}

	dirname, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("saveCredentials: %v", err)
	}
	f, err := os.OpenFile(dirname+"/.netrc", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("saveCredentials: %v", err)
	}
	if _, err := f.Write([]byte(fmt.Sprintf(`
machine %s
login %s
password %s
`, host, username, password))); err != nil {
		return fmt.Errorf("saveCredentials: %v", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("saveCredentials: %v", err)
	}

	if err := os.Setenv("GONOSUMDB", host); err != nil {
		return fmt.Errorf("saveCredentials: %v", err)
	}
	return nil
}
