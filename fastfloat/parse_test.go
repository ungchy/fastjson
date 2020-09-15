package fastfloat

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
)

func TestParseUint64BestEffort(t *testing.T) {
	f := func(s string, expectedNum uint64) {
		t.Helper()

		num := ParseUint64BestEffort(s)
		if num != expectedNum {
			t.Fatalf("unexpected umber parsed from %q; got %v; want %v", s, num, expectedNum)
		}
	}

	// Invalid first char
	f("", 0)
	f("   ", 0)
	f("foo", 0)
	f("-", 0)
	f("-foo", 0)
	f("-123", 0)

	// Invalid suffix
	f("1foo", 0)
	f("13223 ", 0)
	f("1-2", 0)

	// Int
	f("1", 1)
	f("123", 123)
	f("1234567890", 1234567890)
	f("9223372036854775807", 9223372036854775807)
	f("18446744073709551615", 18446744073709551615)

	// Too big int
	f("18446744073709551616", 0)
}

func TestParseInt64BestEffort(t *testing.T) {
	f := func(s string, expectedNum int64) {
		t.Helper()

		num := ParseInt64BestEffort(s)
		if num != expectedNum {
			t.Fatalf("unexpected umber parsed from %q; got %v; want %v", s, num, expectedNum)
		}
	}

	// Invalid first char
	f("", 0)
	f("   ", 0)
	f("foo", 0)
	f("-", 0)
	f("-foo", 0)

	// Invalid suffix
	f("1foo", 0)
	f("13223 ", 0)
	f("1-2", 0)

	// Int
	f("0", 0)
	f("-0", 0)
	f("1", 1)
	f("123", 123)
	f("-123", -123)
	f("1234567890", 1234567890)
	f("9223372036854775807", 9223372036854775807)
	f("-9223372036854775807", -9223372036854775807)

	// Too big int
	f("9223372036854775808", 0)
	f("18446744073709551615", 0)
}

func TestParseBestEffort(t *testing.T) {
	f := func(s string, expectedNum float64) {
		t.Helper()

		num := ParseBestEffort(s)
		if math.IsNaN(expectedNum) {
			if !math.IsNaN(num) {
				t.Fatalf("unexpected number parsed from %q; got %v; want %v", s, num, expectedNum)
			}
		} else if num != expectedNum {
			t.Fatalf("unexpected number parsed from %q; got %v; want %v", s, num, expectedNum)
		}
	}

	// Invalid first char
	f("", 0)
	f("  ", 0)
	f("foo", 0)
	f(" bar ", 0)
	f("-", 0)
	f("--", 0)
	f("-.", 0)
	f("-.e", 0)
	f("+112", 0)
	f("++", 0)
	f("e123", 0)
	f("E123", 0)
	f("-e12", 0)
	f(".", 0)
	f("..34", 0)
	f("-.32", 0)
	f("-.e3", 0)
	f(".e+3", 0)

	// Invalid suffix
	f("1foo", 0)
	f("1  foo", 0)
	f("12.34.56", 0)
	f("13e34.56", 0)
	f("12.34e56e4", 0)
	f("12.", 0)
	f("123..45", 0)
	f("123ee34", 0)
	f("123e", 0)
	f("123e+", 0)
	f("123E-", 0)
	f("123E+.", 0)
	f("-123e-23foo", 0)

	// Integer
	f("0", 0)
	f("-0", 0)
	f("0123", 123)
	f("-00123", -123)
	f("1", 1)
	f("-1", -1)
	f("1234567890123456", 1234567890123456)
	f("12345678901234567", 12345678901234567)
	f("123456789012345678", 123456789012345678)
	f("1234567890123456789", 1234567890123456789)
	f("12345678901234567890", 12345678901234567890)
	f("-12345678901234567890", -12345678901234567890)
	f("9223372036854775807", 9223372036854775807)
	f("18446744073709551615", 18446744073709551615)
	f("-18446744073709551615", -18446744073709551615)

	// Fractional part
	f("0.0", 0)
	f("0.000", 0)
	f("0.1", 0.1)
	f("0.3", 0.3)
	f("-0.1", -0.1)
	f("-0.123", -0.123)
	f("1.66", 1.66)
	f("12345.12345678901", 12345.12345678901)
	f("12345.123456789012", 12345.123456789012)
	f("12345.1234567890123", 12345.1234567890123)
	f("12345.12345678901234", 12345.12345678901234)
	f("12345.123456789012345", 12345.123456789012345)
	f("12345.1234567890123456", 12345.1234567890123456)
	f("12345.12345678901234567", 12345.12345678901234567)
	f("12345.123456789012345678", 12345.123456789012345678)
	f("12345.1234567890123456789", 12345.1234567890123456789)
	f("12345.12345678901234567890", 12345.12345678901234567890)
	f("-12345.12345678901234567890", -12345.12345678901234567890)
	f("18446744073709551615.18446744073709551615", 18446744073709551615.18446744073709551615)
	f("-18446744073709551615.18446744073709551615", -18446744073709551615.18446744073709551615)

	// Exponent part
	f("0e0", 0)
	f("123e+001", 123e1)
	f("0e12", 0)
	f("-0E123", 0)
	f("-0E-123", 0)
	f("-0E+123", 0)
	f("123e12", 123e12)
	f("-123E-12", -123e-12)
	f("-123e-400", 0)
	f("123e456", math.Inf(1))   // too big exponent
	f("-123e456", math.Inf(-1)) // too big exponent
	f("1e4", 1e4)
	f("-1E-10", -1e-10)

	// Fractional + exponent part
	f("0.123e4", 0.123e4)
	f("-123.456E-10", -123.456e-10)
	f("1.e4", 1.e4)
	f("-1.E-10", -1.e-10)

	// inf and nan
	f("inf", math.Inf(1))
	f("-Inf", math.Inf(-1))
	f("INF", math.Inf(1))
	f("nan", math.NaN())
	f("naN", math.NaN())
	f("NaN", math.NaN())
}

func TestParseBestEffortFuzz(t *testing.T) {
	r := rand.New(rand.NewSource(0))
	for i := 0; i < 100000; i++ {
		f := r.Float64()
		s := strconv.FormatFloat(f, 'g', -1, 64)
		numExpected, err := strconv.ParseFloat(s, 64)
		if err != nil {
			t.Fatalf("unexpected error when parsing %q: %s", s, err)
		}
		num := ParseBestEffort(s)
		if num != numExpected {
			t.Fatalf("unexpected number parsed from %q; got %g; want %g", s, num, numExpected)
		}
	}
}
