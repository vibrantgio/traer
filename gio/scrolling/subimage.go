package main

import "image"

type subimage interface {
	SubImage(r image.Rectangle) image.Image
}

func SubImage(img image.Image, x0 int, y0 int, x1 int, y1 int) image.Image {
	return img.(subimage).SubImage(image.Rect(x0, y0, x1, y1))
}
