package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fieldValidation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func InvalidArgument(validation []*errdetails.BadRequest_FieldViolation) error {
	badrequest := &errdetails.BadRequest{FieldViolations: validation}
	statusinvalid := status.New(codes.InvalidArgument, "invalid parameter")

	status, err := statusinvalid.WithDetails(badrequest)
	if err != nil {
		return statusinvalid.Err()
	}
	return status.Err()
}
