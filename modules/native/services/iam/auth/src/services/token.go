package services

type Policy struct {
	Resources []string
	Actions   []string
}

type TokenData struct {
	Uuid     string
	Identity string
	Policies []Policy
}

func (t *TokenData) ToJWT() string {
	return ""
}
