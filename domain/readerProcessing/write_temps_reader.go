package readerProcessing

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/supermegacatdog/queryCounter/entity"
	"github.com/supermegacatdog/queryCounter/pkg"
)

// Write creates and rewrites temp files
func (s *ReaderProcessingService) Write(
	existedTempFiles []*os.File,
	queriesFreq entity.GroupQueryFrequency,
	maxQueries int,
) ([]*os.File, error) {
	var newTempFiles []*os.File

	// need to pass through all created files
	// if these files are not enough, the new are created
	for i := 0; len(queriesFreq) > 0; i++ {
		if i > len(existedTempFiles)-1 {
			tempFile, err := os.Create(fmt.Sprintf(pkg.TempFileNameMask, i))
			if err != nil {
				return nil, fmt.Errorf("failed to create temp file:, %w\n", err)
			}

			existedTempFiles = append(existedTempFiles, tempFile)
			newTempFiles = append(newTempFiles, tempFile)
		}

		err := s.writeUsingReader(existedTempFiles[i], queriesFreq, maxQueries)
		if err != nil {
			return nil, fmt.Errorf("failed to write with reader: %w\n", err)
		}
	}

	return newTempFiles, nil
}

// writeUsingReader rewrites strings of temp file with overlapping values from queriesFreq
func (s *ReaderProcessingService) writeUsingReader(tempFile *os.File, queriesFreq entity.GroupQueryFrequency, maxQueries int) error {
	_, err := tempFile.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to get current position in temp file: %w\n", err)
	}

	reader := bufio.NewReader(tempFile)
	tempQueries := make(entity.GroupQueryFrequency, 0)

	// receiving a batch of parsed into the map strings from current temp file
	tempQueries, err = s.getBlockFromTempFile(reader, maxQueries)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to get block from file: %w\n", err)
	}

	// moving values from maxQueries to temp queries
	// if len(tempQueries) > N, then values are remained in maxQueries, but at the end of the function keys cannot be overlapped
	tempQueries.Add(queriesFreq, maxQueries)

	_, err = tempFile.WriteAt(tempQueries.BytesEncode(), 0)
	if err != nil {
		return fmt.Errorf("failed to write temp file at: %w\n", err)
	}

	return nil
}

// getBlockFromTempFile returns a batch of parsed into the map strings using current reader
func (s *ReaderProcessingService) getBlockFromTempFile(r *bufio.Reader, max int) (entity.GroupQueryFrequency, error) {
	group := make(entity.GroupQueryFrequency, 0)

	for len(group) < max {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return group, err
			}

			return nil, err
		}

		parsedString := &entity.StringQueryFrequency{}
		err = parsedString.Parse(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse string of temp file: %w\n", err)
		}

		group[parsedString.Query] = parsedString.Frequency
	}

	return group, nil
}
