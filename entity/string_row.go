package entity

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/supermegacatdog/queryCounter/pkg"
)

type StringQueryFrequency struct {
	Query     string
	Frequency int64
}

func (s *StringQueryFrequency) Parse(str string) error {
	if str == "" {
		return pkg.StringEmpty
	}

	parsedString := strings.Split(strings.Trim(str, "\n"), " ")
	if len(parsedString) < 2 {
		return pkg.StringIncorrect
	}

	freq, err := strconv.Atoi(parsedString[1])
	if err != nil {
		return fmt.Errorf("failed to parse frequency: %w\n", err)
	}

	s.Query = parsedString[0]
	s.Frequency = int64(freq)

	return nil
}
