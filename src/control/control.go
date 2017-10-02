package control

import (
	"../data"
	"../front"
	"net/http"
)

/*
	Model to view functions never return nil, so that something is always
	rendered
*/

// Build a views map from the active user and a generic view
func MakeViews(view interface{}, activeUser *data.User) *front.Views {
	views := make(front.Views)

	// Content (main) view
	if view != nil {
		views["content"] = view
	}

	// Active user account view
	if activeUser != nil {
		views["active"] = MakeActiveView(activeUser.Handle)
	}

	return &views
}

// MakeViews, but automatically load the active user
func GetUserAndMakeViews(view interface{}, in *http.Request) *front.Views {
	account, _ := data.GetUserFromSession(in) // Nil check done in MakeViews
	return MakeViews(view, account)
}

// Build an ActiveView
func MakeActiveView(handle string) *front.ActiveView {
	return &front.ActiveView{
		Handle: handle,
	}
}

// Build a UserView from a user model
func MakeUserView(
	user *data.User, status string, account *data.User,
) *front.UserView {
	handle := user.Handle

	// Always display something to the frontend
	if handle == "" {
		handle = "Username Invalid"
	}

	name := user.Name

	// Always display something to the frontend
	if name == "" {
		name = "Name Invalid"
	}

	// Compare the active user to the input user
	active := false

	if account != nil && account.Handle == user.Handle {
		active = true
	}

	return &front.UserView{
		Handle:       handle,
		Name:         name,
		Status:       status,
		IsActiveUser: active,
	}
}

// Load the active user and build a UserView
func GetUserAndMakeUserView(
	user *data.User, status string, in *http.Request,
) (*front.UserView, *data.User) {
	account, _ := data.GetUserFromSession(in)
	return MakeUserView(user, status, account), account
}

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
