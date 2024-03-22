package vo

type CreatArticleRequest struct {
	CategoryId uint   `json:"category_id" binging:"required"`
	Title      string `json:"title" binging:"required"`
	Content    string `json:"content" binging:"required"`
	HeadImage  string `json:"head_image"`
}
