package image

// Tagger is a Processor that adds tags to images.
type Tagger struct {
	fn func(Processed) Tags
}

// Tag returns a Tagger that adds the provided tags to images.
func Tag(tags Tags) *Tagger {
	return TagBy(func(Processed) Tags {
		return tags
	})
}

// TagBy returns a Tagger that adds tags to images. For each image, the provided
// function is called to determine which tags to add to the image.
func TagBy(fn func(Processed) Tags) *Tagger {
	return &Tagger{fn}
}

// Process implements [Processor]. It adds the configured tags to the image.
func (tagger *Tagger) Process(ctx ProcessorContext) ([]Processed, error) {
	pimg := ctx.Image()
	pimg.Tags = pimg.Tags.With(tagger.fn(pimg)...)
	return []Processed{pimg}, nil
}
