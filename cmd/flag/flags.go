package flag

import (
	"flag"
)

var (
	ConfigFiles = ListFlag{}
)

func init() {
	flag.Var(&ConfigFiles, "c", "config files")
}
