package miniostorage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/thdoikn/sihp-be/config"
	"github.com/thdoikn/sihp-be/pkg/storage"
)

var (
	allowedContentTypes = map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/webp": ".webp",
	}
)

type minioStorage struct {
	client    *minio.Client
	bucket    string
	publicURL string
}

func NewMinIOStorage(cfg *config.Config) (storage.KomoditasStorage, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	client, err := minio.New(cfg.MinIO.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: cfg.MinIO.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}

	publicURL := strings.TrimRight(cfg.MinIO.PublicURL, "/")

	return &minioStorage{
		client:    client,
		bucket:    cfg.MinIO.Bucket,
		publicURL: publicURL,
	}, nil
}

func (s *minioStorage) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("check bucket: %w", err)
	}
	if exists {
		return nil
	}
	if err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("create bucket: %w", err)
	}
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{
			"Effect": "Allow",
			"Principal": {"AWS": ["*"]},
			"Action": ["s3:GetObject"],
			"Resource": ["arn:aws:s3:::%s/*"]
		}]
	}`, s.bucket)
	if err := s.client.SetBucketPolicy(ctx, s.bucket, policy); err != nil {
		return fmt.Errorf("set bucket policy: %w", err)
	}
	return nil
}

func (s *minioStorage) UploadKomoditasImage(ctx context.Context, komoditasID uuid.UUID, reader io.Reader, size int64, contentType string) (string, error) {
	ext, ok := allowedContentTypes[contentType]
	if !ok {
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}

	objectKey := komoditasID.String() + ext
	_, err := s.client.PutObject(ctx, s.bucket, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload object: %w", err)
	}

	return s.objectURL(objectKey), nil
}

func (s *minioStorage) DeleteKomoditasImage(ctx context.Context, objectKey string) error {
	if objectKey == "" {
		return nil
	}
	if err := s.client.RemoveObject(ctx, s.bucket, objectKey, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("remove object: %w", err)
	}
	return nil
}

func (s *minioStorage) objectURL(objectKey string) string {
	return s.publicURL + "/" + url.PathEscape(objectKey)
}

// ObjectKeyFromURL extracts the object key from a stored public URL.
func ObjectKeyFromURL(publicURL, basePublicURL string) string {
	base := strings.TrimRight(basePublicURL, "/")
	if publicURL == "" || !strings.HasPrefix(publicURL, base) {
		return ""
	}
	trimmed := strings.TrimPrefix(publicURL, base+"/")
	if trimmed == "" {
		return ""
	}
	key, err := url.PathUnescape(trimmed)
	if err != nil {
		return trimmed
	}
	return path.Base(key)
}
