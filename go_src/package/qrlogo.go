package qrlogo

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"log"

	// qr "github.com/skip2/go-qrcode"
	qr "github.com/garaio/qrlogo/package/qrcode"
)

// Encoder defines settings for QR/Overlay encoder.
type Encoder struct {
	AlphaThreshold int
	GreyThreshold  int
	QRLevel        qr.RecoveryLevel
}

// DefaultEncoder is the encoder with default settings.
var DefaultEncoder = Encoder{
	AlphaThreshold: 2000,      // FIXME: don't remember where this came from
	GreyThreshold:  30,        // in percent
	QRLevel:        qr.Medium, // Better would be 'Highest', as logo steals some redundant space - but we need Medium
}

// Encode encodes QR image, adds logo overlay and renders result as PNG.
func Encode(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	return DefaultEncoder.Encode(str, logo, size)
}

// Encode encodes QR image, adds logo overlay and renders result as PNG.
func (e Encoder) Encode(str string, logo image.Image, size int) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	errcheck(nil, "here")
	code, err := qr.New(str, e.QRLevel)

	errcheck(nil, "here-1")
	if err != nil {
		return nil, err
	}

	errcheck(nil, "here-2")
	img := code.Image(size)

	errcheck(nil, "here-3")
	e.overlayLogo(img, logo)

	errcheck(nil, "here-4")
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}

// overlayLogo blends logo to the center of the QR code,
// changing all colors to black.
func (e Encoder) overlayLogo(dst, src image.Image) {
	grey := uint32(^uint16(0)) * uint32(e.GreyThreshold) / 100
	alphaOffset := uint32(e.AlphaThreshold)
	offset := dst.Bounds().Max.X/2 - src.Bounds().Max.X/2
	for x := 0; x < src.Bounds().Max.X; x++ {
		for y := 0; y < src.Bounds().Max.Y; y++ {
			if r, g, b, alpha := src.At(x, y).RGBA(); alpha > alphaOffset {
				col := color.Black
				if r > grey && g > grey && b > grey {
					col = color.White
				}
				dst.(*image.Paletted).Set(x+offset, y+offset, col)
			}
		}
	}
}

func errcheck(err error, str string) {
	f, e := os.OpenFile("qr-encode.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if e != nil {
		log.Fatalf("error opening file: %v", err)
		// fmt.Println(str, err)
		// os.Exit(1)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(str, err)

}
