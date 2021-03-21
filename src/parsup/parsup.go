package parsup

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ParamsSupport 参数辅助类
type ParamsSupport struct {
	IsDeep       *bool // 深度递归
	IsConvOID    *bool // 转化ObjectID
	IsConvTime   *bool // 转化时间对象
	IsDenyInject *bool // 防注入
	IsConvStruct *bool // 转结构
}

// ParSup 工厂方法
func ParSup() *ParamsSupport {
	t := true
	f := false
	return &ParamsSupport{
		IsDeep:       &t,
		IsConvOID:    &t,
		IsConvTime:   &t,
		IsDenyInject: &t,
		IsConvStruct: &f,
	}
}

// SetIsDeep 设置方法
func (p *ParamsSupport) SetIsDeep(b bool) *ParamsSupport {
	p.IsDeep = &b
	return p
}

// SetIsConvOID 设置方法
func (p *ParamsSupport) SetIsConvOID(b bool) *ParamsSupport {
	p.IsConvOID = &b
	return p
}

// SetIsConvTime 设置方法
func (p *ParamsSupport) SetIsConvTime(b bool) *ParamsSupport {
	p.IsConvTime = &b
	return p
}

// SetIsDenyInject 设置方法
func (p *ParamsSupport) SetIsDenyInject(b bool) *ParamsSupport {
	p.IsDenyInject = &b
	return p
}

// SetIsConvStruct 设置方法
func (p *ParamsSupport) SetIsConvStruct(b bool) *ParamsSupport {
	p.IsConvStruct = &b
	return p
}

// ConvBase base handler
func (p *ParamsSupport) ConvBase(i interface{}) (interface{}, error) {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Invalid:
		return nil, nil
	case reflect.Map:
		if *p.IsDeep {
			return p.ConvMap(i.(map[string]interface{}))
		}
		break
	case reflect.Slice:
		if *p.IsDeep {
			return p.ConvSlice(i.([]interface{}))
		}
		break
	case reflect.Struct:
		if *p.IsDeep && *p.IsConvStruct {
			return p.ConvStruct(i)
		}
		break
	case reflect.String:
		return p.ConvStr(i.(string))
	}
	return i, nil

}

// ConvStr string handler
func (p *ParamsSupport) ConvStr(s string) (interface{}, error) {
	if *p.IsDenyInject {
		if strings.ContainsAny(s, "$[]{}()") {
			return nil, errors.New("不能包含$[]{}()等特殊符号")
		}
	}

	if *p.IsConvOID {
		if oid, err := primitive.ObjectIDFromHex(s); err == nil {
			return oid, nil
		}
	}
	if *p.IsConvTime {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t.Local(), nil
		}
	}
	return s, nil
}

// ConvMap map handler
func (p *ParamsSupport) ConvMap(m map[string]interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	for k, v := range m {
		dv, err := p.ConvBase(v)
		if err != nil {
			return nil, err
		}
		res[k] = dv
	}
	return res, nil
}

// ConvSlice slice handler
func (p *ParamsSupport) ConvSlice(s []interface{}) ([]interface{}, error) {
	res := make([]interface{}, len(s))
	for k, v := range s {
		dv, err := p.ConvBase(v)
		if err != nil {
			return nil, err
		}
		res[k] = dv
	}
	return res, nil
}

// ConvJSON byte handler
func (p *ParamsSupport) ConvJSON(s []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(s, &m)
	if err != nil {
		return nil, err
	}
	return p.ConvMap(m)
}

// ConvStruct struct handler
func (p *ParamsSupport) ConvStruct(s interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Struct {
		return nil, errors.New("not a struct")
	}
	t := val.Type()
	o := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		var omitempty, skip, omitzero bool
		tagsName := t.Field(i).Name
		cur := val.Field(i)
		if tags, ok := t.Field(i).Tag.Lookup("parsup"); ok {
			tagsArr := strings.Split(tags, ",")
			tagsName = tagsArr[0]
			for _, v := range tagsArr {
				switch v {
				case "omitempty":
					omitempty = true
				case "omitzero":
					omitzero = true
				case "-":
					skip = true
				}
			}
		}

		if skip || (cur.Kind() == reflect.Ptr && cur.IsNil() && omitempty) || (cur.IsZero() && omitzero) {
			continue
		}

		ele, err := p.ConvBase(p.safeInterface(cur))

		if err != nil {
			return nil, err
		}

		o[tagsName] = ele
	}

	return o, nil
}

func (p *ParamsSupport) safeInterface(v reflect.Value) interface{} {
	if v.IsValid() && v.CanInterface() {
		return v.Interface()
	}
	return nil
}
