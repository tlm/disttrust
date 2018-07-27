package healthz

type Checker interface {
	Check() error
	Name() string
}
