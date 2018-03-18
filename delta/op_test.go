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
