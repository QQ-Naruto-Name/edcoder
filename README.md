# edcoder
edcoder package include encoder and decoder implemented by serval format,such as toml,ymal,ini,xml and json.

#decocder
```
  type User struct {
    Name string
    EName string `json:"company"`
  }
  
  // data
  d, err := NewDecoder(SetDecoderExt('json'), SetDecoderData(`{"name":"june","company":"abc"}`))
  if nil != err {
    ...
  }
  err = d.Decode(&User{})
  ...
  
  fi, err := os.Open("define.json")
	if err != nil {
		return "", err
	}
	defer fi.Close()
  
  d, err := NewDecoder(SetDecoderExt('json'), SetDecoderReader(fi))
  if nil != err {
    ...
  }
  err = d.Decode(&User{})
  ...
```

#encode
```
  type User struct {
    Name string
    EName string `json:"company"`
  }
  
  fi, err = os.OpenFile("define.json", os.O_RDWR|os.O_APPEND|os.O_CREATE, 066) 
  if err != nil {
    return err
  }
  defer fi.Close()
  
  e, err := NewEncoder(SetEncoderExt("json"), SetEncoderObj(&User{"june","abc"},SetEncoderWriter(fi)))
  if nil != err {
    ...
  }
  err = d.Encode()
  ...
```
