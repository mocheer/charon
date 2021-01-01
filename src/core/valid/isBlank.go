package valid

//IsBlank check the string is blank string
func IsBlank(str string) bool {
	for _, r := range str {
		switch r {
		case 9:
		case 10:
		case 13:
		case 32:
			break
		default:
			return false
		}
	}
	return true
}
