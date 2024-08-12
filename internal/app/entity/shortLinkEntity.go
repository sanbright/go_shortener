package entity

import "github.com/google/uuid"

type ShortLinkEntity struct {
	UUID      string `json:"uuid" db:"uuid"`
	ShortLink string `json:"short_link" db:"short_link"`
	URL       string `json:"url" json:"url"`
}

func NewShortLinkEntity(shortLink string, url string) *ShortLinkEntity {
	return &ShortLinkEntity{
		UUID:      uuid.New().String(),
		ShortLink: shortLink,
		URL:       url,
	}
}
