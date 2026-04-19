package auth_request

type LoginDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
