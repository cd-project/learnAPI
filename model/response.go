package model

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Success bool        `json:"success"`
}

type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refreshToken"`
	Data         interface{} `json:"data"`
	Message      string      `json:"message"`
	Code         string      `json:"code"`
	Success      bool        `json:"success"`
}
