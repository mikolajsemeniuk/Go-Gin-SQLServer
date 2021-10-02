package payloads

type Post struct {
	PostId int64  `json:"post_id"`
	UserId int64  `json:"-"`
	Title  string `json:"title"`
}
