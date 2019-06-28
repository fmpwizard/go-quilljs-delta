package delta

// AttrCompose takes two attributes maps and composes (combine) them
func AttrCompose(a, b map[string]interface{}, keepNil bool) map[string]interface{} {
	attributes := make(map[string]interface{})
	if b != nil {
		for k, v := range b {
			attributes[k] = v
		}
	}

	for k := range a {
		aa, aFound := a[k]
		bb, bFound := b[k]
		if !keepNil && bb == nil { // a nil check to match the null special case in quilljs
			delete(attributes, k)
		}

		if aFound && !bFound {
			attributes[k] = aa
		}
	}
	// clean up any nil attributes that were on b but not on a
	for k, v := range attributes {
		if !keepNil && v == nil {
			delete(attributes, k)
		}
	}
	if len(attributes) > 0 {
		return attributes
	}
	return nil
}

// AttrDiff returns the diff between two maps of attributes
func AttrDiff(a, b map[string]interface{}) map[string]interface{} {
	keys := make([]string, len(a)+len(b))
	attributes := make(map[string]interface{})

	if a == nil {
		a = make(map[string]interface{})
	}

	if b == nil {
		b = make(map[string]interface{})
	}
	for k := range a {
		keys = append(keys, k)
	}
	for kk := range b {
		keys = append(keys, kk)
	}

	for _, v := range keys {
		if a[v] != b[v] {
			bb, bFound := b[v]
			if !bFound {
				attributes[v] = nil
			} else {
				attributes[v] = bb
			}
		}
	}
	if len(attributes) > 0 {
		return attributes
	}
	return nil
}

// AttrTransform is used in Detal.transform(), hard to really explain
func AttrTransform(a, b map[string]interface{}, priority bool) map[string]interface{} {
	if a == nil {
		return b
	}
	if b == nil {
		return nil
	}
	// b simply overwrites us without priority
	if !priority {
		return b
	}

	attributes := make(map[string]interface{})

	for k := range b {
		_, aOk := a[k]
		if !aOk {
			// nil is a valid value
			attributes[k] = b[k]
		}
	}

	if len(attributes) > 0 {
		return attributes
	}
	return nil
}

// AttrInvert inverts an attribute map, used in Delta.Invert()
func AttrInvert(attr map[string]interface{}, base map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for k, v := range base {
		v2, exists := attr[k]
		if exists && v2 != v {
			ret[k] = v
		}
	}
	for k, v := range attr {
		v2, exists := base[k]
		if (!exists) && v2 != v {
			ret[k] = nil
		}
	}
	if len(ret) == 0 {
		return nil
	}
	return ret
}

// OpsIterator returns an Iterator wrapping the ops
func OpsIterator(ops []Op) Iterator {
	return NewIterator(ops)
}

// OpsLength returns the length of the string insert, or the numeric value of Delete or Retain
func OpsLength(op Op) int {
	if op.Delete != nil {
		return *op.Delete
	}
	if op.Retain != nil {
		return *op.Retain
	}
	if op.Insert != nil {
		return len(op.Insert)
	}
	if op.InsertEmbed != nil {
		return 1
	}

	return 1
}
