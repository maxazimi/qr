package camera

import (
	"image"
)

func convertBGRAToRGBA(rgbBuffer []byte, width, height int) *image.RGBA {
	rgbaImage := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Calculate the index for the mirrored column
			mirroredX := width - 1 - x
			rgbaIndex := (y*width + x) * 4
			mirroredIndex := (y*width + mirroredX) * 4

			b := rgbBuffer[rgbaIndex]
			g := rgbBuffer[rgbaIndex+1]
			r := rgbBuffer[rgbaIndex+2]
			a := rgbBuffer[rgbaIndex+3]

			rgbaImage.Pix[mirroredIndex] = r
			rgbaImage.Pix[mirroredIndex+1] = g
			rgbaImage.Pix[mirroredIndex+2] = b
			rgbaImage.Pix[mirroredIndex+3] = a
		}
	}

	return rgbaImage
}
