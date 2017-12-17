package smhb

import "time"

const (
	/* Work channel buffer size */

	QUEUE_SIZE = 128

	/* Timeout before closing connection (seconds) */

	IO_TIMEOUT = 10

	/* New file permissions */

	PERM = 0600

	/* Post address identifier format */

	ADDRESS_LENGTH = len(time.RFC3339Nano) + 1 // Includes the leading slash

	/* Content character limits */

	HANDLE_LIMIT  = 24
	TITLE_LIMIT   = 20
	CONTENT_LIMIT = 100

	/* Length of the access token */

	TOKEN_SIZE = 32

	/* Size of target string in header (includes null-terminator) byte */

	TARGET_LENGTH = HANDLE_LIMIT + ADDRESS_LENGTH + 1

	/* Length of message header (includes token handle) */

	HEADER_SIZE = 6 + HANDLE_LIMIT + 1 + TOKEN_SIZE + TARGET_LENGTH

	/* Clock synchronization */

	NTP_SERVER  = "pool.ntp.org"
	NTP_TIMEOUT = 4

	/* Replication timeouts */

	RM_TIMEOUT          = 2
	TRANSACTION_TIMEOUT = 3
	LOG_TIMEOUT         = 4
)
