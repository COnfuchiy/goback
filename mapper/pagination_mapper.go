package mapper

import (
	"goback/api/response"
	"math"
)

type IPaginationMapper interface {
	ToPaginationResponse(totalCount int64, currentPage int) *response.PaginationResponse
}

type PaginationMapper struct {
	entitiesPerPage int64
}

func NewPaginationMapper(entitiesPerPage int64) *PaginationMapper {
	return &PaginationMapper{entitiesPerPage: entitiesPerPage}
}

func (m PaginationMapper) ToPaginationResponse(totalCount int64, currentPage int) *response.PaginationResponse {
	totalPages := int64(math.Ceil(float64(totalCount) / float64(m.entitiesPerPage)))

	pagination := response.PaginationResponse{
		TotalEntities: totalCount,
		TotalPages:    totalPages,
		CurrentPage:   int64(currentPage),
		NextPage:      nil,
		PrevPage:      nil,
	}

	if pagination.CurrentPage+1 > totalPages {
		pagination.NextPage = nil
	} else {
		nextPage := pagination.CurrentPage + 1
		pagination.NextPage = &nextPage
	}

	if pagination.CurrentPage-1 < 1 {
		pagination.PrevPage = nil
	} else {
		prevPage := pagination.CurrentPage - 1
		pagination.PrevPage = &prevPage
	}

	return &pagination

}
