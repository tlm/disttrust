package config

import (
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/file"
)

func DestFileBuilder(path, mode, gid, uid string) (file.File, error) {
	builder := file.New(path)
	if mode != "" {
		conv, err := strconv.ParseUint(mode, 8, 32)
		if err != nil {
			return builder, errors.Wrap(err, "invalid mode uint")
		}
		builder.Mode = os.FileMode(conv)
	}
	builder.Gid = gid
	builder.Uid = uid
	return builder, nil
}
