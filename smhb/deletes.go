package smhb

func (this client) DeleteUser(handle string) error {
	return this.delete(USER, handle)
}
