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

// ColumnFamilyMetricsAllMemtablesLiveDataSizeGetReader is a Reader for the ColumnFamilyMetricsAllMemtablesLiveDataSizeGet structure.
type ColumnFamilyMetricsAllMemtablesLiveDataSizeGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ColumnFamilyMetricsAllMemtablesLiveDataSizeGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewColumnFamilyMetricsAllMemtablesLiveDataSizeGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewColumnFamilyMetricsAllMemtablesLiveDataSizeGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewColumnFamilyMetricsAllMemtablesLiveDataSizeGetOK creates a ColumnFamilyMetricsAllMemtablesLiveDataSizeGetOK with default headers values
func NewColumnFamilyMetricsAllMemtablesLiveDataSizeGetOK() *ColumnFamilyMetricsAllMemtablesLiveDataSizeGetOK {
	return &ColumnFamilyMetricsAllMemtablesLiveDataSizeGetOK{}
}

/*
ColumnFamilyMetricsAllMemtablesLiveDataSizeGetOK handles this case with default header values.

Success
*/
type ColumnFamilyMetricsAllMemtablesLiveDataSizeGetOK struct {
	Payload interface{}
}

func (o *ColumnFamilyMetricsAllMemtablesLiveDataSizeGetOK) GetPayload() interface{} {
	return o.Payload
}

func (o *ColumnFamilyMetricsAllMemtablesLiveDataSizeGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewColumnFamilyMetricsAllMemtablesLiveDataSizeGetDefault creates a ColumnFamilyMetricsAllMemtablesLiveDataSizeGetDefault with default headers values
func NewColumnFamilyMetricsAllMemtablesLiveDataSizeGetDefault(code int) *ColumnFamilyMetricsAllMemtablesLiveDataSizeGetDefault {
	return &ColumnFamilyMetricsAllMemtablesLiveDataSizeGetDefault{
		_statusCode: code,
	}
}

/*
ColumnFamilyMetricsAllMemtablesLiveDataSizeGetDefault handles this case with default header values.

internal server error
*/
type ColumnFamilyMetricsAllMemtablesLiveDataSizeGetDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the column family metrics all memtables live data size get default response
func (o *ColumnFamilyMetricsAllMemtablesLiveDataSizeGetDefault) Code() int {
	return o._statusCode
}

func (o *ColumnFamilyMetricsAllMemtablesLiveDataSizeGetDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *ColumnFamilyMetricsAllMemtablesLiveDataSizeGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *ColumnFamilyMetricsAllMemtablesLiveDataSizeGetDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
