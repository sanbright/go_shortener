package dto

type Request struct {
	Url string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Success bool            `json:"success"`
	Errors  []*CurrentError `json:"errors"`
}

type CurrentError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}
