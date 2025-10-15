package pagination

import (
	"common/pkg/validation"
)

type PageParams struct {
	Page     int `form:"page" binding:"omitempty,min=1" label:"页码"`
	PageSize int `form:"page_size" binding:"omitempty,min=1" label:"每页大小"`
}

func (p *PageParams) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
}

var _ validation.Defaultable = (*PageParams)(nil)
