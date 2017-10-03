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
func MakePoolView(handle string, status string) (*front.PoolView, error) {
	pool, err := data.LoadPool(handle)

	if err != nil {
		// Return empty pool view if pool is not found
		empty := &front.PoolView{
			Handles: []string{},
			Status:  "Error: access failure",
		}

		return empty, err
	}

	if len(pool.Users) <= 1 {
		// Override the empty pool message with the input status message
		if status == "" {
			status = "Your pool is empty!"
		}

		// Return empty pool view
		empty := &front.PoolView{
			Handles: []string{},
			Status:  status,
		}

		return empty, nil
	}

	result := &front.PoolView{
		Handles: []string{},
		Status:  status,
	}

	// Build handles slice from pool users
	for _, value := range pool.Users {
		if value == handle {
			continue
		}

		result.Handles = append(result.Handles, value)
	}

	return result, nil
}
