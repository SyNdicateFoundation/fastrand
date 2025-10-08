package fastrand

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"math/rand"
	"strings"
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
	startTag         = []byte("{RAND")
	startUrlEncoded  = []byte("%7BRAND")
	startHtmlEncoded = []byte("&lbrace;RAND")
	startTagOpt      = []byte("OM")
	endTag           = byte('}')
	endTagUrl        = []byte("%7D")
	endTagHtml       = []byte("&rbrace;")
	sepTag           = byte(';')
	sepTagUrl        = []byte("%3B")
	sepTagHtml       = []byte("&semi;")
	kwABL            = []byte("ABL")
	kwABU            = []byte("ABU")
	kwABR            = []byte("ABR")
	kwDIGIT          = []byte("DIGIT")
	kwHEX            = []byte("HEX")
	kwSPACE          = []byte("SPACE")
	kwUUID           = []byte("UUID")
	kwNULL           = []byte("NULL")
	kwIPV4           = []byte("IPV4")
	kwIPV6           = []byte("IPV6")
	kwBYTES          = []byte("BYTES")
	kwEMAIL          = []byte("EMAIL")
)

const defaultLength = 16

func hasPrefix(slice, prefix []byte, pos int) bool {
	if pos+len(prefix) > len(slice) {
		return false
	}
	return bytes.Equal(slice[pos:pos+len(prefix)], prefix)
}

func RandomizerString(payload string) string {
	return string(Randomizer([]byte(payload)))
}

func Randomizer(payload []byte) []byte {
	if !bytes.ContainsAny(payload, "{%&") {
		return payload
	}

	if bytes.ContainsAny(payload, "%&") {
		var normalizedBuf bytes.Buffer
		normalizedBuf.Grow(len(payload))
		cursor := 0
		for cursor < len(payload) {
			idx := bytes.IndexAny(payload[cursor:], "%&")
			if idx == -1 {
				normalizedBuf.Write(payload[cursor:])
				break
			}
			normalizedBuf.Write(payload[cursor : cursor+idx])
			cursor += idx

			if hasPrefix(payload, startUrlEncoded, cursor) {
				normalizedBuf.Write(startTag)
				cursor += len(startUrlEncoded)
			} else if hasPrefix(payload, startHtmlEncoded, cursor) {
				normalizedBuf.Write(startTag)
				cursor += len(startHtmlEncoded)
			} else if hasPrefix(payload, endTagUrl, cursor) {
				normalizedBuf.WriteByte(endTag)
				cursor += len(endTagUrl)
			} else if hasPrefix(payload, endTagHtml, cursor) {
				normalizedBuf.WriteByte(endTag)
				cursor += len(endTagHtml)
			} else if hasPrefix(payload, sepTagUrl, cursor) {
				normalizedBuf.WriteByte(sepTag)
				cursor += len(sepTagUrl)
			} else if hasPrefix(payload, sepTagHtml, cursor) {
				normalizedBuf.WriteByte(sepTag)
				cursor += len(sepTagHtml)
			} else {
				normalizedBuf.WriteByte(payload[cursor])
				cursor++
			}
		}
		payload = normalizedBuf.Bytes()
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

		parseAndReplaceFast(tag, &buffer)
	}

	return buffer.Bytes()
}

func parseAndReplaceFast(tag []byte, buffer *bytes.Buffer) {
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

	length := defaultLength
	var typeKeyword, lenPart []byte

	sepIndex := bytes.IndexByte(tag, sepTag)

	if sepIndex == -1 {
		lenPart = tag
	} else {
		lenPart = tag[:sepIndex]
		typeKeyword = tag[sepIndex+1:]
	}

	rangeSepIndex := bytes.IndexByte(lenPart, '-')
	if rangeSepIndex != -1 {
		minPart := lenPart[:rangeSepIndex]
		maxPart := lenPart[rangeSepIndex+1:]

		if minX, ok1 := parseLengthFast(minPart); ok1 {
			if maxX, ok2 := parseLengthFast(maxPart); ok2 && minX <= maxX {
				length = rand.Intn(maxX-minX+1) + minX
			}
		}
	} else {
		if l, ok := parseLengthFast(lenPart); ok && l > 0 {
			length = l
		} else if typeKeyword == nil {
			typeKeyword = lenPart
		}
	}

	if bytes.Contains(typeKeyword, []byte(",")) {
		choices := bytes.Split(typeKeyword, []byte(","))
		typeKeyword = choices[rand.Intn(len(choices))]
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

func parseLengthFast(b []byte) (int, bool) {
	if len(b) == 1 {
		c := b[0]
		if c >= '0' && c <= '9' {
			return int(c - '0'), true
		}
	}
	if len(b) == 2 {
		c1, c2 := b[0], b[1]
		if c1 >= '0' && c1 <= '9' && c2 >= '0' && c2 <= '9' {
			return int(c1-'0')*10 + int(c2-'0'), true
		}
	}

	return 0, false
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
