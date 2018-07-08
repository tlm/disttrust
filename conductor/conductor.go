package conductor

type Conductor struct {
	members []Member
}

func NewConductor() *Conductor {
	return &Conductor{}
}

func (c *Conductor) AddMember(mem Member) {
	c.members = append(c.members, mem)
}

func (c *Conductor) Conduct() {
	for _, member := range c.members {
		member.Do()
	}
}
