package services

type HttpHandler interface {
	Setup() error
	Wait() error
}
