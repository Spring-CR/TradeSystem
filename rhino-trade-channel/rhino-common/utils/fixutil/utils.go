package fixutil

import (
	"bytes"
	"rhino-common/utils/byteutils"
)
var(
	fixSep =  []byte{0x01}
	fixNewSep =  []byte{'|'}
)
func ConvertFIXDataToString(fixMsg[]byte) string {
	if len(fixMsg) == 0 {
		return ""
	}
	replaced := bytes.ReplaceAll(fixMsg, fixSep, fixNewSep)
	return byteutils.GetZeroCopyString(replaced)
}