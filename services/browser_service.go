package services

type BrowserService interface {
	Open(url string) error
}
