package main

type ShortenRequest struct {
	LongUrl string `json:"longUrl"`
}

type ShortenResponse struct {
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

type RedirectResponse struct {
	LongUrl string `json:"longUrl"`
}
