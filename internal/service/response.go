package service

type AddNewsResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

type GetNewsFromSitesResponse struct {
	News []News `json:"news"`
}

type News struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Text  string `json:"text"`
}
