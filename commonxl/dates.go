package commonxl

import (
	"regexp"
	"strings"
	"time"
)

// GetDateTime is a wrapper for chained invocation of
// FormatIsDateTime and ConvertToDate.
func GetDateTime(fno uint16, f string, val float64, mode1904 bool) (time.Time, bool) {
	if !FormatIsDateTime(fno, f) {
		return time.Time{}, false
	}
	return ConvertToDate(val, mode1904), true
}

// ConvertToDate converts a floating-point value using the
// Excel date serialization conventions.
func ConvertToDate(val float64, mode1904 bool) time.Time {
	// http://web.archive.org/web/20190808062235/http://aa.usno.navy.mil/faq/docs/JD_Formula.php
	v := int(val)
	if v < 61 {
		jdate := val + 0.5
		if mode1904 {
			jdate += 2416480.5
		} else {
			jdate += 2415018.5
		}
		JD := int(jdate)
		frac := jdate - float64(JD)

		L := JD + 68569
		N := 4 * L / 146097
		L = L - (146097*N+3)/4
		I := 4000 * (L + 1) / 1461001
		L = L - 1461*I/4 + 31
		J := 80 * L / 2447
		day := L - 2447*J/80
		L = J / 11
		month := time.Month(J + 2 - 12*L)
		year := 100*(N-49) + I + L

		t := time.Duration(float64(time.Hour*24) * frac)
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Add(t)
	}
	frac := val - float64(v)
	date := time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC)
	if !mode1904 {
		date = time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	}

	t := time.Duration(float64(time.Hour*24) * frac)
	return date.AddDate(0, 0, v).Add(t)
}

// FormatIsDateTime returns true if the given format number or
// format string contains a date/time formatting instruction.
func FormatIsDateTime(fno uint16, f string) bool {
	if _, ok := builtInDateFormats[fno]; ok {
		return true
	}

	// fast path
	if !strings.ContainsAny(f, "ymdhs") {
		return false
	}

	// remove colors, escaped characters, and quoted text
	f = formatMatchBrackets.ReplaceAllString(f, "")
	f = formatMatchEscaped.ReplaceAllString(f, "")
	f = formatMatchTextLiteral.ReplaceAllString(f, "")

	// if there's still any of ymdhs in there, it's a date
	return strings.ContainsAny(f, "ymdhs")
}

var (
	formatMatchBrackets    = regexp.MustCompile(`\[[^\]]*\]`)
	formatMatchEscaped     = regexp.MustCompile(`\\.`)
	formatMatchTextLiteral = regexp.MustCompile(`"[^"]*"`)
)

// 0x0001 = date   0b0010 = time    0b0011 = date+time
var builtInDateFormats = map[uint16]byte{
	14: 1, 15: 1, 16: 1, 17: 1, 18: 2, 19: 2, 20: 2, 21: 2, 22: 3,
	45: 2, 46: 2, 47: 2, 27: 1, 28: 1, 29: 1, 30: 1, 31: 1, 32: 2,
	33: 2, 34: 2, 35: 2, 36: 1, 50: 1, 51: 1, 52: 1, 53: 1, 54: 1,
	55: 2, 56: 2, 57: 1, 58: 1, 71: 1, 72: 1, 73: 1, 74: 1, 75: 2,
	76: 2, 77: 3, 78: 2, 79: 2, 80: 2, 81: 1,
}

// mapping of standard built-ins to Go date format strings.
var builtInGoFormats = map[uint16]string{
	14: `01-02-06`,
	15: `2-Jan-06`,
	16: `2-Jan`,
	17: `Jan-06`,
	18: `3:04 AM`,
	19: `3:04:05 AM`,
	20: `15:04`,
	21: `15:04:05`,
	22: `1/2/06 15:04`,
	45: `04:05`,
	46: `3:04:05`,
	47: `0405.9`,

	// zh-cn format codes
	27: `2006"年"1"月"`,
	28: `1"月"2"日"`,
	29: `1"月"2"日"`,
	30: `1-2-06`,
	31: `2006"年"1"月"2"日"`,
	32: `15"时"04"分"`,
	33: `15"时"04"分"05"秒"`,
	34: `上午/下午 3"时"04"分"`, // FIXME: am/pm not properly handled here
	35: `上午/下午 3"时"04"分"05"秒"`,
	36: `2006"年"2"月"`,
	50: `2006"年"2"月"`,
	51: `1"月"2"日"`,
	52: `2006"年"1"月"`,
	53: `1"月"2"日"`,
	54: `1"月"2"日"`,
	55: `上午/下午 3"时"04"分"`, // TODO am/pm
	56: `上午/下午 3"时"04"分"05"秒`,
	57: `2006"年"1"月"`,
	58: `1"月"2"日"`,

	71: `2/1/2006`,
	72: `2-Jan-06`,
	73: `2-Jan`,
	74: `Jan-06`,
	75: `15:04`,
	76: `15:04:05`,
	77: `2/1/2006 15:04`,
	78: `04:05`,
	79: `15:04:05`,
	80: `04:05.9`,
	81: `2/1/06`,
}