package extract_data

import "github.com/Sakagam1/parallel-golang/extract_data/models"

type DocsDataType interface {
	models.Character
}

type Response[DDType DocsDataType] struct {
	Docs   []DDType `json:"docs"`
	Total  uint32   `json:"total"`
	Limit  uint16   `json:"limit"`
	Offset uint16   `json:"offset"`
	Page   uint16   `json:"page"`
	Pages  uint16   `json:"pages"`
}
