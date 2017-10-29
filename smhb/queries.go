package smhb

func (this client) GetUser(handle string) (User, error) {
	buffer, err := this.query(USER, handle)

	if err != nil {
		return nil, err
	}

	loaded, err := deserializeUserInfo(handle, buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

func (this client) GetPool(handle string) (Pool, error) {
	buffer, err := this.query(POOL, handle)

	if err != nil {
		return nil, err
	}

	loaded, err := deserializePool(handle, buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

func (this client) GetPost(address string) (Post, error) {
	buffer, err := this.query(POST, address)

	if err != nil {
		return nil, err
	}

	loaded, err := deserializePost(address, buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

func (this client) GetPostAddresses(author string) ([]string, error) {
	buffer, err := this.query(POST_ADDRESSES, author)

	if err != nil {
		return nil, err
	}

	loaded, err := deserializePostAddresses(buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}

func (this client) GetFeed(handle string) (Feed, error) {
	buffer, err := this.query(FEED, handle)

	if err != nil {
		return nil, err
	}

	loaded, err := deserializeFeed(buffer)

	if err != nil {
		return nil, err
	}

	return loaded, nil
}
