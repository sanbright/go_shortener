// Package api Dto объект для обработки входящих запросов
package api

// Request - данные входящего запроса
type Request struct {
	// URL - УРЛ
	URL string `json:"url"`
}

// Response - данные ответа на запрос
type Response struct {
	// Result - результат
	Result string `json:"result"`
}

// ErrorResponse - Дто для форматирования информации об ошибках выполнения запроса
type ErrorResponse struct {
	// Success - успех выполнения запроса
	Success bool `json:"success"`
	// Errors - массив ошибок
	Errors []*CurrentError `json:"errors"`
}

// CurrentError - ошибка по конкретному полю во входящих данных
type CurrentError struct {
	// Path - путь поля
	Path string `json:"path"`
	// Message - текст ошибки
	Message string `json:"message"`
}
