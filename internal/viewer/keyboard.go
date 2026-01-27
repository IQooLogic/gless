package viewer

import (
	"fmt"
	"os"
	"strings"

	"github.com/iqool/gless/internal/ansi"
	"golang.org/x/term"
)

// handleInput handles keyboard input
func (v *Viewer) handleInput() error {
	buf := make([]byte, 16)

	for !v.quit {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return err
		}

		if n == 0 {
			continue
		}

		// Handle input
		v.processInput(buf[:n])

		// Re-render
		v.render()
	}

	return nil
}

// processInput processes keyboard input
func (v *Viewer) processInput(input []byte) {
	// Check for escape sequences (arrow keys, etc.)
	if len(input) >= 3 && input[0] == 0x1b && input[1] == '[' {
		switch input[2] {
		case 'A': // Up arrow
			v.Scroll(-1)
		case 'B': // Down arrow
			v.Scroll(1)
		case 'C': // Right arrow (could be used for horizontal scrolling)
			// Not implemented yet
		case 'D': // Left arrow
			// Not implemented yet
		case '5': // Page Up (ESC[5~)
			if len(input) >= 4 && input[3] == '~' {
				v.Scroll(-(v.height - 2))
			}
		case '6': // Page Down (ESC[6~)
			if len(input) >= 4 && input[3] == '~' {
				v.Scroll(v.height - 2)
			}
		case 'H': // Home
			v.GoToLine(0)
		case 'F': // End
			v.GoToLine(v.fileReader.LineCount() - 1)
		}
		return
	}

	// Handle single character commands
	if len(input) == 1 {
		ch := input[0]
		switch ch {
		case 'q', 'Q': // Quit
			v.quit = true
		case 'h', 'H', '?': // Help
			v.showHelp()
		case 'g': // Go to first line
			v.GoToLine(0)
		case 'G': // Go to last line
			v.GoToLine(v.fileReader.LineCount() - 1)
		case 'j': // Down (vim-style)
			v.Scroll(1)
		case 'k': // Up (vim-style)
			v.Scroll(-1)
		case 'd': // Half page down
			v.Scroll((v.height - 2) / 2)
		case 'u': // Half page up
			v.Scroll(-(v.height - 2) / 2)
		case 'f', ' ': // Page down (space or f)
			v.Scroll(v.height - 2)
		case 'b': // Page up
			v.Scroll(-(v.height - 2))
		case '/': // Search
			v.enterSearchMode()
		case 'n': // Next search result
			v.nextSearchResult()
		case 'N': // Previous search result
			v.previousSearchResult()
		case '#': // Toggle line numbers
			v.showLineNumbers = !v.showLineNumbers
		case 0x03: // Ctrl+C
			v.quit = true
		}
	}
}

// showHelp displays the help screen
func (v *Viewer) showHelp() {
	v.clearScreen()
	fmt.Print("\x1b[H")

	help := []string{
		"",
		"                          GLess - Help",
		"                          ============",
		"",
		"  Navigation:",
		"    ↑, k           Move up one line",
		"    ↓, j           Move down one line",
		"    PageUp, b      Move up one page",
		"    PageDown, f, SPACE  Move down one page",
		"    u              Move up half page",
		"    d              Move down half page",
		"    Home, g        Go to first line",
		"    End, G         Go to last line",
		"",
		"  Search:",
		"    /              Enter search mode",
		"    n              Next search result",
		"    N              Previous search result",
		"",
		"  Display:",
		"    #              Toggle line numbers",
		"",
		"  Other:",
		"    h, ?           Show this help",
		"    q              Quit",
		"    Ctrl+C         Quit",
		"",
		"  GLess displays files with ANSI color codes preserved.",
		"",
		"",
		"                Press any key to continue...",
	}

	for _, line := range help {
		fmt.Print("\x1b[2K")
		fmt.Println(line)
	}

	// Wait for keypress
	buf := make([]byte, 16)
	os.Stdin.Read(buf)

	v.clearScreen()
}

// enterSearchMode prompts for search input
func (v *Viewer) enterSearchMode() {
	// Show cursor and move to bottom
	v.showCursor()
	fmt.Printf("\x1b[%d;1H", v.height)
	fmt.Print("\x1b[2K")
	fmt.Print("/")

	// Restore terminal for search input
	if v.terminalState != nil {
		term.Restore(int(os.Stdin.Fd()), v.terminalState)
	}

	// Read search term
	var searchTerm string
	fmt.Scanln(&searchTerm)

	// Re-enter raw mode
	v.enterRawMode()
	v.hideCursor()

	if searchTerm == "" {
		return
	}

	// Perform search
	v.searchTerm = searchTerm
	v.performSearch()

	// Jump to first result if found
	if len(v.searchResults) > 0 {
		v.currentResult = 0
		v.GoToLine(v.searchResults[0])
	}
}

// performSearch searches for the term in all lines
func (v *Viewer) performSearch() {
	v.searchResults = []int{}
	v.currentResult = -1

	totalLines := v.fileReader.LineCount()
	for i := 0; i < totalLines; i++ {
		line, err := v.fileReader.GetLine(i)
		if err != nil {
			continue
		}

		// Search in the stripped version (without ANSI codes)
		stripped := ansi.StripANSI(line)
		if strings.Contains(strings.ToLower(stripped), strings.ToLower(v.searchTerm)) {
			v.searchResults = append(v.searchResults, i)
		}
	}
}

// nextSearchResult jumps to the next search result
func (v *Viewer) nextSearchResult() {
	if len(v.searchResults) == 0 {
		return
	}

	v.currentResult++
	if v.currentResult >= len(v.searchResults) {
		v.currentResult = 0
	}

	v.GoToLine(v.searchResults[v.currentResult])
}

// previousSearchResult jumps to the previous search result
func (v *Viewer) previousSearchResult() {
	if len(v.searchResults) == 0 {
		return
	}

	v.currentResult--
	if v.currentResult < 0 {
		v.currentResult = len(v.searchResults) - 1
	}

	v.GoToLine(v.searchResults[v.currentResult])
}
