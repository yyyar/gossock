//
// @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
//

package gossock

//
// Frame types
//
const (
	frameTypeBinary = 'b'
	frameTypeJson   = 'j'
)

//
// frame represents message
//
type frame struct {
	name string

	typ byte

	body []byte
}
