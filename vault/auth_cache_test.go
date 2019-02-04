package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"

	"github.com/tlmiller/disttrust/file"
)

func TestFileAuthCacheEmptyReading(t *testing.T) {
	f := file.New("/tmp/disttrust-vaul-noexist.json")
	cache := NewFileAuthCache(f)
	val, err := cache.Read()
	if err != nil {
		t.Errorf("unexpected error when reading non existent file auth cache value: %v", err)
	}
	if val != nil {
		t.Error("recieved auth cache file secret for non existent file")
	}
}

func TestFileAuthCacheReadingAndWriting(t *testing.T) {
	f := file.New("/tmp/disttrust-vault-authcache.json")
	cache := NewFileAuthCache(f)
	err := cache.Write(&api.Secret{
		RequestID: "request-id",
	})
	if err != nil {
		t.Fatalf("unexpected error when writing auth secret value to file cache: %v", err)
	}

	val, err := cache.Read()
	if err != nil {
		t.Fatalf("unexpected error when reading file auth cache value: %v", err)
	}
	if val == nil {
		t.Fatal("recieved empty auth cache secret reading file auth")
	}

	if val.RequestID != "request-id" {
		t.Fatalf("expected request id value to be request-id but got %s",
			val.RequestID)
	}
}
