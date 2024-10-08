package util

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
	"regexp"
)

/**
 * @Author Khr
 * @Description 验证参数是否存在
 * @Date 9:38 2024/2/20
 * @Param
 * @return
 **/
func ValidateExist(a string, b []string) bool {
	for _, s := range b {
		if a == s {
			return true
		}
	}
	return false
}

/**
 * @Author Khr
 * @Description 通过反射获取结构体字段的标签处理验证错误并
 * @Date 14:55 2024/5/15
 * @Param
 * @return
 **/
func GetValidate(err error, obj any) error {
	log.Println(err, "++++++++++++++")
	var invalid *validator.InvalidValidationError
	ok := errors.As(err, &invalid)
	if ok {
		fmt.Println("param error:", invalid)
		return invalid
	}
	//反射获取标签的注释
	getObj := reflect.TypeOf(obj)
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		//return errs
		for _, e := range errs {
			if f, exist := getObj.Elem().FieldByName(e.Field()); exist {
				msg := f.Tag.Get("msg")
				log.Println(msg, "msg********")
				return errors.New(msg)
			}
		}
	}
	log.Println(err, "********err")
	return err
}

/*
 * @MethodName FuzzyMatch
 * @Description 正则模糊匹配路径
 * @Author khr
 * @Date 2023/5/9 16:25
 */
func FuzzyMatch(param string, paths []string) bool {
	for _, y := range paths {
		if regexp.MustCompile(y).MatchString(param) {

			//fmt.Print("匹配道路进了")
			return true
		}

	}
	return false
}
