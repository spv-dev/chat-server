package validator

import "fmt"

const (
	errBodyIsEmpty = "Тело сообщения пустое"
)

// CheckBody проверки тела сообщения
func CheckBody(body string) error {
	if len(body) == 0 {
		return fmt.Errorf(errBodyIsEmpty)
	}

	return nil
}
