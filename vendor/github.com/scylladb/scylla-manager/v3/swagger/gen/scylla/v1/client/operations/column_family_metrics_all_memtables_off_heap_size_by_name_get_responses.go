// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/scylladb/scylla-manager/v3/swagger/gen/scylla/v1/models"
)

// ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetReader is a Reader for the ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGet structure.
type ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetOK creates a ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetOK with default headers values
func NewColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetOK() *ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetOK {
	return &ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetOK{}
}

/*
ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetOK handles this case with default header values.

Success
*/
type ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetOK struct {
	Payload interface{}
}

func (o *ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetOK) GetPayload() interface{} {
	return o.Payload
}

func (o *ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetDefault creates a ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetDefault with default headers values
func NewColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetDefault(code int) *ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetDefault {
	return &ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetDefault{
		_statusCode: code,
	}
}

/*
ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetDefault handles this case with default header values.

internal server error
*/
type ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the column family metrics all memtables off heap size by name get default response
func (o *ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetDefault) Code() int {
	return o._statusCode
}

func (o *ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *ColumnFamilyMetricsAllMemtablesOffHeapSizeByNameGetDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
