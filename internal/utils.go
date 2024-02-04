package internal

func OptionalParam[T any](defaultValue T, s ...T) T {
	if len(s) == 0 {
		return defaultValue
	} else if len(s) == 1 {
		return s[0]
	} else {
		panic("too many parameters")
	}
}

func OptionalPointer[T any](s ...T) *T {
	if len(s) == 0 {
		return nil
	} else if len(s) == 1 {
		return &s[0]
	} else {
		panic("too many parameters")
	}
}
