package impl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func CopyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

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
