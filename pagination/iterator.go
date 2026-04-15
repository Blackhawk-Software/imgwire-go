package pagination

import "context"

type PageFetcher[T any] func(ctx context.Context, page int, limit int) (Page[T], error)

type PageIterator[T any] struct {
	ctx       context.Context
	fetch     PageFetcher[T]
	nextPage  int
	limit     int
	current   Page[T]
	err       error
	started   bool
	exhausted bool
}

func NewPageIterator[T any](
	ctx context.Context,
	page int,
	limit int,
	fetch PageFetcher[T],
) *PageIterator[T] {
	if ctx == nil {
		ctx = context.Background()
	}
	if page <= 0 {
		page = 1
	}
	return &PageIterator[T]{
		ctx:      ctx,
		fetch:    fetch,
		nextPage: page,
		limit:    limit,
	}
}

func (it *PageIterator[T]) Next() bool {
	if it.err != nil || it.exhausted {
		return false
	}

	page, err := it.fetch(it.ctx, it.nextPage, it.limit)
	if err != nil {
		it.err = err
		return false
	}

	it.current = page
	it.started = true
	it.limit = page.Pagination.Limit
	if it.limit == 0 {
		it.limit = page.Pagination.Limit
	}

	if page.Pagination.NextPage == nil {
		it.exhausted = true
	} else {
		it.nextPage = *page.Pagination.NextPage
	}

	return true
}

func (it *PageIterator[T]) Page() Page[T] {
	return it.current
}

func (it *PageIterator[T]) Err() error {
	return it.err
}

type ItemIterator[T any] struct {
	pages *PageIterator[T]
	index int
	items []T
	err   error
}

func NewItemIterator[T any](pages *PageIterator[T]) *ItemIterator[T] {
	return &ItemIterator[T]{pages: pages}
}

func (it *ItemIterator[T]) Next() bool {
	for {
		if it.index < len(it.items) {
			it.index++
			return true
		}

		if !it.pages.Next() {
			it.err = it.pages.Err()
			return false
		}

		page := it.pages.Page()
		it.items = page.Data
		it.index = 0
	}
}

func (it *ItemIterator[T]) Item() T {
	return it.items[it.index-1]
}

func (it *ItemIterator[T]) Err() error {
	return it.err
}
