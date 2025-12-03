package service

type RegisterParams struct {
	Username string
	Email    string
	Password string
}

type LoginParams struct {
	Username string
	Password string
}
