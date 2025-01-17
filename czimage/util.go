package czimage

import (
	"bytes"
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"io"
)

// GetOutputInfo 读取分块信息
//  Description 读取分块信息
//  Param data []byte
//  Return outputInfo
//
func GetOutputInfo(data []byte) (outputInfo *CzOutputInfo) {
	outputInfo = &CzOutputInfo{}
	err := restruct.Unpack(data, binary.LittleEndian, outputInfo)
	if err != nil {
		panic(err)
	}
	for _, block := range outputInfo.BlockInfo {
		outputInfo.TotalRawSize += int(block.RawSize)
		outputInfo.TotalCompressedSize += int(block.CompressedSize)
	}
	outputInfo.Offset = 4 + int(outputInfo.FileCount)*8
	return outputInfo
}

// WriteStruct 写入结构体
//  Description 写入结构体
//  Param writer io.Writer
//  Param list ...interface{}
//  Return error
//
func WriteStruct(writer io.Writer, list ...interface{}) error {
	for _, v := range list {
		temp, err := restruct.Pack(binary.LittleEndian, v)
		if err != nil {
			return err
		}
		writer.Write(temp)
	}
	return nil
}

// Decompress 解压数据
//  Description
//  Param data []byte 压缩的数据
//  Param outputInfo *CzOutputInfo 分块信息
//  Return []byte
//
func Decompress(data []byte, outputInfo *CzOutputInfo) []byte {
	offset := 0

	// fmt.Println("uncompress info", outputInfo)
	outputBuf := &bytes.Buffer{}
	for _, block := range outputInfo.BlockInfo {
		lzwBuf := make([]uint16, int(block.CompressedSize))
		//offsetTemp := offset
		for j := 0; j < int(block.CompressedSize); j++ {
			lzwBuf[j] = binary.LittleEndian.Uint16(data[offset : offset+2])
			offset += 2
		}
		//os.WriteFile("../data/LB_EN/IMAGE/2.ori.lzw", data[offsetTemp:offset], 0666)
		rawBuf := decompressLZW(lzwBuf, int(block.RawSize))
		//os.WriteFile("../data/LB_EN/IMAGE/2.ori", rawBuf, 0666)
		//panic("11")
		outputBuf.Write(rawBuf)
	}
	return outputBuf.Bytes()

}

// Compress 压缩数据
//  Description 压缩数据
//  Param data []byte 未压缩数据
//  Param size int 分块大小
//  Return compressed
//  Return outputInfo
//
func Compress(data []byte, size int) (compressed []byte, outputInfo *CzOutputInfo) {

	if size == 0 {
		size = 0xFEFD
	}
	var partData []uint16
	offset := 0
	count := 0
	last := ""
	tmp := make([]byte, 2)
	outputBuf := &bytes.Buffer{}
	outputInfo = &CzOutputInfo{
		TotalRawSize: len(data),
		BlockInfo:    make([]CzBlockInfo, 0),
	}
	for {
		count, partData, last = compressLZW(data[offset:], size, last)
		if count == 0 {
			break
		}
		offset += count
		for _, d := range partData {
			binary.LittleEndian.PutUint16(tmp, d)
			outputBuf.Write(tmp)
		}

		outputInfo.BlockInfo = append(outputInfo.BlockInfo, CzBlockInfo{
			CompressedSize: uint32(len(partData)),
			RawSize:        uint32(count),
		})
		outputInfo.FileCount++
	}
	outputInfo.TotalCompressedSize = outputBuf.Len() / 2

	return outputBuf.Bytes(), outputInfo
}
