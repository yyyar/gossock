//
// @author Yaroslav Pogrebnyak <yyyaroslav@gmail.com>
//

package gossock

//
// Frame types
//
const (
	FRAME_TYPE_BINARY = 'b'
	FRAME_TYPE_JSON   = 'j'
)

//
// frame represents message
//
type frame struct {

	name string

	typ  byte

	body []byte
}
