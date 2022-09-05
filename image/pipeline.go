package image

import (
	"context"
	"fmt"
	"image"
	"strings"

	"github.com/modernice/media-tools/internal/slices"
)

// Original is the tag that is assigned to the original image that is passed
// to [Pipeline.Run].
const Original = "original"

// Pipeline is a list of Processors that are applied to an image.
type Pipeline []Processor

// Processed is a processed image, except if the Original field is set to true.
// In that case, it is the original image that was passed to [Pipeline.Run].
type Processed struct {
	Image    image.Image
	Tags     Tags
	Original bool
}

// A Processor processes an image and returns possibly multiple processed images.
type Processor interface {
	Process(ProcessorContext) ([]Processed, error)
}

// ProcessorFunc allows functions to be used a Processors.
type ProcessorFunc func(ProcessorContext) ([]Processed, error)

// ProcessorContext is passed to Processors.
type ProcessorContext interface {
	context.Context

	// Image is the input image for the Processor.
	Image() Processed
}

// NewProcessorContext returns a new [ProcessorContext] for a [Processor].
func NewProcessorContext(ctx context.Context, img Processed) ProcessorContext {
	return &processorContext{
		Context: ctx,
		image:   img,
	}
}

type processorContext struct {
	context.Context

	image Processed
}

func (ctx *processorContext) Image() Processed {
	return ctx.image
}

// PipelineResult is the result of running a [Pipeline] on an image.
type PipelineResult struct {
	// Images are the processed images.
	Images []Processed
}

// Tags is a list of tags that Processors assigned to images in a [Pipeline].
type Tags []string

// NewTags returns the given tags as [Tags]. Duplicates are removed.
func NewTags(tags ...string) Tags {
	return slices.Unique(Tags(tags))
}

// Contains returns whether a tag is contained within tags.
func (tags Tags) Contains(tag string) bool {
	return slices.Contains(tag, tags)
}

// With appends additional tags and returns the new [Tags]. Duplicate tags are removed.
func (tags Tags) With(add ...string) Tags {
	return slices.Unique(append(tags, add...))
}

// Without returns a copy of tags without the given tags.
func (tags Tags) Without(remove ...string) Tags {
	out := make(Tags, 0, len(tags))
	for _, tag := range tags {
		if slices.Contains(tag, remove) {
			continue
		}
		out = append(out, tag)
	}
	return out
}

// Run runs the pipeline on an image and returns the [PipelineResult],
// containing the processed images.
func (pipeline Pipeline) Run(ctx context.Context, img image.Image) (PipelineResult, error) {
	previous := []Processed{{Image: img, Original: true, Tags: NewTags(Original)}}

	for _, processor := range pipeline {
		_previous := previous
		previous = previous[:0]

		for _, img := range _previous {
			pctx := NewProcessorContext(ctx, img)

			processed, err := processor.Process(pctx)
			if err != nil {
				return PipelineResult{}, fmt.Errorf("%T processor: %w", processor, err)
			}

			var originalCount int
			for _, pimg := range processed {
				if pimg.Original {
					originalCount++
				}

				if originalCount > 1 {
					return PipelineResult{}, fmt.Errorf("%T processor returned more than one %q image", processor, Original)
				}
			}

			previous = append(previous, processed...)
		}
	}

	return PipelineResult{Images: previous}, nil
}

// DimensionName extracts the dimension name from the tags of a processed image.
func DimensionName(tags Tags) string {
	for _, tag := range tags {
		if strings.HasPrefix(tag, "size=") {
			return tag[5:]
		}
	}
	return ""
}
