package aws

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
)

func TestGetObject(t *testing.T) {
	type args struct {
		objGetFn ObjectGetFunc
	}
	tests := []struct {
		name  string
		args  args
		wantR HTTPResponse
	}{
		{
			name: "given error, expect error only response",
			args: args{
				objGetFn: func() (*s3.GetObjectOutput, error) {
					return &s3.GetObjectOutput{}, errors.New("generic failure")
				},
			},
			wantR: HTTPResponse{
				Code: http.StatusInternalServerError,
				Response: Response{
					Error: "generic failure",
				},
			},
		},
		{
			name: "given object at path, expect object detail response",
			args: args{
				objGetFn: func() (*s3.GetObjectOutput, error) {
					return &s3.GetObjectOutput{
						Body: io.NopCloser(strings.NewReader("hello world")),
					}, nil
				},
			},
			wantR: HTTPResponse{
				Code: http.StatusOK,
				Response: Response{
					Type:    "object",
					Content: "hello world",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR := GetObject(tt.args.objGetFn)
			assert.Equal(t, tt.wantR, gotR)
		})
	}
}

func TestListObjects(t *testing.T) {
	type args struct {
		objListFn ObjectListFunc
	}
	tests := []struct {
		name  string
		args  args
		wantR HTTPResponse
	}{
		{
			name: "given error, expect error only response",
			args: args{
				objListFn: func() (*s3.ListObjectsV2Output, error) {
					return &s3.ListObjectsV2Output{}, errors.New("generic failure")
				},
			},
			wantR: HTTPResponse{
				Code: http.StatusInternalServerError,
				Response: Response{
					Error: "generic failure",
				},
			},
		},
		{
			name: "given 0 objects, expect an empty list of children",
			args: args{
				objListFn: func() (*s3.ListObjectsV2Output, error) {
					return &s3.ListObjectsV2Output{
						Contents: []types.Object{},
					}, nil
				},
			},
			wantR: HTTPResponse{
				Code: http.StatusOK,
				Response: Response{
					Type: "folder",
					// TODO I believe we should be more explicit/forthright and return an empty array
					// To do that, we'd have to go down the path of implementing a custom json unmarshaller
					// Children: []Child{},
				},
			},
		},
		{
			name: "given 1 object, expect a list of 1 children",
			args: args{
				objListFn: func() (*s3.ListObjectsV2Output, error) {
					return &s3.ListObjectsV2Output{
						Contents: []types.Object{
							{
								Key:  aws.String("foo/bar"),
								Size: 123,
							},
						},
					}, nil
				},
			},
			wantR: HTTPResponse{
				Code: http.StatusOK,
				Response: Response{
					Type: "folder",
					Children: []Child{
						{
							Key:  "foo/bar",
							Size: 123,
						},
					},
				},
			},
		},
		{
			name: "given 3 objects, expect a list of 3 children",
			args: args{
				objListFn: func() (*s3.ListObjectsV2Output, error) {
					return &s3.ListObjectsV2Output{
						Contents: []types.Object{
							{
								Key:  aws.String("foo/bar"),
								Size: 123,
							},
							{
								Key:  aws.String("foo/biz"),
								Size: 456,
							},
							{
								Key:  aws.String("foo/baz"),
								Size: 789,
							},
						},
					}, nil
				},
			},
			wantR: HTTPResponse{
				Code: http.StatusOK,
				Response: Response{
					Type: "folder",
					Children: []Child{
						{
							Key:  "foo/bar",
							Size: 123,
						},
						{
							Key:  "foo/biz",
							Size: 456,
						},
						{
							Key:  "foo/baz",
							Size: 789,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR := ListObjects(tt.args.objListFn)
			assert.Equal(t, tt.wantR, gotR)
		})
	}
}
