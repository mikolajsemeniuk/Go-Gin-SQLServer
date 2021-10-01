package payloads

type User struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
}
