package smhb

func (this client) GetUser(handle string) (User, error) {
	buffer, err := this.query(USER, handle)

	if err != nil {
		return nil, err
	}

	account, err := deserializeUser(handle, buffer)

	if err != nil {
		return nil, err
	}

	return account, nil
}
