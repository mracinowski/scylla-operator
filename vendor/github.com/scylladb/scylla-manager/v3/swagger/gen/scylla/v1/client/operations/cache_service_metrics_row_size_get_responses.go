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

// CacheServiceMetricsRowSizeGetReader is a Reader for the CacheServiceMetricsRowSizeGet structure.
type CacheServiceMetricsRowSizeGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CacheServiceMetricsRowSizeGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCacheServiceMetricsRowSizeGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewCacheServiceMetricsRowSizeGetDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCacheServiceMetricsRowSizeGetOK creates a CacheServiceMetricsRowSizeGetOK with default headers values
func NewCacheServiceMetricsRowSizeGetOK() *CacheServiceMetricsRowSizeGetOK {
	return &CacheServiceMetricsRowSizeGetOK{}
}

/*
CacheServiceMetricsRowSizeGetOK handles this case with default header values.

Success
*/
type CacheServiceMetricsRowSizeGetOK struct {
	Payload interface{}
}

func (o *CacheServiceMetricsRowSizeGetOK) GetPayload() interface{} {
	return o.Payload
}

func (o *CacheServiceMetricsRowSizeGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCacheServiceMetricsRowSizeGetDefault creates a CacheServiceMetricsRowSizeGetDefault with default headers values
func NewCacheServiceMetricsRowSizeGetDefault(code int) *CacheServiceMetricsRowSizeGetDefault {
	return &CacheServiceMetricsRowSizeGetDefault{
		_statusCode: code,
	}
}

/*
CacheServiceMetricsRowSizeGetDefault handles this case with default header values.

internal server error
*/
type CacheServiceMetricsRowSizeGetDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the cache service metrics row size get default response
func (o *CacheServiceMetricsRowSizeGetDefault) Code() int {
	return o._statusCode
}

func (o *CacheServiceMetricsRowSizeGetDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *CacheServiceMetricsRowSizeGetDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *CacheServiceMetricsRowSizeGetDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
