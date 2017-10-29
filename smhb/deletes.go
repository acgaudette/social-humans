package smhb

func (this client) DeleteUser(handle string, token Token) error {
	return this.delete(USER, handle, &token)
}

func (this client) DeletePost(address string, token Token) error {
	return this.delete(POST, address, &token)
}
