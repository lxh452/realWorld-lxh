package utils

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"realWorld/global"
	"realWorld/model"
)

func init() {
	registerValidate()
	defineTranslator()
}

func usernameUnique(fieldLevel validator.FieldLevel) bool {
	// title的值
	value := fieldLevel.Field().Interface().(string)
	// id的值
	//id := fieldLevel.Parent().FieldByName("ID").Interface().(uint64)
	// 校验是否重复
	row := model.User{}
	global.DB.Raw("select * from user where username = ?", value).First(&row)
	fmt.Print("输出结果", row)
	// 判断是否查询到了
	return row.ID == 0
}

func registerValidate() {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册
		err := validate.RegisterValidation("usernameUnique", usernameUnique)
		if err != nil {
			panic(err)
		}
	}
}

var translator ut.Translator

func defineTranslator() {
	universalTranslator := ut.New(zh.New())
	// 具体验证引擎
	validate := binding.Validator.Engine().(*validator.Validate)
	translator, _ = universalTranslator.GetTranslator("zh")
	err := zhTranslations.RegisterDefaultTranslations(validate, translator)
	if err != nil {
		panic(err)
	}
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("chinese")
		if label == "-" {
			return ""
		}
		return label
	})
	validate.RegisterTranslation("usernameUnique", translator, func(ut ut.Translator) error {
		return ut.Add("usernameUnique", "{0} {1} 已存在", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("usernameUnique", fe.Field(), fe.Value().(string))
		return t
	})
}

// 翻译错误消息
func Translate(err error) string {
	// 仅翻译验证消息
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return ""
	}
	msg := ""
	for _, err := range errs {
		// 在这里使用定义好的翻译器
		msg += err.Translate(translator)
	}
	return msg
}
