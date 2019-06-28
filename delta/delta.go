// Package delta is a port of quill's Delta package and implements Operational Transformations.
// The idea behind it is to provide a Go backend to the Quill editor
package delta

import (
	"encoding/json"
	"errors"
	"math"
	"reflect"
)

// Delta is the main type representing a QuillJs delta
type Delta struct {
	Ops []Op `json:"ops"`
}

// Embed is the type represending a QuillJs embed object, it should
// have only one key and only one value.
type Embed struct {
	Key   string
	Value interface{}
}

// Op is the smallest "operation"
type Op struct {
	Insert      []rune                 `json:"insert,omitempty"`
	InsertEmbed *Embed                 `json:"-"`
	Retain      *int                   `json:"retain,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
	Delete      *int                   `json:"delete,omitempty"`
}

// IsNil tells you if the current Op is a nil operation
func (o *Op) IsNil() bool {
	return o.Attributes == nil &&
		o.Delete == nil &&
		o.Insert == nil &&
		o.Retain == nil
}

// Length calculates the length of current Op
func (o *Op) Length() int {
	return OpsLength(*o)
}

// New creates a new Delta with the given ops
func New(ops []Op) *Delta {
	return &Delta{
		Ops: ops,
	}
}

// FromJSON takes a list of ops in json format and creates a Delta
func FromJSON(in []byte) (*Delta, error) {
	var ret Delta
	if err := json.Unmarshal(in, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

// UnmarshalJSON let's us unmarshal a string in the `insert` op to a []rune
func (o *Op) UnmarshalJSON(data []byte) error {
	type Alias Op
	aux := &struct {
		Insert *json.RawMessage `json:"insert"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Insert != nil {
		if (*aux.Insert)[0] == '"' {
			var b string
			if err := json.Unmarshal(*aux.Insert, &b); err != nil {
				return err
			}
			o.Insert = []rune(b)
			o.InsertEmbed = nil
		} else {
			var m map[string]interface{}
			if err := json.Unmarshal(*aux.Insert, &m); err != nil {
				return err
			}
			if len(m) != 1 {
				return errors.New("invalid embed")
			}
			var embed Embed
			// There should be only one key in the map, so this operation should
			// work without problems
			for k, v := range m {
				embed.Key = k
				embed.Value = v
			}
			o.InsertEmbed = &embed
			o.Insert = nil
		}
	}
	return nil
}

// MarshalJSON let's us marshal our Insert []rune into a string
func (o *Op) MarshalJSON() ([]byte, error) {
	type Alias Op
	var message *json.RawMessage
	if o.Insert != nil {
		b, err := json.Marshal(string(o.Insert))
		if err != nil {
			return nil, err
		}
		message = (*json.RawMessage)(&b)
	} else if o.InsertEmbed != nil {
		m := make(map[string]interface{})
		m[o.InsertEmbed.Key] = o.InsertEmbed.Value
		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		message = (*json.RawMessage)(&b)
	}
	return json.Marshal(&struct {
		Insert *json.RawMessage `json:"insert,omitempty"`
		*Alias
	}{
		Insert: message,
		Alias:  (*Alias)(o),
	})
}

// Insert takes a string and a map of attributes and adds them to the Delta d
// If the string is empty, we return the original delta
func (d *Delta) Insert(text string, attrs map[string]interface{}) *Delta {
	if len([]rune(text)) == 0 {
		return d
	}
	newOp := Op{
		Insert:     []rune(text),
		Attributes: attrs,
	}
	d.Push(newOp)
	return d
}

// InsertEmbed takes a map of embeds and a map of attributes, adds them to the Delta. This
// can be used to insert images or URLs to the Delta
func (d *Delta) InsertEmbed(embed Embed, attrs map[string]interface{}) *Delta {
	if len(embed.Key) == 0 || embed.Value == nil {
		return d
	}
	newOp := Op{
		InsertEmbed: &embed,
		Attributes:  attrs,
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
				mergedText := append(lastOp.Insert, newOp.Insert...)
				d.Ops[idx-1] = Op{
					Insert: mergedText,
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
		// expand the slice, don't worry about the last element, it will be overridden on the next line
		d.Ops = append(d.Ops, Op{})
		// copy the min of either src or dst elements from the secon part to the first parameter
		// see https://www.reddit.com/r/golang/comments/3z91f1/can_someone_explain_inserting_into_a_slice_with/
		// for more info, specially the playground example
		copy(d.Ops[idx+1:], d.Ops[idx:])
		// finally, insert our desired element at idx position
		d.Ops[idx] = newOp
	}
	return d
}

// Chop removes the last retain operation if it doesn't have any attributes
func (d *Delta) Chop() *Delta {
	x := len(d.Ops)
	if x == 0 {
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
					newOp.Insert = append([]rune(nil), thisOp.Insert...)
					newOp.InsertEmbed = thisOp.InsertEmbed
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

// Length calculates the length of a Delta, whic is the sum of lengths of
// all its operations
func (d *Delta) Length() int {
	length := 0
	for _, op := range d.Ops {
		length += op.Length()
	}
	return length
}

// Slice returns copy of the delta containing the sliced subset of operations
func (d *Delta) Slice(start int, end int) *Delta {
	iter := OpsIterator(d.Ops)
	delta := New(nil)
	index := 0
	for index < end && iter.HasNext() {
		var nextOp Op
		if index < start {
			nextOp = iter.Next(start - index)
		} else {
			nextOp = iter.Next(end - index)
			delta.Push(nextOp)
		}
		index += nextOp.Length()
	}
	return delta
}

// Invert calculates the inverted delta given a base document data, which
// has the opposite effect when applying. In other words,
// base.Compose(delta).Compose(inverted) == base
func (d *Delta) Invert(base *Delta) *Delta {
	inverted := New(nil)
	baseIndex := 0
	for _, op := range d.Ops {
		if op.Insert != nil || op.InsertEmbed != nil {
			inverted.Delete(op.Length())
		} else if op.Retain != nil && op.Attributes == nil {
			inverted.Retain(*op.Retain, nil)
			baseIndex += *op.Retain
		} else if op.Retain != nil && op.Attributes != nil {
			length := *op.Retain
			slice := base.Slice(baseIndex, baseIndex+length)
			for _, baseOp := range slice.Ops {
				inverted.Retain(baseOp.Length(), AttrInvert(op.Attributes, baseOp.Attributes))
			}
			baseIndex += length
		} else if op.Delete != nil {
			length := *op.Delete
			slice := base.Slice(baseIndex, baseIndex+length)
			for _, baseOp := range slice.Ops {
				inverted.Push(baseOp)
			}
			baseIndex += length
		}
	}
	return inverted
}
