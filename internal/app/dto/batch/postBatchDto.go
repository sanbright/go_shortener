// Package batch Dto объект для обработки пачки данных по коротким ссылкам
package batch

// Request - массив входящих элементов
type Request []*ItemRequest

// ItemRequest - элемент входящего запроса
type ItemRequest struct {
	// CorrelationID - идентификатор входящих данных
	CorrelationID string `json:"correlation_id"`
	// OriginalURL - УРЛ для сокащения
	OriginalURL string `json:"original_url"`
}

// Response - массив ответов на запрос
type Response []*ItemResponse

// ItemResponse - элемент ответа на запрос
type ItemResponse struct {
	// CorrelationID - идентификатор входящих данных
	CorrelationID string `json:"correlation_id"`
	// ShortURL - короткая ссылка
	ShortURL string `json:"short_url"`
}

// AddBatchDtoList - массив ДТО на добавление новых коротких ссылок
type AddBatchDtoList []*AddBatchDto

// AddBatchDto - ДТО на добавление краткой ссылки
type AddBatchDto struct {
	// CorrelationID - идентификатор входящих данных
	CorrelationID string
	// OriginalURL - оригинальный УРЛ
	OriginalURL string
	// ShortURL - краткий УРЛ
	ShortURL string
	// UserID - уникальный идентификатор пользователя
	UserID string
}
