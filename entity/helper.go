package entity

func Contains(trumps []Trump, trump Trump) bool {
	for _, t := range trumps {
		if t == trump {
			return true
		}
	}
	return false
}
