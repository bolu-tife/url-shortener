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

type Url struct {
	ID       int    `json:"id"`
	ShortUrl string `json:"shortUrl"`
	LongUrl  string `json:"longURl"`
}

func NewUrl(shortUrl, longUrl string) (*Url, error) {
	return &Url{
		ShortUrl: shortUrl,
		LongUrl:  longUrl,
	}, nil
}

type Config struct {
	Port             string `default:"3000"`
	PostgresUser     string
	PostgresDb       string
	PostgresPassword string
	PostgresSslMode  string
	RedisUrl         string
}
