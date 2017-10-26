package smhb

func (this client) DeleteUser(handle string) error {
	return this.delete(USER, handle)
}

func (this client) DeletePost(address string) error {
	return this.delete(POST, address)
}
