package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func findComma(s string) int {
	return strings.Index(s, ",")
}

func GenerateAllowedFields(structType interface{}) map[string]bool {
	fields := make(map[string]bool)
	t := reflect.TypeOf(structType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			name := jsonTag
			if commaIdx := findComma(jsonTag); commaIdx != -1 {
				name = jsonTag[:commaIdx]
			}
			fields[name] = true
		}
	}
	return fields
}

func BindAndValidateBody(ctx *fiber.Ctx, dest interface{}, allowedFields map[string]bool, validate *validator.Validate) error {
	rawBody := ctx.Body()
	var tempMap map[string]interface{}
	if err := json.Unmarshal(rawBody, &tempMap); err != nil {
		return fmt.Errorf("invalid JSON format: %v", err)
	}
	if len(tempMap) == 0 {
		return fmt.Errorf("at least one field must be provided")
	}
	for field := range tempMap {
		if !allowedFields[field] {
			return fmt.Errorf("field '%s' is not allowed", field)
		}
	}

	decoder := json.NewDecoder(bytes.NewReader(rawBody))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dest); err != nil {
		return fmt.Errorf("invalid or unknown field: %v", err)
	}

	if err := validate.Struct(dest); err != nil {
		var combinedErrMsg string
		for _, e := range err.(validator.ValidationErrors) {
			combinedErrMsg += fmt.Sprintf("field '%s' failed on '%s' validation; ", e.Field(), e.Tag())
		}
		return fmt.Errorf("validation failed: %s", combinedErrMsg)
	}
	return nil
}
