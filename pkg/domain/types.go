package domain

type UserId string

func (userId *UserId) String() string {
	return string(*userId)
}
