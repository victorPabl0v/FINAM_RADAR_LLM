package handlers

type GetNewsReq struct {
	Site  string `json:"site"`
	Limit int    `json:"limit"`
}
