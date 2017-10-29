package smhb

func (this client) EditPoolAdd(owner, handle string, token Token) error {
	return this.edit(POOL_ADD, owner, []byte(handle), &token)
}

func (this client) EditPoolBlock(owner, handle string, token Token) error {
	return this.edit(POOL_BLOCK, owner, []byte(handle), &token)
}

// Data to send along with a post edit
type postEdit struct {
	Title   string
	Content string
}

func (this client) EditPost(
	address, title, content string, token Token,
) error {
	data, err := serialize(postEdit{title, content})

	if err != nil {
		return err
	}

	return this.edit(POST, address, data, &token)
}

func (this client) EditUserName(handle, name string, token Token) error {
	return this.edit(USER_NAME, handle, []byte(name), &token)
}

func (this client) EditUserPassword(
	handle, password string, token Token,
) error {
	return this.edit(USER_PASSWORD, handle, []byte(password), &token)
}
