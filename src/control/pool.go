package control

import (
	"../../smhb"
	"../data"
	"../views"
)

/*
	Model to view functions never return nil so that something is always
	rendered
*/

// Build a Pool view
func MakePoolView(handle string, token smhb.Token) (*views.Pool, error) {
	pool, err := data.Backend.GetPool(handle, token)

	view := &views.Pool{
		Handles: []string{},
	}

	if err != nil {
		// Return empty pool view if pool is not found
		return view, &AccessError{handle}
	}

	if len(pool.Users()) <= 1 {
		// Return empty pool view
		return view, &EmptyPoolError{handle}
	}

	// Build handles slice from pool users
	for _, value := range pool.Users() {
		// Ignore self
		if value == handle {
			continue
		}

		view.Handles = append(view.Handles, value)
	}

	return view, nil
}
