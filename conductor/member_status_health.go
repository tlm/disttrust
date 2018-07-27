package conductor

func (m *MemberStatus) Check() error {
	_, err := m.State()
	return err
}

func (m *MemberStatus) Name() string {
	return m.member.Name()
}
