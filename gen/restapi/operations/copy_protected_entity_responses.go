// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/vmware-tanzu/astrolabe/gen/models"
)

// CopyProtectedEntityAcceptedCode is the HTTP code returned for type CopyProtectedEntityAccepted
const CopyProtectedEntityAcceptedCode int = 202

/*CopyProtectedEntityAccepted Create in progress

swagger:response copyProtectedEntityAccepted
*/
type CopyProtectedEntityAccepted struct {

	/*
	  In: Body
	*/
	Payload *models.CreateInProgressResponse `json:"body,omitempty"`
}

// NewCopyProtectedEntityAccepted creates CopyProtectedEntityAccepted with default headers values
func NewCopyProtectedEntityAccepted() *CopyProtectedEntityAccepted {

	return &CopyProtectedEntityAccepted{}
}

// WithPayload adds the payload to the copy protected entity accepted response
func (o *CopyProtectedEntityAccepted) WithPayload(payload *models.CreateInProgressResponse) *CopyProtectedEntityAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the copy protected entity accepted response
func (o *CopyProtectedEntityAccepted) SetPayload(payload *models.CreateInProgressResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *CopyProtectedEntityAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(202)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
