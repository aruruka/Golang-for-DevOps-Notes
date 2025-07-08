package users

import "fmt"

type User struct {
	Sub               string `json:"sub"`
	Name              string `json:"name"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	PreferredUsername string `json:"preferred_username"`
	Email             string `json:"email"`
	Picture           string `json:"picture"`
}

func Auth(login, password, mfa string) (bool, User, error) {
	if login == "raymond" && password == "password" {
		return true, GetAllUsers()[0], nil
	}
	return false, User{}, fmt.Errorf("invalid login or password")
}

func GetAllUsers() []User {
	return []User{
		{
			Sub:               "9-9-9-9",
			Name:              "Raymond Yan",
			GivenName:         "Raymond",
			FamilyName:        "Yan",
			PreferredUsername: "raymond",
			Email:             "keraymondyan69@gmail.com",
		},
	}
}
