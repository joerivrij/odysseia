package util

import (
	"fmt"
	"github.com/kpango/glg"
	"io"
	"os"
)

func WriteFile(input []byte, outputFile string) {
	openedFile, err := os.Create(outputFile)
	if err != nil {
		glg.Error(err)
	}
	defer openedFile.Close()

	outputFromWrite, err := openedFile.Write(input)
	if err != nil {
		glg.Error(err)
	}

	glg.Info(fmt.Sprintf("finished writing %d bytes", outputFromWrite))
	glg.Info(fmt.Sprintf("file written to %s", outputFile))
}

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

func CopyFile(src, dst string) (err error) {
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
