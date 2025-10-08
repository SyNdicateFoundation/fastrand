package fastrand

import "strings"

type Engine struct {
	defaultLength         int
	minLength             int
	maxLength             int
	inputEncoding         RandomizerEncoding
	outputEncoding        RandomizerEncoding
	rangesEnabled         bool
	keywordChoicesEnabled bool
	lengthChoicesEnabled  bool
	enabledKeywords       map[string]bool
	mailProviders         []string
	customCharsets        map[string][]byte
	customKeywords        map[string]CustomKeywordGenerator
}

type Option func(*Engine)

func NewEngine(opts ...Option) *Engine {
	enabledKeywords := make(map[string]bool, len(allKeywords))
	for _, kw := range allKeywords {
		enabledKeywords[kw] = true
	}

	e := &Engine{
		defaultLength:         16,
		minLength:             1,
		maxLength:             99,
		inputEncoding:         RandomizerEncodingURL | RandomizerEncodingHTML,
		outputEncoding:        RandomizerEncodingNone,
		rangesEnabled:         true,
		keywordChoicesEnabled: true,
		lengthChoicesEnabled:  true,
		enabledKeywords:       enabledKeywords,
		mailProviders:         SafeMailProviders,
		customCharsets:        make(map[string][]byte),
		customKeywords:        make(map[string]CustomKeywordGenerator),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

func (e *Engine) Reset() {
	freshEngine := NewEngine()
	*e = *freshEngine
}

func WithDefaultLength(length int) Option {
	return func(e *Engine) {
		if length > 0 {
			e.defaultLength = length
		}
	}
}

func WithMinLength(length int) Option {
	return func(e *Engine) {
		if length > 0 {
			e.minLength = length
		}
	}
}

func WithMaxLength(length int) Option {
	return func(e *Engine) {
		if length > 0 {
			e.maxLength = length
		}
	}
}

func WithDisabledKeywords(keywords ...string) Option {
	return func(e *Engine) {
		for _, kw := range keywords {
			e.enabledKeywords[strings.ToUpper(kw)] = false
		}
	}
}

func WithMailProviders(providers []string) Option {
	return func(e *Engine) {
		if len(providers) > 0 {
			e.mailProviders = providers
		}
	}
}

func WithCustomCharset(keyword string, charset []byte) Option {
	return func(e *Engine) {
		e.customCharsets[strings.ToUpper(keyword)] = charset
	}
}

func WithCustomKeyword(keyword string, generator CustomKeywordGenerator) Option {
	return func(e *Engine) {
		e.customKeywords[strings.ToUpper(keyword)] = generator
	}
}

func WithInputEncoding(encoding RandomizerEncoding) Option {
	return func(e *Engine) {
		e.inputEncoding = encoding
	}
}

func WithOutputEncoding(encoding RandomizerEncoding) Option {
	return func(e *Engine) {
		e.outputEncoding = encoding
	}
}

func WithRanges(enabled bool) Option {
	return func(e *Engine) {
		e.rangesEnabled = enabled
	}
}

func WithKeywordChoices(enabled bool) Option {
	return func(e *Engine) {
		e.keywordChoicesEnabled = enabled
	}
}

func WithLengthChoices(enabled bool) Option {
	return func(e *Engine) {
		e.lengthChoicesEnabled = enabled
	}
}
