package delta

import (
	//"log"
	"math"
)

// NewIterator takes a []Op and returns an Iterator instance with them
func NewIterator(ops []Op) Iterator {
	return Iterator{
		Ops:    ops,
		Index:  0,
		Offset: 0,
	}
}

// Iterator holds a list of Op, an Index and Offset
type Iterator struct {
	Ops    []Op
	Index  int
	Offset int
}

// HasNext returns true if we have more ops
func (x *Iterator) HasNext() bool {
	return x.PeekLength() < math.MaxInt64
}

// Next moves the index on item
func (x *Iterator) Next(length int) Op {
	// TODO: see where we need to allow passing infinity as length
	// js code does that if length is missing
	if len(x.Ops) <= x.Index {
		m := math.MaxInt64
		return Op{
			Retain: &m,
		}
	}

	nextOp := x.Ops[x.Index]
	offset := x.Offset
	opLength := OpsLength(nextOp)
	if length >= opLength-offset {
		length = opLength - offset
		x.Index++
		x.Offset = 0
	} else {
		x.Offset += length
	}
	if nextOp.Delete != nil {
		return Op{Delete: &length}
	}
	retOp := Op{}
	if nextOp.Attributes != nil {
		retOp.Attributes = nextOp.Attributes
	}
	if nextOp.Retain != nil {
		retOp.Retain = &length
	}
	if nextOp.Insert != nil {
		// when using Go's slice syntax to extract characters from a string, note that the
		// number after the ":" isn't the number of characters to take, but the position, starting from 0
		// to extract. This is different than the way substr in js is implemented
		// Also, using a :length greater than the actual len(str) panics
		l := length + offset // because of how Go's slice[a:b] work
		if xx := len(nextOp.Insert); xx < l {
			l = xx
		}
		str := (nextOp.Insert)[offset:l]
		retOp.Insert = str
	}
	if nextOp.InsertEmbed != nil {
		retOp.InsertEmbed = nextOp.InsertEmbed
	}
	return retOp
}

// Peek returns the current Op without advancing the index
func (x *Iterator) Peek() Op {
	return x.Ops[x.Index]
}

// PeekLength returns the length of the Op at the current index
func (x *Iterator) PeekLength() int {
	if len(x.Ops) > x.Index {
		return OpsLength(x.Ops[x.Index]) - x.Offset
	}
	return math.MaxInt64
}

// PeekType tells you the type of operation at the current Index
func (x *Iterator) PeekType() string {
	if len(x.Ops) > x.Index { //TODO may be a bug getting the last item (highest index)
		if x.Ops[x.Index].Delete != nil {
			return "delete"
		}
		if x.Ops[x.Index].Retain != nil {
			return "retain"
		}
		if x.Ops[x.Index].Insert != nil || x.Ops[x.Index].InsertEmbed != nil {
			return "insert"
		}
	}
	return "retain"
}
