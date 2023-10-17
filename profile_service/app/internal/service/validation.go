package service

func stringContainsChar(str []byte, toFind byte) bool {
	for _, char := range str {
		if char == toFind {
			return true
		}
	}
	return false
}
