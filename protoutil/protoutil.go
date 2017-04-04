package protoutil

import (
	"github.com/gogo/protobuf/types"
	"log"
)

func StructSet(s *types.Struct, key string, value interface{}) {
	vw := WrapValue(value)
	s.Fields[key] = vw
}

func WrapValue(value interface{}) *types.Value {
	switch v := value.(type) {
	case string:
		return &types.Value{Kind: &types.Value_StringValue{v}}
	case int:
		return &types.Value{Kind: &types.Value_NumberValue{float64(v)}}
	case int64:
		return &types.Value{Kind: &types.Value_NumberValue{float64(v)}}
	case float64:
		return &types.Value{Kind: &types.Value_NumberValue{float64(v)}}
	case bool:
		return &types.Value{Kind: &types.Value_BoolValue{v}}
	case *types.Value:
		return v
	case []interface{}:
		o := make([]*types.Value, len(v))
		for i, k := range v {
			wv := WrapValue(k)
			o[i] = wv
		}
		return &types.Value{Kind: &types.Value_ListValue{&types.ListValue{Values: o}}}
	case []string:
		o := make([]*types.Value, len(v))
		for i, k := range v {
			wv := &types.Value{Kind: &types.Value_StringValue{k}}
			o[i] = wv
		}
		return &types.Value{Kind: &types.Value_ListValue{&types.ListValue{Values: o}}}
	case map[string]interface{}:
		o := &types.Struct{Fields: map[string]*types.Value{}}
		for k, v := range v {
			wv := WrapValue(v)
			o.Fields[k] = wv
		}
		return &types.Value{Kind: &types.Value_StructValue{o}}
	default:
		log.Printf("unknown data type: %T", value)
	}
	return nil
}

func CopyToStructSub(s *types.Struct, keys []string, values map[string]interface{}) {
	for _, i := range keys {
		StructSet(s, i, values[i])
	}
}

func CopyToStruct(s *types.Struct, values map[string]interface{}) {
	for i := range values {
		StructSet(s, i, values[i])
	}
}

func CopyStructToStruct(dst *types.Struct, src *types.Struct) {
	for k, v := range src.Fields {
		StructSet(dst, k, v)
	}
}

func CopyStructToStructSub(dst *types.Struct, keys []string, src *types.Struct) {
	for _, k := range keys {
		StructSet(dst, k, src.Fields[k])
	}
}

func AsMap(src *types.Struct) map[string]interface{} {
	out := map[string]interface{}{}
	for k, f := range src.Fields {
		if v, ok := f.Kind.(*types.Value_StringValue); ok {
			out[k] = v.StringValue
		} else if v, ok := f.Kind.(*types.Value_NumberValue); ok {
			out[k] = v.NumberValue
		} else if v, ok := f.Kind.(*types.Value_StructValue); ok {
			out[k] = AsMap(v.StructValue)
		}
	}
	return out
}

func AsStruct(src map[string]interface{}) *types.Struct {
	out := types.Struct{Fields: map[string]*types.Value{}}
	for k, v := range src {
		StructSet(&out, k, v)
	}
	return &out
}
