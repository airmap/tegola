package subject

import (
	"fmt"

	colour "github.com/logrusorgru/aurora"

	"github.com/airmap/tegola/container/singlelist/point/list"
	"github.com/airmap/tegola/maths"
)

func (s *Subject) GoString() string {
	str := fmt.Sprintf("  Subject:(%v)", s.winding)
	s.ForEachIdx(func(idx int, pt list.ElementerPointer) bool {
		str += fmt.Sprintf("[%v](%#v)", idx, colour.Green(pt))
		return true
	})
	return str
}

func (s *Subject) DebugStringAugmented(augmenter func(idx int, e maths.Pt) string) string {
	str := fmt.Sprintf("  Subject:(%v)", s.winding)
	s.ForEachIdx(func(i int, pt list.ElementerPointer) bool {
		str += augmenter(i, pt.Point())
		return true
	})
	return str
}
