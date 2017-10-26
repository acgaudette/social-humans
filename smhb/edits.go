package smhb

func (this client) EditPoolAdd(owner, handle string) error {
	return this.edit(POOL_ADD, owner, []byte(handle))
}

func (this client) EditPoolBlock(owner, handle string) error {
	return this.edit(POOL_BLOCK, owner, []byte(handle))
}
