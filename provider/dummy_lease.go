package provider

import (
	"errors"
	"time"
)

type DummyLease struct {
	IDVal       string
	ResponseVal *Response
	TillVal     func() time.Time
}

func (d *DummyLease) ID() string {
	return d.IDVal
}

func (d *DummyLease) HasResponse() bool {
	return d.ResponseVal != nil
}

func (d *DummyLease) Response() (*Response, error) {
	if d.ResponseVal == nil {
		return nil, errors.New("lease has no response")
	}
	return d.ResponseVal, nil
}

func (d *DummyLease) Start() time.Time {
	return time.Now()
}

func (d *DummyLease) Till() time.Time {
	return d.TillVal()
}
