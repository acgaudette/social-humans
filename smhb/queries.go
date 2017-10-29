package smhb

func (this client) GetUser(handle string) (User, error) {
	buffer, err := this.query(USER, handle, nil, nil)

	if err != nil {
		return nil, err
	}

	loaded, err := deserializeUserInfo(handle, buffer)

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

	loaded, err := deserializePool(handle, buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

func (this client) GetPost(address string, token Token) (Post, error) {
	buffer, err := this.query(POST, address, nil, &token)

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
