package post

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestUpsertPostRequest_Validate_ValidRequest(t *testing.T) {
	req := &UpsertPostRequest{
		Title:   "Valid Title",
		Content: "This is a valid content with more than 10 characters",
	}
	err := req.Validate()
	assert.NoError(t, err)
}

func TestUpsertPostRequest_Validate_TitleTooShort(t *testing.T) {
	req := &UpsertPostRequest{
		Title:   "Hi", // less than 3 characters
		Content: "This is a valid content with more than 10 characters",
	}
	err := req.Validate()
	assert.Error(t, err)

	validationErr, ok := err.(validator.ValidationErrors)
	assert.True(t, ok)

	var errFields []string
	for _, e := range validationErr {
		errFields = append(errFields, e.Field())
	}
	assert.Contains(t, errFields, "Title")
}

func TestUpsertPostRequest_Validate_TitleTooLong(t *testing.T) {
	req := &UpsertPostRequest{
		Title:   "This is an extremely long title that exceeds the maximum allowed length of two hundred characters and should trigger a validation error because it is just too verbose and unnecessary for any practical purpose",
		Content: "This is a valid content with more than 10 characters",
	}
	err := req.Validate()
	assert.Error(t, err)

	validationErr, ok := err.(validator.ValidationErrors)
	assert.True(t, ok)

	var errFields []string
	for _, e := range validationErr {
		errFields = append(errFields, e.Field())
	}
	assert.Contains(t, errFields, "Title")
}

func TestUpsertPostRequest_Validate_EmptyTitle(t *testing.T) {
	req := &UpsertPostRequest{
		Title:   "", // empty title
		Content: "This is a valid content with more than 10 characters",
	}
	err := req.Validate()
	assert.Error(t, err)

	validationErr, ok := err.(validator.ValidationErrors)
	assert.True(t, ok)

	var errFields []string
	for _, e := range validationErr {
		errFields = append(errFields, e.Field())
	}
	assert.Contains(t, errFields, "Title")
}

func TestUpsertPostRequest_Validate_ContentTooShort(t *testing.T) {
	req := &UpsertPostRequest{
		Title:   "Valid Title",
		Content: "Short", // less than 10 characters
	}
	err := req.Validate()
	assert.Error(t, err)

	validationErr, ok := err.(validator.ValidationErrors)
	assert.True(t, ok)

	var errFields []string
	for _, e := range validationErr {
		errFields = append(errFields, e.Field())
	}
	assert.Contains(t, errFields, "Content")
}

func TestUpsertPostRequest_Validate_EmptyContent(t *testing.T) {
	req := &UpsertPostRequest{
		Title:   "Valid Title",
		Content: "", // empty content
	}
	err := req.Validate()
	assert.Error(t, err)

	validationErr, ok := err.(validator.ValidationErrors)
	assert.True(t, ok)

	var errFields []string
	for _, e := range validationErr {
		errFields = append(errFields, e.Field())
	}
	assert.Contains(t, errFields, "Content")
}

func TestUpsertPostRequest_Validate_TitleAtMinimumLength(t *testing.T) {
	req := &UpsertPostRequest{
		Title:   "Min", // exactly 3 characters
		Content: "This is a valid content with more than 10 characters",
	}
	err := req.Validate()
	assert.NoError(t, err, "Title of exactly 3 characters should be valid")
}

func TestUpsertPostRequest_Validate_TitleAtMaximumLength(t *testing.T) {
	longTitle := "This is a title that is exactly two hundred characters long and should pass validation because it meets the maximum length requirement without exceeding it"
	req := &UpsertPostRequest{
		Title:   longTitle,
		Content: "This is a valid content with more than 10 characters",
	}
	err := req.Validate()
	assert.NoError(t, err, "Title of exactly 200 characters should be valid")
}

func TestUpsertPostRequest_Validate_ContentAtMinimumLength(t *testing.T) {
	req := &UpsertPostRequest{
		Title:   "Valid Title",
		Content: "Ten Charc.", // exactly 10 characters
	}
	err := req.Validate()
	assert.NoError(t, err, "Content of exactly 10 characters should be valid")
}
