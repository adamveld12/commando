package commando

import (
	"fmt"
	"strconv"
	"testing"
)

func TestAddAcceptsHandlerObjectOfTypeFunc(t *testing.T) {
	c := New()

	c.Add("add", "adds 2 numbers", func(a, b int) {})

	if len(c.commands) < 1 {
		t.Errorf("expected 1 command actual %d", len(c.commands))
	}
}

func TestAddPanicsOnHandlerObjectNotOfTypeFunc(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fail()
		}
	}()

	c := New()

	c.Add("add", "adds 2 numbers", 1)
}

func TestAddCommandRunsWithCorrectArguments(t *testing.T) {
	c := New()

	c.Add("add", "adds 2 numbers", func(a, b int) {
		if a != 2 {
			t.Errorf("expected 2 actual %d", a)
		}

		if b != 4 {
			t.Errorf("expected 4 actual %d", b)
		}
	})

	if err := c.Execute("add", "2", "4"); err != nil {
		t.Error(err)
	}
}

func TestParseUintTypes(t *testing.T) {
	c := New()

	c.Add("1", "", func(a uint8) {})
	c.Add("2", "", func(a uint16) {})
	c.Add("3", "", func(a uint32) {})
	c.Add("4", "", func(a uint64) {})
	c.Add("5", "", func(a uint) {})

	for i := 1; i <= 5; i++ {
		idx := fmt.Sprintf("%d", i)
		if err := c.Execute(idx, idx); err != nil {
			t.Error(err)
		}
	}
}

func TestParseBoolAndStringTypes(t *testing.T) {
	c := New()

	c.Add("1", "", func(a bool) {
		if !a {
			t.Errorf("expected true actual %v", strconv.FormatBool(a))
		}
	})
	c.Add("2", "", func(a string) {})

	for i := 1; i <= 2; i++ {
		idx := fmt.Sprintf("%d", i)
		if err := c.Execute(idx, "true"); err != nil {
			t.Error(err)
		}
	}
}

func TestParseIntTypes(t *testing.T) {
	c := New()

	c.Add("1", "", func(a int8) {})
	c.Add("2", "", func(a int16) {})
	c.Add("3", "", func(a int32) {})
	c.Add("4", "", func(a int64) {})
	c.Add("5", "", func(a int) {})

	for i := 1; i <= 5; i++ {
		idx := fmt.Sprintf("%d", i)
		if err := c.Execute(idx, idx); err != nil {
			t.Error(err)
		}
	}
}

func TestUsagesIncludeAllItems(t *testing.T) {
	c := New()

	c.Add("add", "adds 2 numbers", func(a, b int) {})

	usages := "Usage:\nhelp, h, --help\tDisplays usages\nadd [int, int]\tadds 2 numbers\n"
	actual := c.Usage()
	if actual != usages {
		t.Errorf("expected \n%s\n actual \n%s\n", usages, actual)
	}
}
