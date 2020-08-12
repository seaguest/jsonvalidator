package jsonvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
)

const (
	validateTypeInt               = "int"
	validateTypeFloat             = "float"
	validateTypeString            = "string"
	validateTypeList              = "list"
	validateTypeRegularExpression = "re"

	typeErrMsg = "tag [%s] expected [%s], but recieved [%+v]"
)

var vd *validator.Validate

func init() {
	vd = validator.New()
}

func Validate(src, tpl []byte) (err error) {
	var srcItf interface{}
	if err := json.Unmarshal(src, &srcItf); err != nil {
		return err
	}

	var tplItf interface{}
	if err := json.Unmarshal(tpl, &tplItf); err != nil {
		return err
	}
	return validate(srcItf, tplItf, "")
}

func validate(src, tpl interface{}, tag string) (err error) {
	if src == nil || tpl == nil {
		return nil

	}

	switch sv := src.(type) {
	case []interface{}:
		switch tv := tpl.(type) {
		case []interface{}:
			if len(sv) != len(tv) {
				return fmt.Errorf("tag [%s] not match", tag)
			}

			for idx, v := range sv {
				err = validate(v, tv[idx], fmt.Sprintf("%s[%d]", tag, idx))
				if err != nil {
					return
				}
			}
		case string:
			return validateVar(tag, tv, sv)
		default:
			return fmt.Errorf("tag [%s] type inconsistent", tag)
		}
	case map[string]interface{}:
		switch tv := tpl.(type) {
		case map[string]interface{}:
			for k, v := range sv {
				err = validate(v, tv[k], k)
				if err != nil {
					return
				}
			}
		default:
			return fmt.Errorf("tag [%s] type inconsistent", tag)
		}
	default:
		return validateVar(tag, tpl.(string), sv)
	}
	return
}

func validateVar(tag, exp string, v interface{}) (err error) {
	ss := strings.Split(exp, "|")
	if len(ss) != 2 {
		return fmt.Errorf("invalid exp [%s]", exp)
	}
	typ := ss[0]
	tags := ss[1]

	wrapErr := func(e error) error {
		if e != nil {
			return fmt.Errorf(strings.Replace(e.Error(), "''", fmt.Sprintf("'%s'", tag), -1))
		}
		return nil
	}

	switch typ {
	case validateTypeInt:
		if reflect.TypeOf(v) != reflect.TypeOf(float64(0)) {
			return fmt.Errorf(typeErrMsg, tag, typ, reflect.TypeOf(v))
		}
		return wrapErr(vd.Var(v, tags))
	case validateTypeFloat:
		if reflect.TypeOf(v) != reflect.TypeOf(float64(0)) {
			return fmt.Errorf(typeErrMsg, tag, typ, reflect.TypeOf(v))
		}
		return wrapErr(vd.Var(v, tags))
	case validateTypeString:
		if reflect.TypeOf(v) != reflect.TypeOf("") {
			return fmt.Errorf(typeErrMsg, tag, typ, reflect.TypeOf(v))
		}
		return wrapErr(vd.Var(v, tags))
	case validateTypeList:
		return wrapErr(vd.Var(v, tags))
	case validateTypeRegularExpression:
		exp := regexp.MustCompile(tags)
		macth := exp.Match([]byte(fmt.Sprint(v)))
		if !macth {
			return fmt.Errorf(typeErrMsg, tag, tags, fmt.Sprint(v))
		}
	}
	return
}
