package parser

import (
	"context"
	"fmt"
	"github.com/openai/openai-go"
	"io"
	"path/filepath"
	"strings"
)

func (p *PrescriptionParser) uploadImageOpenAI(ctx context.Context, fileName string, file io.Reader) (string, error) {
	fileExt := strings.ToLower(filepath.Ext(fileName))
	var contentType string
	switch fileExt {
	case ".pdf":
		contentType = "application/pdf"
	default:
		return "", fmt.Errorf("unsupported file type: %s. file must be PDF", fileExt)
	}

	inputFile := openai.File(file, fileName, contentType)

	storedFile, err := p.openAIClient.Files.New(ctx, openai.FileNewParams{
		File:    inputFile,
		Purpose: openai.FilePurposeUserData,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	return storedFile.ID, nil
}
