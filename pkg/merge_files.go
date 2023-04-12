package pkg

import (
	"fmt"
	"io"
	"os"
)

func MergeFiles(files []*os.File, outFile *os.File) error {
	for _, file := range files {
		_, err := file.Seek(0, 0)
		if err != nil {
			return fmt.Errorf("failed to get current position: %w\n", err)
		}

		_, err = io.Copy(outFile, file)
		if err != nil {
			return fmt.Errorf("failed to copy file into another: %w\n", err)
		}

		err = os.Remove(file.Name())
		if err != nil {
			return fmt.Errorf("failed to remove file: %w\n", err)
		}
	}

	return nil
}
