package delta

import (
	//"log"
	"math"
	"reflect"
)

// Delta is the main type representing a QuillJs delta
type Delta struct {
	Ops []Op
}

// Op is the smallest "operation"
// TODO: Handle embeds
type Op struct {
	Insert     *string
	Retain     *int
	Attributes map[string]interface{}
	Delete     *int
}

// IsNil tells you if the current Op is a nil operation
func (o *Op) IsNil() bool {
	return o.Attributes == nil &&
		o.Delete == nil &&
		o.Insert == nil &&
		o.Retain == nil
}

// New creates a new Delta with the given ops
func New(ops []Op) *Delta {
	return &Delta{
		Ops: ops,
	}
}

// Insert takes a string and a map of attributes and adds them to the Delta d
// If the string is empty, we return the original delta
func (d *Delta) Insert(text string, attrs map[string]interface{}) *Delta {
	if len(text) == 0 {
		return d
	}
	newOp := Op{
		Insert: &text,
	}

	if attrs != nil {
		newOp.Attributes = attrs
	}
	d.Push(newOp)
	return d
}

// Delete deletes `n` characters from the deltal d`
func (d *Delta) Delete(n int) *Delta {
	if n <= 0 {
		return d
	}
	d.Push(Op{Delete: &n})

	return d
}

// Retain keeps n characters and applies the attrs if present
func (d *Delta) Retain(n int, attrs map[string]interface{}) *Delta {
	if n <= 0 {
		return d
	}
	newOp := Op{
		Retain: &n,
	}
	if attrs != nil {
		newOp.Attributes = attrs
	}
	d.Push(newOp)

	return d
}

// Push adds the newOp Operation to the delta, but reorganizes the ops based on certain rules
func (d *Delta) Push(newOp Op) *Delta {
	idx := len(d.Ops)
	var lastOp *Op
	if idx > 0 {
		lastOp = &d.Ops[idx-1]
	}
	// if we don't have any ops, add the newOp as is
	if lastOp == nil {
		d.Ops = append(d.Ops, newOp)
		return d
	}
	if !newOp.IsNil() {

		if newOp.Delete != nil && lastOp.Delete != nil {
			sum := *lastOp.Delete + *newOp.Delete
			d.Ops[idx-1] = Op{Delete: &sum}

			return d
		}

		// Since it does not matter if we insert before or after deleting at the same index,
		// always prefer to insert first
		if lastOp.Delete != nil && newOp.Insert != nil {
			idx--
			if idx < 1 {
				d.Ops = append([]Op{newOp}, d.Ops...)
				return d
			}
			lastOp = &d.Ops[idx-1]
		}
		if reflect.DeepEqual(newOp.Attributes, lastOp.Attributes) {
			if newOp.Insert != nil && lastOp.Insert != nil {
				mergedText := *lastOp.Insert + *newOp.Insert
				d.Ops[idx-1] = Op{
					Insert: &mergedText,
				}
				if newOp.Attributes != nil {
					d.Ops[idx-1].Attributes = newOp.Attributes
				}
				return d
			} else if newOp.Retain != nil && lastOp.Retain != nil {
				x := *lastOp.Retain + *newOp.Retain
				d.Ops[idx-1] = Op{
					Retain: &x,
				}
				if newOp.Attributes != nil {
					d.Ops[idx-1].Attributes = newOp.Attributes
				}
				return d
			}
		}
	}
	if idx == len(d.Ops) {
		d.Ops = append(d.Ops, newOp)
	} else {
		// expand the slice, don't worry about the last element, it will be overriden on the next line
		d.Ops = append(d.Ops, Op{})
		// copy the min of either src or dst elements from the secon part to the first parameter
		// see https://www.reddit.com/r/golang/comments/3z91f1/can_someone_explain_inserting_into_a_slice_with/
		// for mroe info, specialy the playground example
		copy(d.Ops[idx+1:], d.Ops[idx:])
		// finally, insert our desired element at idx position
		d.Ops[idx] = newOp
	}
	return d
}

// Chop removes the last retain operation if it doesn't have any attributes
func (d *Delta) Chop() *Delta {
	x := len(d.Ops)
	if x <= 1 {
		return d
	}
	lastOp := d.Ops[x-1]
	if lastOp.Retain != nil && lastOp.Attributes == nil {
		d.Ops = d.Ops[:x-1]
	}
	return d
}

// Compose returns a Delta that is equivalent to applying the operations of own Delta, followed by another Delta.
func (d *Delta) Compose(other Delta) *Delta {
	thisIter := OpsIterator(d.Ops)
	otherIter := OpsIterator(other.Ops)
	//log.Printf("ss %+v\n", *other.Ops[0].Retain)
	delta := New(nil)
	for thisIter.HasNext() || otherIter.HasNext() {
		if otherIter.PeekType() == "insert" {
			delta.Push(otherIter.Next(math.MaxInt64))
		} else if thisIter.PeekType() == "delete" {
			delta.Push(thisIter.Next(math.MaxInt64))
		} else {
			length := int(math.Min(float64(thisIter.PeekLength()), float64(otherIter.PeekLength())))
			thisOp := thisIter.Next(length)
			otherOp := otherIter.Next(length)
			if otherOp.Retain != nil {
				newOp := Op{}
				if thisOp.Retain != nil {
					newOp.Retain = &length
				} else {
					newOp.Insert = thisOp.Insert
				}
				// Preserve null when composing with a retain, otherwise remove it for inserts
				attributes := AttrCompose(thisOp.Attributes, otherOp.Attributes, thisOp.Retain != nil)
				if attributes != nil {
					newOp.Attributes = attributes
				}
				delta.Push(newOp)
				// Other op should be delete, we could be an insert or retain
				// Insert + delete cancels out
			} else if otherOp.Delete != nil && thisOp.Retain != nil {
				delta.Push(otherOp)
			}
		}
	}
	return delta.Chop()
}

// Concat concatenates two Deltas
func (d *Delta) Concat(other Delta) *Delta {
	delta := New(d.Ops)
	if len(other.Ops) > 0 {
		delta.Push(other.Ops[0])
		delta.Ops = append(delta.Ops, other.Ops[1:]...)
	}
	return delta
}

// TransformPosition returns the new index after applying a list of Ops
func (d *Delta) TransformPosition(index int, priority bool) int {
	thisIter := OpsIterator(d.Ops)
	offset := 0
	for thisIter.HasNext() && offset <= index {
		length := thisIter.PeekLength()
		nextType := thisIter.PeekType()
		thisIter.Next(math.MaxInt64)
		if nextType == "delete" {
			index -= int(math.Min(float64(length), float64(index-offset)))
			continue
		} else if nextType == "insert" && (offset < index || !priority) {
			index += length
		}
		offset += length
	}
	return index
}

// Transform given Delta against own operations
func (d *Delta) Transform(other Delta, priority bool) *Delta {
	thisIter := OpsIterator(d.Ops)
	otherIter := OpsIterator(other.Ops)
	delta := New(nil)
	for thisIter.HasNext() || otherIter.HasNext() {
		if thisIter.PeekType() == "insert" && (priority || otherIter.PeekType() != "insert") {
			delta.Retain(OpsLength(thisIter.Next(math.MaxInt64)), nil)
		} else if otherIter.PeekType() == "insert" {
			delta.Push(otherIter.Next(math.MaxInt64))
		} else {
			length := int(math.Min(float64(thisIter.PeekLength()), float64(otherIter.PeekLength())))
			thisOp := thisIter.Next(length)
			otherOp := otherIter.Next(length)
			if thisOp.Delete != nil {
				// Our delete either makes their delete redundant or removes their retain
				continue
			} else if otherOp.Delete != nil {
				delta.Push(otherOp)
			} else {
				// We retain either their retain or insert
				delta.Retain(length, AttrTransform(thisOp.Attributes, otherOp.Attributes, priority))
			}
		}
	}

	return delta.Chop()
}
