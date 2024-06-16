package utils

import (
	"fmt"
	"mxshs/redis-go/app/types"
	"strings"
)

func Encode(data *types.Data) (*string, error) {
	switch data.T {
	case types.Int:
		i := data.Value.(string)
		res := fmt.Sprintf(":%s\r\n", i)

		return &res, nil
	case types.String:
		s := data.Value.(string)
		res := fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)

		return &res, nil
	case types.Array:
		arr := data.Value.([]*types.Data)

		var sb strings.Builder

		sb.Write([]byte(fmt.Sprintf("*%d\r\n", len(arr))))
		for _, val := range arr {
			resp, err := Encode(val)
			if err != nil {
				return nil, err
			}

			sb.Write([]byte(*resp))
		}

		res := sb.String()

		return &res, nil
	default:
		return nil, types.InvalidType
	}
}
