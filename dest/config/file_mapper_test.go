package config

import (
	"testing"

	"github.com/tlmiller/disttrust/dest"
)

func TestFileMapperWithNoData(t *testing.T) {
	conf := &File{}
	mapDest, err := FileMapper(conf)
	if err != nil {
		t.Fatalf("unexpected error calling file mapper with no data")
	}

	aggDest, ok := mapDest.(*dest.Aggregate)
	if !ok {
		t.Fatal("unexpected dest type for file mapper")
	}

	if len(aggDest.Dests) != 0 {
		t.Fatalf("expected file mapper aggregate dest size to be 0, got %d",
			len(aggDest.Dests))
	}
}

func TestFileMapperProperties(t *testing.T) {
	conf := &File{
		CAFile:             "caFile",
		CAFileMode:         "0700",
		CAFileGid:          "group",
		CAFileUid:          "user",
		CertFile:           "certFile",
		CertFileMode:       "0700",
		CertFileGid:        "group",
		CertFileUid:        "user",
		CertBundleFile:     "certBundleFile",
		CertBundleFileMode: "0700",
		CertBundleFileGid:  "group",
		CertBundleFileUid:  "user",
		PrivKeyFile:        "privKeyFile",
		PrivKeyFileMode:    "0700",
		PrivKeyFileGid:     "group",
		PrivKeyFileUid:     "user",
	}

	mapDest, err := FileMapper(conf)
	if err != nil {
		t.Fatalf("unexpected error calling file mapper with full data")
	}

	aggDest, ok := mapDest.(*dest.Aggregate)
	if !ok {
		t.Fatal("unexpected dest type for file mapper")
	}

	if len(aggDest.Dests) != 4 {
		t.Fatalf("expected file mapper aggregate dest size to be 4, got %d",
			len(aggDest.Dests))
	}
}
