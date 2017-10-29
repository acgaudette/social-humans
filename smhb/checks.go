package smhb

func (this client) Validate(handle, cleartext string) error {
	return this.check(VALIDATE, handle, []byte(cleartext))
}

func (this client) CheckUser(handle string) error {
	return this.check(USER, handle, nil)
}
