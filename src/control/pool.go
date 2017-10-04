package control

import (
	"../data"
	"../front"
)

/*
	Model to view functions never return nil, so that something is always
	rendered
*/

// Build a PoolView
func MakePoolView(handle string) (*front.PoolView, error) {
	pool, err := data.LoadPool(handle)

	view := &front.PoolView{
		Handles: []string{},
	}

	if err != nil {
		// Return empty pool view if pool is not found
		return view, &AccessError{handle}
	}

	if len(pool.Users) <= 1 {
		// Return empty pool view
		return view, &EmptyPoolError{handle}
	}

	// Build handles slice from pool users
	for _, value := range pool.Users {
		if value == handle {
			continue
		}

		view.Handles = append(view.Handles, value)
	}

	return view, nil
}
