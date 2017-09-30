package data

func path(handle, extension string) string {
	return DATA_PATH + "/" + handle + "." + extension
}
