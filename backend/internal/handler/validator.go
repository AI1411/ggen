package handler

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"time"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/ja_JP"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	jatranslations "github.com/go-playground/validator/v10/translations/ja"

	myerrors "g_gen/internal/errors"
)

const (
	passwordMinLength     = 8
	minimumLength         = "8"
	maximumLength         = "20"
	passwordTag           = "password"
	datetimeTag           = "datetime"
	alphaNumUnderscoreTag = "alphanum_underscore"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
	jatrans  ut.Translator
)

func init() {
	ja := ja_JP.New()
	uni = ut.New(ja, ja)
	jatrans, _ = uni.GetTranslator("ja")

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate = v
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			return field.Tag.Get("ja")
		})
	}

	_ = validate.RegisterValidation("vcf", validateVcfFormat)
	_ = validate.RegisterTranslation("vcf", jatrans, func(ut ut.Translator) error {
		return ut.Add("vcf", "{0}は`vcf-`と8文字の文字列である必要があります", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("vcf", fe.Field())
		return t
	})

	_ = validate.RegisterValidation("vcs", validateVcsFormat)
	_ = validate.RegisterTranslation("vcs", jatrans, func(ut ut.Translator) error {
		return ut.Add("vcs", "{0}は`vcs-`と8文字の文字列である必要があります", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("vcs", fe.Field())
		return t
	})

	_ = jatranslations.RegisterDefaultTranslations(validate, jatrans)
	_ = validate.RegisterValidation(passwordTag, validatePassword)
	_ = validate.RegisterTranslation(passwordTag, jatrans, func(ut ut.Translator) error {
		return ut.Add(passwordTag, "{0}は{1}文字以上{2}文字以下で、大文字、小文字、数字、特殊文字をそれぞれ1つ以上含む必要があります", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(passwordTag, fe.Field(), minimumLength, maximumLength)
		return t
	})

	_ = validate.RegisterValidation(datetimeTag, isDateTimeString)
	_ = validate.RegisterTranslation(datetimeTag, jatrans, func(ut ut.Translator) error {
		return ut.Add(datetimeTag, "{0}は日付時刻形式である必要があります", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(datetimeTag, fe.Field())
		return t
	})

	// 半角英数と_のみを許容するバリデーションを登録
	_ = validate.RegisterValidation(alphaNumUnderscoreTag, validateAlphaNumUnderscore)
	_ = validate.RegisterTranslation(alphaNumUnderscoreTag, jatrans, func(ut ut.Translator) error {
		return ut.Add(alphaNumUnderscoreTag, "{0}は半角英数字とアンダースコア(_)のみ使用できます", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(alphaNumUnderscoreTag, fe.Field())
		return t
	})
}

// 日付時刻のバリデーション
func isDateTimeString(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.RFC3339, fl.Field().String())
	return err == nil
}

func createValidateErrorResponse(err error) *ErrorResponseDetail {
	var verr validator.ValidationErrors
	errors.As(err, &verr)

	details := make([]ValidationError, len(verr))
	for i, fe := range verr {
		// エラーメッセージを日本語に翻訳
		details[i] = ValidationError{
			Attribute: fe.Field(),
			Tag:       fe.Tag(),
			Message:   fe.Translate(jatrans),
		}
	}

	return &ErrorResponseDetail{
		Code:    myerrors.ValidationError,
		Message: myerrors.ValidationErrorMessage,
		Details: details,
		err:     err,
		status:  http.StatusBadRequest,
	}
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < passwordMinLength {
		return false
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func validateVcfFormat(fl validator.FieldLevel) bool {
	tagID := fl.Field().String()
	ok, _ := regexp.MatchString(`^vcf-[a-zA-Z0-9]{8}$`, tagID)
	return ok
}

func validateVcsFormat(fl validator.FieldLevel) bool {
	tagID := fl.Field().String()
	ok, _ := regexp.MatchString(`^vcs-[a-zA-Z0-9]{8}$`, tagID)
	return ok
}

func validateAlphaNumUnderscore(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true // 空文字は他のバリデーション（required等）に任せる
	}

	// 正規表現で半角英数とアンダースコアのみかチェック
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, value)
	return matched
}
