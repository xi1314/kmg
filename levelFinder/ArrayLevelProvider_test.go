package levelFinder
import (
	"testing"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestArrayLevelProvider(t *testing.T) {
	arrayLevelProvider := NewArrayLevelProvider(3)
	arrayLevelProvider.SetExpByLevel(1, 0)
	arrayLevelProvider.SetExpByLevel(2, 100)
	arrayLevelProvider.SetExpByLevel(3, 200)
	kmgTest.Equal([]int(arrayLevelProvider), []int{100, 200})
	kmgTest.Equal(arrayLevelProvider.MaxLevel(), 3)

	arrayLevelProvider = NewArrayLevelProvider(2)
	arrayLevelProvider.SetExpByLevel(1, 0)
	arrayLevelProvider.SetExpByLevel(2, 100)
	kmgTest.Equal([]int(arrayLevelProvider), []int{100})
	kmgTest.Equal(arrayLevelProvider.MaxLevel(), 2)

	kmgTest.Equal(arrayLevelProvider.GetExpByLevel(2), 100)
	kmgTest.Equal(arrayLevelProvider.GetExpByLevel(1), 0)

	arrayLevelProvider = NewArrayLevelProvider(1)
	kmgTest.Equal(arrayLevelProvider.MaxLevel(), 1)
	kmgTest.Equal([]int(arrayLevelProvider), []int{})
}
