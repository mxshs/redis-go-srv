package utils

import "mxshs/redis-go/app/types"

func ValidateBaseValue(value *types.Data, t types.DataType) error {
	if value.T != t {
		return types.InvalidType
	}

	switch t {
	case types.String:
		_, ok := value.Value.(string)
		if !ok {
			return types.InvalidType
		}

		return nil
	case types.Int:
		_, ok := value.Value.(int)
		if !ok {
			return types.InvalidType
		}

		return nil
	default:
		return types.InvalidType
	}
}
