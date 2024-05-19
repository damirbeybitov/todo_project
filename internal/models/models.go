package models

type Config struct {
	SqlConnection string `json:"sqlConnection"`
}

type Task struct {
	ID          int64
	Title       string
	Description string
	Status      string
	UserID      int64
}
