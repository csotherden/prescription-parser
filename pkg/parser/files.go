package parser

import (
	"context"
	"fmt"
	"io"
)

func (p *PrescriptionParser) UploadImage(ctx context.Context, fileName string, file io.Reader) (string, error) {
	return p.uploadImageFunc(ctx, fileName, file)
}

func (p *PrescriptionParser) DeleteImage(ctx context.Context, id string) error {
	_, err := p.openAIClient.Files.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}
