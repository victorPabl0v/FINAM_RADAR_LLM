package models

import "time"

type NewsRequest struct {
	News []News `json:"news"`
}

type News struct {
	Headline   string    `json:"headline"`
	Hotness    float64   `json:"hotness"`
	WhyNow     string    `json:"why_now"`
	Entities   []string  `json:"entities"`
	Sources    []string  `json:"sources"`
	Timeline   time.Time `json:"timeline"`
	Draft      draft     `json:"draft"`
	DedupGroup string    `json:"dedup_group"`
}

type draft struct {
	Title   string   `json:"title"`
	Lead    string   `json:"lead"`
	Bullets []string `json:"bullets"`
	Quote   string   `json:"quotes"`
}
