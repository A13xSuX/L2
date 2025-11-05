package models

type Event struct {
	UserId      int    `json:"user_id"` //теги показывают как поля структуры должны называться в json
	Date        string `json:"date"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}
