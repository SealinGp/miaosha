package model

type UserDetails struct {
	UserId int64
	Username string
	Password string
	Authorities []string
}

func (userDetails *UserDetails)IsMatch(username, password string) bool {
	return userDetails.Password == password && userDetails.Username == username
}