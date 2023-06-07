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

func FlagOrConfig[T any](flag T, configVal *T) T {
	if configVal != nil {
		flag = *configVal
	}

	return flag
}
