package delta

import (
//"log"
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
