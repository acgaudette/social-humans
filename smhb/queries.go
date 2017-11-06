package smhb

func (this client) GetToken(handle, cleartext string) (*Token, error) {
	buffer, err := this.query(TOKEN, handle, []byte(cleartext), nil)

	if err != nil {
		return nil, err
	}

	token := NewToken(string(buffer), handle)
	return &token, nil
}

func (this client) GetUser(handle string) (User, error) {
	buffer, err := this.query(USER, handle, nil, nil)

	if err != nil {
		return nil, err
	}

	loaded := userInfo{InfoHandle: handle}
	err = loaded.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

func (this client) GetPool(handle string, token Token) (Pool, error) {
	buffer, err := this.query(POOL, handle, nil, &token)

	if err != nil {
		return nil, err
	}

	loaded := &pool{handle: handle}
	err = loaded.UnmarshalBinary(buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

func (this client) GetPost(
	requester, address string, token Token,
) (Post, error) {
	buffer, err := this.query(POST, requester, []byte(address), &token)

	if err != nil {
		return nil, err
	}

	loaded, err := deserializePost(address, buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

func (this client) GetPostAddresses(
	author string, token Token,
) ([]string, error) {
	buffer, err := this.query(POST_ADDRESSES, author, nil, &token)

	if err != nil {
		return nil, err
	}

	loaded, err := deserializePostAddresses(buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

func (this client) GetFeed(handle string, token Token) (Feed, error) {
	buffer, err := this.query(FEED, handle, nil, &token)

	if err != nil {
		return nil, err
	}

	loaded, err := deserializeFeed(buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}
