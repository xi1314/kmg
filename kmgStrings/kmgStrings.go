package kmgStrings

func IsInSlice(slice []string, s string) bool {
	for _, thisS := range slice {
		if thisS == s {
			return true
		}
	}
	return false
}
