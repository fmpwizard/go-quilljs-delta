package delta

import (
	"reflect"
	"testing"
)

func TestAttrComposeLeftNil(t *testing.T) {
	attr := make(map[string]interface{})
	attr["bold"] = true
	attr["color"] = "red"

	if !reflect.DeepEqual(attr, AttrCompose(nil, attr, false)) {
		t.Errorf("failed to compose attr map, got: %+v\n", AttrCompose(nil, attr, false))
	}

}

func TestAttrComposeRightNil(t *testing.T) {
	attr := make(map[string]interface{})
	attr["bold"] = true
	attr["color"] = "red"

	if !reflect.DeepEqual(attr, AttrCompose(attr, nil, false)) {
		t.Errorf("failed to compose attr map, got: %+v\n", AttrCompose(attr, nil, false))
	}
}

func TestAttrComposeBothNil(t *testing.T) {
	if AttrCompose(nil, nil, false) != nil {
		t.Errorf("failed to compose attr map, got: %+v\n", AttrCompose(nil, nil, false))
	}
}

func TestAttrComposeMissingAttr(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	attr1["color"] = "red"

	attr2 := make(map[string]interface{})
	attr2["italic"] = true

	ret := make(map[string]interface{})
	ret["bold"] = true
	ret["color"] = "red"
	ret["italic"] = true

	if !reflect.DeepEqual(ret, AttrCompose(attr1, attr2, false)) {
		t.Errorf("failed to compose attr map, got: %+v\n", AttrCompose(attr1, attr2, false))
	}
}
func TestAttrComposeOverrideAttr(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	attr1["color"] = "red"

	attr2 := make(map[string]interface{})
	attr2["color"] = "blue"

	ret := make(map[string]interface{})
	ret["bold"] = true
	ret["color"] = "blue"

	if !reflect.DeepEqual(ret, AttrCompose(attr1, attr2, false)) {
		t.Errorf("failed to compose attr map, got: %+v\n", AttrCompose(attr1, attr2, false))
	}
}
func TestAttrComposeRemoveAttr(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	attr1["color"] = "red"

	attr2 := make(map[string]interface{})
	attr2["bold"] = nil

	ret := make(map[string]interface{})
	ret["color"] = "red"

	if !reflect.DeepEqual(ret, AttrCompose(attr1, attr2, false)) {
		t.Errorf("failed to compose attr map, got: %+v\n", AttrCompose(attr1, attr2, false))
	}
}
func TestAttrComposeRemoveAllAttr(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	attr1["color"] = "red"

	attr2 := make(map[string]interface{})
	attr2["bold"] = nil
	attr2["color"] = nil

	if AttrCompose(attr1, attr2, false) != nil {
		t.Errorf("failed to compose attr map, got: %+v\n", AttrCompose(attr1, attr2, false))
	}
}

func TestAttrComposeRemoveMissingAttr(t *testing.T) {
	attr1 := make(map[string]interface{})
	attr1["bold"] = true
	attr1["color"] = "red"

	attr2 := make(map[string]interface{})
	attr2["italic"] = nil

	ret := make(map[string]interface{})
	ret["color"] = "red"
	ret["bold"] = true

	if !reflect.DeepEqual(ret, AttrCompose(attr1, attr2, false)) {
		t.Errorf("failed to compose attr map, got: %+v\n", AttrCompose(attr1, attr2, false))
	}
}

func TestAttrDiffLeftNil(t *testing.T) {
	format := make(map[string]interface{})
	format["bold"] = true
	format["color"] = "red"
	if !reflect.DeepEqual(format, AttrDiff(nil, format)) {
		t.Errorf("failed to diff attr map, got: %+v\n", AttrDiff(nil, format))
	}
}
func TestAttrDiffRightNil(t *testing.T) {
	format := make(map[string]interface{})
	format["bold"] = true
	format["color"] = "red"

	expected := make(map[string]interface{})
	expected["bold"] = nil
	expected["color"] = nil

	if !reflect.DeepEqual(expected, AttrDiff(format, nil)) {
		t.Errorf("failed to diff attr map, got: %+v\n", AttrDiff(format, nil))
	}
}
func TestAttrDiffSame(t *testing.T) {
	format := make(map[string]interface{})
	format["bold"] = true
	format["color"] = "red"

	if AttrDiff(format, format) != nil {
		t.Errorf("failed to diff attr map, got: %+v\n", AttrDiff(format, format))
	}
}

func TestAttrDiffAddFormat(t *testing.T) {
	format := make(map[string]interface{})
	format["bold"] = true
	format["color"] = "red"

	added := make(map[string]interface{})
	added["bold"] = true
	added["italic"] = true
	added["color"] = "red"

	expected := make(map[string]interface{})
	expected["italic"] = true

	if !reflect.DeepEqual(expected, AttrDiff(format, added)) {
		t.Errorf("failed to diff attr map, got: %+v\n", AttrDiff(format, added))
	}
}
func TestAttrDiffRemoveFormat(t *testing.T) {
	format := make(map[string]interface{})
	format["bold"] = true
	format["color"] = "red"

	removed := make(map[string]interface{})
	removed["bold"] = true

	expected := make(map[string]interface{})
	expected["color"] = nil

	if !reflect.DeepEqual(expected, AttrDiff(format, removed)) {
		t.Errorf("failed to diff attr map, got: %+v\n", AttrDiff(format, removed))
	}
}

func TestAttrDiffOverrideFormat(t *testing.T) {
	format := make(map[string]interface{})
	format["bold"] = true
	format["color"] = "red"

	overwritten := make(map[string]interface{})
	overwritten["bold"] = true
	overwritten["color"] = "blue"

	expected := make(map[string]interface{})
	expected["color"] = "blue"

	if !reflect.DeepEqual(expected, AttrDiff(format, overwritten)) {
		t.Errorf("failed to diff attr map, got: %+v\n", AttrDiff(format, overwritten))
	}
}

func TestAttrTransformLeftNil(t *testing.T) {
	left := make(map[string]interface{})
	left["bold"] = true
	left["color"] = "red"
	left["font"] = nil

	right := make(map[string]interface{})
	right["font"] = "serif"
	right["color"] = "blue"
	right["italic"] = true

	if !reflect.DeepEqual(left, AttrTransform(nil, left, false)) {
		t.Errorf("failed to transform attr map, got: %+v\n", AttrTransform(nil, left, false))
	}
}
func TestAttrTransformRightNil(t *testing.T) {
	left := make(map[string]interface{})
	left["bold"] = true
	left["color"] = "red"
	left["font"] = nil

	right := make(map[string]interface{})
	right["font"] = "serif"
	right["color"] = "blue"
	right["italic"] = true

	if AttrTransform(left, nil, false) != nil {
		t.Errorf("failed to transform attr map, got: %+v\n", AttrTransform(left, nil, false))
	}
}
func TestAttrTransformBothtNil(t *testing.T) {
	if AttrTransform(nil, nil, false) != nil {
		t.Errorf("failed to transform attr map, got: %+v\n", AttrTransform(nil, nil, false))
	}
}
func TestAttrTransformWithPriority(t *testing.T) {
	left := make(map[string]interface{})
	left["bold"] = true
	left["color"] = "red"
	left["font"] = nil

	right := make(map[string]interface{})
	right["font"] = "serif"
	right["color"] = "blue"
	right["italic"] = true

	expected := make(map[string]interface{})
	expected["italic"] = true

	if !reflect.DeepEqual(expected, AttrTransform(left, right, true)) {
		t.Errorf("failed to transform attr map, got: %+v\n", AttrTransform(left, right, true))
	}
}
func TestAttrTransformWithoutPriority(t *testing.T) {
	left := make(map[string]interface{})
	left["bold"] = true
	left["color"] = "red"
	left["font"] = nil

	right := make(map[string]interface{})
	right["font"] = "serif"
	right["color"] = "blue"
	right["italic"] = true

	if !reflect.DeepEqual(right, AttrTransform(left, right, false)) {
		t.Errorf("failed to transform attr map, got: %+v\n", AttrTransform(left, right, false))
	}
}

func TestAttrLengthDelete(t *testing.T) {
	a := 5
	r := OpsLength(Op{
		Delete: &a,
	})

	if r != 5 {
		t.Error("failed to get length 5 for delete")
	}
}
func TestAttrLengthRetain(t *testing.T) {
	a := 4
	r := OpsLength(Op{
		Retain: &a,
	})

	if r != 4 {
		t.Error("failed to get length 4 for retain")
	}
}
func TestAttrLengthInsert(t *testing.T) {
	a := []rune("text")
	r := OpsLength(Op{
		Insert: a,
	})

	if r != 4 {
		t.Error("failed to get length 4 for insert")
	}
}

func TestAttrInvertUndefined(t *testing.T) {
	base := map[string]interface{}{"bold": true}

	if AttrInvert(nil, base) != nil {
		t.Errorf("Invalid inverted nil map")
	}
}

func TestAttrInvertBaseUndefined(t *testing.T) {
	attr := map[string]interface{}{"bold": true}
	expected := map[string]interface{}{"bold": nil}

	ret := AttrInvert(attr, nil)
	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("Wrong inverted attribute map, got: %+v\n", ret)
	}
}

func TestAttrInvertBothUndefined(t *testing.T) {
	if AttrInvert(nil, nil) != nil {
		t.Errorf("Invalid inverted nil map")
	}
}

func TestAttrInvertMerge(t *testing.T) {
	attr := map[string]interface{}{"bold": true}
	base := map[string]interface{}{"italic": true}
	expected := map[string]interface{}{"bold": nil}

	ret := AttrInvert(attr, base)
	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("Wrong inverted attribute map, got: %+v\n", ret)
	}
}

func TestAttrInvertNull(t *testing.T) {
	attr := map[string]interface{}{"bold": nil}
	base := map[string]interface{}{"bold": true}
	expected := map[string]interface{}{"bold": true}

	ret := AttrInvert(attr, base)
	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("Wrong inverted attribute map, got: %+v\n", ret)
	}
}

func TestAttrInvertReplace(t *testing.T) {
	attr := map[string]interface{}{"color": "red"}
	base := map[string]interface{}{"color": "blue"}
	expected := map[string]interface{}{"color": "blue"}

	ret := AttrInvert(attr, base)
	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("Wrong inverted attribute map, got: %+v\n", ret)
	}
}

func TestAttrInvertNoop(t *testing.T) {
	attr := map[string]interface{}{"color": "red"}
	base := map[string]interface{}{"color": "red"}

	ret := AttrInvert(attr, base)
	if ret != nil {
		t.Errorf("Wrong inverted attribute map, got: %+v\n", ret)
	}
}

func TestAttrInvertCombined(t *testing.T) {
	attr := map[string]interface{}{"bold": true, "italic": nil, "color": "red", "size": "12px"}
	base := map[string]interface{}{"font": "serif", "italic": true, "color": "blue", "size": "12px"}
	expected := map[string]interface{}{"bold": nil, "italic": true, "color": "blue"}

	ret := AttrInvert(attr, base)
	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("Wrong inverted attribute map, got: %+v\n", ret)
	}
}
