package fastrand

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"strconv"
	"strings"
	"unsafe"
)

var (
	CharsNull         = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	SafeMailProviders []string
)

//go:embed mail_providers.txt
var mailProviders string

func init() {
	lines := strings.Split(mailProviders, "\n")
	SafeMailProviders = make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			SafeMailProviders = append(SafeMailProviders, trimmed)
		}
	}
}

var (
	startTag    = []byte("{RAND")
	startTagOpt = []byte("OM")
	endTag      = byte('}')
	sepTag      = byte(';')
	kwABL       = []byte("ABL")
	kwABU       = []byte("ABU")
	kwABR       = []byte("ABR")
	kwDIGIT     = []byte("DIGIT")
	kwHEX       = []byte("HEX")
	kwSPACE     = []byte("SPACE")
	kwUUID      = []byte("UUID")
	kwNULL      = []byte("NULL")
	kwIPV4      = []byte("IPV4")
	kwIPV6      = []byte("IPV6")
	kwBYTES     = []byte("BYTES")
	kwEMAIL     = []byte("EMAIL")
)

const defaultLength = 16
const maxLen = 99

func RandomizerString(payload string) string {
	return string(Randomizer([]byte(payload)))
}

func Randomizer(payload []byte) []byte {
	if !bytes.Contains(payload, startTag) {
		return payload
	}
	var buffer bytes.Buffer
	buffer.Grow(len(payload) + defaultLength*4)
	cursor := 0
	for {
		startIndex := bytes.Index(payload[cursor:], startTag)
		if startIndex == -1 {
			buffer.Write(payload[cursor:])
			break
		}
		startIndex += cursor
		buffer.Write(payload[cursor:startIndex])
		cursor = startIndex
		endIndex := bytes.IndexByte(payload[cursor:], endTag)
		if endIndex == -1 {
			buffer.Write(payload[cursor:])
			break
		}
		endIndex += cursor
		tag := payload[cursor:endIndex]
		cursor = endIndex + 1
		parseAndReplace(tag, &buffer)
	}
	return buffer.Bytes()
}

func parseAndReplace(tag []byte, buffer *bytes.Buffer) {
	tag = tag[len(startTag):]
	if bytes.HasPrefix(tag, startTagOpt) {
		tag = tag[len(startTagOpt):]
	}
	if len(tag) == 0 {
		buffer.WriteString(String(defaultLength, CharsAll))
		return
	}
	if tag[0] != sepTag {
		buffer.Write(startTag)
		if bytes.HasPrefix(tag, startTagOpt) {
			buffer.Write(startTagOpt)
		}
		buffer.Write(tag)
		return
	}
	tag = tag[1:]
	parts := bytes.SplitN(tag, []byte{sepTag}, 2)
	length := defaultLength
	var typeKeyword []byte
	if len(parts) == 1 {
		if l, err := parseLength(parts[0]); err == nil {
			length = l
		} else {
			typeKeyword = parts[0]
		}
	} else if len(parts) == 2 {
		if l, err := parseLength(parts[0]); err == nil {
			length = l
		}
		typeKeyword = parts[1]
	}
	switch {
	case bytes.Equal(typeKeyword, kwABL):
		buffer.WriteString(String(length, CharsAlphabetLower))
	case bytes.Equal(typeKeyword, kwABU):
		buffer.WriteString(String(length, CharsAlphabetUpper))
	case bytes.Equal(typeKeyword, kwABR):
		buffer.WriteString(String(length, CharsAlphabet))
	case bytes.Equal(typeKeyword, kwDIGIT):
		buffer.WriteString(String(length, CharsDigits))
	case bytes.Equal(typeKeyword, kwNULL):
		for i := 0; i < length; i++ {
			buffer.WriteByte(Choice(CharsNull))
		}
	case bytes.Equal(typeKeyword, kwSPACE):
		buffer.Write(bytes.Repeat([]byte(" "), length))
	case bytes.Equal(typeKeyword, kwUUID):
		buffer.Write(generateUUID())
	case bytes.Equal(typeKeyword, kwBYTES):
		buffer.Write(Bytes(length))
	case bytes.Equal(typeKeyword, kwIPV4):
		buffer.WriteString(IPv4().String())
	case bytes.Equal(typeKeyword, kwIPV6):
		buffer.WriteString(IPv6().String())
	case bytes.Equal(typeKeyword, kwEMAIL):
		buffer.Write(generateRandomEmail(length))
	case bytes.Equal(typeKeyword, kwHEX):
		buffer.Write(generateRandomHex(length))
	default:
		buffer.WriteString(String(length, CharsAll))
	}
}

func generateUUID() []byte {
	uuid := MustFastUUID()
	b := make([]byte, 36)
	hex.Encode(b[0:8], uuid[0:4])
	b[8] = '-'
	hex.Encode(b[9:13], uuid[4:6])
	b[13] = '-'
	hex.Encode(b[14:18], uuid[6:8])
	b[18] = '-'
	hex.Encode(b[19:23], uuid[8:10])
	b[23] = '-'
	hex.Encode(b[24:], uuid[10:])
	return b
}

func parseLength(b []byte) (int, error) {
	if len(b) == 0 || len(b) > 2 {
		return 0, strconv.ErrSyntax
	}
	val, err := strconv.Atoi(*(*string)(unsafe.Pointer(&b)))
	if err != nil {
		return 0, err
	}
	if val > 0 && val <= maxLen {
		return val, nil
	}
	return 0, strconv.ErrRange
}

func generateRandomEmail(userLength int) []byte {
	if userLength <= 0 {
		userLength = 8
	}
	user := String(userLength, CharsAlphabetLower)
	provider := "gmail.com"
	if len(SafeMailProviders) > 0 {
		provider = Choice(SafeMailProviders)
	}
	return []byte(user + "@" + provider)
}

func generateRandomHex(byteLength int) []byte {
	if byteLength <= 0 {
		byteLength = defaultLength
	}
	srcBytes := Bytes(byteLength)
	hexBytes := make([]byte, byteLength*2)
	hex.Encode(hexBytes, srcBytes)
	return hexBytes
}
