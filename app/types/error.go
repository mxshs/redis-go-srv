package types

type RedisError string

func (re RedisError) Error() string {
	return string(re)
}

const (
	Expired      RedisError = "requested value has expired"
	NotFound     RedisError = "requested value was not present or got deleted due to expiration"
	InvalidKey   RedisError = "unexpected type for key"
	InvalidType  RedisError = "unexpected internal type"
	EmptyMessage RedisError = "empty message or uncaught error during parsing"
)
