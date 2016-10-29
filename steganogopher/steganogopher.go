package steganogopher

import (
	"bufio"
	"errors"
	"image"
	"io"
)

// Encoder is the signature of the Encode function. Specs TBD
type Encoder func(w io.Writer, i image.Image, m string, o *Options) error

// Decoder is the signature of the Decode function. Specs TBD.
type Decoder func(r io.Reader) (image.Image, error)

// Decode reads a JPEG image from r and returns it as an image.Image.
func Decode(r io.Reader) (image.Image, error) {
	var d decoder
	return d.decode(r, false)
}

// Encode writes the Image m to w in JPEG 4:2:0 baseline format with the given
// options. Default parameters are used if a nil *Options is passed.
func Encode(w io.Writer, m image.Image, message string, o *Options) error {
	b := m.Bounds()
	if b.Dx() >= 1<<16 || b.Dy() >= 1<<16 {
		return errors.New("jpeg: image is too large to encode")
	}
	var e encoder
	e.data = []byte(message)
	if ww, ok := w.(writer); ok {
		e.w = ww
	} else {
		e.w = bufio.NewWriter(w)
	}
	// Clip quality to [1, 100].
	quality := DefaultQuality
	if o != nil {
		quality = o.Quality
		if quality < 1 {
			quality = 1
		} else if quality > 100 {
			quality = 100
		}
	}
	// Convert from a quality rating to a scaling factor.
	var scale int
	if quality < 50 {
		scale = 5000 / quality
	} else {
		scale = 200 - quality*2
	}
	// Initialize the quantization tables.
	for i := range e.quant {
		for j := range e.quant[i] {
			x := int(unscaledQuant[i][j])
			x = (x*scale + 50) / 100
			if x < 1 {
				x = 1
			} else if x > 255 {
				x = 255
			}
			e.quant[i][j] = uint8(x)
		}
	}
	// Write the Start Of Image marker.
	e.buf[0] = 0xff
	e.buf[1] = 0xd8
	e.write(e.buf[:2])
	// Write the quantization tables.
	e.writeDQT()
	// Write the image dimensions.
	e.writeSOF0(b.Size())
	// Write the Huffman tables.
	e.writeDHT()
	// Write the image data.
	e.writeSOS(m)
	// Write the End Of Image marker.
	e.buf[0] = 0xff
	e.buf[1] = 0xd9
	e.write(e.buf[:2])
	e.flush()
	return e.err
}

//
func DecodeAndRead(r io.Reader) (image.Image, string, error) {
	var d decoder
	i, err := d.decode(r, false)
	return i, string(d.data), err
}

// DecodeConfig returns the color model and dimensions of a JPEG image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	c, _, err := image.DecodeConfig(r)
	if err != nil {
		return image.Config{}, err
	}
	return c, nil
}
