package storage

import (
	"context"
	"io"
	"github.com/google/uuid"
)

// KomoditasStorage handles object storage for komoditas images and pengumpulan signatures.
type KomoditasStorage interface {
	UploadKomoditasImage(ctx context.Context, komoditasID uuid.UUID, reader io.Reader, size int64, contentType string) (publicURL string, err error)
	UploadPengumpulanSignature(ctx context.Context, pengumpulanDataID uuid.UUID, reader io.Reader, size int64, contentType string) (publicURL string, err error)
	DeleteKomoditasImage(ctx context.Context, objectKey string) error
	DeletePengumpulanSignature(ctx context.Context, objectKey string) error
	EnsureBucket(ctx context.Context) error
}
