package images

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

func TestComposeAvatarAtRect(t *testing.T) {
	dstFile, err := os.Create("dst.jpg")
	if err != nil {
		t.Fatalf("Create dstFile Error: %s\n", err.Error())
	}

	defer dstFile.Close()
	imgFile, err := os.Open("google.png")
	if err != nil {
		t.Fatalf("Open ImgFile Error: %s\n", err.Error())
	}
	defer imgFile.Close()
	img, err := png.Decode(imgFile)
	if err != nil {
		t.Fatalf("Decode ImgFile Error: %s\n", err.Error())
	}
	avatarsCount := 81
	var avatars []image.Image
	for i := 0; i < avatarsCount; i++ {
		avatars = append(avatars, img)
	}
	ret := ComposeAvatarAtRect(300, 10, 2, avatars...)
	jpeg.Encode(dstFile, ret, nil)
}
