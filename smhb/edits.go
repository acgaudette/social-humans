package smhb

func (this client) EditPoolAdd(owner, handle string) error {
	return this.edit(POOL_ADD, owner, []byte(handle))
}

func (this client) EditPoolBlock(owner, handle string) error {
	return this.edit(POOL_BLOCK, owner, []byte(handle))
}

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