package main

import (
	"testing"
)

// TestStruct — простая структура для проверки пула.
type TestStruct struct {
	Value int
	Slice []int
}

// Reset реализует интерфейс Resetter.
func (t *TestStruct) Reset() {
	if t == nil {
		return
	}
	t.Value = 0
	t.Slice = t.Slice[:0]
}

func TestPool_BasicFlow(t *testing.T) {
	p := New(func() *TestStruct {
		return &TestStruct{
			Slice: make([]int, 0, 10),
		}
	})

	obj1 := p.Get()
	obj1.Value = 123
	obj1.Slice = append(obj1.Slice, 1, 2, 3)

	if obj1.Value != 123 || len(obj1.Slice) != 3 {
		t.Fatalf("unexpected state: %+v", obj1)
	}

	p.Put(obj1)

	obj2 := p.Get()

	if obj2.Value != 0 {
		t.Errorf("expected Value 0, got %d", obj2.Value)
	}
	if len(obj2.Slice) != 0 {
		t.Errorf("expected Slice len 0, got %d", len(obj2.Slice))
	}

	if cap(obj2.Slice) != 10 {
		t.Errorf("expected Slice cap 10, got %d", cap(obj2.Slice))
	}
}
