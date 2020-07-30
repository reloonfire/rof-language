package helpers

func ContainsValue(m map[string]interface{}, value interface{}) bool {
	for _, v := range m {
		if value == v {
			return true
		}
	}
	return false
}

func ContainsKey(m map[string]interface{}, key string) bool {
	for k, _ := range m {
		if key == k {
			return true
		}
	}
	return false
}
