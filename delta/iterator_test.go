package delta

import (
	"math"
	"testing"
)

func TestPeekType(t *testing.T) {
	attr := make(map[string]interface{})
	attr["bold"] = true
	delta := New(nil).Insert("Hello", attr).Retain(3, nil).Insert("some link I need to fix", nil).Delete(4)

	iter := NewIterator(delta.Ops)
	if s := iter.PeekType(); s != "insert" {
		t.Errorf("iter.PeekType() expected 'insert' but got '%s'\n", s)
	}
	iter.Next(math.MaxInt64)
	if s := iter.PeekType(); s != "retain" {
		t.Errorf("iter.PeekType() expected 'retain' but got '%s'\n", s)
	}
	iter.Next(math.MaxInt64)
	if s := iter.PeekType(); s != "insert" {
		t.Errorf("iter.PeekType() expected 'insert' but got '%s'\n", s)
	}
	iter.Next(math.MaxInt64)
	if s := iter.PeekType(); s != "delete" {
		t.Errorf("iter.PeekType() expected 'delete' but got '%s'\n", s)
	}
	iter.Next(math.MaxInt64)
	if s := iter.PeekType(); s != "retain" {
		t.Errorf("iter.PeekType() expected 'retain' but got '%s'\n", s)
	}
}

func TestNext(t *testing.T) {
	attr := make(map[string]interface{})
	attr["bold"] = true
	delta := New(nil).Insert("Hello", attr).Retain(3, nil).Insert("some link I need to fix", nil).Delete(4)

	iter := NewIterator(delta.Ops)
	for x := 0; x < len(delta.Ops); x++ {
		n := iter.Next(math.MaxInt64)
		nn := delta.Ops[x]
		if n.Insert != nil && *n.Insert != *nn.Insert {
			t.Errorf("failed to call Next(), '%+v' diff than '%+v'\n", *n.Insert, *nn.Insert)
		}
		if n.Retain != nil && *n.Retain != *nn.Retain {
			t.Errorf("failed to call Next(), '%+v' diff than '%+v'\n", *n.Retain, *nn.Retain)
		}
		if n.Delete != nil && *n.Delete != *nn.Delete {
			t.Errorf("failed to call Next(), '%+v' diff than '%+v'\n", *n.Delete, *nn.Delete)
		}
	}
	n := iter.Next(math.MaxInt64)
	if *n.Retain != math.MaxInt64 {
		t.Error("didn't get MaxInt64, got: ", *n.Retain)
	}
	n = iter.Next(4)
	if *n.Retain != math.MaxInt64 {
		t.Error("didn't get MaxInt64, got: ", *n.Retain)
	}
	n = iter.Next(math.MaxInt64)
	if *n.Retain != math.MaxInt64 {
		t.Error("didn't get MaxInt64, got: ", *n.Retain)
	}
}

func TestNext2(t *testing.T) {
	attr := make(map[string]interface{})
	attr["bold"] = true
	delta := New(nil).Insert("Hello", attr).Retain(3, nil).Insert("some link I need to fix", nil).Delete(4)
	iter := NewIterator(delta.Ops)

	n := iter.Next(2)
	if *n.Insert != "He" {
		t.Error("didn't get 'He', got: ", *n.Insert)
	}
	if n.Attributes["bold"] != true {
		t.Errorf("didn't get correct attr, got: %+v\n", n.Attributes)
	}

	n = iter.Next(10)
	if *n.Insert != "llo" {
		t.Error("didn't get 'llo', got: ", *n.Insert)
	}
	if n.Attributes["bold"] != true {
		t.Errorf("didn't get correct attr, got: %+v\n", n.Attributes)
	}

	n = iter.Next(1)
	if *n.Retain != 1 {
		t.Error("didn't get retain 1, got: ", *n.Retain)
	}

	n = iter.Next(2)
	if *n.Retain != 2 {
		t.Error("didn't get retain 2, got: ", *n.Retain)
	}

}
