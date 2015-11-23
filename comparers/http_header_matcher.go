package comparers

func containsAllHeaders(expected, actual map[string][]string) bool {
	if len(expected) > len(actual) {
		return false
	}

	for key, evals := range expected {
		avals := actual[key]
		for _, eval := range evals {
			if !contains(avals, eval) {
				return false
			}
		}
	}
	return true
}

func contains(src []string, lookup string) bool {
	for i := 0; i < len(src); i++ {
		if src[i] == lookup {
			return true
		}
	}
	return false
}
