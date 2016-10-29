// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package steganogopher

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"testing"
)

var testCase = []struct {
	filename  string
	output    string
	quality   int
	tolerance int64
}{

	{"_test/twitter.png", "_test/twitter_out-20.jpg", 20, 12 << 8},
	{"_test/twitter.png", "_test/twitter_out-60.jpg", 60, 8 << 8},
	{"_test/twitter.png", "_test/twitter_out-80.jpg", 80, 6 << 8},
	{"_test/twitter.png", "_test/twitter_out-90.jpg", 90, 4 << 8},
	{"_test/twitter.png", "_test/twitter_out-100.jpg", 100, 2 << 8},

	{"_test/evilcat.png", "_test/evlicat_out-20.jpg", 20, 12 << 8},
	{"_test/evilcat.png", "_test/evlicat_out-60.jpg", 60, 8 << 8},
	{"_test/evilcat.png", "_test/evlicat_out-80.jpg", 80, 6 << 8},
	{"_test/evilcat.png", "_test/evlicat_out-90.jpg", 90, 4 << 8},
	{"_test/evilcat.png", "_test/evlicat_out-100.jpg", 100, 2 << 8},
}

func delta(u0, u1 uint32) int64 {
	d := int64(u0) - int64(u1)
	if d < 0 {
		return -d
	}
	return d
}

func decodePng(f *os.File) (image.Image, error) {

	defer f.Close()
	return png.Decode(f)
}

func TestWriteToDisk(t *testing.T) {
	for _, tc := range testCase {

		msg := "Hello world!"

		file, err := os.Open(tc.filename)
		if err != nil {
			t.Errorf("Error while opening the file %s", tc.filename)
		}

		i, err := decodePng(file)
		if err != nil {
			t.Error(tc.filename, err)
		}
		f, err := os.Create(tc.output)
		if err != nil {
			t.Error(tc.filename, err)
		}
		err = Encode(f, i, msg, &Options{Quality: tc.quality})
		if err != nil {
			t.Error(tc.filename, err)
		}
	}
}

func TestWriter(t *testing.T) {
	for _, tc := range testCase {
		msg := "Hello world; or こんにちは 世界!"

		file, err := os.Open(tc.filename)
		if err != nil {
			t.Errorf("Error while opening file %s", tc.filename)
		}

		// Read the image.
		m0, err := decodePng(file)
		if err != nil {
			t.Error(tc.filename, err)
			continue
		}
		// Encode that image as JPEG.
		buf := bytes.NewBuffer(nil)
		err = Encode(buf, m0, msg, &Options{Quality: tc.quality})
		if err != nil {
			t.Error(tc.filename, err)
			continue
		}
		// Decode that JPEG.
		m1, data, err := DecodeAndRead(buf)

		if data != msg {
			t.Error("Got wrong message back:", data)
		}

		if err != nil {
			t.Error(tc.filename, err)
			continue
		}
		// Compute the average delta in RGB space.
		b := m0.Bounds()
		var sum, n int64
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				c0 := m0.At(x, y)
				c1 := m1.At(x, y)
				r0, g0, b0, _ := c0.RGBA()
				r1, g1, b1, _ := c1.RGBA()
				sum += delta(r0, r1)
				sum += delta(g0, g1)
				sum += delta(b0, b1)
				n += 3
			}
		}
		// Compare the average delta to the tolerance level.
		if sum/n > tc.tolerance {
			t.Errorf("%s, quality=%d: average delta is too high", tc.filename, tc.quality)
			continue
		}
	}
}
