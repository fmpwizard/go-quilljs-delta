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

func TestInsertAfterDeleteWithMerge(t *testing.T) {
	n := New(nil)
	n.Insert("a", nil).Delete(1).Insert("b", nil)
	exp := New(nil)
	exp.Insert("ab", nil).Delete(1)

	if len(n.Ops) != 2 {
		t.Errorf("failed to create Delta with delete and insert merge, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Insert != *exp.Ops[0].Insert {
		t.Logf("n.Ops and exp.Ops are not equal.\nn: %+v\nexp: %+v\n", n.Ops, exp.Ops)
		t.Errorf("n.Ops and exp.Ops are not equal.\nn: %+v\n", *n.Ops[0].Insert)
		t.Errorf("n.Ops and exp.Ops are not equal.\nn: %+v\n", *n.Ops[1].Insert)
	}
}

func TestDelete(t *testing.T) {
	n := New(nil)
	n.Delete(0)

	if len(n.Ops) != 0 {
		t.Errorf("failed to create Delta with delete(0), got: %+v\n", n.Ops)
	}
}
func TestDeletePositive(t *testing.T) {
	n := New(nil)
	n.Delete(10)

	if len(n.Ops) != 1 {
		t.Errorf("failed to create Delta with delete(10), got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Delete != 10 {
		t.Errorf("failed to create Delta with delete(10), got: %+v\n", n.Ops)
	}
}

func TestRetain(t *testing.T) {
	n := New(nil)
	n.Retain(0, nil)

	if len(n.Ops) != 0 {
		t.Errorf("failed to create Delta with retain(0), got: %+v\n", n.Ops)
	}
}

func TestRetainPositive(t *testing.T) {
	n := New(nil)
	n.Retain(2, nil)

	if len(n.Ops) != 1 {
		t.Errorf("failed to create Delta with retain(2), got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Retain != 2 {
		t.Errorf("failed to create Delta with retain(2), got: %+v\n", n.Ops)
	}
}
func TestRetainPositiveAndAttrs(t *testing.T) {
	n := New(nil)
	attr := make(map[string]interface{})
	attr["bold"] = true
	n.Retain(2, attr)

	if len(n.Ops) != 1 {
		t.Errorf("failed to create Delta with retain(2, {bold: true}), got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Retain != 2 {
		t.Errorf("failed to create Delta with retain(2, {bold: true}), got: %+v\n", n.Ops)
	}
	if n.Ops[0].Attributes["bold"] != true {
		t.Errorf("failed to create Delta with retain(2, {bold: true}), got: %+v\n", n.Ops)
	}
}

func TestPush(t *testing.T) {
	n := New(nil)
	x := "test"
	n.Push(Op{Insert: &x})
	if len(n.Ops) != 1 {
		t.Errorf("failed to Push({insert: 'test'}) to Delta, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Insert != x {
		t.Errorf("failed to Push to Delta, got: %+v\n", n.Ops)
	}
}
func TestPushMultiDelete(t *testing.T) {
	n := New(nil)
	n.Delete(2)
	x := 1
	n.Push(Op{Delete: &x})
	if len(n.Ops) != 1 {
		t.Errorf("failed to Push({insert: 'test'}) to Delta, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Delete != 3 {
		t.Errorf("failed to Push to Delta, got: %+v\n", *n.Ops[0].Delete)
	}
}

func TestPushMultiInsert(t *testing.T) {
	n := New(nil)
	n.Insert("Diego ", nil)
	x := "Smith"
	n.Push(Op{Insert: &x})
	if len(n.Ops) != 1 {
		t.Errorf("failed to Push({insert: 'test'}) to Delta, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Insert != "Diego Smith" {
		t.Errorf("failed to Push to Delta, got: %+v\n", *n.Ops[0].Insert)
	}
}
func TestPushMultiInsertMathingAttrs(t *testing.T) {
	n := New(nil)
	attr := make(map[string]interface{})
	attr["bold"] = true
	n.Insert("Diego ", attr)
	x := "Smith"
	n.Push(Op{Insert: &x, Attributes: attr})
	if len(n.Ops) != 1 {
		t.Errorf("failed to Push({insert: 'test'}) to Delta, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Insert != "Diego Smith" {
		t.Errorf("failed to Push to Delta, got: %+v\n", *n.Ops[0].Insert)
	}
}

func TestPushMultiInsertNonMathingAttrs(t *testing.T) {
	n := New(nil)
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	n.Insert("Diego ", attr1)
	x := "Smith"
	attr2 := make(map[string]interface{})
	attr2["bold"] = false
	n.Push(Op{Insert: &x, Attributes: attr2})
	if len(n.Ops) != 2 {
		t.Errorf("failed to Push multi insert, diff attributes, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Insert != "Diego " {
		t.Errorf("failed to Push multi insert, diff attributes, got: %+v\n", *n.Ops[0].Insert)
	}
	if *n.Ops[1].Insert != "Smith" {
		t.Errorf("failed to Push multi insert, diff attributes, got: %+v\n", *n.Ops[1].Insert)
	}
}

func TestPushMultiRetainMathingAttrs(t *testing.T) {
	n := New(nil)
	attr := make(map[string]interface{})
	attr["bold"] = true
	n.Retain(5, attr)
	x := 3
	n.Push(Op{Retain: &x, Attributes: attr})
	if len(n.Ops) != 1 {
		t.Errorf("failed to multi retain with matching attr, got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Retain != 8 {
		t.Errorf("failed to multi retain with matching attr, got: %+v\n", *n.Ops[0].Retain)
	}
}

func TestPushMultiRetainNonMathingAttrs(t *testing.T) {
	n := New(nil)
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	n.Retain(5, attr1)
	x := 3
	attr2 := make(map[string]interface{})
	attr2["bold"] = false

	n.Push(Op{Retain: &x, Attributes: attr2})
	if len(n.Ops) != 2 {
		t.Errorf("failed to multi retain without matching attr, expected 2 ops got: %+v\n", n.Ops)
	}
	if *n.Ops[0].Retain != 5 {
		t.Errorf("failed to multi retain without matching attr, expected 5, got: %+v\n", *n.Ops[0].Retain)
	}
	if *n.Ops[1].Retain != 3 {
		t.Errorf("failed to multi retain without matching attr, expected 3, got: %+v\n", *n.Ops[1].Retain)
	}
}
