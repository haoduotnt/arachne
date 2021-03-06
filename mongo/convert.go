package mongo

import (
	//"fmt"
	"github.com/bmeg/arachne/aql"
	"github.com/bmeg/arachne/protoutil"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"gopkg.in/mgo.v2/bson"
)

var FIELD_SRC string = "from"
var FIELD_DST string = "to"
var FIELD_BUNDLE string = "bundle"

func PackVertex(v aql.Vertex) map[string]interface{} {
	p := map[string]interface{}{}
	if v.Properties != nil {
		p = protoutil.AsMap(v.Properties)
	}
	//fmt.Printf("proto:%s\nmap:%s\n", v.Properties, p)
	return map[string]interface{}{
		"_id":        v.Gid,
		"label":      v.Label,
		"properties": p,
	}
}

func PackEdge(e aql.Edge) map[string]interface{} {
	p := map[string]interface{}{}
	if e.Properties != nil {
		p = protoutil.AsMap(e.Properties)
	}
	o := map[string]interface{}{
		FIELD_SRC:    e.From,
		FIELD_DST:    e.To,
		"label":      e.Label,
		"properties": p,
	}
	if e.Gid != "" {
		o["_id"] = e.Gid
	}
	return o
}

func PackBundle(e aql.Bundle) map[string]interface{} {
	m := map[string]interface{}{}
	for k, v := range e.Bundle {
		m[k] = protoutil.AsMap(v)
	}
	o := map[string]interface{}{
		FIELD_SRC:    e.From,
		FIELD_BUNDLE: m,
		"label":      e.Label,
	}
	if e.Gid != "" {
		o["_id"] = e.Gid
	}
	return o
}

func UnpackVertex(i map[string]interface{}) aql.Vertex {
	o := aql.Vertex{}
	o.Gid = i["_id"].(string)
	o.Label = i["label"].(string)
	o.Properties = protoutil.AsStruct(i["properties"].(map[string]interface{}))
	return o
}

func UnpackEdge(i map[string]interface{}) aql.Edge {
	o := aql.Edge{}
	id := i["_id"]
	if idb, ok := id.(bson.ObjectId); ok {
		o.Gid = idb.String()
	} else {
		o.Gid = id.(string)
	}
	o.Label = i["label"].(string)
	o.From = i[FIELD_SRC].(string)
	o.To = i[FIELD_DST].(string)
	o.Properties = protoutil.AsStruct(i["properties"].(map[string]interface{}))
	return o
}

func UnpackBundle(i map[string]interface{}) aql.Bundle {
	o := aql.Bundle{}
	id := i["_id"]
	if idb, ok := id.(bson.ObjectId); ok {
		o.Gid = idb.Hex()
	} else {
		o.Gid = id.(string)
	}
	o.Label = i["label"].(string)
	o.From = i[FIELD_SRC].(string)
	m := map[string]*structpb.Struct{}
	for k, v := range i[FIELD_BUNDLE].(map[string]interface{}) {
		m[k] = protoutil.AsStruct(v.(map[string]interface{}))
	}
	o.Bundle = m
	return o
}
