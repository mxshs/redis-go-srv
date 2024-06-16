package types

type StandardResponse string

const (
	NULL_BULK_STRING StandardResponse = "$-1\r\n"
	OK               StandardResponse = "+OK\r\n"
)
