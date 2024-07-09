package common

type Pagging struct {
	Total int64 `json:"total" form:"-"`
	Limit int   `json:"limit" form:"limit"`
	Page  int   `json:"page" form:"page"`
}

func (p *Pagging) Process() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
}
