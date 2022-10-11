package aws

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	_ "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewClient(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(cfg), nil
}

type HTTPResponse struct {
	Code     int
	Response Response
}

type Response struct {
	Type     string  `json:"type,omitempty"`
	Content  string  `json:"content,omitempty"`
	Children []Child `json:"children,omitempty"`
	Error    string  `json:"error,omitempty"`
}

type Child struct {
	Key  string
	Size int64
}

type ObjectGetFunc func() (*s3.GetObjectOutput, error)

func ToObjectGetFunc(ctx context.Context, client *s3.Client, bucket, path string) ObjectGetFunc {
	return func() (*s3.GetObjectOutput, error) {
		return client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(path),
		})
	}
}

func GetObject(objGetFn ObjectGetFunc) (r HTTPResponse) {
	output, err := objGetFn()

	r.Code = http.StatusInternalServerError
	if err != nil {
		r.Response.Error = err.Error()
		return
	}

	var b []byte
	if b, err = io.ReadAll(output.Body); err != nil {
		r.Response.Error = err.Error()
		return
	}

	r.Code = http.StatusOK
	r.Response.Type = "object"
	r.Response.Content = string(b)
	return
}

type ObjectListFunc func() (*s3.ListObjectsV2Output, error)

func ToObjectListFunc(ctx context.Context, client *s3.Client, bucket, path string) ObjectListFunc {
	return func() (*s3.ListObjectsV2Output, error) {
		return client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(path),
		})
	}
}

func ListObjects(objListFn ObjectListFunc) (r HTTPResponse) {
	output, err := objListFn()

	r.Code = http.StatusInternalServerError
	if err != nil {
		r.Response = Response{Error: err.Error()}
		return
	}

	r.Response.Type = "folder"
	for _, object := range output.Contents {
		r.Response.Children = append(r.Response.Children, Child{
			Key:  *object.Key,
			Size: object.Size,
		})
	}

	r.Code = http.StatusOK
	return
}

func PathDecider(path string, listFunc ObjectListFunc, getFunc ObjectGetFunc) func() HTTPResponse {
	if strings.HasSuffix(path, "/") {
		return func() HTTPResponse {
			return ListObjects(listFunc)
		}
	}

	return func() HTTPResponse {
		return GetObject(getFunc)
	}
}
