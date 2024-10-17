package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ContextKey string

const (
	UserContextKey ContextKey = "user"
)

type Response struct {
	Message string  `json:"message"`
	Error   *string `json:"error"`
	Data    any     `json:"data"`
}

func NewResponse(ctx *fiber.Ctx, message string, data any) error {
	return ctx.JSON(&Response{
		Message: message,
		Data:    data,
	})
}

func NewErrorResponse(ctx *fiber.Ctx, status int, err error) error {
	msg := err.Error()
	ctx.JSON(&Response{
		Error: &msg,
	})

	return ctx.SendStatus(status)
}

type SQLMap[K comparable, V any] map[K]V

// Scan scan value into Jsonb, implements sql.Scanner interface
func (v *SQLMap[K, V]) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := SQLMap[K, V]{}
	err := json.Unmarshal(bytes, &result)

	return err
}

// Value return json value, implement driver.Valuer interface
func (v SQLMap[K, V]) Value() (driver.Value, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}
