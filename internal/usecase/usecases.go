package usecase

type UseCases struct {
	Shorten Shorten
	Ping    Ping
	User    User
}

func New(shorten Shorten, ping Ping, user User) *UseCases {
	return &UseCases{
		Shorten: shorten,
		Ping:    ping,
		User:    user,
	}
}
