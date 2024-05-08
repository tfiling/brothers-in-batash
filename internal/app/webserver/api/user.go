package api

type UserRegistrationReqBody struct {
	Username string `json:"username" validate:"alphanum,ascii,min=4,max=100"`
	Password string `json:"password" validate:"ascii,min=4,max=100"`
}
