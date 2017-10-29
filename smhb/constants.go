package smhb

const (
	/* Work channel buffer size */

	QUEUE_SIZE = 128

	/* Timeout before closing connection (seconds) */

	IO_TIMEOUT = 10

	/* Post address identifier format */

	TIMESTAMP_LAYOUT = "20060102150405"
	ADDRESS_LENGTH   = len(TIMESTAMP_LAYOUT) + 1 // Includes the leading slash

	/* Length of message header */

	HEADER_SIZE = 64

	/* Size of target string in header */

	TARGET_LENGTH = HEADER_SIZE - 6 - 1 // Includes the null-terminator byte

	/* Length of the access token */

	TOKEN_SIZE = 64

	/* Content character limits */

	HANDLE_LIMIT  = TARGET_LENGTH - ADDRESS_LENGTH
	TITLE_LIMIT   = 20
	CONTENT_LIMIT = 100
)
