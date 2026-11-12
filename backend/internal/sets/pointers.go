package sets

import (
	"bytes"
	"encoding/json"
)

// Pointers are the named marks under the cells of a figure row, in the order of the content.
// The JSON form is an object from a name to a cell index. The web client draws the names in the order of the object.
type Pointers []Pointer

// Pointer is one named mark under the cell Cell.
type Pointer struct {
	Name string
	Cell int
}

// MarshalJSON writes the pointers as one object with the names in order.
func (p Pointers) MarshalJSON() ([]byte, error) {
	var out bytes.Buffer
	out.WriteByte('{')
	for i, pointer := range p {
		if i > 0 {
			out.WriteByte(',')
		}
		name, err := json.Marshal(pointer.Name)
		if err != nil {
			return nil, err
		}
		out.Write(name)
		out.WriteByte(':')
		cell, err := json.Marshal(pointer.Cell)
		if err != nil {
			return nil, err
		}
		out.Write(cell)
	}
	out.WriteByte('}')
	return out.Bytes(), nil
}

// UnmarshalJSON reads an object and keeps the order of its names.
func (p *Pointers) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	if _, err := decoder.Token(); err != nil {
		return err
	}
	*p = Pointers{}
	for decoder.More() {
		name, err := decoder.Token()
		if err != nil {
			return err
		}
		var cell int
		if err := decoder.Decode(&cell); err != nil {
			return err
		}
		*p = append(*p, Pointer{Name: name.(string), Cell: cell})
	}
	return nil
}
