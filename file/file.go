package file

import (
	"os"
	"os/user"
	"strconv"

	"github.com/pkg/errors"
)

type File struct {
	Gid  string
	Path string
	Mode os.FileMode
	Uid  string
}

func (f File) Chown() error {
	var gid, uid int
	if f.Gid == "" {
		gid = -1
	} else if conv, err := strconv.ParseInt(f.Gid, 10, 32); err == nil {
		gid = int(conv)
	} else {
		group, err := user.LookupGroup(f.Gid)
		if err != nil {
			return errors.Wrapf(err, "looking up group for id %s", f.Gid)
		}
		conv, err := strconv.ParseInt(group.Gid, 10, 32)
		if err != nil {
			return errors.Wrapf(err, "parsing group lookup gid for id %s", group.Gid)
		}
		gid = int(conv)
	}

	if f.Uid == "" {
		uid = -1
	} else if conv, err := strconv.ParseInt(f.Uid, 10, 32); err == nil {
		uid = int(conv)
	} else {
		user, err := user.Lookup(f.Uid)
		if err != nil {
			return errors.Wrapf(err, "looking up user for id %s", f.Uid)
		}
		conv, err := strconv.ParseInt(user.Uid, 10, 32)
		if err != nil {
			return errors.Wrapf(err, "parsing user lookup gid for id %s", user.Uid)
		}
		uid = int(conv)
	}

	return os.Chown(f.Path, uid, gid)
}

func (f File) HasPath() bool {
	return f.Path != ""
}

func New(path string) File {
	return File{
		Path: path,
		Mode: os.FileMode(0644),
		Gid:  "",
		Uid:  "",
	}
}
