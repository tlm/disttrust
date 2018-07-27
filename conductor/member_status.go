package conductor

import "sync"

type MemberStatus struct {
	err     error
	lock    sync.Mutex
	member  Member
	running bool
}

func NewMemberStatus(member Member) *MemberStatus {
	return &MemberStatus{
		err:     nil,
		lock:    sync.Mutex{},
		member:  member,
		running: false,
	}
}

func (m *MemberStatus) setState(running bool, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.running, m.err = running, err
}

func (m *MemberStatus) State() (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.running, m.err
}
