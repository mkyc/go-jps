package jps

func maxInt(v ...int) int {
	m := v[0]
	for _, i := range v {
		if i > m {
			m = i
		}
	}
	return m
}

func absInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func signInt(i int) int {
	switch {
	case i > 0:
		return 1
	case i < 0:
		return -1
	}
	return 0
}
