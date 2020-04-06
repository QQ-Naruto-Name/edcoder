package edcoder

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type DMocker struct {
	ConfDecoder //有时候只为测试某一个方法，则DMocker实现即可，避开实现完整的接口方法集
	m           mock.Mock
}

type EMocker struct {
	ConfEncoder //有时候只为测试某一个方法，则EMocker实现即可，避开实现完整的接口方法集
	m           mock.Mock
}

func TestED2Map(t *testing.T) {
	tests := []struct {
		ext string
		str string
		obj map[string]interface{}
	}{
		{
			ext: "json",
			str: `
                {
                    "title":"Json Example",
                    "owner":{
                        "name" : "Tom Preston-Werner",
                        "dob" : "1979-05-27T07:32:00-08:00"
                    },
                    "database":{
                        "server" : "192.168.1.1",
                        "ports" : [ 8001, 8001, 8002 ],
                        "connection_max" : 5000,
                        "enabled" : true
                    },
                    "servers":{
                        "alpha":{
                            "ip": "10.0.0.1",
                            "dc": "eqdc10"
                        },
                        "beta":{
                            "ip": "10.0.0.2",
                            "dc": "eqdc10"
                        }
                    },
                    "hosts":[
                        "alpha",
                        "omega"
                    ]
                }
                `,
			obj: make(map[string]interface{}),
		},
		{
			ext: "toml",
			str: `
                # This is a TOML document.

                title = "TOML Example"

                [owner]
                name = "Tom Preston-Werner"
                dob = 1979-05-27T07:32:00-08:00 # First class dates

                [database]
                server = "192.168.1.1"
                ports = [ 8001, 8001, 8002 ]
                connection_max = 5000
                enabled = true

                [servers]

                # Indentation (tabs and/or spaces) is allowed but not required
                [servers.alpha]
                IP = "10.0.0.1"
                dc = "eqdc10"

                [servers.beta]
                Ip = "10.0.0.2"
                dc = "eqdc10"

                [clients]
                data = [ ["gamma", "delta"], [1, 2] ]

                # Line breaks are OK when inside arrays
                hosts = [
                "alpha",
                "omega"
                ]
                `,
			obj: make(map[string]interface{}),
		},
		{
			ext: "yaml",
			str: `
                # This is a YAML document.

                title: "YAML Example"
                owner:
                    name: Tom Preston-Werner
                    dob: 1979-05-27T07:32:00-08:00 # First class dates
                database:
                    server: 192.168.1.1
                    ports: [8001,8001,8002]
                    enabled: true
                    connection_max: 5000

                servers:
                    alpha:
                        IP: 10.0.0.1
                        dc: &tag eqdc10
                    beta:
                        IP: 10.0.0.2
                        dc: *tag
                clients:
                    data:
                        -
                            - gamma
                            - delta
                        -
                            - 1
                            - 2

                hosts:
                    - alpha
                    - omega
                `,
			obj: make(map[string]interface{}),
		},
		{
			ext: "xml", // xml不支持map类型的序列化和反序列化，但可自己实现该功能 https://stackoverflow.com/questions/30928770/marshall-map-to-xml-in-go
			str: `
                <persons>
                    <person name="studygolang" age="27">
                        <career>码农</career>
                        <interest>编程</interest>
                        <interesta>编程a</interesta>
                    </person>
                </persons>
                `,
			obj: make(map[string]interface{}),
		},
		{
			ext: "ini", //不支持map
			str: `
                ; A comment line
                outter = string
                [Section]
                enabled = true
                path = /usr/local
            `,
			obj: make(map[string]interface{}),
		},
	}

	for _, test := range tests {
		d, err := NewDecoder(SetDecoderExt(test.ext), SetDecoderData(test.str))
		assert.NoError(t, err, "new decoder failed %s error:%v", test.ext, err)
		dmk := &DMocker{ConfDecoder: d}
		err = dmk.Decode(&test.obj)
		assert.NoError(t, err, "decode failed %s error:%v", test.ext, err)

		e, err := NewEncoder(SetEncoderExt(test.ext), SetEncoderObj(test.obj))
		assert.NoError(t, err, "new encoder failed %s error:%v", test.ext, err)
		emk := &EMocker{ConfEncoder: e}
		err = emk.Encode()
		assert.NoError(t, err, "encode failed %s error:%v", test.ext, err)
	}
}

func TestED2Struct(t *testing.T) {
	type Sub struct {
		AA string `yaml:"AA" xml:"AA"`
		BB string //无标签，照样解析该字段
	}

	type Dst struct {
		A     string `xml:"a" `
		B     string `yaml:"B" xml:"B"`
		Slice []Sub  `xml:"slice>sub"`
	}

	type ResultD struct {
		XMLName xml.Name `xml:"Dst"` //节点persons
		A       string   `xml:"a" `
		B       string   `yaml:"B" xml:"B"`
		Slice   []Sub    `xml:"slice>sub"`
	}

	type Person struct {
		Name      string   `xml:"name,attr"` //person标签属性名为name的属性值
		Age       int      `xml:"age,attr"`
		Career    string   `xml:"career"`             //person中标签名为career的值 若不定义标签则该字段为空，不填充
		Interests []string `xml:"interests>interest"` //节点interests下的interest数组
	}
	type Result struct {
		XMLName xml.Name `xml:"persons"` //节点persons
		Persons []Person `xml:"person"`  //多个person节点
	}

	type Section struct {
		Enable bool
		A      string `gcfg:"aa"`
		B      string
	}

	type Ini struct { //字段名称必须大写，视为导出字段 标签gcfg
		Section //不能定义命名字段，必须匿名，否则解析报错can't store data at section "section"   ( expected section header)
	}

	tests := []struct {
		ext string
		str string
		obj interface{}
	}{
		{
			ext: "json",
			str: `
						{
							"a":"a2A",
							"B":"B2B",
							"Slice":[
								{
									"Aa":"Aa2AA",
									"Bb":"Bb2BB"
								},
								{
									"aa":"aa2AA",
									"bb":"bb2BB"
								},
								{
									"AA":"AA2AA",
									"BB":"BB2BB"
								}
							]
						}
						`,
			obj: &Dst{},
		},
		{
			ext: "toml",
			str: `
				a = "a2A"
				B = "B2B"
				
				[[Slice]]
						Aa = "Aa2AA"
						bB = "bB2BB"
				
				[[Slice]]
						aa = "aa2AA"
						bb = "bb2BB"

				[[Slice]]
						AA = "AA2AA"
						BB = "BB2BB"
			`,
			obj: &Dst{},
		},
		{
			ext: "yaml", //默认只支持yaml文件中小写的key，除非struct中对key进行特定标签设置
			str: `
                a: a2A
                B: B2B
                slice: 
                    -
                        Aa: &tag Aa2AA
                        Bb: Bb2BB
                    -
                        aa: aa2AA
                        bb: bb2BB
                    -
                        AA: AA2AA
                        BB: BB2BB
                    -
                        Aa: *tag
                        bB: bB2BB
            `,
			obj: &Dst{},
		},
		{
			ext: "xml", //只支持struct解析及编码；struct中字段想要与xml中对应输出，则必须进行xml标签设定，否则不可decode到结构体中
			str: `
                <?xml version="1.0" encoding="UTF-8"?>
                <persons>
                    <person name="polaris" age="28">
                        <career>无业游民</career>
                        <interests>
                            <interest>编程</interest>
                            <interest>下棋</interest>
                        </interests>
                    </person>
                    <person name="studygolang" age="27">
                        <career>码农</career>
                        <interests>
                            <interest>编程</interest>
                            <interest>下棋</interest>
                        </interests>
                    </person>
                </persons>
            `,
			obj: &Result{},
		},
		{
			ext: "xml",
			str: `
                <?xml version="1.0" encoding="UTF-8"?>
                <Dst>
                    <a>a2A</a>
                    <B>B2B</B>
                    <slice>
                        <sub>
                            <AA>无业游民</AA>
                            <BB>码农</BB>
                        </sub>
                        <sub>
                            <AA>水电工</AA>
                            <BB>物业</BB>
                        </sub>
                    </slice>
                </Dst>
            `,
			obj: &ResultD{},
		},
		{
			ext: "ini", //https://godoc.org/gopkg.in/gcfg.v1
			str: `
                ; A comment line
                [section]
                enable = true
                aa = a_string
                b = b_string
                
                
            `,
			obj: &Ini{},
		},
	}

	for _, test := range tests {
		d, err := NewDecoder(SetDecoderExt(test.ext), SetDecoderData(test.str))
		assert.NoError(t, err, "new decoder failed %s error:%v", test.ext, err)
		dmk := &DMocker{ConfDecoder: d}
		err = dmk.Decode(test.obj)
		assert.NoError(t, err, "decode failed %s error:%v", test.ext, err)

		e, err := NewEncoder(SetEncoderExt(test.ext), SetEncoderObj(test.obj))
		assert.NoError(t, err, "new encoder failed %s error:%v", test.ext, err)
		emk := &EMocker{ConfEncoder: e}
		err = emk.Encode()
		assert.NoError(t, err, "encode failed %s error:%v", test.ext, err)
	}
}
