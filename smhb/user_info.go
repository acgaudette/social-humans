package smhb

// User data wrapper for transmission
type userInfo struct {
	InfoHandle string
	InfoName   string
}

/* Interface implementation */

func (this userInfo) Handle() string {
	return this.InfoHandle
}

func (this userInfo) Name() string {
	return this.InfoName
}

// Compare two users
func (this userInfo) Equals(other User) bool {
	return this.InfoHandle == other.Handle()
}

func (this userInfo) GetPath() string {
	return this.InfoHandle + ".user"
}

func (this userInfo) String() string {
	return "user info \"" + this.InfoHandle + "\""
}

func (this userInfo) UnmarshalBinary(buffer []byte) error {
	// Deserialize user
	loaded := &user{}
	err := loaded.UnmarshalBinary(buffer)

	if err != nil {
		return err
	}

	// Strip out hash and load into info struct
	info := &userInfo{
		InfoName: loaded.name,
	}

	return nil
}

// Load user info with lookup handle
func getUserInfo(
	handle string, context serverContext, access Access,
) (*userInfo, error) {
	info := &userInfo{InfoHandle: handle}
	err := access.Load(info, context)

	if err != nil {
		return nil, err
	}

	return info, nil
}

func getRawUserInfo(
	handle string, context serverContext, access Access,
) ([]byte, error) {
	info := &userInfo{InfoHandle: handle}
	buffer, err := access.LoadRaw(info, context)

	if err != nil {
		return nil, err
	}

	return buffer, nil
}
