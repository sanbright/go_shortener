package error

import "fmt"

type NotUniqShortLinkError struct {
	URL string
	Err error
}

func (e *NotUniqShortLinkError) Error() string {
	return fmt.Sprintf("Попытка добавить cуществующую ссылку '%s'", e.URL)
}

func NewNotUniqShortLinkError(URL string, err error) error {
	return &NotUniqShortLinkError{
		URL: URL,
		Err: err,
	}
}
