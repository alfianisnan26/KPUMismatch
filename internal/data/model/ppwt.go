package model

import (
	"fmt"
	"strings"
)

var canonicalIndex = []int{2, 4, 6, 10, 13}

type PPWTEntity struct {
	ID      int64
	Kode    string
	Nama    string
	Tingkat int

	Parent *PPWTEntity
}

func (e PPWTEntity) GetCanonicalName() []string {
	if e.Tingkat <= 0 {
		return nil
	}
	canonicalName := make([]string, 5)
	canonicalName[e.Tingkat-1] = e.Nama

	node := e.Parent
	for node != nil {
		canonicalName[node.Tingkat-1] = node.Nama
		node = node.Parent
	}

	return canonicalName
}

func (e PPWTEntity) String() string {
	if e.Tingkat == 0 {
		return fmt.Sprintf("NaN")
	}

	if e.Nama == "" {
		return e.Kode
	}

	return fmt.Sprintf("%v | %v", e.Kode, e.Nama)
}

func (e PPWTEntity) LowestLevel() bool {
	return e.Tingkat == 5
}

func NewPPWT(kode string) PPWTEntity {
	tingkat := 0
	for i, index := range canonicalIndex {
		if len(kode) == index {
			tingkat = i + 1
			break
		}
	}

	return PPWTEntity{
		Kode:    kode,
		Tingkat: tingkat,
	}
}

func (e PPWTEntity) GetCanonicalCode() string {
	if e.Tingkat == 0 {
		return "0"
	}

	return strings.Join(e.GetCanonicalCodeAll(), "/")
}

func (e PPWTEntity) GetCanonicalCodeAll() []string {
	if e.Tingkat == 0 {
		return nil
	}

	canonical := make([]string, e.Tingkat)
	for i := 0; i < e.Tingkat; i++ {
		canonical[i] = e.Kode[:canonicalIndex[i]]
	}

	return canonical
}
