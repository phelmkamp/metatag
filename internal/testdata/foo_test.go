package testdata

import "testing"

import "reflect"

type Resizer interface {
	Size() int
	SetSize(int)
}

func TestFoo_Name(t *testing.T) {
	f := Foo{name: "bar"}
	if got := f.Name(); got != f.name {
		t.Errorf("Name() = %v, want %v", got, f.name)
	}
}

func TestFoo_GetDesc(t *testing.T) {
	f := Foo{Desc: "bar"}
	if got := f.GetDesc(); got != f.Desc {
		t.Errorf("GetDesc() = %v, want %v", got, f.Desc)
	}
}

func TestFoo_Size(t *testing.T) {
	f := Foo{size: 42}
	var r Resizer = &f
	if got := r.Size(); got != f.size {
		t.Errorf("Size() = %v, want %v", got, f.size)
	}
}

func TestFoo_SetSize(t *testing.T) {
	var r Resizer = &Foo{}
	r.SetSize(42)
	if got := r.Size(); got != 42 {
		t.Errorf("SetSize() = %v, want %v", got, 42)
	}
}

func TestFoo_FindLabel(t *testing.T) {
	f := Foo{labels: []string{"a", "b", "c"}}
	if got := f.FindLabel("b"); got != 1 {
		t.Errorf("FindLabel() = %d, want %d", got, 1)
	}
}

func TestFoo_FilterLabels(t *testing.T) {
	f := Foo{labels: []string{"a", "aa", "b", "bb"}}
	isMultiByte := func(s string) bool { return len(s) > 1 }
	if got := f.FilterLabels(isMultiByte); !reflect.DeepEqual(got, []string{"aa", "bb"}) {
		t.Errorf("FilterLabels() = %v, want %v", got, []string{"aa", "bb"})
	}
}
