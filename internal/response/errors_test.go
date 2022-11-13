package response_test

import (
	"testing"

	"github.com/eduardohoraciosanto/users/internal/response"
	"github.com/stretchr/testify/assert"
)

func TestError_OK(t *testing.T) {
	err := response.Error{
		Code:        "someCode",
		Description: "Some Cool Description",
	}

	expected := "someCode:Some Cool Description"

	assert.Equal(t, expected, err.Error())
}
