package conductor

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type Conductor struct {
	healthErr error
	members   []*MemberStatus
	memCount  int32
	stopCh    chan struct{}
	watchCh   chan *MemberStatus
	waitGroup sync.WaitGroup
}

func NewConductor() *Conductor {
	return &Conductor{
		healthErr: nil,
		stopCh:    make(chan struct{}),
		watchCh:   make(chan *MemberStatus),
		waitGroup: sync.WaitGroup{},
	}
}

func (c *Conductor) AddMember(member Member) *MemberStatus {
	mstatus := NewMemberStatus(member)
	c.members = append(c.members, mstatus)
	return mstatus
}

func (c *Conductor) AddMembers(members ...Member) []*MemberStatus {
	statuses := make([]*MemberStatus, len(members))
	for i, member := range members {
		statuses[i] = c.AddMember(member)
	}
	return statuses
}

func (c *Conductor) Play() *Conductor {
	for _, fmstatus := range c.members {
		if running, err := fmstatus.State(); running || err != nil {
			continue
		}
		fmstatus.setState(true, nil)
		gmstatus := fmstatus

		c.waitGroup.Add(1)
		go func() {
			log := logrus.WithFields(logrus.Fields{
				"member": gmstatus.member.Name(),
			})
			log.Info("playing member")
			go gmstatus.member.Play()

		Outer:
			for {
				select {
				case err := <-gmstatus.member.DoneCh():
					log.Info("member stopped")
					gmstatus.setState(false, err)
					c.watchCh <- gmstatus
					c.waitGroup.Done()
					break Outer
				case <-c.stopCh:
					gmstatus.member.Stop()
				}
			}
		}()
	}
	go c.Watch()
	return c
}

func (c *Conductor) Stop() {
	close(c.stopCh)
	c.waitGroup.Wait()
}

func (c *Conductor) Watch() {
	for {
		select {
		case mstatus := <-c.watchCh:
			if _, err := mstatus.State(); err != nil {
				log := logrus.WithFields(logrus.Fields{
					"member": mstatus.member.Name(),
				})
				log.Errorf("member failed: %v", err)
			}
		}
	}
}
