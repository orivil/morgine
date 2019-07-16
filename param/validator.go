// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package param

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

// 参数分割符
const Separator = " "

// 参数标签
const (
	TagRequired = "required"
	TagNum      = "num"
	TagLen      = "len"
	TagItem     = "item"
	TagEmail    = "email"
	TagEnum     = "enum"
	TagRegexp   = "reg"
	TagFileByte = "size-Byte"
	TagFileKB   = "size-KB"
	TagFileMB   = "size-MB"
	TagFileExt  = "exts"
	TagFileType = "mime"
)

var emailPattern = "[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?"
var emailExp = regexp.MustCompile(emailPattern)

type condition struct {
	required *string // nil means not required

	// number condition
	maxNum   *float64
	eqMaxNum *struct{}
	minNum   *float64
	eqMinNum *struct{}
	numMsgID *string

	// string condition
	maxLen   *int // max []rune length
	minLen   *int // min []rune length
	lenMsgID *string

	// file size condition
	maxFileByte   *int64
	minFileByte   *int64
	fileSizeMsgID *string

	// file extension condition
	fileExtensions []string
	fileExtMsgID   *string

	// file mime type condition
	fileMimeTypes []string
	fileMimeMsgID *string

	// slice condition
	maxItem   *int
	minItem   *int
	itemMsgID *string

	// enum condition
	enums     []string
	enumMsgID *string

	// regular condition
	pattern  *string
	regexp   *regexp.Regexp
	regMsgID *string
}

type info struct {
	Required *string `json:",omitempty"`

	// number condition
	MaxNum   *float64  `json:",omitempty"`
	EqMaxNum *struct{} `json:",omitempty"`
	MinNum   *float64  `json:",omitempty"`
	EqMinNum *struct{} `json:",omitempty"`
	NumMsgID *string   `json:",omitempty"`

	// string condition
	MaxLen   *int    `json:",omitempty"`
	MinLen   *int    `json:",omitempty"`
	LenMsgID *string `json:",omitempty"`

	// file size condition
	MaxFileByte   *int64  `json:",omitempty"`
	MinFileByte   *int64  `json:",omitempty"`
	FileSizeMsgID *string `json:",omitempty"`

	// file extension condition
	FileExtensions []string `json:",omitempty"`
	FileExtMsgID   *string  `json:",omitempty"`

	// file mime type condition
	FileMimeTypes []string `json:",omitempty"`
	FileMimeMsgID *string  `json:",omitempty"`

	// slice condition
	MaxItem   *int    `json:",omitempty"`
	MinItem   *int    `json:",omitempty"`
	ItemMsgID *string `json:",omitempty"`

	// enum condition
	Enums     []string `json:",omitempty"`
	EnumMsgID *string  `json:",omitempty"`

	// regular condition
	Pattern  *string `json:",omitempty"`
	regexp   *regexp.Regexp
	RegMsgID *string `json:",omitempty"`
}

func (c *condition) getInfo() *info {
	pointer := reflect.ValueOf(c).Pointer()
	i := (*info)(unsafe.Pointer(pointer))
	return i
}

func (c *condition) Required(msg string) *condition {
	if msg == "" {
		msg = "required"
	}
	c.required = &msg
	return c
}

// 字符串长度, 按 []rune 长度计算
func (c *condition) Len(min, max int, msg string) *condition {
	if msg == "" {
		msg = fmt.Sprintf("string length %d-%d", min, max)
	}
	c.minLen = &min
	c.maxLen = &max
	c.lenMsgID = &msg
	return c
}

func (c *condition) Email(msg string) *condition {
	if msg == "" {
		msg = "email address format incorrect"
	}
	c.pattern = &emailPattern
	c.regexp = emailExp
	c.regMsgID = &msg
	return c
}

func (c *condition) Regexp(pattern, msg string) *condition {
	if msg == "" {
		msg = "string format incorrect"
	}
	c.pattern = &pattern
	c.regexp = regexp.MustCompile(pattern)
	c.regMsgID = &msg
	return c
}

func (c *condition) Item(min, max int, msg string) *condition {
	if msg == "" {
		msg = fmt.Sprintf("%d-%d items", min, max)
	}
	c.minItem = &min
	c.maxItem = &max
	c.itemMsgID = &msg
	return c
}

func (c *condition) Enums(elements []string, msg string) *condition {
	if msg == "" {
		msg = fmt.Sprintf("one of %v", elements)
	}
	c.enums = elements
	c.enumMsgID = &msg
	return c
}

// set MIME type condition, the condition could be a specific
// type(e.g. image/png image/jpeg text/plain application/octet-stream)
// or just the main type(e.g. image text application)
func (c *condition) FileMimeTypes(types []string, msg string) *condition {
	if msg == "" {
		msg = fmt.Sprintf("file types %v", types)
	}
	c.fileMimeTypes = types
	c.fileMimeMsgID = &msg
	return c
}

// 限制文件大小, 单位: MB
func (c *condition) FileSizeMB(min, max int64, msg string) *condition {
	if msg == "" {
		msg = fmt.Sprintf("%d-%dMB", min, max)
	}
	min = min << 10
	max = max << 10
	return c.FileSizeKB(min, max, msg)
}

// 限制文件大小, 单位: KB, tag 标签中设置的文件大小限制默认调用此方法.
// (e.g. size:"20-500" size-msg:"20-500KB")
func (c *condition) FileSizeKB(min, max int64, msg string) *condition {
	if msg == "" {
		msg = fmt.Sprintf("%d-%dKB", min, max)
	}
	min = min << 10
	max = max << 10
	return c.FileSizeByte(min, max, msg)
}

// 限制文件大小, 单位: Byte
func (c *condition) FileSizeByte(min, max int64, msg string) *condition {
	if msg == "" {
		msg = fmt.Sprintf("%d-%dByte", min, max)
	}
	c.minFileByte = &min
	c.maxFileByte = &max
	c.fileSizeMsgID = &msg
	return c
}

// 限制文件后缀名
func (c *condition) FileExts(exts []string, msg string) *condition {
	if msg == "" {
		msg = fmt.Sprintf("file extensions %v", exts)
	}
	c.fileExtensions = exts
	c.fileExtMsgID = &msg
	return c
}

// 限制数字最大值
func (c *condition) MaxNum(max float64, eq bool, msg string) *condition {
	if c.numMsgID == nil || *c.numMsgID != msg { // 检测是否是同一条消息(标签消息只能有一条)
		if msg == "" {
			if eq {
				msg = fmt.Sprintf("<=%v", max)
			} else {
				msg = fmt.Sprintf("<%v", max)
			}
		}
		if c.numMsgID != nil {
			msg = fmt.Sprintf("%s & %s", *c.numMsgID, msg)
		}
		c.numMsgID = &msg
	}
	c.maxNum = &max
	if eq {
		c.eqMaxNum = &struct{}{}
	}
	return c
}

// 限制数字最小值
func (c *condition) MinNum(min float64, eq bool, msg string) *condition {
	if c.numMsgID == nil || *c.numMsgID != msg { // 检测是否是同一条消息
		if msg == "" {
			if eq {
				msg = fmt.Sprintf("%v<=", min)
			} else {
				msg = fmt.Sprintf("%v<", min)
			}
		}
		if c.numMsgID != nil {
			msg = fmt.Sprintf("%s & %s", msg, *c.numMsgID)
		}
		c.numMsgID = &msg
	}
	c.minNum = &min
	if eq {
		c.eqMinNum = &struct{}{}
	}
	return c
}

func isCondition(tag reflect.StructTag) bool {
	tags := []string{
		TagRequired,
		TagNum,
		TagLen,
		TagItem,
		TagEmail,
		TagEnum,
		TagRegexp,
		TagFileByte,
		TagFileKB,
		TagFileMB,
		TagFileExt,
		TagFileType,
	}
	return containsTags(tag, tags)
}

func containsTags(tag reflect.StructTag, tags []string) bool {
	for _, t := range tags {
		if _, ok := tag.Lookup(t); ok {
			return true
		}
	}
	return false
}

func (c *condition) Syntax(syntax string) error {
	tag := reflect.StructTag(syntax)
	if required, ok := tag.Lookup(TagRequired); ok {
		c.Required(required)
	}
	if syntax, ok := tag.Lookup(TagLen); ok {
		min, max, err := getBetween(syntax)
		if err != nil {
			return err
		}
		msg := tag.Get(MsgName(TagLen))
		c.Len(*min, *max, msg)
	}
	if pattern, ok := tag.Lookup(TagRegexp); ok {
		msg := tag.Get(MsgName(TagRegexp))
		c.Regexp(pattern, msg)
	}
	if msg, ok := tag.Lookup(TagEmail); ok {
		c.Email(msg)
	}
	if syntax, ok := tag.Lookup(TagFileByte); ok {
		min, max, err := getBetween(syntax)
		if err != nil {
			return err
		}
		msg := tag.Get(MsgName(TagFileByte))
		c.FileSizeByte(int64(*min), int64(*max), msg)
	}
	if syntax, ok := tag.Lookup(TagFileKB); ok {
		min, max, err := getBetween(syntax)
		if err != nil {
			return err
		}
		msg := tag.Get(MsgName(TagFileKB))
		c.FileSizeKB(int64(*min), int64(*max), msg)
	}
	if syntax, ok := tag.Lookup(TagFileMB); ok {
		min, max, err := getBetween(syntax)
		if err != nil {
			return err
		}
		msg := tag.Get(MsgName(TagFileMB))
		c.FileSizeMB(int64(*min), int64(*max), msg)
	}
	if syntax, ok := tag.Lookup(TagFileType); ok {
		types := strings.Split(syntax, Separator)
		for key, typ := range types {
			types[key] = strings.TrimSpace(typ)
		}
		msg := tag.Get(MsgName(TagFileType))
		c.FileMimeTypes(types, msg)
	}
	if syntax, ok := tag.Lookup(TagFileExt); ok {
		exts := strings.Split(syntax, Separator)
		for key, ext := range exts {
			exts[key] = strings.TrimSpace(ext)
		}
		msg := tag.Get(MsgName(TagFileExt))
		c.FileExts(exts, msg)
	}
	if syntax, ok := tag.Lookup(TagNum); ok {
		min, max, eqMin, eqMax, e := readArea(syntax)
		if e != nil {
			return e
		}
		msg := tag.Get(MsgName(TagNum))
		if min != nil {
			c.MinNum(*min, eqMin, msg)
		}
		if max != nil {
			c.MaxNum(*max, eqMax, msg)
		}
	}
	if syntax, ok := tag.Lookup(TagItem); ok {
		min, max, err := getBetween(syntax)
		if err != nil {
			return err
		}
		msg := tag.Get(MsgName(TagItem))
		c.Item(*min, *max, msg)
	}
	if syntax, ok := tag.Lookup(TagEnum); ok {
		enums := strings.Split(syntax, Separator)
		for key, enum := range enums {
			enums[key] = strings.TrimSpace(enum)
		}
		msg := tag.Get(MsgName(TagEnum))
		c.Enums(enums, msg)
	}
	return nil
}

type Validator struct {
	ptr        uintptr
	conditions map[uintptr]*condition
}

func (v *Validator) offsetField(offset uintptr) *condition {

	c := &condition{}
	v.conditions[offset] = c
	return c
}

func (v *Validator) Field(pointer interface{}) *condition {
	ptr := reflect.ValueOf(pointer).Pointer()
	offset := ptr - v.ptr
	return v.offsetField(offset)
}

func NewValidator(schema interface{}) *Validator {
	return &Validator{
		ptr:        reflect.ValueOf(schema).Pointer(),
		conditions: make(map[uintptr]*condition),
	}
}

type ValidatorErr struct {
	Field   string
	Kind    ConditionKind
	Enums   []string `json:",omitempty"`
	Max     *float64 `json:",omitempty"`
	EQMax   *bool    `json:",omitempty"`
	Min     *float64 `json:",omitempty"`
	EQMin   *bool    `json:",omitempty"`
	Message string
}

type ConditionKind int

func (re *ValidatorErr) Error() string {
	return fmt.Sprintf("%s: %s", re.Field, re.Message)
}

const (
	ConditionRequired ConditionKind = iota
	ConditionInvalidNumber
	ConditionInvalidBoolean
	ConditionStringRegexp
	ConditionStringLength
	ConditionNumber
	ConditionItem
	ConditionFileSize
	ConditionFileExtensions
	ConditionFileMimeTypes
	ConditionEnums
)

var Conditions = map[ConditionKind]string{
	ConditionRequired:       "required",
	ConditionInvalidNumber:  "invalid-number",
	ConditionInvalidBoolean: "invalid-boolean",
	ConditionStringRegexp:   "string-regexp",
	ConditionStringLength:   "string-length",
	ConditionNumber:         "number",
	ConditionItem:           "item-length",
	ConditionFileSize:       "file-size",
	ConditionFileExtensions: "file-extensions",
	ConditionFileMimeTypes:  "file-mime-types",
	ConditionEnums:          "enums",
}

// 用非空指针表示 equal, 空指针表示 not equal
func Equal() *bool {
	b := true
	return &b
}

func Number(n float64) *float64 {
	return &n
}

func Integer(n int) *int {
	return &n
}

func (c *condition) validStr(field string, value string) (err error) {
	if c.required == nil && value == "" {
		return nil
	}

	if c.required != nil && value == "" {
		return &ValidatorErr{Field: field, Message: *c.required, Kind: ConditionRequired}
	}

	if reg := c.regexp; reg != nil {
		if !reg.MatchString(value) {
			return &ValidatorErr{Field: field, Message: *c.regMsgID, Kind: ConditionStringRegexp}
		}
	}

	ln := len([]rune(value))
	if max, min := c.maxLen, c.minLen; max != nil && min != nil {
		if ln < *min || *max < ln {
			return &ValidatorErr{
				Field:   field,
				Message: *c.lenMsgID,
				Kind:    ConditionStringLength,
				Max:     ptrIntToFloat(max), EQMax: Equal(),
				Min: ptrIntToFloat(min), EQMin: Equal(),
			}
		}
	}
	return nil
}

func ptrIntToFloat(n *int) *float64 {
	f := float64(*n)
	return &f
}

func ptrInt64ToFloat(n *int64) *float64 {
	f := float64(*n)
	return &f
}

func (c *condition) validNum(field string, valueStr string, value float64) (err error) {
	if c.required == nil && valueStr == "" {
		return nil
	}

	if c.required != nil && valueStr == "" {
		return &ValidatorErr{Field: field, Message: *c.required, Kind: ConditionRequired}
	}

	if c.minNum != nil {
		if c.eqMinNum != nil {
			if value < *c.minNum {
				return &ValidatorErr{Field: field, Message: *c.numMsgID, Kind: ConditionNumber, Min: c.minNum}
			}
		} else {
			if value <= *c.minNum {
				return &ValidatorErr{Field: field, Message: *c.numMsgID, Kind: ConditionNumber, Min: c.minNum, EQMin: Equal()}
			}
		}
	}
	if c.maxNum != nil {
		if c.eqMaxNum != nil {
			if value > *c.maxNum {
				return &ValidatorErr{Field: field, Message: *c.numMsgID, Kind: ConditionNumber, Max: c.maxNum}
			}
		} else {
			if value >= *c.maxNum {
				return &ValidatorErr{Field: field, Message: *c.numMsgID, Kind: ConditionNumber, Max: c.maxNum, EQMax: Equal()}
			}
		}
	}
	return nil
}

func (c *condition) validFile(field string, form *multipart.Form) (err error) {
	fileItem := len(form.File[field])
	if c.required == nil && fileItem == 0 {
		return nil
	}
	if c.required != nil && fileItem == 0 {
		return &ValidatorErr{Field: field, Message: *c.required, Kind: ConditionRequired}
	}
	// 验证文件个数
	if min, max := c.minItem, c.maxItem; min != nil && max != nil {
		if *min > fileItem || *max < fileItem {
			return &ValidatorErr{
				Field:   field,
				Message: *c.itemMsgID,
				Kind:    ConditionItem,
				Max:     ptrIntToFloat(max), EQMax: Equal(),
				Min: ptrIntToFloat(min), EQMin: Equal(),
			}
		}
	}
	// 验证文件大小
	if min, max := c.minFileByte, c.maxFileByte; min != nil && max != nil {
		for _, header := range form.File[field] {
			if header.Size > *max || header.Size < *min {
				return &ValidatorErr{
					Field:   field,
					Message: *c.fileSizeMsgID,
					Kind:    ConditionFileSize,
					Max:     ptrInt64ToFloat(max), EQMax: Equal(),
					Min: ptrInt64ToFloat(min), EQMin: Equal(),
				}
			}
		}
	}
	// 验证后缀名
	if len(c.fileExtensions) > 0 {
		for _, header := range form.File[field] {
			exist := false
			for _, ext := range c.fileExtensions {
				if filepath.Ext(header.Filename) == ext {
					exist = true
					break
				}
			}
			if !exist {
				return &ValidatorErr{Field: field, Message: *c.fileExtMsgID, Kind: ConditionFileExtensions, Enums: c.fileExtensions}
			}
		}
	}
	// 验证 Mime type
	if len(c.fileMimeTypes) > 0 {
		for _, header := range form.File[field] {
			exist := false
			for _, need := range c.fileMimeTypes {
				got := header.Header.Get("Content-Type")
				if got == need {
					exist = true
					break
				} else {
					if !strings.Contains(need, "/") { // e.g. need = "image" got = "image/png"
						if need == got[:strings.Index(got, "/")] {
							exist = true
							break
						}
					}
				}
			}
			if !exist {
				return &ValidatorErr{Field: field, Message: *c.fileMimeMsgID, Kind: ConditionFileMimeTypes, Enums: c.fileMimeTypes}
			}
		}
	}
	return nil
}

func (c *condition) validItem(field string, lenVs int) (err error) {
	if c.required == nil && lenVs == 0 {
		return nil
	}
	if c.required != nil && lenVs == 0 {
		return &ValidatorErr{Field: field, Message: *c.required, Kind: ConditionRequired}
	}
	if min, max := c.minItem, c.maxItem; min != nil && max != nil {
		if *min > lenVs || *max < lenVs {
			return &ValidatorErr{
				Field:   field,
				Message: *c.itemMsgID,
				Kind:    ConditionItem,
				Max:     ptrIntToFloat(max), EQMax: Equal(),
				Min: ptrIntToFloat(min), EQMin: Equal(),
			}
		}
	}
	return nil
}

func (c *condition) validTime(field string, form *multipart.Form) (err error) {
	item := len(form.Value[field])
	if c.required == nil && item == 0 {
		return nil
	}
	if c.required != nil && item == 0 {
		return &ValidatorErr{Field: field, Message: *c.required, Kind: ConditionRequired}
	}
	return nil
}

func (c *condition) validEnum(field string, values []string) (err error) {
	if c.required == nil && len(values) == 0 {
		return nil
	}
	if c.required != nil && len(values) == 0 {
		return &ValidatorErr{Field: field, Message: *c.required, Kind: ConditionRequired}
	}
	if len(c.enums) > 0 {
		for _, value := range values {
			got := false
			for _, e := range c.enums {
				if value == e {
					got = true
				}
			}
			if !got {
				return &ValidatorErr{Field: field, Message: *c.enumMsgID, Kind: ConditionEnums, Enums: c.enums}
			}
		}
	}
	return
}

// 20 <= x
var rightGt = regexp.MustCompile(`(-?\d+\.?\d*)(<[=]*)[a-zA-Z]+`)

// x > 20
var leftGt = regexp.MustCompile(`[a-zA-Z]+(>[=]*)(-?\d+\.?\d*)`)

// 20 >= x
var rightLt = regexp.MustCompile(`(-?\d+\.?\d*)(>[=]*)[a-zA-Z]+`)

// x > 20
var leftLt = regexp.MustCompile(`[a-zA-Z]+(<[=]*)(-?\d+\.?\d*)`)

func gt(str string) (n float64, eq, match bool, err error) {
	rms := rightGt.FindAllStringSubmatch(str, 1)
	if len(rms) == 1 {
		if len(rms[0]) == 3 {
			i64, e := strconv.ParseFloat(rms[0][1], 64)
			if e != nil {
				err = e
				return
			}
			return i64, rms[0][2] == "<=", true, nil
		}
	}

	lms := leftGt.FindAllStringSubmatch(str, 1)
	if len(lms) == 1 {
		if len(lms[0]) == 3 {
			i64, e := strconv.ParseFloat(lms[0][2], 64)
			if e != nil {
				err = e
				return
			}
			return i64, lms[0][1] == ">=", true, nil
		}
	}
	return
}

func lt(str string) (n float64, eq, match bool, err error) {
	rms := rightLt.FindAllStringSubmatch(str, 1)
	if len(rms) == 1 {
		if len(rms[0]) == 3 {
			i64, e := strconv.ParseFloat(rms[0][1], 64)
			if e != nil {
				err = e
				return
			}
			return i64, rms[0][2] == ">=", true, nil
		}
	}

	lms := leftLt.FindAllStringSubmatch(str, 1)
	if len(lms) == 1 {
		if len(lms[0]) == 3 {
			i64, e := strconv.ParseFloat(lms[0][2], 64)
			if e != nil {
				err = e
				return
			}
			return i64, lms[0][1] == "<=", true, nil
		}
	}
	return
}

// -1.02<=x<99.99
func readArea(str string) (min, max *float64, eqmin, eqmax bool, err error) {
	str = strings.Replace(str, " ", "", -1)
	if n, eq, ok, e := gt(str); e == nil {
		if ok {
			min = &n
			eqmin = eq
		}
	} else {
		err = e
		return
	}
	if n, eq, ok, e := lt(str); e == nil {
		if ok {
			max = &n
			eqmax = eq
		}
	} else {
		err = e
		return
	}
	return
}

func MsgName(op string) string {
	return op + "-msg"
}

var between = regexp.MustCompile("^([\\d]+)-([\\d]+)$")

// match: "6-12"
func getBetween(str string) (min, max *int, err error) {
	str = strings.Replace(str, " ", "", -1)
	data := between.FindAllStringSubmatch(str, -1)
	if len(data) == 1 && len(data[0]) == 3 {
		imin, _ := strconv.Atoi(data[0][1])
		imax, _ := strconv.Atoi(data[0][2])
		if imin > imax { // 12-6 => 6-12
			imin, imax = imax, imin
		}
		min = &imin
		max = &imax
		return
	} else {
		if imax, _ := strconv.Atoi(str); imax > 0 {
			imin := 1
			min = &imin
			max = &imax
			return
		} else {
			return nil, nil, errors.New(`pattern should be like "1-18" or just "18"`)
		}
	}
}
