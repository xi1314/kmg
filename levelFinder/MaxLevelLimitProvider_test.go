package levelFinder
import (
	"testing"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestMaxLevelLimitProvider(t *testing.T) {
	ExpTable, err := MaxLevelLimitProvider(ArrayLevelProvider{12, 29, 57}, 4)
	kmgTest.Equal(err, nil)
	levelTestTable := []struct {
		exp        int
		actualExp  int
		lv         int
		excess     int
		nextAll    int
		isMaxLevel bool
	}{
		{11, 11, 1, 11, 12, false},
		{12, 12, 2, 0, 17, false},
		{13, 13, 2, 1, 17, false},
		{-1, 0, 1, 0, 12, false},
		{56, 56, 3, 27, 28, false},
		{57, 57, 4, 0, 0, true},
		{58, 57, 4, 0, 0, true},
		{1000, 57, 4, 0, 0, true},
	}
	for _, c := range levelTestTable {
		result := GetLevelByExp(ExpTable, c.exp)
		kmgTest.Equal(result.Exp, c.actualExp, "actualExp at exp:%d", c.exp)
		kmgTest.Equal(result.Level, c.lv, "Level at exp:%d", c.exp)
		kmgTest.Equal(result.CurrentLevelExcessExp, c.excess, "CurrentLevelExcessExp at exp:%d", c.exp)
		kmgTest.Equal(result.NextLevelAllNeedExp, c.nextAll, "NextLevelAllNeedExp at exp:%d", c.exp)
		kmgTest.Equal(result.IsMaxLevel, c.isMaxLevel, "IsMaxLevel at exp:%d", c.exp)
	}

	ExpTable, err = MaxLevelLimitProvider(ArrayLevelProvider{12, 29, 57}, 1)
	kmgTest.Equal(err, nil)
	levelTestTable = []struct {
		exp        int
		actualExp  int
		lv         int
		excess     int
		nextAll    int
		isMaxLevel bool
	}{
		{11, 0, 1, 0, 0, true},
	}
	for _, c := range levelTestTable {
		result := GetLevelByExp(ExpTable, c.exp)
		kmgTest.Equal(result.Exp, c.actualExp, "actualExp at exp:%d", c.exp)
		kmgTest.Equal(result.Level, c.lv, "Level at exp:%d", c.exp)
		kmgTest.Equal(result.CurrentLevelExcessExp, c.excess, "CurrentLevelExcessExp at exp:%d", c.exp)
		kmgTest.Equal(result.NextLevelAllNeedExp, c.nextAll, "NextLevelAllNeedExp at exp:%d", c.exp)
		kmgTest.Equal(result.IsMaxLevel, c.isMaxLevel, "IsMaxLevel at exp:%d", c.exp)
	}
}
