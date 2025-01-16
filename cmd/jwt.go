package cmd

import "gateway/app/jwt"

func InitJwt() (jwt.Service, error) {
	return jwt.NewJwt()
}
