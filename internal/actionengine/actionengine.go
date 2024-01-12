package actionengine

import (
	"image"
	"image/color"
	"sync"
)

/*ActOnImagePixel runs provided function on image. Acts a worker pool*/
func ActOnImagePixel(imageOld image.Image, pixelAction func(p image.Point, imageOld image.Image) color.Color, poolSize int) *image.RGBA {

	/* Initialize need variables */
	var wg sync.WaitGroup
	wg.Add(poolSize)
	imageNew := image.NewRGBA(image.Rect(imageOld.Bounds().Min.X, imageOld.Bounds().Min.Y, imageOld.Bounds().Max.X, imageOld.Bounds().Max.Y))
	ch := make(chan image.Point)

	/* set up worker pool */
	for i := 0; i < poolSize; i++ {
		go func() {
			defer wg.Done()
			for p := range ch {
				/* set pixel in new image based on provided action function */
				imageNew.Set(p.X, p.Y, pixelAction(p, imageOld))
			}
		}()
	}

	/* feed each pixel to workers */
	for y := imageOld.Bounds().Min.Y; y < imageOld.Bounds().Max.Y; y++ {
		for x := imageOld.Bounds().Min.X; x < imageOld.Bounds().Max.X; x++ {
			ch <- image.Point{X: x, Y: y}
		}
	}
	close(ch)
	wg.Wait()
	return imageNew
}

/*ActOnImageKernal runs provided function on image. Acts a worker pool*/
func ActOnImageKernal(imageOld image.Image, kernalAction func(p image.Point, imageOld image.Image, kernalSize int) [][]color.Color, kernalSize int, poolSize int) *image.RGBA {

	/* Initialize need variables */
	var wg sync.WaitGroup
	wg.Add(poolSize)
	imageNew := image.NewRGBA(image.Rect(imageOld.Bounds().Min.X, imageOld.Bounds().Min.Y, imageOld.Bounds().Max.X, imageOld.Bounds().Max.Y))
	ch := make(chan image.Point)

	/* set up worker pool */
	for i := 0; i < poolSize; i++ {
		go func() {
			defer wg.Done()
			for p := range ch {

				kernalResult := kernalAction(p, imageOld, kernalSize)
				for y := 0; y < kernalSize; y++ {
					for x := 0; x < kernalSize; x++ {
						col := kernalResult[x][y]
						imageNew.Set(x+p.X, y+p.Y, col)
					}
				}
			}
		}()
	}

	/* feed each pixel to workers */
	for y := imageOld.Bounds().Min.Y; y < imageOld.Bounds().Max.Y-kernalSize; y += 1 {
		for x := imageOld.Bounds().Min.X; x < imageOld.Bounds().Max.X-kernalSize; x += 1 {
			ch <- image.Point{X: x, Y: y}
		}
	}
	close(ch)
	wg.Wait()
	return imageNew
}
