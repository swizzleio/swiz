package configutil

func getZero[T comparable]() T {
	var result T

	return result
}

func SetOrDefault[T comparable](val T, def T) T {
	if val == getZero[T]() {
		return def
	}
	return val
}
