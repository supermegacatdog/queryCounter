package main

import (
	"log"
	"os"
	"strconv"

	"github.com/supermegacatdog/queryCounter/domain/readerProcessing"
	"github.com/supermegacatdog/queryCounter/domain/scannerProcessing"
	"github.com/supermegacatdog/queryCounter/pkg"
	"github.com/supermegacatdog/queryCounter/presenters"
)

func main() {
	if len(os.Args) < 4 {
		log.Print(pkg.NotEnoughArgs)
		return
	}

	n, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Printf("failed to parse n:\n%s", err.Error())
		return
	}

	// default mode is scan because it is faster
	mode := pkg.CounterModeScan
	if len(os.Args) == 5 {
		m, err := strconv.Atoi(os.Args[4])
		if err != nil {
			log.Printf("failed to parse mode:\n%s", err.Error())
			return
		}

		mode = m
	}

	scanner := scannerProcessing.NewScannerProcessingService()
	reader := readerProcessing.NewReaderProcessingService()

	switch mode {
	case pkg.CounterModeScan:
		presenter := presenters.NewStringCounterPresenter(scanner)

		presenter.CountByScan(os.Args[1], os.Args[2], n)
		return

	case pkg.CounterModeRead:
		presenter := presenters.NewStringCounterPresenter(reader)

		presenter.CountByRead(os.Args[1], os.Args[2], n)
		return

	default:
		log.Print(pkg.InvalidMode.Error())
		return
	}
}
