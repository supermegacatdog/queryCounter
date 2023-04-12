package scannerProcessing

import (
	"bufio"
	"fmt"
	"os"

	"github.com/supermegacatdog/queryCounter/entity"
	"github.com/supermegacatdog/queryCounter/pkg"
)

// Write creates and rewrites temp files
func (s *ScannerProcessingService) Write(
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

		err := s.writeUsingScanner(existedTempFiles[i], queriesFreq, maxQueries)
		if err != nil {
			return nil, fmt.Errorf("failed to write with reader: %w\n", err)
		}
	}

	return newTempFiles, nil
}

// writeUsingScanner rewrites strings of temp file with overlapping values from queriesFreq
func (s *ScannerProcessingService) writeUsingScanner(tempFile *os.File, queriesFreq entity.GroupQueryFrequency, maxQueries int) error {
	stat, err := tempFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get temp file stat: %w\n", err)
	}

	// if a temp file is new, its is much more convenient to save data here and return
	if stat.Size() == 0 {
		_, er := tempFile.Write(queriesFreq.BytesEncode())
		if er != nil {
			return fmt.Errorf("failed to write to new temp file: %w\n", er)
		}

		queriesFreq.Clear()
		return nil
	}

	_, err = tempFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to set start position in temp file: %w\n", err)
	}

	tempScanner := bufio.NewScanner(tempFile)
	tempQueries := make(entity.GroupQueryFrequency, 0)

	for tempScanner.Scan() {
		parsedString := &entity.StringQueryFrequency{}
		er := parsedString.Parse(tempScanner.Text())
		if er != nil {
			return fmt.Errorf("failed to parse string of temp file: %w\n", er)
		}

		if len(tempQueries) < maxQueries {
			tempQueries[parsedString.Query] = parsedString.Frequency
			continue
		}

		// moving values from maxQueries to temp queries
		// if len(tempQueries) > N, then values are remained in maxQueries, but at the end of the function keys cannot be overlapped
		tempQueries.Add(queriesFreq, maxQueries)

		_, er = tempFile.WriteAt(tempQueries.BytesEncode(), 0)
		if er != nil {
			return fmt.Errorf("failed to write temp file at: %w\n", er)
		}

		tempQueries = make(entity.GroupQueryFrequency)
		tempQueries[parsedString.Query] = parsedString.Frequency
	}

	if len(tempQueries) > 0 {
		tempQueries.Add(queriesFreq, maxQueries)

		_, er := tempFile.WriteAt(tempQueries.BytesEncode(), 0)
		if er != nil {
			return fmt.Errorf("failed to write temp file at: %w\n", er)
		}
	}

	return nil
}
