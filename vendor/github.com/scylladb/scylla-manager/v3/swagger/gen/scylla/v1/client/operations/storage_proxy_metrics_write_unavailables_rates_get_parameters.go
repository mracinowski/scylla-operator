// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewStorageProxyMetricsWriteUnavailablesRatesGetParams creates a new StorageProxyMetricsWriteUnavailablesRatesGetParams object
// with the default values initialized.
func NewStorageProxyMetricsWriteUnavailablesRatesGetParams() *StorageProxyMetricsWriteUnavailablesRatesGetParams {

	return &StorageProxyMetricsWriteUnavailablesRatesGetParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewStorageProxyMetricsWriteUnavailablesRatesGetParamsWithTimeout creates a new StorageProxyMetricsWriteUnavailablesRatesGetParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewStorageProxyMetricsWriteUnavailablesRatesGetParamsWithTimeout(timeout time.Duration) *StorageProxyMetricsWriteUnavailablesRatesGetParams {

	return &StorageProxyMetricsWriteUnavailablesRatesGetParams{

		timeout: timeout,
	}
}

// NewStorageProxyMetricsWriteUnavailablesRatesGetParamsWithContext creates a new StorageProxyMetricsWriteUnavailablesRatesGetParams object
// with the default values initialized, and the ability to set a context for a request
func NewStorageProxyMetricsWriteUnavailablesRatesGetParamsWithContext(ctx context.Context) *StorageProxyMetricsWriteUnavailablesRatesGetParams {

	return &StorageProxyMetricsWriteUnavailablesRatesGetParams{

		Context: ctx,
	}
}

// NewStorageProxyMetricsWriteUnavailablesRatesGetParamsWithHTTPClient creates a new StorageProxyMetricsWriteUnavailablesRatesGetParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewStorageProxyMetricsWriteUnavailablesRatesGetParamsWithHTTPClient(client *http.Client) *StorageProxyMetricsWriteUnavailablesRatesGetParams {

	return &StorageProxyMetricsWriteUnavailablesRatesGetParams{
		HTTPClient: client,
	}
}

/*
StorageProxyMetricsWriteUnavailablesRatesGetParams contains all the parameters to send to the API endpoint
for the storage proxy metrics write unavailables rates get operation typically these are written to a http.Request
*/
type StorageProxyMetricsWriteUnavailablesRatesGetParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the storage proxy metrics write unavailables rates get params
func (o *StorageProxyMetricsWriteUnavailablesRatesGetParams) WithTimeout(timeout time.Duration) *StorageProxyMetricsWriteUnavailablesRatesGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the storage proxy metrics write unavailables rates get params
func (o *StorageProxyMetricsWriteUnavailablesRatesGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the storage proxy metrics write unavailables rates get params
func (o *StorageProxyMetricsWriteUnavailablesRatesGetParams) WithContext(ctx context.Context) *StorageProxyMetricsWriteUnavailablesRatesGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the storage proxy metrics write unavailables rates get params
func (o *StorageProxyMetricsWriteUnavailablesRatesGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the storage proxy metrics write unavailables rates get params
func (o *StorageProxyMetricsWriteUnavailablesRatesGetParams) WithHTTPClient(client *http.Client) *StorageProxyMetricsWriteUnavailablesRatesGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the storage proxy metrics write unavailables rates get params
func (o *StorageProxyMetricsWriteUnavailablesRatesGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *StorageProxyMetricsWriteUnavailablesRatesGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
