package delta

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
	lastOp := Op{}
	if idx > 0 {
		lastOp = d.Ops[idx-1]
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
			if idx == 0 {
				d.Ops = append([]Op{newOp}, d.Ops...)
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

//  var index = this.ops.length;
//  var lastOp = this.ops[index - 1];
//  newOp = extend(true, {}, newOp);
//  if (typeof lastOp === 'object') {
//    if (typeof newOp['delete'] === 'number' && typeof lastOp['delete'] === 'number') {
//      this.ops[index - 1] = { 'delete': lastOp['delete'] + newOp['delete'] };
//      return this;
//    }
//    // Since it does not matter if we insert before or after deleting at the same index,
//    // always prefer to insert first
//    if (typeof lastOp['delete'] === 'number' && newOp.insert != null) {
//      index -= 1;
//      lastOp = this.ops[index - 1];
//      if (typeof lastOp !== 'object') {
//        this.ops.unshift(newOp);
//        return this;
//      }
//    }
//   // skip for now if (equal(newOp.attributes, lastOp.attributes)) {
//   // skip for now   if (typeof newOp.insert === 'string' && typeof lastOp.insert === 'string') {
//   // skip for now     this.ops[index - 1] = { insert: lastOp.insert + newOp.insert };
//   // skip for now     if (typeof newOp.attributes === 'object') this.ops[index - 1].attributes = newOp.attributes
//   // skip for now     return this;
//   // skip for now   } else if (typeof newOp.retain === 'number' && typeof lastOp.retain === 'number') {
//   // skip for now     this.ops[index - 1] = { retain: lastOp.retain + newOp.retain };
//   // skip for now     if (typeof newOp.attributes === 'object') this.ops[index - 1].attributes = newOp.attributes
//   // skip for now     return this;
//   // skip for now   }
//   // skip for now }
//  }
//  if (index === this.ops.length) {
//    this.ops.push(newOp);
//  } else {
//    this.ops.splice(index, 0, newOp);
//  }
//  return this;
//
