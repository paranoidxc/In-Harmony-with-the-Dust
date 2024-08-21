package utils

import (
	"bytes"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"os"
)

func QrcodeWeChat(content string) ([]byte, error) {
	qr, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	qrImg, err := png.Decode(bytes.NewReader(qr))
	if err != nil {
		return nil, err
	}

	iconFile, err := os.Open("./resource/wechat.png")
	if err != nil {
		return nil, err
	}
	defer iconFile.Close()

	iconImg, _, err := image.Decode(iconFile)
	if err != nil {
		return nil, err
	}

	qrImgWithIcon := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(qrImgWithIcon, qrImgWithIcon.Bounds(), qrImg, image.Point{}, draw.Src)

	iconWidth := iconImg.Bounds().Dx()
	iconHeight := iconImg.Bounds().Dy()
	iconX := (256 - iconWidth) / 2
	iconY := (256 - iconHeight) / 2
	draw.Draw(qrImgWithIcon, image.Rect(iconX, iconY, iconX+iconWidth, iconY+iconHeight), iconImg, image.Point{}, draw.Over)

	var qrImgWithIconBuf bytes.Buffer
	err = png.Encode(&qrImgWithIconBuf, qrImgWithIcon)
	if err != nil {
		return nil, err
	}

	qrImgWithIconBytes := qrImgWithIconBuf.Bytes()
	return qrImgWithIconBytes, nil
}
