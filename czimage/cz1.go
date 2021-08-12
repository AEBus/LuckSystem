package czimage

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

type Cz1Image struct {
	CzHeader
	ColorPanel [][]byte // []BGRA
	Image      image.Image
}

func (cz *Cz1Image) Load(header CzHeader, data []byte) {
	cz.CzHeader = header
	pic := image.NewRGBA(image.Rect(0, 0, int(header.Width), int(header.Heigth)))
	offset := int(cz.HeaderLength)
	switch cz.Colorbits {
	case 4:
		// TODO 未测试
		cz.ColorPanel = make([][]byte, 16)
		for i := 0; i < 16; i++ {
			cz.ColorPanel[i] = data[offset : offset+4]
			offset += 4
		}
		fmt.Println("cz1 colorPanel", len(cz.ColorPanel))
		buf := Decompress(data[offset:])
		fmt.Println("uncompress size", len(buf))
		i := 0
		var index uint8
		for y := 0; y < int(header.Heigth); y++ {
			for x := 0; x < int(header.Width); x++ {
				if i%2 == 0 {
					index = buf[i/2] & 0x0F // low4bit
				} else {
					index = (buf[i/2] & 0xF0) >> 4 // high4bit
				}
				pic.SetRGBA(x, y, color.RGBA{
					R: cz.ColorPanel[index][2],
					G: cz.ColorPanel[index][1],
					B: cz.ColorPanel[index][0],
					A: cz.ColorPanel[index][3],
				})
				i++
			}
		}
	case 8:
		cz.ColorPanel = make([][]byte, 256)
		for i := 0; i < 256; i++ {
			cz.ColorPanel[i] = data[offset : offset+4]
			offset += 4
		}
		fmt.Println("cz1 colorPanel", len(cz.ColorPanel))
		buf := Decompress(data[offset:])
		fmt.Println("uncompress size", len(buf))
		// B,G,R,A
		// 0,1,2,3
		i := 0
		for y := 0; y < int(header.Heigth); y++ {
			for x := 0; x < int(header.Width); x++ {
				pic.SetRGBA(x, y, color.RGBA{
					R: cz.ColorPanel[buf[i]][2],
					G: cz.ColorPanel[buf[i]][1],
					B: cz.ColorPanel[buf[i]][0],
					A: cz.ColorPanel[buf[i]][3],
				})
				i++
			}
		}
	case 24:
		// TODO 未测试
		// RGB
		buf := Decompress(data[offset:])
		fmt.Println("uncompress size", len(buf))
		i := 0
		for y := 0; y < int(header.Heigth); y++ {
			for x := 0; x < int(header.Width); x++ {
				pic.SetRGBA(x, y, color.RGBA{
					R: buf[i+0],
					G: buf[i+1],
					B: buf[i+2],
					A: 0xFF,
				})
				i += 3
			}
		}
	case 32:
		// TODO 未测试
		// RGBA
		buf := Decompress(data[offset:])
		fmt.Println("uncompress size", len(buf))
		pic.Pix = buf
	}

	cz.Image = pic
}
func (cz *Cz1Image) Save(path string) {
	f, _ := os.Create(path)
	defer f.Close()
	png.Encode(f, cz.Image)
}

func (cz *Cz1Image) Get() image.Image {
	return cz.Image
}