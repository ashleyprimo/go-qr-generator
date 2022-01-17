package qr

import (
	"github.com/boombuler/barcode"
	"image"
	"image/color"
)

type WhiteBorder struct {
	width   int
	barcode barcode.Barcode
}

func (wb WhiteBorder) At(x, y int) color.Color {
	bounds := wb.barcode.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if x < wb.width || x >= w+wb.width || y < wb.width || y >= h+wb.width {
		return color.White
	}
	return wb.barcode.At(x-wb.width, y-wb.width)
}

func (wb WhiteBorder) Bounds() image.Rectangle {
	b := wb.barcode.Bounds()
	return image.Rect(0, 0, b.Dx()+2*wb.width, b.Dy()+2*wb.width)
}

func (wb WhiteBorder) ColorModel() color.Model {
	return wb.barcode.ColorModel()
}
