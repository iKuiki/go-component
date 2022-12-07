package compress

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

// Zlib 以zlib压缩数据
func Zlib(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// DeZlib 以zlib解压数据
func DeZlib(data []byte) ([]byte, error) {
	dataReader := bytes.NewReader(data)
	gzipReader, err := zlib.NewReader(dataReader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()
	result, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}
	return result, nil
}
