package user

type Response []*ItemResponse

type ItemResponse struct {
	OriginalUrl string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}
