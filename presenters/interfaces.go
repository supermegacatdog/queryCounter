package presenters

import (
	"os"

	"github.com/supermegacatdog/queryCounter/entity"
)

type TempFilesWriter interface {
	Write(existedTempFiles []*os.File, queriesFreq entity.GroupQueryFrequency, maxQueries int) ([]*os.File, error)
}
