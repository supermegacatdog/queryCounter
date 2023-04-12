package presenters

type QueryCounterPresenter struct {
	tempFilesWriter TempFilesWriter
}

func NewStringCounterPresenter(tempFilesWriter TempFilesWriter) *QueryCounterPresenter {
	return &QueryCounterPresenter{
		tempFilesWriter: tempFilesWriter,
	}
}
