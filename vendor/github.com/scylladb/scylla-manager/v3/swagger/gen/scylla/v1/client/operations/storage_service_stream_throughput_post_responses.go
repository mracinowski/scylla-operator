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

// StorageServiceStreamThroughputPostReader is a Reader for the StorageServiceStreamThroughputPost structure.
type StorageServiceStreamThroughputPostReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *StorageServiceStreamThroughputPostReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewStorageServiceStreamThroughputPostOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewStorageServiceStreamThroughputPostDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewStorageServiceStreamThroughputPostOK creates a StorageServiceStreamThroughputPostOK with default headers values
func NewStorageServiceStreamThroughputPostOK() *StorageServiceStreamThroughputPostOK {
	return &StorageServiceStreamThroughputPostOK{}
}

/*
StorageServiceStreamThroughputPostOK handles this case with default header values.

Success
*/
type StorageServiceStreamThroughputPostOK struct {
}

func (o *StorageServiceStreamThroughputPostOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewStorageServiceStreamThroughputPostDefault creates a StorageServiceStreamThroughputPostDefault with default headers values
func NewStorageServiceStreamThroughputPostDefault(code int) *StorageServiceStreamThroughputPostDefault {
	return &StorageServiceStreamThroughputPostDefault{
		_statusCode: code,
	}
}

/*
StorageServiceStreamThroughputPostDefault handles this case with default header values.

internal server error
*/
type StorageServiceStreamThroughputPostDefault struct {
	_statusCode int

	Payload *models.ErrorModel
}

// Code gets the status code for the storage service stream throughput post default response
func (o *StorageServiceStreamThroughputPostDefault) Code() int {
	return o._statusCode
}

func (o *StorageServiceStreamThroughputPostDefault) GetPayload() *models.ErrorModel {
	return o.Payload
}

func (o *StorageServiceStreamThroughputPostDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorModel)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func (o *StorageServiceStreamThroughputPostDefault) Error() string {
	return fmt.Sprintf("agent [HTTP %d] %s", o._statusCode, strings.TrimRight(o.Payload.Message, "."))
}
