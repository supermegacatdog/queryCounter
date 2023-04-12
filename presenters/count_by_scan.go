package presenters

import (
	"bufio"
	"log"
	"os"

	"github.com/supermegacatdog/queryCounter/entity"
	"github.com/supermegacatdog/queryCounter/pkg"
)

func (s *QueryCounterPresenter) CountByScan(inputFilePath string, outputFilePath string, maxElemNum int) {
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
		inQueries = make(entity.GroupQueryFrequency, 0)
		inScanner = bufio.NewScanner(inputFile)
		files     []*os.File
	)

	for inScanner.Scan() {
		// filling inQueries map using scan
		// map cannot include more than N/2 elements
		if len(inQueries) < maxElemNum/2 {
			inQueries[inScanner.Text()]++
			continue
		}

		// newFiles are files that were created because the data could not fit into already created files
		newFiles, er := s.tempFilesWriter.Write(files, inQueries, maxElemNum/2)
		if err != nil {
			log.Printf("failed to write temp files:\n%s", er.Error())
			return
		}

		files = append(files, newFiles...)
		inQueries[inScanner.Text()]++
	}

	// for cases when the map inQueries was full, but the new line was already scanned
	if len(inQueries) > 0 {
		newFiles, er := s.tempFilesWriter.Write(files, inQueries, maxElemNum/2)
		if er != nil {
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
			log.Printf("failed to close input file:\n%s", er.Error())
			return
		}

		return
	}()

	err = pkg.MergeFiles(files, outFile)
	if err != nil {
		log.Printf("failed to merge temp:\n%s", err.Error())
		return
	}
}
