package utils

func Has[T interface{}](data []T, condition func(item *T) bool) bool {
	for _, item := range data {
		if condition(&item) {
			return true
		}
	}

	return false
}
