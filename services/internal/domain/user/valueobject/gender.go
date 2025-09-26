package valueobject

type Gender int

const (
	GenderMale   Gender = 100 // 男性
	GenderFemale Gender = 200 // 女性
	GenderOther  Gender = 300 // 其他
)

func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderOther:
		return true
	}
	return false
}

func (g Gender) Int() int {
	return int(g)
}
