package v1

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/golang/glog"
	numv1 "github.com/tuanden0/generate_number/proto/gen/go/v1/number"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	vd *validator.Validate
)

type validate struct {
	Validate *validator.Validate
	Trans    ut.Translator
}

func NewValidate() validate {
	vd = validator.New()
	return validate{
		Validate: vd,
	}
}

func (v *validate) Init() error {
	v.initValidate()
	errTranslator := v.initTranslator()
	if errTranslator != nil {
		return errTranslator
	}

	return nil
}

func (v *validate) initValidate() {}

func (v *validate) initTranslator() error {

	en := en.New()
	uni := ut.New(en, en)

	trans, found := uni.GetTranslator("en")
	if !found {
		return fmt.Errorf("translator not found")
	}

	v.Trans = trans

	if err := en_translations.RegisterDefaultTranslations(v.Validate, trans); err != nil {
		return err
	}

	// Get lower-case field name
	v.Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom error msg translator
	v.registerCustomTranslate(trans, "required", "{0} is a required field")

	return nil
}

func (v *validate) parseError(err error) error {

	if _, ok := err.(*validator.InvalidValidationError); ok {
		glog.Error(err)
		return nil
	}

	errs := err.(validator.ValidationErrors)
	st := status.New(codes.InvalidArgument, "invalid_argument")
	br := &errdetails.BadRequest{}

	for _, e := range errs {
		v := &errdetails.BadRequest_FieldViolation{
			Field:       e.Field(),
			Description: e.Translate(v.Trans),
		}
		br.FieldViolations = append(br.FieldViolations, v)
	}

	st, err = st.WithDetails(br)
	if err != nil {
		glog.Errorf("Unexpected error attaching metadata %v", err)
		return err
	}

	return st.Err()
}

func (v *validate) registerCustomTranslate(trans ut.Translator, tag, msg string) {

	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, msg, true)
	}

	translationFn := func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.StructField())
		return t
	}

	v.Validate.RegisterTranslation(tag, trans, registerFn, translationFn)
}

// Validate service functions
func (v *validate) Get(ctx context.Context, in *numv1.NumberServiceGetRequest) error {

	if err := v.Validate.Struct(in); err != nil {
		return v.parseError(err)
	}

	return nil
}
