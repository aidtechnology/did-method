package resolver

import (
	"log"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("bryk", func(t *testing.T) {
		t.Skip()
		contents, err := Get("did:bryk:7889c965-4644-44ff-b760-f396f1d11444")
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("%s", contents)
	})
	t.Run("stack", func(t *testing.T) {
		contents, err := Get("did:stack:v0:15gxXgJyT5tM5A4Cbx99nwccynHYsBouzr-0")
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("%s", contents)
	})
	t.Run("ccp", func(t *testing.T) {
		contents, err := Get("did:ccp:ceNobbK6Me9F5zwyE3MKY88QZLw")
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("%s", contents)
	})
}
