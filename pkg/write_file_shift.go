package pkg

import (
	"fmt"
	"io"
	"os"
)

// NOTE: not needed

// WriteFileWithShift writes byteGroup into the file at current position with shift
func WriteFileWithShift(file *os.File, byteGroup []byte, shift int64) error {
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stats: %w\n", err)
	}

	pos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("failed to temp file seek when shift: %w\n", err)
	}

	if shift > 0 {
		remainingBytes := make([]byte, stat.Size()-pos)

		_, err = file.ReadAt(remainingBytes, pos)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read temp file at when shift: %w\n", err)
		}

		byteGroup = append(byteGroup, remainingBytes...)
	}

	_, err = file.WriteAt(byteGroup, pos-shift-int64(len(byteGroup)))
	if err != nil {
		return fmt.Errorf("failed to write temp file at when shift: %w\n", err)
	}

	return nil
}
