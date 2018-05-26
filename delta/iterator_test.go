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
		if n.Insert != nil && len(n.Insert) != len(nn.Insert) {
			t.Errorf("failed to call Next(), '%+v' diff than '%+v'\n", n.Insert, nn.Insert)
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
	if string(n.Insert) != "He" {
		t.Error("didn't get 'He', got: ", string(n.Insert))
	}
	if n.Attributes["bold"] != true {
		t.Errorf("didn't get correct attr, got: %+v\n", n.Attributes)
	}

	n = iter.Next(10)
	if string(n.Insert) != "llo" {
		t.Error("didn't get 'llo', got: ", string(n.Insert))
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
func TestHasNext(t *testing.T) {
	attr := make(map[string]interface{})
	attr["bold"] = true
	delta := New(nil).Insert("Hello", attr).Retain(3, nil).Insert("some link I need to fix", nil).Delete(4)
	iter := NewIterator(delta.Ops)

	if !iter.HasNext() {
		t.Error("didn't get true to HasNext")
	}
}
func TestHasNext2(t *testing.T) {
	iter := NewIterator(nil)

	if iter.HasNext() {
		t.Error("got true to HasNext but expected false")
	}
}

func TestPeekLength(t *testing.T) {
	attr := make(map[string]interface{})
	attr["bold"] = true
	delta := New(nil).Insert("Hello", attr).Retain(3, nil).Insert("some link I need to fix", nil).Delete(4)
	iter := NewIterator(delta.Ops)

	if x := iter.PeekLength(); x != 5 {
		t.Errorf("expected 5 but got '%d'\n", x)
	}

	iter.Next(math.MaxInt64)
	if x := iter.PeekLength(); x != 3 {
		t.Errorf("expected 3 but got '%d'\n", x)
	}
	iter.Next(math.MaxInt64)
	if x := iter.PeekLength(); x != 23 {
		t.Errorf("expected 1 but got '%d'\n", x)
	}
	iter.Next(math.MaxInt64)
	if x := iter.PeekLength(); x != 4 {
		t.Errorf("expected 4 but got '%d'\n", x)
	}
}

func TestPeekLengthPositiveOffset(t *testing.T) {
	attr := make(map[string]interface{})
	attr["bold"] = true
	delta := New(nil).Insert("Hello", attr).Retain(3, nil).Insert("some link I need to fix", nil).Delete(4)
	iter := NewIterator(delta.Ops)

	iter.Next(2)
	if x := iter.PeekLength(); x != 5-2 {
		t.Errorf("expected '5-3' but got '%d'\n", x)
	}
}
func TestPeekLengthNoOpsLeft(t *testing.T) {
	iter := NewIterator(nil)

	if x := iter.PeekLength(); x != math.MaxInt64 {
		t.Errorf("expected 'math.MaxInt64' but got '%d'\n", x)
	}
}
