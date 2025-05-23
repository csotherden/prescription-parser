package parser

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"

	"github.com/csotherden/prescription-parser/pkg/jobs"
)

const JobTypeParsePrescription = "parse_prescription"

func (p *PrescriptionParser) ParseImage(ctx context.Context, fileName string, file io.Reader) (string, error) {
	if p.parseImageFunc == nil {
		return "", fmt.Errorf("no parser backend enabled")
	}

	// Create a job for asynchronous processing
	jobID := jobs.GlobalTracker.CreateJob(
		JobTypeParsePrescription,
		fmt.Sprintf("Processing image: %s", fileName),
	)

	p.logger.Info("starting image parsing", zap.String("job_id", jobID), zap.String("file_name", fileName))

	go p.parseImageFunc(context.Background(), jobID, fileName, file)

	return jobID, nil
}
