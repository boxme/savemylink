package models

type ViewerContext struct {
	Id    uint64
	Token string
}

func NewViewerContext(id uint64, token string) *ViewerContext {
	return &ViewerContext{
		Id:    id,
		Token: token,
	}
}
