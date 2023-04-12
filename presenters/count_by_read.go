package presenters

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/supermegacatdog/queryCounter/entity"
	"github.com/supermegacatdog/queryCounter/pkg"
)

func (s *QueryCounterPresenter) CountByRead(inputFilePath string, outputFilePath string, maxElemNum int) {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		log.Printf("failed to open input file:\n%s", err.Error())
		return
	}

	defer func() {
		er := inputFile.Close()
		if er != nil {
			log.Printf("failed to close input file:\n%s", er.Error())
			return
		}

		return
	}()

	var (
		inReader  = bufio.NewReader(inputFile)
		files     []*os.File
		offset    int64
		inQueries entity.GroupQueryFrequency
	)

	for {
		var n int64
		var er error
		inQueries, n, er = getBlockFromFile(inReader, maxElemNum/2)
		if er != nil {
			if er == io.EOF {
				break
			}

			log.Printf("failed to get block from file:\n%s", er.Error())
			return
		}

		offset += n

		pos, er := inputFile.Seek(0, io.SeekCurrent)
		if err != nil {
			log.Printf("failed to get current position:\n%s", er.Error())
			return
		}

		pos, er = inputFile.Seek(offset-pos, io.SeekCurrent)
		if err != nil {
			log.Printf("failed to change position with offset:\n%s", er.Error())
			return
		}

		_, er = inReader.Discard(inReader.Buffered())
		if err != nil {
			log.Printf("failed to discard buffer:\n%s", er.Error())
			return
		}

		newFiles, er := s.tempFilesWriter.Write(files, inQueries, maxElemNum/2)
		if er != nil {
			log.Printf("failed to write temp files:\n%s", er.Error())
			return
		}

		if len(newFiles) > 0 {
			files = append(files, newFiles...)
		}
	}

	if len(inQueries) > 0 {
		newFiles, er := s.tempFilesWriter.Write(files, inQueries, maxElemNum/2)
		if err != nil {
			log.Printf("failed to write temp files:\n%s", er.Error())
			return
		}

		if len(newFiles) > 0 {
			files = append(files, newFiles...)
		}
	}

	outFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Printf("failed to create output file:\n%s", err.Error())
		return
	}

	defer func() {
		er := outFile.Close()
		if er != nil {
			log.Printf("failed to close out file:\n%s", err.Error())
			return
		}
	}()

	err = pkg.MergeFiles(files, outFile)
	if err != nil {
		log.Printf("failed to merge temp:\n%s", err.Error())
		return
	}
}

func getBlockFromFile(r *bufio.Reader, max int) (entity.GroupQueryFrequency, int64, error) {
	var bytesRead int64

	queries := make(entity.GroupQueryFrequency, 0)

	for len(queries) < max {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				bytesRead += int64(len([]byte(line)))
				queries[strings.Split(line, "\n")[0]]++
				return queries, bytesRead, err
			}

			return nil, 0, fmt.Errorf("failed to read string using reader: %w\n", err)
		}

		bytesRead += int64(len([]byte(line)))
		queries[strings.Split(line, "\n")[0]]++
	}

	return queries, bytesRead, nil
}
