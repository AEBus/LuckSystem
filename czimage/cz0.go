package czimage

import (
	"encoding/binary"
	"github.com/golang/glog"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/go-restruct/restruct"
)

type Cz0Header struct {
	Flag    uint8
	X       uint16
	Y       uint16
	Width1  uint16
	Heigth1 uint16

	Width2  uint16
	Heigth2 uint16
}

// Cz0Image
//  Description Cz0.Load() 载入数据
//  Description Cz0.Export() 解压数据，转化成Image并导出
//  Description Cz0.GetImage() 解压数据，转化成Image
type Cz0Image struct {
	CzHeader
	Cz0Header
	CzData
}

func (cz *Cz0Image) Load(header CzHeader, data []byte) {
	cz.CzHeader = header
	cz.Raw = data
	err := restruct.Unpack(cz.Raw[15:], binary.LittleEndian, &cz.Cz0Header)
	if err != nil {
		panic(err)
	}
	glog.V(6).Infoln("cz0 header ", cz.Cz0Header)
	cz.OutputInfo = nil
}

func (cz *Cz0Image) decompress() {
	//os.WriteFile("../data/LB_EN/IMAGE/2.lzw", cz.Raw[int(cz.HeaderLength)+cz.OutputInfo.Offset:], 0666)
	buf := Decompress(cz.Raw[int(cz.HeaderLength)+cz.OutputInfo.Offset:], cz.OutputInfo)
	glog.V(6).Infoln("uncompress size ", len(buf))
	cz.Image = LineDiff(&cz.CzHeader, buf)
}

func (cz *Cz0Image) Export(path string) {
	pic := image.NewRGBA(image.Rect(0, 0, int(cz.Width), int(cz.Heigth)))
	offset := int(cz.HeaderLength)
	switch cz.Colorbits {
	case 32:
		for y := 0; y < int(cz.Heigth); y++ {
			for x := 0; x < int(cz.Width); x++ {
				pic.SetRGBA(x, y, color.RGBA{
					R: cz.Raw[offset],
					G: cz.Raw[offset+1],
					B: cz.Raw[offset+2],
					A: cz.Raw[offset+3]},
				)
				offset += 4
			}
		}
	}
	cz.Image = pic
	f, _ := os.Create(path)
	defer f.Close()
	png.Encode(f, cz.Image)
}

func (cz *Cz0Image) GetImage() image.Image {
	if cz.Image == nil {
		cz.decompress()
	}
	return cz.Image
}

func (cz *Cz0Image) Import(file string) {
	pngFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer pngFile.Close()
	cz.PngImage, err = png.Decode(pngFile)
	if err != nil {
		panic(err)
	}

	cz0File, _ := os.Create(file + ".cz0")
	defer cz0File.Close()
	glog.V(6).Infoln(cz.CzHeader)
	err = WriteStruct(cz0File, &cz.CzHeader, &cz.Cz0Header)
	if err != nil {
		panic(err)
	}
	pic := cz.PngImage.(*image.RGBA)
	switch cz.Colorbits {
	case 32:
		cz0File.Write(pic.Pix)
	}

}
