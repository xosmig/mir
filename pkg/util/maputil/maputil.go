package maputil

func GetKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func GetValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func GetValuesOf[K comparable, V any](m map[K]V, keys []K) []V {
	values := make([]V, 0, len(keys))
	for _, k := range keys {
		values = append(values, m[k])
	}
	return values
}
