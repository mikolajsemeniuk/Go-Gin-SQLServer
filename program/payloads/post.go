package payloads

type Post struct {
	PostId int64      `json:"post_id"`
	Title  string     `json:"title"`
	Likes  []PostLike `json:"likes"`
}
