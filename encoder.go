package edcoder

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// Encoder 配置
type Encoder struct {
	obj    interface{}
	ext    string
	writer io.Writer
}

// EncoderOption 配置参数接口
type EncoderOption interface {
	apply(*Encoder)
}

type opf func(*Encoder)

func (invoke opf) apply(r *Encoder) {
	invoke(r)
}

func newEncoderOption(t func(o *Encoder)) EncoderOption {
	return opf(t)
}

// SetEncoderExt 设置配置格式
func SetEncoderExt(ext string) EncoderOption {
	return newEncoderOption(func(o *Encoder) {
		o.ext = ext
	})
}

// SetEncoderObj 设置配置对象
func SetEncoderObj(obj interface{}) EncoderOption {
	return newEncoderOption(func(o *Encoder) {
		o.obj = obj
	})
}

// SetEncoderWriter 设置配置输出
func SetEncoderWriter(w io.Writer) EncoderOption {
	return newEncoderOption(func(o *Encoder) {
		o.writer = w
	})
}

// NewEncoder 创建解析器
func NewEncoder(opts ...EncoderOption) (*Encoder, error) {
	d := &Encoder{}
	for _, v := range opts {
		o := v.(opf)
		o.apply(d)
	}

	if "" == d.ext {
		return nil, errors.New("ext must be set")
	}
	if nil == d.obj {
		return nil, errors.New("obj must be set")
	}

	return d, nil
}

// ConfEncoder 配置编码器
type ConfEncoder interface {
	Encode() error
}

// Encode 通用编码
func (d *Encoder) Encode() error {
	var err error
	var fi *os.File
	if nil == d.writer {
		fi, err = os.OpenFile("default."+d.ext, os.O_RDWR|os.O_APPEND|os.O_CREATE, 066)
		if err != nil {
			return err
		}
		defer fi.Close()

		d.writer = fi
	}

	switch d.ext {
	case "toml": //https://github.com/BurntSushi/toml
		e := toml.NewEncoder(d.writer)
		e.Indent = "    "
		return e.Encode(d.obj)
	case "yaml": //https://gopkg.in/yaml.v2
		e := yaml.NewEncoder(d.writer)
		return e.Encode(d.obj)
	case "xml":
		switch d.obj.(type) {
		case map[string]interface{}:
			return errors.New("type not supported by encoder")
		default:
			e := xml.NewEncoder(d.writer)
			e.Indent("", "    ")
			return e.Encode(d.obj)
		}

	case "json":
		e := json.NewEncoder(d.writer)
		e.SetIndent("", "    ")
		return e.Encode(d.obj)
	}

	return errors.New("ext not be supported by encoder")
}
