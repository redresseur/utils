package charset

import (
	"errors"
	"regexp"
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
