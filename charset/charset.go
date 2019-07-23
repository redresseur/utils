package charset

import (
	"errors"
	"regexp"
	"unicode/utf8"
	"strings"
)

/*
	instruction: the charset operation
	author: wangzhipengtest@163.com
	date: 2019/07/06
*/

//判断身份证是否符合
func CheckPIdRight(pId string) error {

	if len(pId) == 15 {
		//验证15位身份证，15位的是全部数字
		if m, _ := regexp.MatchString(`^(\d{15})$`, pId); !m {
			return errors.New("身份证号不正确")
		}
	} else {
		//验证18位身份证，18位前17位为数字，最后一位是校验位，可能为数字或字符X。
		if m, _ := regexp.MatchString(`^(\d{17})([0-9]|X)$`, pId); !m {
			return errors.New("身份证号不正确")
		}
	}

	return nil
}

var ErrSpecialCharacter = errors.New("Special characters are included")

// special charset check
func CheckSpecialCharacter(str string)bool{
	//re :=regexp.MustCompile("[~!@#$%^&*(){}|<>\\\\/+\\-【】:\"?'：；‘’“”，。、《》\\]\\[`]")
	//re :=regexp.MustCompile("[!-/]|[:-@]|[\\[-`]")
	//re :=regexp.MustCompile("[~!@#$%^&*(){}|<>\\\\/+\\-【】:\"?'：；‘’“”，。、《》\\]\\[`]")
	// fixMe 取消空格限制,取消下划线限制
	// 空格： \u0020
	// - ：\u002D
	// _ : \u005F
	re :=regexp.MustCompile("[\u0021-\u002F]|[\u003A-\u0040]|[\u005B-\u0060]|[\u00A0-\u00BF]")
	return re.MatchString(str)
}

var ErrNotLetter = errors.New("the word is not letter")

// 小写转大写
func capitals(c rune) (rune, error){
	//unicode A-Z 0x0041 ~0x005a a-z 0x0061 ~0x007a
	if c >= 0x0041 && c <= 0x005a{
		return c, nil
	}

	if c >= 0x0061 && c <= 0x007a{
		return c - 0x0020 , nil
	}

	return c, ErrNotLetter
}

func CamelCaseFormatMust(firstWordCapital bool, str string) string {
	// 剔除 _ 字符
	res, err := CamelCaseFormat(firstWordCapital, strings.Split(str, "_")...)
	if err != nil{
		panic(err)
	}

	return res
}

// 驼峰格式化
func CamelCaseFormat(firstWordCapital bool, strs... string) (string, error){
	res := ""
	// var err error
	if len(strs) == 0{
		return "", nil
	}

	// TODO: 过滤掉所有特殊字符
	for i, str := range strs {
		if CheckSpecialCharacter(str){
			return "", ErrSpecialCharacter
		}

		if !firstWordCapital && i == 0{
			res += strings.ToLower(str)
		}else {
			data := make([]byte, len(str))
			for i, c := range str {
				// 首字符必须转为大写
				if i == 0{
					//if c, err = capitals(c); err != nil{
					//	return "", err
					//}
					c, _ = capitals(c)
				}

				utf8.EncodeRune(data[i:], c)
			}
			res += string(data)
		}
	}

	return res, nil
}