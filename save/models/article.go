package save

type Article struct {
	Id      uint64
	Url     string
	Title   string
	Content string
	Image   string
}

func NewArticle(id uint64, url, title, content, image string) *Article {
	return &Article{
		Id:      id,
		Url:     url,
		Title:   title,
		Content: content,
		Image:   image,
	}
}
