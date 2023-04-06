package lib

import (
	"log"
	"os/exec"
)

func CheckDependencies() (string, error) {
	dependencies := []string{"g++", "go"}

	for _, dep := range dependencies {
		if _, err := exec.LookPath(dep); err != nil {
			return dep, err
		}
	}
	log.Println("All dependency install")
	return "", nil
}
