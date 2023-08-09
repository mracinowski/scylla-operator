// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Mapper mapper
//
// # A key value mapping
//
// swagger:model mapper
type Mapper struct {

	// The key
	Key string `json:"key,omitempty"`

	// The value
	Value string `json:"value,omitempty"`
}

// Validate validates this mapper
func (m *Mapper) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this mapper based on context it is used
func (m *Mapper) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Mapper) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Mapper) UnmarshalBinary(b []byte) error {
	var res Mapper
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}