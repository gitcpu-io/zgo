package zgoutils

import (
  "crypto/tls"
  "encoding/xml"
  "io"
  "math/rand"
  "reflect"
  "sort"
  "strconv"
  "strings"
  "time"

  "github.com/parnurzeal/gorequest"
)

type BodyMap map[string]interface{}

//设置参数
//    value：仅支持类型 string,int,int64,float32,float64,ptr,struct,slice,map 类型，其他类型一律设置空字符串
func (bm BodyMap) Set(key string, value interface{}) {
  vKind := reflect.ValueOf(value).Kind()
  switch vKind {
  case reflect.String:
    bm[key] = value.(string)
  case reflect.Int:
    bm[key] = Int2String(value.(int))
  case reflect.Int64:
    bm[key] = Int642String(value.(int64))
  case reflect.Float32:
    bm[key] = Float32ToString(value.(float32))
  case reflect.Float64:
    bm[key] = Float64ToString(value.(float64))
  case reflect.Ptr:
    bm[key] = value
  case reflect.Struct:
    bm[key] = value
  case reflect.Map:
    bm[key] = value
  case reflect.Slice:
    bm[key] = value
  default:
    bm[key] = ""
  }
}

//获取参数
func (bm BodyMap) Get(key string) string {
  if bm == nil {
    return ""
  }
  var (
    value interface{}
    ok    bool
    v     string
  )
  if value, ok = bm[key]; !ok {
    return ""
  }
  if v, ok = value.(string); ok {
    return v
  }
  return jsonToString(value)
}

//删除参数
func (bm BodyMap) Remove(key string) {
  delete(bm, key)
}

type xmlMapEntry struct {
  XMLName xml.Name
  Value   string `xml:",chardata"`
}

func (bm BodyMap) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
  if len(bm) == 0 {
    return nil
  }
  var (
    value string
    vKind reflect.Kind
  )
  if err = e.EncodeToken(start); err != nil {
    return
  }
  for k, v := range bm {
    vKind = reflect.ValueOf(v).Kind()
    switch vKind {
    case reflect.String:
      value = v.(string)
    case reflect.Int:
      value = Int2String(v.(int))
    case reflect.Int64:
      value = Int642String(v.(int64))
    case reflect.Float32:
      value = Float32ToString(v.(float32))
    case reflect.Float64:
      value = Float64ToString(v.(float64))
    default:
      value = ""
    }
    err := e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: value})
    if err != nil {
      return err
    }
  }
  return e.EncodeToken(start.End())
}

func (bm *BodyMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
  for {
    var e xmlMapEntry
    err = d.Decode(&e)
    if err == io.EOF {
      break
    } else if err != nil {
      return
    }
    bm.Set(e.XMLName.Local, e.Value)
  }
  return
}

// ("bar=baz&foo=quux") sorted by key.
func (bm BodyMap) EncodeWechatSignParams(apiKey string) string {
  var (
    buf     strings.Builder
    keyList []string
  )
  keyList = make([]string, 0, len(bm))
  for k := range bm {
    keyList = append(keyList, k)
  }
  sort.Strings(keyList)
  for _, k := range keyList {
    buf.WriteString(k)
    buf.WriteByte('=')
    buf.WriteString(bm.Get(k))
    buf.WriteByte('&')
  }
  buf.WriteString("key")
  buf.WriteByte('=')
  buf.WriteString(apiKey)
  return buf.String()
}

// ("bar=baz&foo=quux") sorted by key.
func (bm BodyMap) EncodeAliPaySignParams() string {
  var (
    buf     strings.Builder
    keyList []string
  )
  keyList = make([]string, 0, len(bm))
  for k := range bm {
    keyList = append(keyList, k)
  }
  sort.Strings(keyList)
  for _, k := range keyList {
    buf.WriteString(k)
    buf.WriteByte('=')
    buf.WriteString(bm.Get(k))
    buf.WriteByte('&')
  }
  return buf.String()[:buf.Len()-1]
}

//HttpAgent
func HttpAgent() (agent *gorequest.SuperAgent) {
  return gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

//获取随机字符串
//    length：字符串长度
func GetRandomString(length int) string {
  str := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
  var (
    result []byte
    b      []byte
    r      *rand.Rand
  )
  b = []byte(str)
  r = rand.New(rand.NewSource(time.Now().UnixNano()))
  for i := 0; i < length; i++ {
    result = append(result, b[r.Intn(len(b))])
  }
  return string(result)
}

//Float64转字符串
//    floatNum：float64数字
//    prec：精度位数（不传则默认float数字精度）
func Float64ToString(floatNum float64, prec ...int) (floatStr string) {
  if len(prec) > 0 {
    floatStr = strconv.FormatFloat(floatNum, 'f', prec[0], 64)
    return
  }
  floatStr = strconv.FormatFloat(floatNum, 'f', -1, 64)
  return
}

//Float32转字符串
//    floatNum：float32数字
//    prec：精度位数（不传则默认float数字精度）
func Float32ToString(floatNum float32, prec ...int) (floatStr string) {
  if len(prec) > 0 {
    floatStr = strconv.FormatFloat(float64(floatNum), 'f', prec[0], 32)
    return
  }
  floatStr = strconv.FormatFloat(float64(floatNum), 'f', -1, 32)
  return
}

//Int转字符串
func Int2String(intNum int) (intStr string) {
  intStr = strconv.Itoa(intNum)
  return
}

//Int64转字符串
func Int642String(intNum int64) (int64Str string) {
  //10, 代表10进制
  int64Str = strconv.FormatInt(intNum, 10)
  return
}

//解密填充模式（去除补全码） PKCS7UnPadding
//解密时，需要在最后面去掉加密时添加的填充byte
func PKCS7UnPadding(origData []byte) (bs []byte) {
  length := len(origData)
  unPaddingNumber := int(origData[length-1]) //找到Byte数组最后的填充byte 数字
  if unPaddingNumber <= 16 {
    bs = origData[:(length - unPaddingNumber)] //只截取返回有效数字内的byte数组
  } else {
    bs = origData
  }
  return
}

//解密填充模式（去除补全码） PKCS5UnPadding
//解密时，需要在最后面去掉加密时添加的填充byte
func PKCS5UnPadding(origData []byte) (bs []byte) {
  length := len(origData)
  unPaddingNumber := int(origData[length-1]) //找到Byte数组最后的填充byte
  if unPaddingNumber <= 16 {
    bs = origData[:(length - unPaddingNumber)] //只截取返回有效数字内的byte数组
  } else {
    bs = origData
  }
  return
}

func jsonToString(v interface{}) (str string) {
  if v == nil {
    return ""
  }
  var (
    bs  []byte
    err error
  )
  if bs, err = Utils.Marshal(v); err != nil {
    return ""
  }

  if str = string(bs); str == "" {
    return ""
  }
  return
}
