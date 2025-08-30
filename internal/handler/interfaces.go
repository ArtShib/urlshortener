package handler


type URLService interface {
	Shorten(url string) (string, error)
	GetID(shortCode string) (string, error)
}
