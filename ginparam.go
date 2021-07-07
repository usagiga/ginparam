package ginparam

import (
	"golang.org/x/xerrors"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Read `out` members from (*gin.Context).Params.
// `query` tagged members(without `query:"-"`) will be read, and others will be ignored.
//
// ## Compatible types
// int, string, bool, their slice and struct contained them.
//
// Slice is separated its own value with ",".
// Unfortunately, it has no way to escape ",", pls take care!
func Read(ctx *gin.Context, out interface{}) (err error) {
	// Drill down pointer or interface
	outValue := reflect.ValueOf(out)
	outType := outValue.Type()
	outKind := outType.Kind()
	for outKind == reflect.Ptr || outKind == reflect.Interface {
		outValue = outValue.Elem()
		outType = outValue.Type()
		outKind = outType.Kind()
	}

	// Now, is it struct?
	if outKind != reflect.Struct {
		return xerrors.Errorf("passed incompatible type(%s): `out` must be assignable struct", outKind)
	}

	// For all struct fields,
	for i := 0; i < outType.NumField(); i++ {
		fieldVal := outValue.Field(i)
		fieldType := outType.Field(i)
		fieldKind := fieldType.Type.Kind()

		// Can set?
		if !fieldVal.CanSet() {
			return xerrors.Errorf("passed incompatible type(%s): `out`'s field(%s) must be assignable", fieldKind, fieldType.Name)
		}

		// Get `query` struct tag
		paramKey := fieldType.Tag.Get("query")

		// If "-", skip it
		if paramKey == "-" {
			continue
		}

		// If nested, process it recursive
		if fieldKind == reflect.Struct {
			err = Read(ctx, fieldVal.Addr().Interface())
			if err != nil {
				return xerrors.Errorf("error raised in nested value: %w", err)
			}
			continue
		}

		// If `query` struct tag not set, skip it
		if paramKey == "" {
			continue
		}

		// Look up param. If can't, skip it
		paramVal, ok := ctx.GetQuery(paramKey)
		if !ok || paramVal == "" {
			continue
		}

		// Set the value
		switch fieldKind {
		case reflect.Slice:
			arrVal := fieldVal
			arrType := fieldType.Type
			arrElemType := arrType.Elem()
			arrElemKind := arrElemType.Kind()

			paramVals := strings.Split(paramVal, ",")
			switch arrElemKind {
			case reflect.String:
				newFieldVal := reflect.ValueOf(paramVals)
				arrVal.Set(newFieldVal)
			case reflect.Int:
				s := make([]int, 0)
				for _, v := range paramVals {
					iv, err := strconv.Atoi(v)
					if err != nil {
						return xerrors.Errorf("can't cast int value: %w", err)
					}
					s = append(s, iv)
				}

				newFieldVal := reflect.ValueOf(s)
				arrVal.Set(newFieldVal)
			case reflect.Bool:
				s := make([]bool, 0)
				for _, v := range paramVals {
					bv := v == "true"
					s = append(s, bv)
				}

				newFieldVal := reflect.ValueOf(s)
				arrVal.Set(newFieldVal)
			}
		case reflect.String:
			newFieldVal := paramVal
			fieldVal.SetString(newFieldVal)
		case reflect.Int:
			newFieldVal, err := strconv.Atoi(paramVal)
			if err != nil {
				return xerrors.Errorf("can't cast int value: %s", err)
			}

			fieldVal.SetInt(int64(newFieldVal))
		case reflect.Bool:
			newFieldVal := paramVal == "true"
			fieldVal.SetBool(newFieldVal)
		}
	}

	return nil
}
