package repository

type IPaginateRepository interface {
	GetOffsetAndLimitFromPage(page int) (int, int)
}

type PaginateRepository struct {
	entitiesPerPage int64
}

func NewPaginateRepository(entitiesPerPage int64) *PaginateRepository {
	return &PaginateRepository{entitiesPerPage: entitiesPerPage}
}

func (r PaginateRepository) GetOffsetAndLimitFromPage(page int) (int, int) {
	if page == 0 {
		return -1, -1
	}
	return (page - 1) * int(r.entitiesPerPage), int(r.entitiesPerPage)
}
