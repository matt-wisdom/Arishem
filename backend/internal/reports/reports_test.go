package reports

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestInitS3(t *testing.T) {
	// Setup a mock S3 server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Mock Server Received: %s %s", r.Method, r.URL.String())
		if r.URL.Query().Get("location") != "" || strings.Contains(r.URL.RawQuery, "location") {
			w.Header().Set("Content-Type", "application/xml")
			w.Write([]byte(`<LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/"></LocationConstraint>`))
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	endpoint := strings.TrimPrefix(ts.URL, "http://")

	t.Setenv("S3_ENDPOINT", endpoint)
	t.Setenv("S3_BUCKET", "arishem-reports")
	t.Setenv("S3_USE_SSL", "false")

	err := InitS3()
	if err != nil {
		t.Fatalf("InitS3 failed: %v", err)
	}

	if s3Client == nil {
		t.Error("expected non-nil s3Client")
	}
}

func TestGetS3Client_and_BucketName(t *testing.T) {
	s3Client = &minio.Client{}
	bucketName = "test-bucket"

	if GetS3Client() != s3Client {
		t.Error("GetS3Client mismatch")
	}
	if GetBucketName() != "test-bucket" {
		t.Errorf("expected test-bucket, got %s", GetBucketName())
	}
}

func TestUploadReport_and_GetSignedURL(t *testing.T) {
	var putReceived bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("Mock Server Received: %s %s", r.Method, r.URL.String())
		if r.URL.Query().Get("location") != "" || strings.Contains(r.URL.RawQuery, "location") {
			w.Header().Set("Content-Type", "application/xml")
			w.Write([]byte(`<LocationConstraint xmlns="http://s3.amazonaws.com/doc/2006-03-01/"></LocationConstraint>`))
			return
		}
		if r.Method == "PUT" && strings.HasPrefix(r.URL.Path, "/arishem-reports/reports/") {
			putReceived = true
			w.Header().Set("ETag", `"mock-etag"`)
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	endpoint := strings.TrimPrefix(ts.URL, "http://")
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("key", "secret", ""),
		Secure: false,
	})
	if err != nil {
		t.Fatalf("failed to create minio client: %v", err)
	}

	s3Client = client
	bucketName = "arishem-reports"

	ctx := context.Background()
	err = UploadReport(ctx, "org_123", "report_123", "html", "text/html", "<html></html>")
	if err != nil {
		t.Errorf("UploadReport failed: %v", err)
	}
	if !putReceived {
		t.Error("expected PUT request to mock S3 server, but none was received")
	}

	url, err := GetSignedURL(ctx, "org_123", "report_123", "html")
	if err != nil {
		t.Errorf("GetSignedURL failed: %v", err)
	}
	if !strings.Contains(url, "reports/report_123.html") {
		t.Errorf("expected URL to contain object path, got %s", url)
	}
}