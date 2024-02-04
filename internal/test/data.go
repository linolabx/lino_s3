package test

import (
	_ "embed"
)

type User struct {
	Name     string `json:"name" cbor:"1,keyasint,omitempty"`
	Email    string `json:"mail" cbor:"2,keyasint,omitempty"`
	Password string `json:"password" cbor:"3,keyasint,omitempty"`
	Entropy  uint64 `json:",omitempty,string" cbor:"4,keyasint,omitempty"`
}

var Violet = &User{Name: "Violet Evergarden", Email: "violet@ch-postal.com", Password: "password", Entropy: 18_446_744_073_709_551_615}
var Gilbert = &User{Name: "Gilbert Bougainvillea", Email: "", Password: "", Entropy: 0}

type Book struct {
	Title   string `json:"title" cbor:"1,keyasint,omitempty"`
	Author  string `json:"author" cbor:"2,keyasint,omitempty"`
	Content string `json:"content" cbor:"3,keyasint,omitempty"`
	Size    int64  `json:"size" cbor:"4,keyasint,omitempty"`
}

//go:embed gzip_test.txt
var bookContent []byte

func GetBook() (*Book, error) {

	return &Book{
		Title:   "帝国主义是资本主义的最高阶段",
		Author:  "列宁",
		Content: string(bookContent),
		Size:    int64(len(bookContent)),
	}, nil
}
