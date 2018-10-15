package healthz

type Checker interface {
	Check() error
	Name() string
}

type ChecksFetcher func() []Checker
