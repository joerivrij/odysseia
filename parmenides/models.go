// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    logos, err := UnmarshalLogos(bytes)
//    bytes, err = logos.Marshal()

package main

import "encoding/json"

type Logos []Word

func UnmarshalLogos(data []byte) (Logos, error) {
	var r Logos
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Logos) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Word struct {
	Greek   string `json:"greek"`
	Dutch   string `json:"dutch"`
	Chapter int64  `json:"chapter"`
}

func (r *Word) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
