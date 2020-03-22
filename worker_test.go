package hive

import (
	"log"
	"testing"
)

func TestHiveJobWithPool(t *testing.T) {
	h := New()

	doGeneric := h.Handle("generic", generic{}, PoolSize(3))

	grp := NewGroup()
	grp.Add(doGeneric("first"))
	grp.Add(doGeneric("first"))
	grp.Add(doGeneric("first"))

	if err := grp.Wait(); err != nil {
		log.Fatal(err)
	}
}
