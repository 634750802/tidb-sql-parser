package wasm

import (
	"reflect"
	"strconv"
)

type typeInfo struct {
	id    ObjectType
	sid   ObjectId
	tp    reflect.Type
	cache map[ObjectId]interface{}
}

type pool struct {
	tid         ObjectType
	typeMap     map[reflect.Type]ObjectType
	typeInfoMap map[ObjectType]*typeInfo
}

func newPool() *pool {
	return &pool{
		tid:         ObjectType(1),
		typeMap:     map[reflect.Type]ObjectType{},
		typeInfoMap: map[ObjectType]*typeInfo{},
	}
}

func (p *pool) defineType(otp ObjectType, tp reflect.Type) {
	if p.tid <= otp {
		p.tid = otp + 1
	}
	if _, ok := p.typeInfoMap[otp]; ok {
		panic("type " + strconv.Itoa(int(otp)) + " already defined")
	}
	tpInfo := &typeInfo{tp: tp, sid: 2, id: otp, cache: map[ObjectId]interface{}{}}
	p.typeInfoMap[otp] = tpInfo
	p.typeMap[tp] = otp
}

func (p *pool) getDefinedType(tp reflect.Type) (ObjectType, bool) {
	otp, ok := p.typeMap[tp]
	return otp, ok
}

func (p *pool) get(handle ObjectHandle) (interface{}, reflect.Type) {
	tpId := ObjectType(uint32(handle) & typeMask)
	objId := ObjectId(uint32(handle) & idMask)

	if tp, ok := p.typeInfoMap[tpId]; ok {
		if obj, ok := tp.cache[objId]; ok {
			return obj, tp.tp
		}
	}

	return nil, nil
}

func (p *pool) getType(otp ObjectType) reflect.Type {
	if tp, ok := p.typeInfoMap[otp]; ok {
		return tp.tp
	}

	return nil
}

func (p *pool) new(obj interface{}) ObjectHandle {
	tp := reflect.TypeOf(obj).Elem()
	var otp ObjectType
	var ok bool
	if otp, ok = p.typeMap[tp]; ok {
		if tpInfo, ok := p.typeInfoMap[otp]; ok {
			oid := tpInfo.sid
			tpInfo.sid += 1
			tpInfo.cache[oid] = obj
			return ObjectHandle(uint32(otp) | uint32(oid))
		}
	} else {
		otp = p.tid << 24
		p.tid += 1
	}
	oid := ObjectId(1)
	cache := map[ObjectId]interface{}{}
	cache[oid] = obj
	tpInfo := &typeInfo{tp: tp, sid: 2, id: otp, cache: cache}
	p.typeInfoMap[otp] = tpInfo
	p.typeMap[tp] = otp
	return ObjectHandle(uint32(otp) | uint32(oid))
}

func (p *pool) delete(handle ObjectHandle) ObjectId {
	tpId := ObjectType(uint32(handle) & typeMask)
	objId := ObjectId(uint32(handle) & idMask)

	if tp, ok := p.typeInfoMap[tpId]; ok {
		if _, ok := tp.cache[objId]; ok {
			tp.cache[objId] = nil
			if len(tp.cache) == 0 {
				tp.sid = 1
			}
			return objId
		}
	}

	return 0
}
