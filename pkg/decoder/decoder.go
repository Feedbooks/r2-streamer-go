package decoder

import (
	"errors"
	"io"

	"github.com/readium/r2-streamer-go/pkg/pub"
)

// missingOrBadKey error return when the key is missing or not correct
const missingOrBadKey = "missing or bad key"

// List TODO add doc
type List struct {
	decoderAlgorithm string
	decoderScheme    string // only for lcp or other encrypted resource
	decoder          (func(*pub.Publication, pub.Link, io.ReadSeeker) (io.ReadSeeker, error))
}

var decoderList []List

// Decode decode the ressource
func Decode(publication *pub.Publication, link pub.Link, reader io.ReadSeeker) (io.ReadSeeker, error) {

	for _, decoderFunc := range decoderList {
		if link.Properties != nil && link.Properties.Encrypted != nil && link.Properties.Encrypted.Algorithm == decoderFunc.decoderAlgorithm && decoderFunc.decoderScheme == link.Properties.Encrypted.Scheme {
			return decoderFunc.decoder(publication, link, reader)
		}
	}

	return nil, errors.New("can't find fetcher")
}

// NeedToDecode check if there a decoder for this resource
func NeedToDecode(publication *pub.Publication, link pub.Link) bool {
	for _, decoderFunc := range decoderList {
		if link.Properties != nil && link.Properties.Encrypted != nil && link.Properties.Encrypted.Algorithm == decoderFunc.decoderAlgorithm && decoderFunc.decoderScheme == link.Properties.Encrypted.Scheme {
			return true
		}
	}

	return false
}
