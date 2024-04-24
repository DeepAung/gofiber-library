package storer

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	"github.com/DeepAung/gofiber-library/pkg/configs"
	"github.com/DeepAung/gofiber-library/pkg/utils"
)

type GCPStorer struct {
	cfg *configs.Config
}

func NewGCPStorer(cfg *configs.Config) Storer {
	return &GCPStorer{
		cfg: cfg,
	}
}

func (s *GCPStorer) UploadFiles(files []*multipart.FileHeader, dir string) ([]FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return []FileRes{}, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	length := len(files)
	jobCh := make(chan FileReq, length)
	resCh := make(chan FileRes, length)
	errCh := make(chan error, length)

	numWorkers := 5
	for range numWorkers {
		go func() {
			for job := range jobCh {
				select {
				case <-ctx.Done():
					errCh <- fmt.Errorf("upload file cancelled")
					return
				default:
				}

				res, err := s.uploadFile(ctx, cancel, client, job.file, dir)
				if err != nil {
					errCh <- err
					return
				}

				errCh <- nil
				resCh <- res
			}
		}()
	}

	for _, file := range files {
		jobCh <- FileReq{file: file, dir: dir}
	}

	var results []FileRes
	for range length {
		err = <-errCh
		if err != nil {
			return results, err
		}

		results = append(results, <-resCh)
	}

	return results, nil
}

func (s *GCPStorer) DeleteFiles(dests []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	length := len(dests)
	jobCh := make(chan string, length)
	errCh := make(chan error, length)

	numWorkers := 5
	for range numWorkers {
		go func() {
			for dest := range jobCh {
				select {
				case <-ctx.Done():
					errCh <- fmt.Errorf("delete file cancelled")
					return
				default:
				}

				err := s.deleteFile(ctx, cancel, client, dest)
				if err != nil {
					errCh <- err
					return
				}

				errCh <- nil
			}
		}()
	}

	for _, dest := range dests {
		jobCh <- dest
	}

	for range length {
		err = <-errCh
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *GCPStorer) uploadFile(
	ctx context.Context,
	cancel func(),
	client *storage.Client,
	file *multipart.FileHeader,
	dir string,
) (FileRes, error) {
	f, err := file.Open()
	if err != nil {
		cancel()
		return FileRes{}, fmt.Errorf("file.Open: %w", err)
	}
	defer f.Close()

	bucket := s.cfg.App.GCPBucket
	dest := utils.Join(dir, file.Filename)
	o := client.Bucket(bucket).Object(dest)

	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		cancel()
		return FileRes{}, fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		cancel()
		return FileRes{}, fmt.Errorf("Writer.Close: %w", err)
	}

	fmt.Printf("Blob %v uploaded.\n", dest)
	return FileRes{
		Filename: file.Filename,
		Url:      utils.MakeGCPUrl(bucket, dest),
	}, nil
}

func (s *GCPStorer) deleteFile(
	ctx context.Context,
	cancel func(),
	client *storage.Client,
	dest string,
) error {
	bucket := s.cfg.App.GCPBucket
	o := client.Bucket(bucket).Object(dest)

	attrs, err := o.Attrs(ctx)
	if err != nil {
		cancel()
		return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		cancel()
		return fmt.Errorf("Object(%q).Delete: %w", dest, err)
	}

	fmt.Printf("Blob %v deleted.\n", dest)
	return nil
}
