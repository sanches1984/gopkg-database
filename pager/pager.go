package pager

import (
	"context"
	"math"

	"github.com/severgroup-tt/gopkg-database/repository"
	dbpager "github.com/severgroup-tt/gopkg-database/repository/pager"
	"github.com/severgroup-tt/gopkg-errors"
)

// pager internal implementation of facade
type pager struct {
	Page       int32
	PageSize   int32
	TotalItems *int32
}

// ErrLastPage is a tag for last page error
var ErrLastPage = errors.NewTag()

// GetOffset compute data offset
func (p *pager) GetOffset() int32 {
	return (p.Page - 1) * p.PageSize
}

// GetLimit compute data limit
func (p *pager) GetLimit() int32 {
	return p.PageSize
}

// GetPage get current page number
func (p *pager) GetPage() int32 {
	return p.Page
}

// GetPageSize get page size
func (p *pager) GetPageSize() int32 {
	return p.PageSize
}

// GetTotalPages compute total pages for all items
func (p *pager) GetTotalPages() int32 {
	if p.TotalItems == nil || *p.TotalItems < 1 {
		return 0
	}

	return int32(math.Ceil(float64(*p.TotalItems) / float64(p.PageSize)))
}

// GetTotalItems get total items
func (p *pager) GetTotalItems() int32 {
	if p.TotalItems == nil {
		return 0
	}
	return *p.TotalItems
}

// SetTotalItems set total items
func (p *pager) SetTotalItems(total int32) {
	p.TotalItems = &total
}

// NextPage try to step onto next page
func (p *pager) NextPage() error {
	if p.TotalItems == nil {
		return errors.Internal.Err(context.TODO(), "Total items must be set").WithTag(errors.LogicError)
	}

	if p.Page*p.PageSize >= *p.TotalItems {
		return errors.Internal.Err(context.TODO(), "Last page").WithTag(ErrLastPage).WithPayloadKV("page", p.Page)
	}

	p.Page++
	return nil
}

// GetApplyFn convert pager to query apply function
func (p *pager) GetApplyFn() repository.QueryApply {
	return dbpager.New(int(p.GetOffset()), int(p.GetLimit()))
}
