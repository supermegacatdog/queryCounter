package main

import (
	"os"
	"testing"

	"github.com/supermegacatdog/queryCounter/domain/readerProcessing"
	"github.com/supermegacatdog/queryCounter/domain/scannerProcessing"
	"github.com/supermegacatdog/queryCounter/presenters"
)

const N = 10000

const inputFileName = "input.txt"
const outputFileName = "output.txt"

func BenchmarkScan(b *testing.B) {
	presenter := presenters.NewStringCounterPresenter(readerProcessing.NewReaderProcessingService())

	for i := 0; i < b.N; i++ {
		presenter.CountByScan(inputFileName, outputFileName, N)
	}

	os.Remove(outputFileName)
}

func BenchmarkRead(b *testing.B) {
	presenter := presenters.NewStringCounterPresenter(scannerProcessing.NewScannerProcessingService())

	for i := 0; i < b.N; i++ {
		presenter.CountByRead(inputFileName, outputFileName, N)
	}

	os.Remove(outputFileName)
}
