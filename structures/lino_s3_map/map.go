package lino_s3_map

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fxamacker/cbor/v2"
	"github.com/linolabx/lino_s3"
)

// file structure:

// | filehead magic_number_23238 | version        | offset         | ->
//   uint16 (2bytes)               uint8 (1 byte)   uint64 (8byte)

// -> index_size     | index | items... |
//    uint64 (8byte)   json         bytes

// filehead: magic_number 23238 to identify the file format
// version: file format version
// offset: size of everything before items, used to random access position of items
// index_size: size of index
// index: cbor([[key(string), start(unit64), end(unit64)]...])
// items: | bytes | bytes | ... |
// * numbers are in little endian

var filehead = []byte{
	// magic_number: 5AC6
	0x5A, 0xC6,

	// version: 01
	0x01,
}

const uint64_max = uint64(0xFFFF_FFFF_FFFF_FFFF)
const filehead_start = uint64(0)
const filehead_end = uint64(3)
const offset_length = 8
const offset_start = filehead_end
const offset_end = offset_start + offset_length
const index_size_length = 8
const index_size_start = offset_end
const index_size_end = index_size_start + index_size_length

const file_probe_size = uint64(1024 * 64) // 64kb

type IndexItem struct {
	key   string
	start uint64
	end   uint64
}

type Map struct {
	obj *lino_s3.LinoS3Object

	rw     sync.RWMutex
	frozen bool // if locked, no more Set allowed

	offset     uint64
	index_size uint64

	index       map[string]IndexItem
	index_bytes []byte

	items_bytes []byte
}

func NewMap(obj *lino_s3.LinoS3Object) *Map {
	if obj.HasInterceptors() {
		panic("interceptors not supported for LinoS3List")
	}

	return &Map{
		obj: obj,

		index:       map[string]IndexItem{},
		items_bytes: []byte{},
	}
}

func loadMap(obj *lino_s3.LinoS3Object, loadAll bool) (*Map, error) {
	if obj.HasInterceptors() {
		panic("interceptors not supported for LinoS3List")
	}

	var filepart = []byte{}
	var err error = nil
	if loadAll {
		filepart, err = obj.ReadBuffer()
	} else {
		filepart, err = obj.Piece(0, int64(file_probe_size)).ReadBuffer()
	}

	if err != nil {
		return nil, err
	}

	if !bytes.Equal(filepart[filehead_start:filehead_end], filehead) {
		return nil, errors.New("invalid filetype of version")
	}

	offset := binary.LittleEndian.Uint64(filepart[offset_start:offset_end])
	index_size := binary.LittleEndian.Uint64(filepart[index_size_start:index_size_end])

	if offset > file_probe_size {
		nextpart, err := obj.Piece(int64(file_probe_size), int64(offset)).ReadBuffer()
		if err != nil {
			return nil, err
		}

		filepart = append(filepart, nextpart...)
	}

	indexBody := make([][]interface{}, 0)
	if err := cbor.Unmarshal(filepart[index_size_end:offset], &indexBody); err != nil {
		return nil, err
	}

	index := map[string]IndexItem{}
	for _, v := range indexBody {
		index[v[0].(string)] = IndexItem{v[0].(string), v[1].(uint64), v[2].(uint64)}
	}

	var items_bytes []byte = nil
	if loadAll {
		items_bytes = filepart[offset:]
	}

	return &Map{
		obj: obj,

		offset:     offset,
		index_size: index_size,
		frozen:     true,

		index:       index,
		items_bytes: items_bytes,
	}, nil
}

func LoadMap(obj *lino_s3.LinoS3Object) (*Map, error) {
	return loadMap(obj, false)
}

func LoadEntireMap(obj *lino_s3.LinoS3Object) (*Map, error) {
	return loadMap(obj, true)
}

func (l *Map) Set(key string, data []byte) error {
	l.rw.Lock()
	defer l.rw.Unlock()

	if l.frozen {
		return errors.New("frozen list, no more Set allowed")
	}

	if _, ok := l.index[key]; ok {
		return nil
	}

	start := uint64(len(l.items_bytes))
	end := start + uint64(len(data))

	l.index[key] = IndexItem{key, start, end}
	l.items_bytes = append(l.items_bytes, data...)

	return nil
}

func (l *Map) Get(key string) ([]byte, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()

	item, ok := l.index[key]
	if !ok {
		return nil, errors.New("key not found")
	}

	if l.items_bytes != nil && len(l.items_bytes) != 0 {
		return l.items_bytes[item.start:item.end], nil
	}

	return l.obj.Piece(int64(item.start+l.offset), int64(item.end+l.offset)).ReadBuffer()
}

func (l *Map) ValidateAndFreeze() error {
	l.rw.Lock()
	defer l.rw.Unlock()

	if l.frozen {
		return nil
	}

	index_body := make([][]interface{}, 0, len(l.index))
	for _, v := range l.index {
		index_body = append(index_body, []interface{}{v.key, v.start, v.end})
	}

	index_bytes, err := cbor.Marshal(index_body)
	if err != nil {
		return err
	}

	index_size := uint64(len(index_bytes))
	if index_size > uint64_max {
		return errors.New("index size too big, max 4 bytes")
	}

	offset := index_size_end + index_size
	if offset > uint64_max {
		return errors.New("offset too big, max 4 bytes")
	}

	l.index_bytes = index_bytes
	l.index_size = index_size
	l.offset = offset
	l.frozen = true

	return nil
}

func (l *Map) Save() error {
	if err := l.ValidateAndFreeze(); err != nil {
		return err
	}

	l.rw.RLock()
	defer l.rw.RUnlock()

	index_size := make([]byte, index_size_length)
	binary.LittleEndian.PutUint64(index_size, l.index_size)

	offset := make([]byte, offset_length)
	binary.LittleEndian.PutUint64(offset, l.offset)

	file := []byte{}
	file = append(file, filehead...)
	file = append(file, offset...)
	file = append(file, index_size...)
	file = append(file, l.index_bytes...)
	file = append(file, l.items_bytes...)

	return l.obj.WriteBuffer(file, "lino/s3-map-v1")
}

func (l *Map) PieceKey(key string) (string, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()

	if !l.frozen {
		return "", errors.New("not frozen")
	}

	item, ok := l.index[key]
	if !ok {
		return "", errors.New("key not found")
	}

	return l.obj.Piece(int64(item.start+l.offset), int64(item.end+l.offset)).Key(), nil
}

func (l *Map) MapAllRangeKey() (map[string]string, error) {
	l.rw.RLock()
	defer l.rw.RUnlock()

	if !l.frozen {
		return nil, errors.New("not frozen")
	}

	result := map[string]string{}
	for k, v := range l.index {
		result[k] = l.obj.Piece(int64(v.start+l.offset), int64(v.end+l.offset)).Key()
	}

	return result, nil
}

func (l *Map) Delete() (*s3.DeleteObjectOutput, error) {
	return l.obj.Delete()
}
