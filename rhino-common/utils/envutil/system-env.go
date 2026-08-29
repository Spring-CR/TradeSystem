package envutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const BASHRC = ".bashrc"
const BASHPROFILE = ".bash_profile"
const ZSHRC = ".zshrc"
const MACPROFILE = ".profile"

func detectProfile() (string, error) {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return "", fmt.Errorf("home directory does not exist")
	}
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "", fmt.Errorf("shell does not exist")
	}
	cmd := exec.Command(shell, "-c", "echo -n \\\"${BASH_VERSION}\\\"")
	bashVer, err := cmd.CombinedOutput()
	if err == nil {
		if len(bashVer) > 2 {
			_, err := os.Stat(filepath.Join(homeDir, BASHRC))
			if err == nil {
				return filepath.Join(homeDir, BASHRC), nil
			}
			_, err = os.Stat(filepath.Join(homeDir, BASHPROFILE))
			if err == nil {
				return filepath.Join(homeDir, BASHPROFILE), nil
			}
		}
	}
	cmd = exec.Command(shell, "-c", "echo -n \\\"${ZSH_VERSION}\\\"")
	zshVer, err := cmd.CombinedOutput()
	if err == nil {
		if len(zshVer) > 2 {
			_, err := os.Stat(filepath.Join(homeDir, ZSHRC))
			if err == nil {
				return filepath.Join(homeDir, ZSHRC), nil
			}
		}
	}
	_, err = os.Stat(filepath.Join(homeDir, MACPROFILE))
	if err == nil {
		return filepath.Join(homeDir, MACPROFILE), nil
	}
	return "", fmt.Errorf("profile does not exist")
}

func setEnvOnUnixLike(key, value string, addToPath bool) error {
	profile, err := detectProfile()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(profile, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	r, err := regexp.Compile(fmt.Sprintf("export *%s=.*", key))
	if err != nil {
		return err
	}
	matchGroups := r.FindAll(content, -1)
	if len(matchGroups) == 0 {
		// env not exist in profile, write export line to file
		_, err = f.WriteString(fmt.Sprintf("\nexport %s=%s\n", key, value))
		if err != nil {
			return err
		}
		if addToPath {
			if !strings.Contains(string(content), fmt.Sprintf("export PATH=$%s:$PATH", key)) {
				_, err = f.WriteString(fmt.Sprintf("\nexport PATH=$%s:$PATH\n", key))
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	modifiedContent := r.ReplaceAll(content, []byte(fmt.Sprintf("export %s=%s", key, value)))
	if err != nil {
		return err
	}
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = f.Write(modifiedContent)
	if err != nil {
		return err
	}
	if addToPath {
		if !strings.Contains(string(content), fmt.Sprintf("export PATH=$%s:$PATH", key)) {
			_, err = f.WriteString(fmt.Sprintf("\nexport PATH=$%s:$PATH\n", key))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SetSystemEnv(key, value string, addToPath bool) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("setx.exe", key, value)
		_, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("cannot set system env, err %v", err)
		}
		if !addToPath {
			return nil
		}
		pathEnv := os.Getenv("PATH")
		for _, p := range strings.Split(pathEnv, ";") {
			if value == p {
				return nil
			}
		}
		pathModified := fmt.Sprintf("%s;%s", value, pathEnv)
		cmd = exec.Command("setx.exe", "PATH", pathModified)
		_, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("cannot set system env, err %v", err)
		}
		return nil
	}
	err := setEnvOnUnixLike(key, value, addToPath)
	if err != nil {
		return fmt.Errorf("cannot set system env, err %v", err)
	}
	return nil
}

func LocateFile(fileName string) ([]string, error) {
	var execCmd = "which"
	if runtime.GOOS == "windows" {
		execCmd = "where.exe"
	}

	cmd := exec.Command(execCmd, fileName)
	output, err := cmd.Output()
	if err != nil {
		if !strings.Contains(err.Error(), "exit status 1") {
			return nil, err
		}
	}

	if runtime.GOOS != "windows" {
		return []string{strings.TrimSpace(string(output))}, nil
	}

	foundPath := make([]string, 0)

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	pathEnv := os.Getenv("PATH")
	var cwdInPath bool
	for _, p := range strings.Split(pathEnv, ";") {
		if cwd == p {
			cwdInPath = true
			break
		}
	}

	for _, f := range strings.Split(strings.TrimSpace(string(output)), "\r\n") {
		if filepath.Base(f) == fileName {
			if strings.Contains(f, cwd) && !cwdInPath {
				continue
			}
			foundPath = append(foundPath, f)
		}
	}
	return foundPath, nil
}
