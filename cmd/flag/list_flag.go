package flag

import (
	"fmt"
)

type ListFlag []string

func (l *ListFlag) Set(i string) error {
	*l = append(*l, i)
	return nil
}

func (l *ListFlag) String() string {
	return fmt.Sprintf("%v", *l)
}
