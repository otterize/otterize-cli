package must

func MustRet[T any](item T, err error) T {
	Must(err)
	return item
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
