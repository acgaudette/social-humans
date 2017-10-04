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
func MakePoolView(
	handle string, status string,
) (*front.PoolView, *front.StatusView, error) {
	pool, err := data.LoadPool(handle)

	view := &front.PoolView{
		Handles: []string{},
	}

	if err != nil {
		// Return empty pool view if pool is not found
		return view, MakeStatusView("Error: access failure"), err
	}

	if len(pool.Users) <= 1 {
		// Override the empty pool message with the input status message
		if status == "" {
			status = "Your pool is empty!"
		}

		// Return empty pool view
		return view, MakeStatusView(status), nil
	}

	// Build handles slice from pool users
	for _, value := range pool.Users {
		if value == handle {
			continue
		}

		view.Handles = append(view.Handles, value)
	}

	return view, MakeStatusView(status), nil
}
