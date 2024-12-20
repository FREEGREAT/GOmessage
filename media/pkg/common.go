package pkg

import (
	"bytes"
	"mime/multipart"

	"github.com/spf13/viper"
)

type CustomFile struct {
	*bytes.Reader
}

const error_file_size = 0

func (f *CustomFile) Close() error {
	return nil
}
func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func NewCustomFile(data []byte) *CustomFile {
	return &CustomFile{bytes.NewReader(data)}
}

func ByteToMultipart(data []byte, filename string) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)

	if err != nil {

		return nil, err

	}

	_, err = part.Write(data)

	if err != nil {

		return nil, err

	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
