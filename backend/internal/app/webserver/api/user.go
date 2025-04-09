package api

type UserRegistrationReqBody struct {
	Username string `json:"username" validate:"ascii,min=4,max=100"`
	Password string `json:"password" validate:"ascii,min=4,max=100"`
}

type UserLoginReqBody struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginRespBody struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenReqBody struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type LogoutReqBody struct {
	Token string `json:"token" validate:"required"`
}
