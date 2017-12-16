package smhb

func (this *Transaction) GetDir() string {
	return "transactions/"
}

func (this *Transaction) GetPath() string {
	return this.GetDir() + this.timestamp
}

func (this *Transaction) String() string {
	return this.timestamp
}

func (this *Transaction) MarshalBinary() ([]byte, error) {
	return serialize(*this)
}

func (this *Transaction) UnmarshalBinary(buffer []byte) error {
	err := deserialize(this, buffer)
	if err != nil {
		return err
	}
	return nil
}

func logTransaction(transaction *Transaction, access Access, context ServerContext) {
	access.SaveWithDir(transaction, transaction.GetDir(), false, context)
}
