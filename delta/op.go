package delta

import (
	"log"
)

// AttrCompose takes two attributes maps and composes (combine) them
func AttrCompose(a, b map[string]interface{}, keepNil bool) map[string]interface{} {
	attributes := make(map[string]interface{})
	if b != nil {
		attributes = b
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
		log.Println("key ", v)
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