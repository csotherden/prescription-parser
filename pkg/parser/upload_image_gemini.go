package parser

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"google.golang.org/genai"
)

func (p *PrescriptionParser) uploadImageGemini(ctx context.Context, fileName string, file io.Reader) (string, error) {
	fileExt := strings.ToLower(filepath.Ext(fileName))
	var contentType string
	switch fileExt {
	case ".pdf":
		contentType = "application/pdf"
	default:
		return "", fmt.Errorf("unsupported file type: %s. file must be PDF", fileExt)
	}

	resp, err := p.geminiClient.Files.Upload(ctx, file, &genai.UploadFileConfig{
		MIMEType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	return resp.URI, nil
}
