package utils

// Type erasure primarily for redis client wrapper
// Initially made to pass []string as variadic any (...any)
// Dunno how to do it better yet
func Erase[T any, U comparable](vals any) any {
    switch vals := vals.(type) {
    case []T:
        res := make([]any, len(vals))

        for idx := range res {
            res[idx] = vals[idx]
        }

        return res
    case map[U]T:
        res := make(map[U]any)

        for key, value := range vals {
            res[key] = value
        }

        return res
    default:
        return vals
    }
}
