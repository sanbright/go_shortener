// Package user Dto объект для обработки запрос на получение коротких ссылок
package user

// Response - список коротких ссылок
type Response []*ItemResponse

// ItemResponse - элемент ответа
type ItemResponse struct {
	// OriginalURL - оригинальный УРЛ
	OriginalURL string `json:"original_url"`
	// ShortURL - краткий УРЛ
	ShortURL string `json:"short_url"`
}
