package smhb

func (this client) GetUser(handle string) (User, error) {
	buffer, err := this.query(USER, handle)

	if err != nil {
		return nil, err
	}

	user := &user{handle: handle}
	err = user.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	return user, nil
}
