package shelly

type Node interface {
	isNode()
}

type Command struct {
	Name   string
	Args   []string
	Redirs []Redirect
}

func (*Command) isNode() {}

type Pipe struct {
	Left  Node
	Right Node
}

func (*Pipe) isNode() {}
