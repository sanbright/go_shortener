// Package error пакет с ошибками приложения
package error

import "fmt"

// NotUniqShortLinkError - не уникальный ссылка
type NotUniqShortLinkError struct {
	URL string
	Err error
}

// Error - формирование текста ошибки
func (e *NotUniqShortLinkError) Error() string {
	return fmt.Sprintf("Попытка добавить cуществующую ссылку '%s'", e.URL)
}

// NewNotUniqShortLinkError - конструктор ошибки в случае конфликта добавления ссылки
func NewNotUniqShortLinkError(URL string, err error) error {
	return &NotUniqShortLinkError{
		URL: URL,
		Err: err,
	}
}
