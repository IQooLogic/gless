package ansi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Style represents the current text styling
type Style struct {
	FgColor   string // Foreground color code
	BgColor   string // Background color code
	Bold      bool
	Dim       bool
	Italic    bool
	Underline bool
	Reverse   bool
}

// Reset resets all styles to default
func (s *Style) Reset() {
	s.FgColor = ""
	s.BgColor = ""
	s.Bold = false
	s.Dim = false
	s.Italic = false
	s.Underline = false
	s.Reverse = false
}

// Segment represents a text segment with its styling
type Segment struct {
	Text  string
	Style Style
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// ParseLine parses a line containing ANSI codes and returns segments with styles
func ParseLine(line string) []Segment {
	var segments []Segment
	currentStyle := Style{}
	lastIndex := 0

	matches := ansiRegex.FindAllStringIndex(line, -1)

	for _, match := range matches {
		// Add text before this ANSI code if any
		if match[0] > lastIndex {
			text := line[lastIndex:match[0]]
			if text != "" {
				segments = append(segments, Segment{
					Text:  text,
					Style: currentStyle,
				})
			}
		}

		// Parse and apply the ANSI code
		ansiCode := line[match[0]:match[1]]
		applyANSICode(&currentStyle, ansiCode)

		lastIndex = match[1]
	}

	// Add remaining text after last ANSI code
	if lastIndex < len(line) {
		text := line[lastIndex:]
		if text != "" {
			segments = append(segments, Segment{
				Text:  text,
				Style: currentStyle,
			})
		}
	}

	// If no ANSI codes found, return the whole line as one segment
	if len(segments) == 0 && line != "" {
		segments = append(segments, Segment{
			Text:  line,
			Style: currentStyle,
		})
	}

	return segments
}

// applyANSICode applies an ANSI escape code to the current style
func applyANSICode(style *Style, code string) {
	// Extract the numeric codes from the sequence
	// Format: \x1b[<codes>m
	code = strings.TrimPrefix(code, "\x1b[")
	code = strings.TrimSuffix(code, "m")

	if code == "" || code == "0" {
		style.Reset()
		return
	}

	codes := strings.Split(code, ";")
	for i := 0; i < len(codes); i++ {
		num, err := strconv.Atoi(codes[i])
		if err != nil {
			continue
		}

		switch num {
		case 0: // Reset
			style.Reset()
		case 1: // Bold
			style.Bold = true
		case 2: // Dim
			style.Dim = true
		case 3: // Italic
			style.Italic = true
		case 4: // Underline
			style.Underline = true
		case 7: // Reverse
			style.Reverse = true
		case 22: // Normal intensity (not bold, not dim)
			style.Bold = false
			style.Dim = false
		case 23: // Not italic
			style.Italic = false
		case 24: // Not underlined
			style.Underline = false
		case 27: // Not reversed
			style.Reverse = false
		case 30, 31, 32, 33, 34, 35, 36, 37: // Foreground colors
			style.FgColor = code
		case 39: // Default foreground
			style.FgColor = ""
		case 40, 41, 42, 43, 44, 45, 46, 47: // Background colors
			style.BgColor = code
		case 49: // Default background
			style.BgColor = ""
		case 90, 91, 92, 93, 94, 95, 96, 97: // Bright foreground colors
			style.FgColor = code
		case 100, 101, 102, 103, 104, 105, 106, 107: // Bright background colors
			style.BgColor = code
		case 38: // Extended foreground color
			if i+1 < len(codes) {
				if codes[i+1] == "5" && i+2 < len(codes) {
					// 256-color mode: 38;5;<color>
					style.FgColor = fmt.Sprintf("38;5;%s", codes[i+2])
					i += 2
				} else if codes[i+1] == "2" && i+4 < len(codes) {
					// RGB mode: 38;2;<r>;<g>;<b>
					style.FgColor = fmt.Sprintf("38;2;%s;%s;%s", codes[i+2], codes[i+3], codes[i+4])
					i += 4
				}
			}
		case 48: // Extended background color
			if i+1 < len(codes) {
				if codes[i+1] == "5" && i+2 < len(codes) {
					// 256-color mode: 48;5;<color>
					style.BgColor = fmt.Sprintf("48;5;%s", codes[i+2])
					i += 2
				} else if codes[i+1] == "2" && i+4 < len(codes) {
					// RGB mode: 48;2;<r>;<g>;<b>
					style.BgColor = fmt.Sprintf("48;2;%s;%s;%s", codes[i+2], codes[i+3], codes[i+4])
					i += 4
				}
			}
		}
	}
}

// StripANSI removes all ANSI codes from a string
func StripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// RenderSegment converts a segment back to a string with ANSI codes
func RenderSegment(seg Segment) string {
	if seg.Text == "" {
		return ""
	}

	var codes []string

	if seg.Style.Bold {
		codes = append(codes, "1")
	}
	if seg.Style.Dim {
		codes = append(codes, "2")
	}
	if seg.Style.Italic {
		codes = append(codes, "3")
	}
	if seg.Style.Underline {
		codes = append(codes, "4")
	}
	if seg.Style.Reverse {
		codes = append(codes, "7")
	}
	if seg.Style.FgColor != "" {
		codes = append(codes, seg.Style.FgColor)
	}
	if seg.Style.BgColor != "" {
		codes = append(codes, seg.Style.BgColor)
	}

	if len(codes) == 0 {
		return seg.Text
	}

	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", strings.Join(codes, ";"), seg.Text)
}

// HighlightText highlights occurrences of searchTerm in segments
// Returns new segments with highlighting applied (yellow background)
func HighlightText(segments []Segment, searchTerm string) []Segment {
	if searchTerm == "" {
		return segments
	}

	searchLower := strings.ToLower(searchTerm)
	var result []Segment

	for _, seg := range segments {
		if seg.Text == "" {
			continue
		}

		// Search for matches in this segment
		textLower := strings.ToLower(seg.Text)
		lastIndex := 0

		for {
			// Find next occurrence
			idx := strings.Index(textLower[lastIndex:], searchLower)
			if idx == -1 {
				// No more matches, add remaining text
				if lastIndex < len(seg.Text) {
					result = append(result, Segment{
						Text:  seg.Text[lastIndex:],
						Style: seg.Style,
					})
				}
				break
			}

			idx += lastIndex

			// Add text before match
			if idx > lastIndex {
				result = append(result, Segment{
					Text:  seg.Text[lastIndex:idx],
					Style: seg.Style,
				})
			}

			// Add highlighted match
			highlightStyle := seg.Style
			highlightStyle.BgColor = "43" // Yellow background
			highlightStyle.FgColor = "30" // Black foreground for better contrast
			result = append(result, Segment{
				Text:  seg.Text[idx : idx+len(searchTerm)],
				Style: highlightStyle,
			})

			lastIndex = idx + len(searchTerm)
		}
	}

	return result
}
