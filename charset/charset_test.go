package charset

import "testing"

func TestCamelCaseFormat(t *testing.T) {
	a := "hello"
	b := "word"

	res, err := CamelCaseFormat(true, a, b)
	t.Logf("res: %s , err: %v", res, err)
	//t.Logf("res: %s , err: %v", HumpFormat(a, b))
}

func TestByteToHexString(t *testing.T) {
	t.Log(ByteToHexString([]byte{0x1f, 0x2e, 0xf7, 0x69}))
}
