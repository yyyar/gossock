//
// @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
//

package gossock

import (
	"io"
	"encoding/binary"
)

//
// serializer does serialization of frame to writer
//
type serializer struct {

	// writer for writing resulting bytes
	writer io.Writer
}

//
// newSerializer creates new instance of serializer
//
func newSerializer(writer io.Writer) *serializer {

	serializer := serializer{
		writer: writer,
	}

	return &serializer
}

//
// serialize writes frame to writer
//
func (s *serializer) serialize(frame* frame) error {

	if err := binary.Write(s.writer, binary.BigEndian, frame.typ); err != nil {
		return err
	}

	if err := binary.Write(s.writer, binary.BigEndian, byte(len(frame.name))); err != nil {
		return err
	}

	if err := binary.Write(s.writer, binary.BigEndian, []byte(frame.name)); err != nil {
		return err
	}

	if err := binary.Write(s.writer, binary.BigEndian, uint32(len(frame.body))); err != nil {
		return err
	}

	if err := binary.Write(s.writer, binary.BigEndian, frame.body); err != nil {
		return err
	}

	return nil
}
