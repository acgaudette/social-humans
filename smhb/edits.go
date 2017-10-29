package smhb

func (this client) EditPoolAdd(owner, handle string) error {
	return this.edit(POOL_ADD, owner, []byte(handle))
}

func (this client) EditPoolBlock(owner, handle string) error {
	return this.edit(POOL_BLOCK, owner, []byte(handle))
}

// Data to send along with a post edit
type postEdit struct {
	Title   string
	Content string
}

func (this client) EditPost(address, title, content string) error {
	data, err := serialize(postEdit{title, content})

	if err != nil {
		return err
	}

	return this.edit(POST, address, data)
}

func (this client) EditUserName(handle, name string) error {
	return this.edit(USER_NAME, handle, []byte(name))
}

func (this client) EditUserPassword(handle, password string) error {
	return this.edit(USER_PASSWORD, handle, []byte(password))
}
