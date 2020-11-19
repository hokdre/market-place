package helper

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	"github.com/market-place/domain"
	"github.com/market-place/usecase/usecase_error"
)

var (
	validation *validator10
)

type ValidationEntity interface {
	Validate(value interface{}) error
}

type validator10 struct {
	validation *validator.Validate
	translator locales.Translator
	trans      ut.Translator
}

func NewValidationEntity() ValidationEntity {
	if validation == nil {
		translator := en.New()
		uni := ut.New(translator, translator)
		trans, _ := uni.GetTranslator("en")

		v := &validator10{
			validation: validator.New(),
			translator: translator,
			trans:      trans,
		}
		v.RegisterTranslation()
		v.RegisterValidation()

		return v
	}

	return validation
}

func (v *validator10) Validate(value interface{}) error {
	log.SetOutput(os.Stdout)

	if err := v.validation.Struct(value); err != nil {
		var entityErrs usecase_error.ErrBadEntityInput
		errs := err.(validator.ValidationErrors)
		for _, err := range errs {
			log.Printf("Validation : field : %s, value : %s, details : %s \n", err.Field(), err.Value(), err.ActualTag())

			e := usecase_error.ErrEntityField{
				Field:   err.Field(),
				Message: err.Translate(v.trans),
			}
			entityErrs = append(entityErrs, e)
		}
		return entityErrs
	}
	return nil
}

func (v *validator10) RegisterTranslation() {
	v.validation.RegisterTranslation("unique", v.trans, func(ut ut.Translator) error {
		return ut.Add("unique", "{0} is not unique", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		field := fe.Field()
		if fe.Field() == "BankAccounts" {
			field = "Number"
		}
		t, _ := ut.T("unique", field)
		return t
	})

	v.validation.RegisterTranslation("required", v.trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is must be filled", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("email", v.trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} is not valid email format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("phone", v.trans, func(ut ut.Translator) error {
		return ut.Add("phone", "{0} is not valid phone number format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("gender", v.trans, func(ut ut.Translator) error {
		return ut.Add("gender", "{0} is not valid", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("gender", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("min", v.trans, func(ut ut.Translator) error {
		return ut.Add("min", "{0} is must greater than equal {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("min", fe.Field(), fe.Param())
		return t
	})

	v.validation.RegisterTranslation("clower", v.trans, func(ut ut.Translator) error {
		return ut.Add("clower", "{0} is must contains a lowercase caracter", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("clower", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("cupper", v.trans, func(ut ut.Translator) error {
		return ut.Add("cupper", "{0} is must contains a uppercase caracter", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("cupper", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("cnumeric", v.trans, func(ut ut.Translator) error {
		return ut.Add("cnumeric", "{0} is must contains a numeric caracter", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("cnumeric", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("csymbol", v.trans, func(ut ut.Translator) error {
		return ut.Add("csymbol", "{0} is must contains a symbol caracter", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("csymbol", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("ltfield", v.trans, func(ut ut.Translator) error {
		return ut.Add("ltfield", "{0} is must less than {1}", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		var param string
		if fe.Param() == "CreatedAt" {
			param = "Today"
		} else {
			param = fe.Param()
		}
		t, _ := ut.T("ltfield", fe.Field(), param)
		return t
	})

	v.validation.RegisterTranslation("bank_provider", v.trans, func(ut ut.Translator) error {
		return ut.Add("bank_provider", "{0} is not valid", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("bank_provider", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("bank_number", v.trans, func(ut ut.Translator) error {
		return ut.Add("bank_number", "{0} is not valid", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("bank_number", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("unique_addresses", v.trans, func(ut ut.Translator) error {
		return ut.Add("unique_addresses", "{0} is not unique", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique_addresses", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("unique_bank_accounts", v.trans, func(ut ut.Translator) error {
		return ut.Add("unique_bank_accounts", "{0} is not unique", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique_bank_accounts", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("unique_shippings", v.trans, func(ut ut.Translator) error {
		return ut.Add("unique_shippings", "{0} is not unique", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique_shippings", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("unique_colors", v.trans, func(ut ut.Translator) error {
		return ut.Add("unique_colors", "{0} is not unique", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique_colors", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("unique_sizes", v.trans, func(ut ut.Translator) error {
		return ut.Add("unique_sizes", "{0} is not unique", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique_sizes", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("unique_items", v.trans, func(ut ut.Translator) error {
		return ut.Add("unique_items", "{0} is not unique", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique_items", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("category", v.trans, func(ut ut.Translator) error {
		return ut.Add("category", "category is not valid", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("category", fe.Field())
		return t
	})

	v.validation.RegisterTranslation("unique_etalase", v.trans, func(ut ut.Translator) error {
		return ut.Add("unique_etalase", "Etalase is not unique", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique_etalase", fe.Field())
		return t
	})
}

func (v *validator10) RegisterValidation() {
	v.validation.RegisterValidation("phone", phoneValidation)
	v.validation.RegisterValidation("gender", genderValidation)
	v.validation.RegisterValidation("clower", containsLowerCaseValidation)
	v.validation.RegisterValidation("cupper", containsUpperCaseValidation)
	v.validation.RegisterValidation("cnumeric", containsNumeric)
	v.validation.RegisterValidation("csymbol", containsSymbl)
	v.validation.RegisterValidation("bank_number", bankAccountNumber)
	v.validation.RegisterValidation("bank_provider", bankProvider)
	v.validation.RegisterValidation("unique_addresses", uniqueAddresses)
	v.validation.RegisterValidation("unique_bank_accounts", uniqueBankAccounts)
	v.validation.RegisterValidation("unique_shippings", uniqueShippings)
	v.validation.RegisterValidation("unique_colors", uniqueColors)
	v.validation.RegisterValidation("unique_sizes", uniqueSizes)
	v.validation.RegisterValidation("unique_items", uniqueItems)
	v.validation.RegisterValidation("category", registeredCategories)
	v.validation.RegisterValidation("unique_etalase", uniqueEtalase)
}

func uniqueEtalase(fl validator.FieldLevel) bool {
	counter := map[string]int{}
	etalase := fl.Field().Interface().([]string)
	for _, etalaseName := range etalase {
		if counter[etalaseName] != 0 {
			return false
		} else {
			counter[etalaseName] = 1
		}
	}

	return true
}

func registeredCategories(fl validator.FieldLevel) bool {
	categories := map[string]map[string]map[string]bool{
		"elektronik": map[string]map[string]bool{
			"dapur": map[string]bool{
				"blender":        true,
				"juicer":         true,
				"kompor listrik": true,
				"kulkas":         true,
				"microwave":      true,
				"mixer":          true,
			},
			"kantor": map[string]bool{
				"mesin fax":         true,
				"mesin fotocopy":    true,
				"mesin hitung uang": true,
				"mesin kasir":       true,
			},
			"rumah": map[string]bool{
				"mesin cuci": true,
				"setrika":    true,
			},
			"lainnya": map[string]bool{
				"lain-lain": true,
			},
		},
		"komputer & laptop": map[string]map[string]bool{
			"komputer & laptop": map[string]bool{
				"komputer": true,
				"laptop":   true,
			},
			"aksesoris": map[string]bool{
				"keyboard":  true,
				"mouse":     true,
				"lain lain": true,
			},
		},
		"handphone & tablet": map[string]map[string]bool{
			"handphone & tablet": map[string]bool{
				"handpone": true,
				"tablet":   true,
			},
			"aksesoris": map[string]bool{
				"charger":    true,
				"casing":     true,
				"anti gores": true,
			},
		},
		"fashion pria": map[string]map[string]bool{
			"atasan pria": map[string]bool{
				"kaos pria":      true,
				"kaos polo pria": true,
				"kemeja pria":    true,
			},
			"celana pria": map[string]bool{
				"celana jeans pria":  true,
				"celana pendek pria": true,
				"celana chino pria":  true,
			},
		},
		"fashion wanita": map[string]map[string]bool{
			"atasan wanita": map[string]bool{
				"kaos wanita":      true,
				"kaos polo wanita": true,
				"kemeja wanita":    true,
			},
			"celana wanita": map[string]bool{
				"celana jeans wanita":  true,
				"celana pendek wanita": true,
				"celana chino wanita":  true,
			},
		},
	}

	category := fl.Field().Interface().(domain.Category)
	if categories[category.Top] != nil && categories[category.Top][category.SecondSub] != nil && categories[category.Top][category.SecondSub][category.ThirdSub] {
		return true
	}

	return false
}

func uniqueBankAccounts(fl validator.FieldLevel) bool {
	accounts := fl.Field().Interface().([]domain.BankAccount)
	counter := map[string]int{}
	for _, account := range accounts {
		key := fmt.Sprintf("%s-%s", account.Number, account.BankCode)
		if counter[key] != 0 {
			return false
		}

		counter[key] = 1

	}
	return true
}

func uniqueAddresses(fl validator.FieldLevel) bool {
	addresses := fl.Field().Interface().([]domain.Address)
	counter := map[string]int{}
	for _, add := range addresses {
		key := fmt.Sprintf("%s-%s-%s", add.City, add.Street, add.Number)
		if counter[key] != 0 {
			return false
		}

		counter[key] = 1

	}
	return true
}

func uniqueShippings(fl validator.FieldLevel) bool {
	shippings := fl.Field().Interface().([]domain.ShippingProvider)
	counter := map[string]int{}
	for _, shipping := range shippings {
		key := fmt.Sprintf("%s-%s", shipping.ID, shipping.Name)
		if counter[key] != 0 {
			return false
		}

		counter[key] = 1

	}
	return true
}

func uniqueColors(fl validator.FieldLevel) bool {
	colors := fl.Field().Interface().([]string)
	counter := map[string]int{}
	for _, color := range colors {
		key := fmt.Sprintf("%s", color)
		if counter[key] != 0 {
			return false
		}

		counter[key] = 1

	}
	return true
}

func uniqueSizes(fl validator.FieldLevel) bool {
	sizes := fl.Field().Interface().([]string)
	counter := map[string]int{}
	for _, size := range sizes {
		key := fmt.Sprintf("%s", size)
		if counter[key] != 0 {
			return false
		}

		counter[key] = 1

	}
	return true
}

func uniqueItems(fl validator.FieldLevel) bool {
	items := fl.Field().Interface().([]domain.Item)
	counter := map[string]int{}
	for _, item := range items {
		key := fmt.Sprintf("%s", item.Product.ID)
		if counter[key] != 0 {
			return false
		}

		counter[key] = 1

	}
	return true
}

func bankProvider(fl validator.FieldLevel) bool {
	const (
		BCA     = "014"
		MANDIRI = "008"
		BNI     = "009"
		BRI     = "002"
	)

	val := fl.Field().String()
	if val != BCA && val != MANDIRI && val != BNI && val != BRI {
		return false
	}

	return true
}

func bankAccountNumber(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) != 10 {
		return false
	}

	return true
}

func phoneValidation(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	var minLengthPhoneNumber int = 7
	if len(val) < minLengthPhoneNumber {
		return false
	}

	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

	if isValid := re.MatchString(val); !isValid {
		return false
	}

	return true
}

func genderValidation(fl validator.FieldLevel) bool {

	if isValid := (fl.Field().String() != "M") && (fl.Field().String() != "F"); isValid {
		return false
	}

	return true
}

func containsLowerCaseValidation(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("^(?:.*[a-z].*)$")

	if match := re.MatchString(fl.Field().String()); !match {
		return false
	}
	return true
}

func containsUpperCaseValidation(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("^(?:.*[A-Z].*)$")

	if match := re.MatchString(fl.Field().String()); !match {
		return false
	}
	return true
}

func containsNumeric(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("^(?:.*[0-9].*)$")

	if match := re.MatchString(fl.Field().String()); !match {
		return false
	}
	return true
}

func containsSymbl(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("^(?:.*[-!$%#^&*()_+|~=`{}\\[\\]:\";'<>?,.\\/].*)$")

	if match := re.MatchString(fl.Field().String()); !match {
		return false
	}
	return true
}
