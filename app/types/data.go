package types

type DataType int

const (
	_ DataType = iota
	String
	Int
	Array
)

type Data struct {
	Value any
	T     DataType
	Sz    int
}

func FastCommand(vals ...any) *Data {
	wrapper := Data{}
	wrapper.T = Array

	args := vals[0].([]any)
	arr := make([]*Data, len(args))

	for idx, val := range args {
		value := Data{}

		switch val := val.(type) {
		case string:
			value.Value = val
			value.T = String
		case int:
			value.Value = val
			value.T = Int
		default:
			return nil
		}

		arr[idx] = &value
	}

	wrapper.Value = arr

	return &wrapper
}
