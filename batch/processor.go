package batch

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"

	base "github.com/falouu/go-libs-public/b"
	"github.com/sirupsen/logrus"
)

type BatchInput map[string]interface{}
type BatchOutput map[string]interface{}
type InputId string

type BatchJob struct {
	Name     string
	Inputs   map[InputId]BatchInput
	Job      func(input BatchInput) (BatchOutput, error)
	Interval time.Duration
	Version  int
}

type InputResult struct {
	err     error
	skipped bool
	output  BatchOutput
}

type BatchJobResult struct {
	Successes int
	Fails     int
	Skipped   int
	Outputs   map[InputId]InputResult
}

type BatchProcessor interface {
	RunJob(job BatchJob) (BatchJobResult, error)
}

func NewBatchProcessor(dir string) BatchProcessor {
	return &batchProcessor{
		dir: dir,
		log: logrus.WithField("src", "batch"),
	}
}

type batchProcessor struct {
	dir string
	log *logrus.Entry
}

func (b *batchProcessor) RunJob(job BatchJob) (BatchJobResult, error) {
	b.log.Infof("Starting Job %v with %v inputs. Interval = %v", job.Name, len(job.Inputs), job.Interval)

	jobResult := BatchJobResult{
		Outputs: map[InputId]InputResult{},
	}

	for inputId, input := range job.Inputs {
		result := InputResult{}
		skip, err := b.shouldSkip(&job, inputId)
		if err != nil {
			return jobResult, base.Wrap(err, "cannot determine if input should be skipped")
		}

		if skip {
			b.log.Infof("Skipping input %v... ", inputId)
			result.skipped = true
			jobResult.Skipped++
		} else {
			time.Sleep(job.Interval)
			b.log.Infof("Processing input %v (%+v)...", inputId, input)
			output, err := job.Job(input)
			result.output = output

			if err != nil {
				jobResult.Fails++
				b.log.WithError(err).Warnf("Input %v failed", inputId)
				result.err = err
			} else {
				b.log.Infof("Input %v processed", inputId)
				jobResult.Successes++
				err := b.saveSuccess(&job, inputId)
				if err != nil {
					return jobResult, err
				}
			}
		}
		jobResult.Outputs[inputId] = result
	}

	b.log.Infof("Batch job %v is done. Succcesses = %v, Fails = %v, Skipped = %v",
		job.Name, jobResult.Successes, jobResult.Fails, jobResult.Skipped)

	return jobResult, nil
}

func (b *batchProcessor) shouldSkip(job *BatchJob, id InputId) (bool, error) {
	return b.isInputProcessed(job, id)
}

func (b *batchProcessor) isInputProcessed(job *BatchJob, id InputId) (bool, error) {
	processedInputs, err := b.getProcessedInputs(job)
	if err != nil {
		return false, err
	}
	_, processed := processedInputs[id]
	return processed, nil
}

func (b *batchProcessor) getProcessedInputs(job *BatchJob) (map[InputId]bool, error) {
	result := map[InputId]bool{}
	bytes, err := os.ReadFile(b.getJobFilepath(job))
	if err != nil {
		if os.IsNotExist(err) {
			return result, nil
		} else {
			return nil, err
		}
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *batchProcessor) getJobFilepath(job *BatchJob) string {
	return filepath.Join(b.dir, job.Name+"-"+strconv.Itoa(job.Version))
}

func (b *batchProcessor) saveSuccess(job *BatchJob, id InputId) error {
	results, err := b.getProcessedInputs(job)
	if err != nil {
		return err
	}

	results[id] = true
	bytes, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		return err
	}
	err = os.MkdirAll(b.dir, 0700)
	if err != nil {
		return err
	}
	err = os.WriteFile(b.getJobFilepath(job), bytes, 0600)
	return err
}
