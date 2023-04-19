package validate

import (
	"log"
	"regexp"

	// "github.com/go-playground/validator/v10"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	id_translations "gopkg.in/go-playground/validator.v9/translations/id"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
)

func IsFloat(fl validator.FieldLevel) bool {
	regexExp := regexp.MustCompile(`^(\d|[1-9]+\d*|\.\d+|0\.\d+|[1-9]+\d*\.\d+)$`)
	checking := regexExp.MatchString(fl.Field().String())

	log.Println("String float : ", fl.Field().String())
	log.Println("Result float : ", checking)

	return checking
}

func IsDateTime(fl validator.FieldLevel) bool {
	regexExp := regexp.MustCompile(`^((19|20)\d\d)[-](0?[1-9]|1[012])[-](0?[1-9]|[12][0-9]|3[01])$`)
	return regexExp.MatchString(fl.Field().String())
}

func IsTextAlphaNumSpecial(fl validator.FieldLevel) bool {
	regexExp := regexp.MustCompile(`^[a-zA-Z]+\-[0-9]+$`)
	return regexExp.MatchString(fl.Field().String())
}

func IsText(fl validator.FieldLevel) bool {
	regexExp := regexp.MustCompile(`^[a-z A-Z \- ]*$`)
	return regexExp.MatchString(fl.Field().String())
}

func IsName(fl validator.FieldLevel) bool {
	regexExp := regexp.MustCompile(`^[a-z A-Z]*$`)
	return regexExp.MatchString(fl.Field().String())
}

func IsAddress(fl validator.FieldLevel) bool {
	regexExp := regexp.MustCompile(`^[a-z A-Z 0-9\. ]*$`)
	checking := regexExp.MatchString(fl.Field().String())

	log.Println("String Value :", fl.Field().String())
	log.Println("String Result :", checking)
	return checking
}

func IsNumberMoreThanZero(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile("^[1-90-9]*$")
	return reg.MatchString(fl.Field().String())
}

func AddCustomeErrorMessage(Tag string, messageError string, validate *validator.Validate, trans ut.Translator) {
	validate.RegisterTranslation(Tag, trans, func(ut ut.Translator) error {
		return ut.Add(Tag, "{0} "+messageError, true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		var t string
		switch fe.Field() {
		case "UserId", "PrakarsaId", "PartnershipMemberId", "PipelineId", "EstimationLoan", "RecordCount", "StartIndex":
			t, _ = ut.T(Tag, "Maaf Parameter "+fe.Field())
		default:
			t, _ = ut.T(Tag, fe.Field())
		}

		return t
	})
}

func CustomeValidateRequest(model interface{}) (status int64, errMsg string) {
	log.Println("MASUK")
	status = 200
	validate = validator.New()

	validate.RegisterValidation("name", IsName)
	validate.RegisterValidation("text", IsText)
	validate.RegisterValidation("float", IsFloat)
	validate.RegisterValidation("dateString", IsDateTime)
	validate.RegisterValidation("numnonzero", IsNumberMoreThanZero)
	validate.RegisterValidation("alnumspecial", IsTextAlphaNumSpecial)
	validate.RegisterValidation("address", IsAddress)

	err := validate.Struct(model)
	log.Println("Model : ", model)
	id := id.New()
	uni = ut.New(id, id)
	trans, _ := uni.GetTranslator("id")

	id_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		for _, errFiled := range err.(validator.ValidationErrors) {
			switch errFiled.Tag() {
			case "required":
				AddCustomeErrorMessage(errFiled.Tag(), "Belum di isi", validate, trans)

				status = 400
				errMsg = errFiled.Translate(trans)

				return status, errMsg

			case "numeric":
				AddCustomeErrorMessage(errFiled.Tag(), "Harus Berupa Numeric", validate, trans)

				status = 400
				errMsg = errFiled.Translate(trans)

				return status, errMsg

			case "number":
				AddCustomeErrorMessage(errFiled.Tag(), "Harus Berupa Angka dan tidak kecil dari 0.", validate, trans)

				status = 400
				errMsg = errFiled.Translate(trans)

				return status, errMsg

			case "name":
				AddCustomeErrorMessage(errFiled.Tag(), "Harus Berupa huruf", validate, trans)

				status = 400
				errMsg = errFiled.Translate(trans)

				return status, errMsg

			case "float":
				AddCustomeErrorMessage(errFiled.Tag(), "Harus berupa Angka atau Angka Decimal", validate, trans)

				status = 400
				errMsg = errFiled.Translate(trans)

				return status, errMsg

			case "text":
				AddCustomeErrorMessage(errFiled.Tag(), "Harus berupa huruf", validate, trans)

				status = 400
				errMsg = errFiled.Translate(trans)

				return status, errMsg

			case "address":
				AddCustomeErrorMessage(errFiled.Tag(), "Tidak boleh memuat karakter spesial", validate, trans)

				status = 400
				errMsg = errFiled.Translate(trans)

				return status, errMsg

			case "alnumspecial":
				AddCustomeErrorMessage(errFiled.Tag(), "Harus berupa huruf, karakter spesial (-), dan angka. Contoh: (abc-123)", validate, trans)

				status = 400
				errMsg = errFiled.Translate(trans)

				return status, errMsg

			case "dateString":
				AddCustomeErrorMessage(errFiled.Tag(), "Format tanggal tidak sesuai. (Format : Y-m-d)", validate, trans)

				status = 400
				errMsg = errFiled.Translate(trans)

				return status, errMsg

			case "numnonzero":
				AddCustomeErrorMessage(errFiled.Tag(), "Harus Berupa Angka dan Tidak Sama dengan 0", validate, trans)

				status = 400
				errMsg = errFiled.Translate(trans)

				return status, errMsg

			case "ne": // ne = not equal
				if errFiled.Field() == "UserId" || errFiled.Field() == "PrakarsaId" || errFiled.Field() == "PartnershipMemberId" || errFiled.Field() == "PipelineId" {
					AddCustomeErrorMessage(errFiled.Tag(), "Tidak Ditemukan", validate, trans)

					status = 400
					errMsg = errFiled.Translate(trans)

					return status, errMsg
				}

			}
		}

	}
	return status, errMsg
}
