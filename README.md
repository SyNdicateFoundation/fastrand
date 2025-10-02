# FastRand for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/SyNdicateFoundation/fastrand)](https://goreportcard.com/report/github.com/SyNdicateFoundation/fastrand)
[![GoDoc](https://godoc.org/github.com/SyNdicateFoundation/fastrand?status.svg)](https://godoc.org/github.com/SyNdicateFoundation/fastrand)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**FastRand** is a powerful, high-performance, and developer-friendly Go library for generating random data. It provides a comprehensive suite of tools, from simple random numbers to complex, templated data generation. The library offers both a high-speed non-cryptographically secure generator (using PCG) and a cryptographically secure generator (using ChaCha8), making it suitable for a wide range of applications.

## Features

-   **Dual Random Sources**: Choose between the blazing-fast PCG algorithm for general-purpose randomness and the secure ChaCha8 for cryptographic needs.
-   **Simple & Idiomatic API**: Functions are designed to be intuitive and easy to remember (e.g., `IntN`, `String`, `Bytes`).
-   **Type-Safe Generics**: Generate random numbers for any standard integer or float type using `Number[T]()` without runtime reflection.
-   **Rich Helper Functions**: One-line functions for common data types like `IPv4`, `IPv6`, `UUID`, `Hex`, and more.
-   **Slice Manipulation**: Easily `Choice` an element, select `ChoiceMultiple` elements, or `Shuffle` a slice.
-   **Powerful Randomizer Engine**: Generate complex data from a template string using placeholders like `{RAND;10;HEX}`.
-   **Concurrency Safe**: All functions are safe for concurrent use in goroutines.
-   **Zero Dependencies**: Relies only on the Go standard library.

## Installation

```sh
go get github.com/SyNdicateFoundation/fastrand
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/SyNdicateFoundation/fastrand"
)

func main() {
	// Generate a random integer between 10 and 100
	fmt.Println("Random Int:", fastrand.Int(10, 100))

	// Generate a 12-character alphanumeric string
	fmt.Println("Random String:", fastrand.String(12, fastrand.CharsAlphabetDigits))

	// Generate 16 cryptographically secure random bytes as a hex string
	secureHex, err := fastrand.SecureHex(16)
	if err != nil {
		panic(err)
	}
	fmt.Println("Secure Hex:", secureHex)

	// Use the Randomizer to generate data from a template
	templatedData := fastrand.RandomizerString("Session ID: {RAND;16;HEX}, User: {RAND;8;ABL}")
	fmt.Println("Templated Data:", templatedData)
}
```

## Full API Examples

### Numeric Types

#### Integers and Floats

```go
// Integer between -50 and 50
n1 := fastrand.Int(-50, 50)

// Integer between 0 and 99
n2 := fastrand.IntN(100)

// Standard float64 between 0.0 and 1.0
f1 := fastrand.Float64()

// A single random byte
b := fastrand.Byte()

fmt.Printf("Int: %d, IntN: %d, Float64: %f, Byte: %d\n", n1, n2, f1, b)
```

#### Generic Numbers

The `Number` and `NumberN` functions work with any standard integer or float type.

```go
// Generate a random int32 between -1000 and 1000
num_i32 := fastrand.Number[int32](-1000, 1000)

// Generate a random uint64 up to 50000
num_u64 := fastrand.NumberN[uint64](50000)

// Generate a random float32 between 0.0 and 10.5
num_f32 := fastrand.NumberN[float32](10.5)

fmt.Printf("int32: %d, uint64: %d, float32: %f\n", num_i32, num_u64, num_f32)
```

### Strings and Bytes

#### Predefined Charsets

Use built-in charsets for common use-cases.
- `CharsDigits`
- `CharsAlphabetLower`
- `CharsAlphabetUpper`
- `CharsAlphabet`
- `CharsAlphabetDigits`
- `CharsSymbolChars`
- `CharsAll`

```go
// 10-digit PIN
pin := fastrand.String(10, fastrand.CharsDigits)

// 8-character lowercase username
user := fastrand.String(8, fastrand.CharsAlphabetLower)

// 32-character API key
apiKey := fastrand.String(32, fastrand.CharsAlphabetDigits)

fmt.Printf("PIN: %s\nUser: %s\nAPI Key: %s\n", pin, user, apiKey)
```

#### Bytes and Hex

```go
// Generate 8 random bytes
rawBytes := fastrand.Bytes(8)

// Generate a 16-byte hex string (results in 32 characters)
hexStr := fastrand.Hex(16)

fmt.Printf("Bytes: %v\nHex: %s\n", rawBytes, hexStr)
```

### Cryptographically Secure Functions

For sensitive data like passwords, session tokens, or cryptographic keys, use the `Secure*` variants.

```go
// Secure integer between 1000 and 2000
s_int, _ := fastrand.SecureInt(1000, 2000)

// Generate a secure 32-character password with symbols
password, _ := fastrand.SecureString(32, fastrand.CharsAll)

// Generate 64 secure bytes
s_bytes, _ := fastrand.SecureBytes(64)

fmt.Printf("Secure Int: %d\nSecure Password: %s\nSecure Bytes: %v\n", s_int, password, s_bytes)
```

### Slices and Choices

```go
// Choose one element from a slice
names := []string{"Alice", "Bob", "Charlie", "Diana"}
chosenName := fastrand.Choice(names)
fmt.Printf("Chosen Name: %s\n", chosenName)

// Choose 3 unique elements from a slice
chosenMultiple := fastrand.ChoiceMultiple(names, 3)
fmt.Printf("Chosen Multiple: %v\n", chosenMultiple)

// Shuffle a slice in-place
numbers := []int{1, 2, 3, 4, 5, 6, 7}
fastrand.Shuffle(len(numbers), func(i, j int) {
	numbers[i], numbers[j] = numbers[j], numbers[i]
})
fmt.Printf("Shuffled: %v\n", numbers)
```

### Network and IDs

```go
// Generate a random IPv4 address
ip4 := fastrand.IPv4()

// Generate a random IPv6 address
ip6 := fastrand.IPv6()

// Generate a fast (non-secure) V4 UUID
uuid1 := fastrand.MustFastUUID()

// Generate a secure V4 UUID
uuid2, _ := fastrand.SecureUUID()

fmt.Printf("IPv4: %s\nIPv6: %s\nFast UUID: %x\nSecure UUID: %x\n", ip4, ip6, uuid1, uuid2)
```

### The `Randomizer` Engine

The `Randomizer` is the most powerful feature for complex data generation. It parses a byte slice or string and replaces placeholders with random data.

**Placeholder Format:** `{RAND[OM];[LENGTH];[TYPE]}`
- `OM`: Optional, `{RANDOM...}` works too.
- `LENGTH`: Optional integer. Defaults vary by type.
- `TYPE`: Optional keyword. Defaults to `CharsAll`.

**Available Types:**
`ABL`, `ABU`, `ABR`, `DIGIT`, `HEX`, `UUID`, `IPV4`, `IPV6`, `EMAIL`, `BYTES`, `SPACE`, `NULL`.

```go
// Template with multiple placeholders
template := "POST /api/v1/data HTTP/1.1\n" +
	"Host: example.com\n" +
	"User-Agent: Client/{RAND;3;DIGIT}\n" +
	"X-Request-ID: {RAND;UUID}\n" +
	"X-Forwarded-For: {RAND;IPV4}\n" +
	"Authorization: Bearer {RANDOM;40;HEX}\n\n" +
	"{\"user_id\": \"{RAND;8;ABL}\", \"data\": \"{RAND;20}\"}"

// Generate the final data
randomizedRequest := fastrand.RandomizerString(template)

fmt.Println(randomizedRequest)
```

**Output of the `Randomizer` example:**
```
POST /api/v1/data HTTP/1.1
Host: example.com
User-Agent: Client/821
X-Request-ID: 7d3c9a1b-9e4f-4a8d-8c1b-2f3a4e5d6f7a
X-Forwarded-For: 198.51.100.27
Authorization: Bearer 8a3f2b1e9c4d58a3b1e9f4a2c8d7e6f5d4c3b2a1a0b9c8d7e6f5a4b3c2d1e0f9a8b7c6d5
{"user_id": "qwerasdf", "data": "A9s(d!f@g#h$j%k^l&"}
```

## Performance

-   **Fast Source (PCG)**: The default non-secure generator is based on the PCG algorithm, which is statistically excellent and significantly faster than the standard library's `math/rand`.
-   **Secure Source (ChaCha8)**: The secure generator uses `ChaCha8`, which is widely recognized for its strong security guarantees and excellent performance, often outperforming AES-based CSPRNGs.
-   **Optimized Algorithms**: Functions like `String` and `Bytes` are optimized to reduce allocations and improve throughput.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
