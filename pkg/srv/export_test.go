package srv

import "encoding/json"

// Decoder - provide read access to private decoder field.
// NOTE: test only func.
func (s *JSONPlugin) Decoder() *json.Decoder {
	return s.decoder
}
