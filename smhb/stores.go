package smhb

type userStore struct {
	Password string
	Name     string
}

func (this client) AddUser(handle, password, name string) (User, error) {
	data, err := serialize(userStore{password, name})

	if err != nil {
		return nil, err
	}

	err = this.store(USER, handle, data)

	if err != nil {
		return nil, err
	}

	return this.GetUser(handle)
}

type postStore struct {
	Content string
	Author  string
}

func (this client) AddPost(title, content, author string) error {
	data, err := serialize(postStore{content, author})

	if err != nil {
		return err
	}

	return this.store(POST, title, data)
}
