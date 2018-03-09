package client

type Client interface {
	Auth(name, pass string)
	Login() error
	Submit(contest, task, code) error
}
