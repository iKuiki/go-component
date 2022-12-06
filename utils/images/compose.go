package images

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/nfnt/resize"
)

// ComposeAvatarAtRect 将给定的一组头像组合为一个群头像
func ComposeAvatarAtRect(size, padding, margin int, avatars ...image.Image) image.Image {
	background := image.NewRGBA(image.Rect(0, 0, size, size))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(background, background.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)
	avatarsCount := len(avatars)
	if avatarsCount == 0 {
		return background
	}
	rowCount := int(math.Sqrt(float64(avatarsCount)))
	var avatarsRect [][]image.Image
	columnCount := avatarsCount / rowCount
	if avatarsCount%rowCount != 0 { // 说明比3行有多出，则前两行每行加1
		columnCount++
	}
	// 组装头像矩阵
	for i := 0; i < rowCount; i++ {
		avatarsRect = append(avatarsRect, avatars[i*columnCount:(i+1)*columnCount])
	}
	// 计算画布边距
	// canvasX1 := padding
	// canvasX2 := size - padding
	canvasSize := size - 2*padding
	itemSize := canvasSize / columnCount
	for rowNo, row := range avatarsRect {
		rowOrigin := padding + rowNo*canvasSize/rowCount
		for columnNo, item := range row {
			columnOrigin := padding + columnNo*itemSize
			// 计算当前图片绘制参数
			itemOriginX := columnOrigin + margin
			itemOriginY := rowOrigin + margin
			itemDrawSize := itemSize - 2*margin
			if item != nil {
				img := resize.Resize(uint(itemDrawSize), uint(itemDrawSize), item, resize.Lanczos3)
				draw.Draw(background, background.Bounds(), img, image.Pt(-itemOriginX, -itemOriginY), draw.Over)
			}
		}
	}
	return background
}
