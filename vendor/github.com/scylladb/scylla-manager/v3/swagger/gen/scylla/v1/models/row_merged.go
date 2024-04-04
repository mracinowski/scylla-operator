// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// RowMerged row_merged
//
// # A row merged information
//
// swagger:model row_merged
type RowMerged struct {

	// The number of sstable
	Key int32 `json:"key,omitempty"`

	// The number or row compacted
	Value interface{} `json:"value,omitempty"`
}

// Validate validates this row merged
func (m *RowMerged) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RowMerged) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RowMerged) UnmarshalBinary(b []byte) error {
	var res RowMerged
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
