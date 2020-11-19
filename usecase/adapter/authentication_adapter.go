package adapter

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticationAdapter interface {
	DecodeLoginInput([]byte) (LoginInput, error)
}
