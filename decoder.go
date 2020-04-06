package edcoder

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/gcfg.v1"
	"gopkg.in/yaml.v2"
)

// Decoder 配置
type Decoder struct {
	ext    string
	data   string
	reader io.Reader
}

// Reader 读取器
//type Reader func(fileName string) (string, error)

// DecoderOption 配置参数接口
type DecoderOption interface {
	apply(*Decoder)
}

type op func(*Decoder)

func (invoke op) apply(r *Decoder) {
	invoke(r)
}

func newDecoderOption(t func(o *Decoder)) DecoderOption {
	return op(t)
}

// SetDecoderData 设置配置数据
func SetDecoderData(data string) DecoderOption {
	return newDecoderOption(func(o *Decoder) {
		o.data = data
	})
}

// SetDecoderExt 设置配置格式
func SetDecoderExt(ext string) DecoderOption {
	return newDecoderOption(func(o *Decoder) {
		o.ext = ext
	})
}

// SetDecoderReader 设置配置读取器
func SetDecoderReader(r io.Reader) DecoderOption {
	return newDecoderOption(func(o *Decoder) {
		o.reader = r
	})
}

// NewDecoder 创建解析器
func NewDecoder(opts ...DecoderOption) (*Decoder, error) {
	d := &Decoder{}
	for _, v := range opts {
		o := v.(op)
		o.apply(d)
	}

	if "" == d.ext {
		return nil, errors.New("ext must be set")
	}

	if "" != d.data {
		return d, nil
	} else if nil != d.reader {
		data, err := readConf(d.reader)
		if nil != err {
			return nil, err
		}
		d.data = data

		return d, nil
	} else {
		_, err := os.Stat("default." + d.ext)
		if nil != err {
			return nil, err
		}

		// 设置data
		data, err := readFile("default." + d.ext)
		if nil != err {
			return nil, err
		}
		d.data = data

		return d, nil
	}
}

// ConfDecoder 配置解析器
type ConfDecoder interface {
	Decode(v interface{}) error
}

// Decode 通用解析
func (d *Decoder) Decode(v interface{}) error {
	switch d.ext {
	case "toml": //https://github.com/BurntSushi/toml
		_, err := toml.Decode(d.data, v)
		return err
	case "yaml": //https://gopkg.in/yaml.v2
		return yaml.Unmarshal([]byte(d.data), v)
	case "ini": //gopkg.in/gcfg.v1
		switch v.(type) {
		case *map[string]interface{}:
			return errors.New("type not supported by decoder")
		default:
			return gcfg.ReadStringInto(v, d.data)
		}
	case "xml":
		switch v.(type) {
		case *map[string]interface{}:
			return errors.New("type not supported by decoder")
		default:
			return xml.Unmarshal([]byte(d.data), v)
		}
	case "json":
		return json.Unmarshal([]byte(d.data), v)
	}

	return errors.New("ext not be support")
}
