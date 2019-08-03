// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/orivil/morgine/log"
	"github.com/orivil/morgine/router"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"
)

// 上传文件时最大内存使用量, 超过最大内存则暂存于硬盘中
var MaxUploadMemory = int64(10 << 20) // 10MB

// 默认消息格式
var MessageDataType = DataTypeJson

type Context struct {
	Writer        http.ResponseWriter
	Request       *http.Request
	rvs           router.Values
	paths         url.Values
	query         url.Values
	form          url.Values
	Values        map[string]interface{}
	multipartForm *multipart.Form
	handler       *Handler
	idx           int
}

func initContext(ctx *Context, res http.ResponseWriter, req *http.Request, rvs router.Values, h *Handler) *Context {
	ctx.Writer = res
	ctx.Request = req
	ctx.rvs = rvs
	ctx.paths = nil
	ctx.query = nil
	ctx.form = nil
	ctx.Values = make(map[string]interface{})
	ctx.multipartForm = nil
	ctx.handler = h
	ctx.idx = 0
	return ctx
}

func (c *Context) handle() {
	if ln := len(c.handler.middles); c.idx == ln {
		c.handler.HandleFunc(c)
	} else if c.idx < ln {
		c.handler.middles[c.idx].HandleFunc(c)
		c.idx++
		c.handle()
	}
}

// 立即执行下一个处理函数
func (c *Context) HandleNext() {
	c.idx++
	c.handle()
}

// 结束处理链, 用于中间件中结束请求处理并立即返回结果, 用在 action 中没有任何效果
func (c *Context) Abort() {
	c.idx = len(c.handler.middles) + 1
}

func (c *Context) abortWithError(depth int, err error) error {
	http.Error(c.Writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	c.Abort()
	return log.Error.Output(depth+1, err.Error())
}

// Error 方法会将错误信息记录到 log.Error 中, 并返回 500 状态码到客户端, 且会调用 Abort 方法结束处理链
func (c *Context) Error(err error) error {
	return c.abortWithError(2, err)
}

func (c *Context) Set(key string, value interface{}) {
	if c.Values == nil {
		c.Values = make(map[string]interface{})
	}
	c.Values[key] = value
}

func (c *Context) Get(key string) (value interface{}) {
	return c.Values[key]
}

func (c *Context) Del(key string) {
	delete(c.Values, key)
}

// 获得 path 中的参数
func (c *Context) Path() url.Values {
	if c.paths == nil {
		c.paths = c.rvs()
	}
	return c.paths
}

// 获得 query 参数
func (c *Context) Query() url.Values {
	if c.query == nil {
		c.query = c.Request.URL.Query()
	}
	return c.query
}

// 获得 form 参数, 编码格式: "application/x-www-form-urlencoded"
func (c *Context) Form() url.Values {
	if c.form == nil {
		_, err := c.parseForm()
		if err != nil {
			c.abortWithError(2, err)
		}
	}
	return c.form
}

func (c *Context) parseForm() (url.Values, error) {
	if c.form == nil {
		err := c.Request.ParseForm()
		if err != nil {
			return nil, err
		} else {
			c.form = c.Request.Form
		}
	}
	return c.form, nil
}

// 获得 form 参数, 编码格式: "multipart/form-data"
func (c *Context) MultipartForm() *multipart.Form {
	if c.multipartForm == nil {
		_, err := c.parseMultipartForm()
		if err != nil {
			c.abortWithError(2, err)
		}
	}
	return c.multipartForm
}

func (c *Context) parseMultipartForm() (*multipart.Form, error) {
	if c.multipartForm == nil {
		err := c.Request.ParseMultipartForm(MaxUploadMemory)
		if err != nil {
			return nil, err
		} else {
			c.multipartForm = c.Request.MultipartForm
		}
	}
	return c.multipartForm, nil
}

func (c *Context) NotFound() {
	http.NotFound(c.Writer, c.Request)
}

func (c *Context) Redirect(url string, code int) {
	http.Redirect(c.Writer, c.Request, url, code)
}

func (c *Context) WriteString(str string) (int, error) {
	return c.Writer.Write([]byte(str))
}

func (c *Context) Write(data []byte) (int, error) {
	return c.Writer.Write(data)
}

// 发送 JSON 数据, 不管是否出现错误都会调用 Abort 方法
func (c *Context) SendJSON(v interface{}) error {
	return c.sendData(DataTypeJson, v)
}

// 发送 XML 数据, 不管是否出现错误都会调用 Abort 方法
func (c *Context) SendXML(v interface{}) error {
	return c.sendData(DataTypeXml, v)
}

// 发送自定义类型的消息, 不管是否出现错误都会调用 Abort 方法
func (c *Context) Message(t MsgType, msg string) error {
	return c.message(t, msg)
}

// 获得 message 数据, 用于生成文档数据
func MsgData(t MsgType, msg string) MAP {
	return MAP{"message": &Message{Type: t, Content: msg}}
}

// 发送 info 类型的消息, 不管是否出现错误都会调用 Abort 方法
func (c *Context) MsgInfo(msg string) error {
	return c.message(MsgTypeInfo, msg)
}

// 发送 error 类型的消息, 不管是否出现错误都会调用 Abort 方法
func (c *Context) MsgError(msg string) error {
	return c.message(MsgTypeError, msg)
}

// 发送 success 类型的消息, 不管是否出现错误都会调用 Abort 方法
func (c *Context) MsgSuccess(msg string) error {
	return c.message(MsgTypeSuccess, msg)
}

// 发送 warning 类型的消息, 不管是否出现错误都会调用 Abort 方法
func (c *Context) MsgWarning(msg string) error {
	return c.message(MsgTypeWarning, msg)
}

func (c *Context) message(t MsgType, msg string) error {
	return c.sendData(MessageDataType, MsgData(t, msg))
}

// 验证参数并将参数解析到 v 中, v 必须经过注册
//
// 如果设置了参数验证且参数验证失败, 将会返回 param.ValidatorErr 错误
func (c *Context) Unmarshal(v ...interface{}) error {
	if ln := len(c.handler.middles); c.idx < ln {
		return c.handler.middles[c.idx].Doc.parser.unmarshal(v, c)
	} else if c.idx == ln {
		return c.handler.Doc.parser.unmarshal(v, c)
	}
	return errors.New("context aborted")
}

var dataBuffer = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

type dataType int

func (dt dataType) contentType() string {
	switch dt {
	case DataTypeXml:
		return "application/xml;charset=UTF-8"
	default:
		return "application/json;charset=UTF-8"
	}
}

const (
	DataTypeJson dataType = 1 + iota
	DataTypeXml
)

type encoder interface {
	Encode(v interface{}) error
}

func (c *Context) sendData(dt dataType, v interface{}) error {
	buf := dataBuffer.Get().(*bytes.Buffer)
	defer dataBuffer.Put(buf)
	defer c.Abort()
	buf.Reset()
	var encoder encoder
	switch dt {
	case DataTypeXml:
		encoder = json.NewEncoder(buf)
	default:
		encoder = xml.NewEncoder(buf)
	}
	err := encoder.Encode(v)
	if err != nil {
		return err
	} else {
		c.Writer.Header().Set("Content-Type", dt.contentType())
		_, err = io.Copy(c.Writer, buf)
		if err != nil {
			return err
		}
		return nil
	}
}
