package levelFinder
import (
	"testing"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestSearchLevelFromProvider(t *testing.T) {
	ExpTable := NewArrayLevelProvider(4)
	ExpTable.SetExpByLevel(1, 0)
	ExpTable.SetExpByLevel(2, 12)
	ExpTable.SetExpByLevel(3, 29)
	ExpTable.SetExpByLevel(4, 57)

	for _, c := range []struct {
		exp int
		lv  int
	}{
		{0, 1},
		{11, 1},
		{12, 2},
		{13, 2},
		{28, 2},
		{29, 3},
		{56, 3},
		{57, 4},
		{58, 4},
	} {
		kmgTest.Equal(searchLevelFromProvider(ExpTable, c.exp), c.lv, "Level at exp:%d", c.exp)
	}
	ExpTable = NewArrayLevelProvider(1)
	kmgTest.Equal(searchLevelFromProvider(ExpTable, 0), 1)
	kmgTest.Equal(searchLevelFromProvider(ExpTable, 2), 1)

	ExpTable = NewArrayLevelProvider(2)
	ExpTable.SetExpByLevel(2, 100)
	kmgTest.Equal(searchLevelFromProvider(ExpTable, 0), 1)
	kmgTest.Equal(searchLevelFromProvider(ExpTable, 99), 1)
	kmgTest.Equal(searchLevelFromProvider(ExpTable, 100), 2)
}

type getLevelByExpTestCase struct {
	exp        int
	actualExp  int
	lv         int
	excess     int
	nextAll    int
	isMaxLevel bool
}

func TestLevelFinder(t *testing.T) {
	ExpTable := NewArrayLevelProvider(4)
	ExpTable.SetExpByLevel(1, 0)
	ExpTable.SetExpByLevel(2, 12)
	ExpTable.SetExpByLevel(3, 29)
	ExpTable.SetExpByLevel(4, 57)
	checkTestTable := func(levelTestTable []getLevelByExpTestCase) {
		for _, c := range levelTestTable {
			result := GetLevelByExp(ExpTable, c.exp)
			kmgTest.Equal(result.Exp, c.actualExp, "actualExp at exp:%d", c.exp)
			kmgTest.Equal(result.Level, c.lv, "Level at exp:%d", c.exp)
			kmgTest.Equal(result.CurrentLevelExcessExp, c.excess, "CurrentLevelExcessExp at exp:%d", c.exp)
			kmgTest.Equal(result.NextLevelAllNeedExp, c.nextAll, "NextLevelAllNeedExp at exp:%d", c.exp)
			kmgTest.Equal(result.IsMaxLevel, c.isMaxLevel, "IsMaxLevel at exp:%d", c.exp)
		}
	}

	checkTestTable([]getLevelByExpTestCase{
		{11, 11, 1, 11, 12, false},
		{12, 12, 2, 0, 17, false},
		{13, 13, 2, 1, 17, false},
		{-1, 0, 1, 0, 12, false},
		{56, 56, 3, 27, 28, false},
		{57, 57, 4, 0, 0, true},
		{58, 57, 4, 0, 0, true},
		{1000, 57, 4, 0, 0, true},
	})

	ExpTable = NewArrayLevelProvider(1)
	checkTestTable([]getLevelByExpTestCase{
		{-1, 0, 1, 0, 0, true},
		{0, 0, 1, 0, 0, true},
		{1, 0, 1, 0, 0, true},
		{1000, 0, 1, 0, 0, true},
	})
}
