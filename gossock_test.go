//
// @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
//

package gossock

import (
	"io"
	"log"
	"testing"
	"sync"
)

//
// MockConnection for emulation through pipes
//
type MockConnection struct {
	io.ReadCloser
	io.Writer
}

//
// BinaryMessage (will be sent in binary form)
//
type BinaryMessage []byte

//
// HelloMessage (will be serialized to json)
//
type HelloMessage struct {
	Content string `json:"content"`
}


//
// TestBasic tests that receiver can get what 
// sender sent
//
func TestBasic(t *testing.T) {

	var wg sync.WaitGroup

	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()

	registry := NewRegistry()
	registry.Register("hello", HelloMessage{})
	registry.Register("binary", BinaryMessage{})

	//
	// Launch receiver party
	//
	go func() {
		s := New(MockConnection{r2, w1}, registry)
		var err error

		err = s.On(func(hello *HelloMessage) {
			t.Log("On Hello:", hello)
			wg.Done()
		})

		if err != nil {
			log.Println("Error On(Hello)", err)
		}

		err = s.On(func(binary *BinaryMessage) {
			t.Log("On Binary:", string(*binary))
			wg.Done()
		})

		if err != nil {
			log.Println("Error On(Binary)", err)
		}

	}()

	//
	// Launch sender party
	//
	func() {
		c := New(MockConnection{r1, w2}, registry)
		var err error

		wg.Add(1)
		err = c.Send(HelloMessage{
			"Hello, World!",
		})

		if err != nil {
			log.Println("Error c.Send(Hello):", err)
		}

		wg.Add(1)
		err = c.Send(BinaryMessage{
			0x62, 0x69, 0x6e, 0x61, 0x72, 0x79, 0xFF,
		})

		if err != nil {
			log.Println("Error c.Send(Binary)", err)
		}
	}()

	wg.Wait()

}
