package model

type ContributorData struct {
	Email               string `json:"email"`
	RedisHost           string `json:"redis_host"`
	PostgresTableRecord string `json:"postgres_table_record"`
	PostgresTableStats  string `json:"postgres_table_stats"`
	PostgresUrl         string `json:"postgres_url"`
}
