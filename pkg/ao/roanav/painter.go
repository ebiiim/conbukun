package roanav

type Painter interface {
	Paint(n Navigation) (path string, err error)
}
