package terraform

type Ops interface {
	Apply() error
	Init() error
}
