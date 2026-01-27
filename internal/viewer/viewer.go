package viewer

import (
	"fmt"
	"os"
	"strings"

	"github.com/iqool/gless/internal/ansi"
	"github.com/iqool/gless/internal/reader"
	"golang.org/x/term"
)

// Viewer manages the terminal UI for viewing files
type Viewer struct {
	fileReader      *reader.FileReader
	currentLine     int // Current top line being displayed (0-based)
	width           int
	height          int
	terminalState   *term.State
	searchTerm      string
	searchResults   []int // Line numbers containing search term
	currentResult   int   // Index in searchResults
	showLineNumbers bool
	quit            bool
}

// NewViewer creates a new viewer for the given file
func NewViewer(fileReader *reader.FileReader) *Viewer {
	return &Viewer{
		fileReader:      fileReader,
		currentLine:     0,
		searchResults:   []int{},
		currentResult:   -1,
		showLineNumbers: false,
		quit:            false,
	}
}

// Run starts the viewer
func (v *Viewer) Run() error {
	// Load file content
	if err := v.fileReader.Load(); err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	// Enter raw mode
	if err := v.enterRawMode(); err != nil {
		return fmt.Errorf("failed to enter raw mode: %w", err)
	}
	defer v.exitRawMode()

	// Clear screen and hide cursor
	v.clearScreen()
	v.hideCursor()
	defer v.showCursor()

	// Initial render
	v.updateSize()
	v.render()

	// Main loop
	return v.handleInput()
}

// enterRawMode puts the terminal into raw mode
func (v *Viewer) enterRawMode() error {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	v.terminalState = oldState
	return nil
}

// exitRawMode restores the terminal to normal mode
func (v *Viewer) exitRawMode() {
	if v.terminalState != nil {
		term.Restore(int(os.Stdin.Fd()), v.terminalState)
	}
}

// updateSize updates the terminal dimensions
func (v *Viewer) updateSize() {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		v.width = 80
		v.height = 24
	} else {
		v.width = width
		v.height = height
	}
}

// clearScreen clears the terminal screen
func (v *Viewer) clearScreen() {
	fmt.Print("\x1b[2J")
}

// hideCursor hides the terminal cursor
func (v *Viewer) hideCursor() {
	fmt.Print("\x1b[?25l")
}

// showCursor shows the terminal cursor
func (v *Viewer) showCursor() {
	fmt.Print("\x1b[?25h")
}

// render draws the current view
func (v *Viewer) render() {
	// Move cursor to home position
	fmt.Print("\x1b[H")

	totalLines := v.fileReader.LineCount()
	displayHeight := v.height - 1 // Reserve last line for status bar

	// Ensure currentLine is within bounds
	if v.currentLine < 0 {
		v.currentLine = 0
	}
	if v.currentLine >= totalLines {
		v.currentLine = totalLines - 1
	}
	if v.currentLine < 0 {
		v.currentLine = 0
	}

	// Get lines to display
	endLine := v.currentLine + displayHeight
	if endLine > totalLines {
		endLine = totalLines
	}

	lines, err := v.fileReader.GetLines(v.currentLine, endLine)
	if err != nil {
		lines = []string{fmt.Sprintf("Error reading lines: %v", err)}
	}

	// Display lines
	for i, line := range lines {
		lineNum := v.currentLine + i + 1 // 1-based for display

		// Clear line
		fmt.Print("\x1b[2K")

		// Line number prefix
		if v.showLineNumbers {
			fmt.Printf("\x1b[90m%6d\x1b[0m ", lineNum)
		}

		// Parse and render line with ANSI codes
		segments := ansi.ParseLine(line)

		// Apply search highlighting if we have a search term and this line is a match
		if v.searchTerm != "" && v.isSearchMatch(lineNum-1) {
			segments = ansi.HighlightText(segments, v.searchTerm)
		}

		for _, seg := range segments {
			fmt.Print(ansi.RenderSegment(seg))
		}

		// Move to next line if not the last line we're rendering
		if i < len(lines)-1 || i < displayHeight-1 {
			fmt.Print("\r\n")
		}
	}

	// Fill remaining lines if file is shorter than screen
	for i := len(lines); i < displayHeight; i++ {
		fmt.Print("\x1b[2K")
		fmt.Print("~")
		if i < displayHeight-1 {
			fmt.Print("\r\n")
		}
	}

	// Render status bar
	v.renderStatusBar()

	// Flush output
	os.Stdout.Sync()
}

// renderStatusBar renders the status bar at the bottom
func (v *Viewer) renderStatusBar() {
	fmt.Print("\r\n")
	fmt.Print("\x1b[2K") // Clear line
	fmt.Print("\x1b[7m") // Reverse video

	totalLines := v.fileReader.LineCount()
	filename := v.fileReader.Filename()

	var percentage int
	if totalLines > 0 {
		percentage = ((v.currentLine + 1) * 100) / totalLines
	}

	status := fmt.Sprintf(" %s | Line %d-%d/%d (%d%%)",
		filename,
		v.currentLine+1,
		min(v.currentLine+v.height-1, totalLines),
		totalLines,
		percentage)

	// Add search info if searching
	if v.searchTerm != "" {
		if len(v.searchResults) > 0 {
			status += fmt.Sprintf(" | Search: \"%s\" (%d/%d)",
				v.searchTerm,
				v.currentResult+1,
				len(v.searchResults))
		} else {
			status += fmt.Sprintf(" | Search: \"%s\" (no matches)", v.searchTerm)
		}
	}

	// Add help hint
	status += " | Press 'h' for help, 'q' to quit"

	// Calculate visual length (without ANSI codes) for proper padding/truncation
	visualLen := len(ansi.StripANSI(status))

	// Truncate if too long
	if visualLen > v.width {
		// Need to truncate - strip ANSI, truncate, then re-apply formatting
		stripped := ansi.StripANSI(status)
		status = stripped[:v.width]
	} else {
		// Pad to full width
		status += strings.Repeat(" ", v.width-visualLen)
	}

	fmt.Print(status)
	fmt.Print("\x1b[0m") // Reset
}

// Scroll scrolls the view by the specified number of lines
func (v *Viewer) Scroll(delta int) {
	v.currentLine += delta
	totalLines := v.fileReader.LineCount()

	if v.currentLine < 0 {
		v.currentLine = 0
	}
	maxLine := totalLines - (v.height - 1)
	if maxLine < 0 {
		maxLine = 0
	}
	if v.currentLine > maxLine {
		v.currentLine = maxLine
	}
}

// GoToLine moves to a specific line
func (v *Viewer) GoToLine(line int) {
	v.currentLine = line
	v.Scroll(0) // Normalize bounds
}

// isSearchMatch checks if a line number is in the search results
func (v *Viewer) isSearchMatch(lineNum int) bool {
	for _, matchLine := range v.searchResults {
		if matchLine == lineNum {
			return true
		}
	}
	return false
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
