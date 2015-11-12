package kmgType

import (
	"testing"

	"github.com/bronze1man/kmg/kmgTest"
)

type T struct {
	String1 string
	Map1    map[string]string
	Map2    map[string]*string
	Map3    map[string]T2
	Map4    map[string]map[string]string
	Map5    map[string][]string
	Slice1  []string
	Ptr1    *string
	Ptr2    *T2
	Array1  [5]string
}

type T2 struct {
	A string
	B string
}

func TestPtrType(ot *testing.T) {
	var data **string
	data = new(*string)
	m, err := NewContext(data)
	kmgTest.Equal(err, nil)

	err = m.SaveByPath(Path{"ptr", "ptr"}, "")
	kmgTest.Equal(err, nil)
	kmgTest.Ok(data != nil)
	kmgTest.Ok(*data != nil)
	kmgTest.Equal(**data, "")
}

func TestStringType(ot *testing.T) {
	var data *string
	data = new(string)
	m, err := NewContext(data)
	kmgTest.Equal(err, nil)

	err = m.SaveByPath(Path{"ptr"}, "123")
	kmgTest.Equal(err, nil)
	kmgTest.Ok(data != nil)
	kmgTest.Equal(*data, "123")
}

func TestStructType(ot *testing.T) {
	data := &struct {
		A string
	}{}
	m, err := NewContext(data)
	kmgTest.Equal(err, nil)

	err = m.SaveByPath(Path{"ptr", "A"}, "123")
	kmgTest.Equal(err, nil)
	kmgTest.Ok(data != nil)
	kmgTest.Equal(data.A, "123")
}

func TestType(ot *testing.T) {
	data := &T{}
	m, err := NewContext(data)
	kmgTest.Equal(err, nil)

	err = m.SaveByPath(Path{"ptr", "String1"}, "B")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(data.String1, "B")

	m.SaveByPath(Path{"ptr", "Map1", "A"}, "1123")
	_, ok := data.Map1["A"]
	kmgTest.Equal(ok, true)
	kmgTest.Equal(data.Map1["A"], "1123")

	err = m.SaveByPath(Path{"ptr", "Map1", "A"}, "1124")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(data.Map1["A"], "1124")

	err = m.DeleteByPath(Path{"ptr", "Map1", "A"})
	kmgTest.Equal(err, nil)
	_, ok = data.Map1["A"]
	kmgTest.Equal(ok, false)

	err = m.SaveByPath(Path{"ptr", "Map2", "B", "ptr"}, "1")
	kmgTest.Equal(err, nil)
	rpString, ok := data.Map2["B"]
	kmgTest.Equal(ok, true)
	kmgTest.Equal(*rpString, "1")

	err = m.SaveByPath(Path{"ptr", "Map2", "B", "ptr"}, "2")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(*rpString, "2")

	err = m.DeleteByPath(Path{"ptr", "Map2", "B", "ptr"})
	kmgTest.Equal(err, nil)
	_, ok = data.Map2["B"]
	kmgTest.Equal(ok, true)
	kmgTest.Equal(data.Map2["B"], nil)

	err = m.DeleteByPath(Path{"ptr", "Map2", "B"})
	kmgTest.Equal(err, nil)
	_, ok = data.Map2["B"]
	kmgTest.Equal(ok, false)

	err = m.SaveByPath(Path{"ptr", "Map3", "C", "A"}, "1")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(data.Map3["C"].A, "1")

	err = m.DeleteByPath(Path{"ptr", "Map3", "C"})
	kmgTest.Equal(err, nil)
	kmgTest.Ok(data.Map3 != nil)
	_, ok = data.Map3["C"]
	kmgTest.Equal(ok, false)

	err = m.SaveByPath(Path{"ptr", "Map4", "D", "F"}, "1234")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(data.Map4["D"]["F"], "1234")

	err = m.SaveByPath(Path{"ptr", "Map4", "D", "H"}, "12345")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(data.Map4["D"]["H"], "12345")

	err = m.SaveByPath(Path{"ptr", "Map4", "D", "H"}, "12346")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(data.Map4["D"]["H"], "12346")

	err = m.DeleteByPath(Path{"ptr", "Map4", "D", "F"})
	kmgTest.Equal(err, nil)
	kmgTest.Ok(data.Map4["D"] != nil)
	_, ok = data.Map4["D"]["F"]
	kmgTest.Equal(ok, false)

	_, ok = data.Map4["D"]["H"]
	kmgTest.Equal(ok, true)

	err = m.SaveByPath(Path{"ptr", "Map5", "D", ""}, "1234")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(len(data.Map5["D"]), 1)
	kmgTest.Equal(data.Map5["D"][0], "1234")

	err = m.DeleteByPath(Path{"ptr", "Map5", "D", "0"})
	kmgTest.Equal(err, nil)
	kmgTest.Equal(len(data.Map5["D"]), 0)

	err = m.SaveByPath(Path{"ptr", "Slice1", ""}, "1234")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(len(data.Slice1), 1)
	kmgTest.Equal(data.Slice1[0], "1234")

	err = m.SaveByPath(Path{"ptr", "Slice1", ""}, "12345")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(data.Slice1[1], "12345")
	kmgTest.Equal(len(data.Slice1), 2)

	err = m.DeleteByPath(Path{"ptr", "Slice1", "0"})
	kmgTest.Equal(err, nil)
	kmgTest.Equal(len(data.Slice1), 1)
	kmgTest.Equal(data.Slice1[0], "12345")

	err = m.SaveByPath(Path{"ptr", "Ptr1", "ptr"}, "12345")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(*data.Ptr1, "12345")

	err = m.SaveByPath(Path{"ptr", "Ptr2", "ptr"}, "")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(data.Ptr2.A, "")

	err = m.SaveByPath(Path{"ptr", "Array1", "1"}, "12345")
	kmgTest.Equal(err, nil)
	kmgTest.Equal(data.Array1[1], "12345")
}
