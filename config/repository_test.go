package config

import (
	"fmt"
	"testing"
)

func TestVersion(t *testing.T) {
	v := &Version{
		Index: 1,
	}
	wanted := "v1"
	if got := v.Name(); got != wanted {
		t.Errorf("wanted %s, got %s", wanted, got)
	}

}

func TestSchemaVersions(t *testing.T) {
	r := &Repository{}
	r.Namespace = "test.schema.net"
	err := r.Create()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	v := &Version{
		Index:         1,
		RootNamespace: r.Namespace,
	}
	r.Versions = append(r.Versions, v)
	err = r.CreateManifest()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	versions := r.SchemaVersions()
	if len(versions) != 1 {
		t.Errorf("expected %d, got %d", 1, len(versions))
	}
	v1 := versions[0]
	expected := "v1"
	if v1 != expected {

		t.Errorf("expected %s, got %s", expected, v1)
	}
	v1ns := v.Namespace()
	expected = "test.schema.net/v1"
	if v1ns != expected {

		t.Errorf("expected %s, got %s", expected, v1ns)
	}
}
