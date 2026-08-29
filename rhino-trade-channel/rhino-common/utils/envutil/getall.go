package envutil

import (
	"os"
	"strings"
)

func Env()map[string]string{
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		envs[pair[0]]=pair[1]
	}
	return envs
}
