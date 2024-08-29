package batch

type Request []*ItemRequest

type ItemRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type Response []*ItemResponse

type ItemResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type AddBatchDtoList []*AddBatchDto

type AddBatchDto struct {
	CorrelationID string
	OriginalURL   string
	ShortURL      string
	UserId        string
}
