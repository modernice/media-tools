package image

import (
	"context"
	"fmt"
	"image"

	"github.com/jaevor/go-nanoid"
	"github.com/modernice/media-tools/image/internal"
	"github.com/modernice/media-tools/internal/slices"
)

var newID func() string

func init() {
	var err error
	if newID, err = nanoid.Standard(21); err != nil {
		panic(fmt.Errorf("nanoid: %w", err))
	}
}

// Pipeline is a list of Processors that are applied to an image.
type Pipeline []Processor

// Processed is a processed image.
type Processed struct {
	ID    string
	Image *image.NRGBA
}

// A Processor processes an image and returns possibly multiple processed images.
type Processor interface {
	Process(ProcessorContext) ([]*image.NRGBA, error)
}

// ProcessorContext is passed to Processors.
type ProcessorContext interface {
	context.Context

	// Image is the input input for the processor.
	Image() *image.NRGBA

	// Original returns whether the image is the original image that was
	// provided to [Pipeline.Run].
	Original() bool
}

// NewProcessorContext returns a new [ProcessorContext] for a [Processor].
func NewProcessorContext(ctx context.Context, img *image.NRGBA, original bool) ProcessorContext {
	return &processorContext{
		Context:  ctx,
		image:    img,
		original: original,
	}
}

type processorContext struct {
	context.Context

	image    *image.NRGBA
	original bool
}

func (ctx *processorContext) Image() *image.NRGBA {
	return ctx.image
}

func (ctx *processorContext) Original() bool {
	return ctx.original
}

// PipelineResult is the result of running a [Pipeline] on an image.
type PipelineResult struct {
	// Images are the processed images.
	Images []Processed
}

// Run runs the pipeline on an image and returns the [PipelineResult],
// containing the processed images.
func (pipeline Pipeline) Run(ctx context.Context, img image.Image) (PipelineResult, error) {
	nimg := internal.ToNRGBA(img)

	previous := []*image.NRGBA{nimg}

	for i, processor := range pipeline {
		_previous := previous
		previous = previous[:0]

		for _, img := range _previous {
			pctx := NewProcessorContext(ctx, img, i == 0)

			processed, err := processor.Process(pctx)
			if err != nil {
				return PipelineResult{}, fmt.Errorf("%T processor: %w", processor, err)
			}

			previous = append(previous, processed...)
		}
	}

	processed := slices.Map(func(img *image.NRGBA) Processed {
		return Processed{
			ID:    newID(),
			Image: img,
		}
	}, previous)

	return PipelineResult{Images: processed}, nil
}
