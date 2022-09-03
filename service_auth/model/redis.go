package model

type RedisToken struct {
	RefreshToken string
	DeviceID     string
	Username     string
	ActiveToken  string
	Email        string
}
