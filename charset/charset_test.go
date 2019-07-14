package charset


import "testing"

func TestHumpFormat(t *testing.T) {
	a := "hello"
	b := "word"

	res, err := HumpFormat(a, b)
	t.Logf("res: %s , err: %v", res, err)
	//t.Logf("res: %s , err: %v", HumpFormat(a, b))
}

