package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// Gzip 以gzip压缩数据
func Gzip(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
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

// DeGzip 以gzip解压数据
func DeGzip(data []byte) ([]byte, error) {
	dataReader := bytes.NewReader(data)
	gzipReader, err := gzip.NewReader(dataReader)
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
