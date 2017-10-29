package smhb

// Data to send along with a user store
type userStore struct {
	Password string
	Name     string
}

func (this client) AddUser(
	handle, password, name string,
) (User, *Token, error) {
	data, err := serialize(userStore{password, name})

	if err != nil {
		return nil, nil, err
	}

	err = this.store(USER, handle, data, nil)

	if err != nil {
		return nil, nil, err
	}

	token, err := this.GetToken(handle, password)

	if err != nil {
		return nil, nil, err
	}

	loaded, err := this.GetUser(handle)

	if err != nil {
		return nil, nil, err
	}

	return loaded, token, nil
}

// Data to send along with a post store
type postStore struct {
	Content string
	Author  string
}

func (this client) AddPost(title, content, author string, token Token) error {
	data, err := serialize(postStore{content, author})

	if err != nil {
		return err
	}

	return this.store(POST, title, data, &token)
}
