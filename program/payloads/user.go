package payloads

type User struct {
	UserId    int64      `json:"user_id"`
	Username  string     `json:"username"`
	Posts     []Post     `json:"posts"`
	Following []UserLike `json:"following"`
	Followers []UserLike `json:"followers"`
}
