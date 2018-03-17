package delta

import (
	//"log"
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
func New(ops []Op) Delta {
	return Delta{
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
		x := append(d.Ops[:idx], newOp)
		d.Ops = append(x, d.Ops[idx+1:]...)
	}
	return d
}
