package validator

import "fmt"

const (
	errTitleIsEmpty = "Пустое наименование чата"
)

// CheckTitle проверки наименования чата
func CheckTitle(name string) error {
	if len(name) == 0 {
		return fmt.Errorf(errTitleIsEmpty)
	}
	return nil
}
