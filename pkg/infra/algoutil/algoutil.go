package algoutil

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must1[T any](v T, err error) T {
	Must(err)
	return v
}

func Must2[T, T2 any](v T, v2 T2, err error) (T, T2) {
	Must(err)
	return v, v2
}
