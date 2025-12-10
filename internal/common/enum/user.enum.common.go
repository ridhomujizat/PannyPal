package enum

type UserStatus string
type UserType string

const (
	ACTIVE   UserStatus = "active"
	INCATIVE UserStatus = "inactive"
	BLCOKE   UserStatus = "blocked"
)

const (
	SAAS UserType = "saas"
	LITE UserType = "lite"
)

func (e UserStatus) ToString() string {
	switch e {
	case ACTIVE:
		return "active"
	case INCATIVE:
		return "inactive"
	case BLCOKE:
		return "blocked"
	default:
		return ""
	}
}
func (e UserStatus) IsValid() bool {
	switch e {
	case ACTIVE, INCATIVE, BLCOKE:
		return true
	}

	return false
}
func (e UserType) ToString() string {
	switch e {
	case SAAS:
		return "saas"
	case LITE:
		return "lite"
	default:
		return ""
	}
}
func (e UserType) IsValid() bool {
	switch e {
	case SAAS, LITE:
		return true
	}

	return false
}
