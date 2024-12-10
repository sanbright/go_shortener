package stats

// GetResponse - Ответ на запрос статистики
type GetResponse struct {
	// Urls - количество сокращённых URL в сервисе
	Urls int `json:"urls"`
	// Users - количество пользователей в сервисе
	Users int `json:"users"`
}
