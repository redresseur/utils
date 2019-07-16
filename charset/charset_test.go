package charset


import "testing"

func TestCamelCaseFormat(t *testing.T) {
	a := "hello"
	b := "word"

	res, err := CamelCaseFormat(true, a, b)
	t.Logf("res: %s , err: %v", res, err)
	//t.Logf("res: %s , err: %v", HumpFormat(a, b))
}

