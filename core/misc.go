package core

func StringToBool(s string) bool {
	if s == "true" {
		return true
	}
	return false
}

func BoolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
