package interpreter

func StringToBool(s string) bool {
	if s == "True" {
		return true
	}
	return false
}

func BoolToString(b bool) string {
	if b {
		return "True"
	}
	return "False"
}
