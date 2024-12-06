package domain

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	PassHash string
}

type User struct {
	Credentials
}
