package main

import (
	"database/sql"
	"fmt"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateUrlTable() error
	CreateUrl(*Url) error
	GetUrlByLongUrl(string) (*Url, error)
	GetUrlByShortUrl(string) (*Url, error)
	GetUrls(int, int) ([]*Url, error)

}

type ShortenerStorage struct {
	db    PostgresStore
	cache CacheStore
}

type PostgresStore struct {
	db *sql.DB
}

type CacheStore struct {
	cache *redis.Client
}

func (s *PostgresStore) Init() error {
	return s.CreateUrlTable()
}

func NewRedisStore() (*redis.Client, error) {
	redisUrl :=  GetConfig().RedisUrl
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)

	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		return nil, err
	}

	return client, err

}

// func (s *CacheStore) GetUrlByShortUrl(url string) (*Url, error) {
// 	ctx := context.Background()
// 	cacheKey := "url:" + url
// 	shortURL, err := s.cache.Get(ctx, cacheKey).Result()
// 	if err == nil {
// 		return shortURL, nil
// 	}
// }

func NewPostgresStore() (*PostgresStore, error) {
	username, dbName, password, sslMode := GetConfig().PostgresUser, GetConfig().PostgresDb, GetConfig().PostgresPassword, GetConfig().PostgresSslMode

	conn := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", username, dbName, password, sslMode)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) CreateUrlTable() error {
	
	query := `create table if not exists url(
		id serial primary key,
		short_url varchar(50) unique,
		long_url varchar(255)
		)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) GetUrlByID(id int) (*Url, error) {
	query := fmt.Sprintf("select * from url where id = %d", id)
	return s.GetUrlByQuery(query)
}

func (s *PostgresStore) GetUrlByShortUrl(shortUrl string) (*Url, error) {
	query := "select * from url where short_url = $1"
	return s.GetUrlByQuery(query, shortUrl)
}

func (s *PostgresStore) GetUrlByLongUrl(longUrl string) (*Url, error) {
	query := "select * from url where long_url = $1"
	return s.GetUrlByQuery(query, longUrl)
}

func (s *PostgresStore) GetUrlByQuery(query string, args ...any) (*Url, error) {

	url := new(Url)
	err := s.db.QueryRow(query, args...).Scan(&url.ID, &url.ShortUrl, &url.LongUrl)
	
	if err == sql.ErrNoRows{

		return nil, fmt.Errorf("url not found")
	}
	if err != nil{
		return nil, err
	}

	return url,nil 
	 
}

func (s *PostgresStore) GetUrls(skip, limit int) ([]*Url, error) {
	query := "select * from url limit $1 offset $2"
	rows, err := s.db.Query(query, limit, skip)

	if err != nil {
		return nil, err
	}
	urls := []*Url{}

	for rows.Next() {
		url, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}
	return urls, nil
}

func (s *PostgresStore) CreateUrl(url *Url) error {
	query := `insert into url
	(short_url, long_url)
	values ($1, $2)`
	resp, err := s.db.Query(
		query,
		url.ShortUrl,
		url.LongUrl,
	)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

func scanIntoAccount(rows *sql.Rows) (*Url, error) {
	url := new(Url)
	err := rows.Scan(
		&url.ID,
		&url.ShortUrl,
		&url.LongUrl,
	)

	return url, err
}
