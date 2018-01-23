//
// @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
//

package gossock

import (
	"encoding/binary"
	"io"
	"sync"
)

//
// serializer does serialization of frame to writer
//
type serializer struct {

	// writer for writing resulting bytes
	writer io.Writer
	mutex  sync.Mutex
}

//
// newSerializer creates new instance of serializer
//
func newSerializer(writer io.Writer) *serializer {

	s := serializer{
		writer: writer,
		mutex:  sync.Mutex{},
	}

	return &s
}

//
// serialize writes frame to writer
//
func (s *serializer) serialize(frame *frame) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

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
