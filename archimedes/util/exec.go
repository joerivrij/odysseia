package util

import (
	"bufio"
	"fmt"
	"os/exec"
)

func ExecCommand(command, filePath string) error {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = filePath

	stdOut, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdOut)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)
	}
	cmd.Wait()

	return nil
}

func ExecCommandWithReturn(command, filePath string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Dir = filePath

	stdOut, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		return "", err
	}

	var text string
	scanner := bufio.NewScanner(stdOut)
	for scanner.Scan() {
		text = scanner.Text()
	}
	cmd.Wait()

	return text, nil
}
