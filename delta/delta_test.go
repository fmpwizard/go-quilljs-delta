package delta

import (
	"testing"
)

func TestEmptyDelta(t *testing.T) {
	n := New(nil)
	if n.Ops != nil {
		t.Error("failed to create Delta with nil ops")
	}
}
func TestNoopOps(t *testing.T) {
	n := New(nil)
	n.Insert("", nil).Delete(0).Retain(0, nil)
	if n.Ops != nil {
		t.Error("failed to create Delta with nil ops")
	}
}
func TestInsert1(t *testing.T) {
	n := New(nil)
	n.Insert("test", nil)
	if len(n.Ops) != 1 {
		t.Errorf("failed to create Delta with insert, got: %+v\n", n.Ops)
	}

	if *n.Ops[0].Insert != "test" {
		t.Error("failed to create Delta with test insert, got: ", *n.Ops[0].Insert)
	}
	if n.Ops[0].Attributes != nil {
		t.Error("failed to create Delta with only insert but no att, att was: ", n.Ops[0].Attributes)
	}
}
func TestInsertWithAttr(t *testing.T) {
	n := New(nil)
	attr := make(map[string]interface{})
	attr["bold"] = true
	n.Insert("test", attr)
	if len(n.Ops) != 1 {
		t.Errorf("failed to create Delta with insert, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Insert != "test" {
		t.Error("failed to create Delta with test insert, got: ", *n.Ops[0].Insert)
	}
	if n.Ops[0].Attributes == nil {
		t.Error("failed to create Delta with only insert with att, attr was nil")
	}
	if n.Ops[0].Attributes["bold"] != true {
		t.Errorf("failed to create Delta with only insert with att, attr was: %+v\n", n.Ops[0].Attributes)
	}
}

func TestInsertAfterDelete(t *testing.T) {
	n := New(nil)
	n.Delete(1).Insert("a", nil)
	exp := New(nil)
	exp.Insert("a", nil).Delete(1)

	if len(n.Ops) != 2 {
		t.Errorf("failed to create Delta with delete and insert, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Insert != *exp.Ops[0].Insert {
		t.Errorf("n.Ops and exp.Ops are not equal.\nn: %+v\nexp: %+v\n", n.Ops, exp.Ops)
	}
}
