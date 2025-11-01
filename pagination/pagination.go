package pagination

type Pagination struct {
	Page      int `form:"page" json:"page"`
	PageSize  int `form:"page_size" json:"page_size"`
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
}

func (p *Pagination) SetDefault() Pagination {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = 20
	}
	return *p
}

func (p *Pagination) Limit() int {
	return p.PageSize
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) SetTotal(totalData int) Pagination {
	p.TotalData = totalData
	p.TotalPage = (totalData + p.PageSize - 1) / p.PageSize
	return *p
}
