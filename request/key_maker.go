package request

type Key interface {
	PKCS8() ([]byte, error)
	Raw() interface{}
}

type KeyMaker interface {
	MakeKey() (Key, error)
}
