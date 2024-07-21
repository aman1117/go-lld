package classes

import (
	"fmt"
	"strings"
)

type Notepad struct {
	allContent []string
	undoStack  [][]string
	redoStack  [][]string
	buffer     []string
}

func NewNotepad(text string) *Notepad {
	return &Notepad{allContent: split(text, '\n')}
}

func split(text string, delim rune) []string {
	return strings.Split(text, string(delim))
}

func (n *Notepad) display() {
	for _, line := range n.allContent {
		fmt.Println(line)
	}
}

func (n *Notepad) displayRange(start, end int) bool {
	if start > len(n.allContent) || end > len(n.allContent) || start > end {
		fmt.Println("Invalid range")
		return false
	}

	for i := start - 1; i < end; i++ {
		fmt.Println(n.allContent[i])
	}
	return true
}

func (n *Notepad) insert(line int, text string) bool {
	if line > len(n.allContent) {
		fmt.Println("The value of line exceeds number of lines in the file")
		return false
	}

	//creating a deep copy of the allContent slice and then append this copy to the undoStack.
	//If you only stored a reference to allContent without making a deep copy, changes to allContent would also affect the stored states in the undoStack, because slices in Go are reference types. Thus, you wouldn’t be able to revert accurately.
	n.undoStack = append(n.undoStack, append([]string{}, n.allContent...))
	n.allContent[line-1] += text
	return true
}

func (n *Notepad) deleteLine(line int) bool {
	if line > len(n.allContent) {
		fmt.Println("The value of line exceeds number of lines in the file")
		return false
	}

	n.undoStack = append(n.undoStack, append([]string{}, n.allContent...))

	// [:line - 1] is not included in the new slice, [line:] is included
	n.allContent = append(n.allContent[:line-1], n.allContent[line:]...)
	return true
}

func (n *Notepad) deleteRange(start, end int) bool {
	if start > len(n.allContent) || end > len(n.allContent) || start > end {
		fmt.Println("Invalid range")
		return false
	}

	n.undoStack = append(n.undoStack, append([]string{}, n.allContent...))
	n.allContent = append(n.allContent[:start-1], n.allContent[end:]...)
	return true
}

func (n *Notepad) copyRange(start, end int) bool {
	if start > len(n.allContent) || end > len(n.allContent) || start > end {
		fmt.Println("Invalid range")
		return false
	}

	n.buffer = append([]string{}, n.allContent[start-1:end]...)
	return true
}

func (n *Notepad) paste(line int) bool {
	if line > len(n.allContent) {
		fmt.Println("The value of line exceeds number of lines in the file")
		return false
	}

	n.undoStack = append(n.undoStack, append([]string{}, n.allContent...))
	// line - 1 is not included in the new slice,
	n.allContent = append(n.allContent[:line-1], append(n.buffer, n.allContent[line-1:]...)...)
	return true
}

func (n *Notepad) undo() bool {
	if len(n.undoStack) == 0 {
		fmt.Println("Nothing to undo!")
		return false
	}

	n.redoStack = append(n.redoStack, append([]string{}, n.allContent...))
	n.allContent = n.undoStack[len(n.undoStack)-1]
	n.undoStack = n.undoStack[:len(n.undoStack)-1]
	return true
}

func (n *Notepad) redo() bool {
	if len(n.redoStack) == 0 {
		fmt.Println("Nothing to redo!")
		return false
	}

	n.undoStack = append(n.undoStack, append([]string{}, n.allContent...))
	n.allContent = n.redoStack[len(n.redoStack)-1]
	n.redoStack = n.redoStack[:len(n.redoStack)-1]
	return true
}

func NotePad() {
	notepad := NewNotepad("At the starting of the week\nI want to kiss my girl Aakanksha and then love her\nIt's the start of the week")
	notepad.display()
	fmt.Println("**************************** 0 ***********************************")
	fmt.Println("** Displaying content: only first two lines **")
	notepad.displayRange(1, 2)
	fmt.Println("**************************** 1 ***********************************")
	fmt.Println("** Inserting yeah to the first line **")
	notepad.insert(1, ", Yeah")
	notepad.display()
	fmt.Println("**************************** 2 ***********************************")
	fmt.Println("** Undoing last move **")
	notepad.undo()
	notepad.display()
	fmt.Println("***************************** 3 *********************************")
	fmt.Println("** Redoing last move **")
	notepad.redo()
	notepad.display()
	fmt.Println("****************************** 4 *******************************")
	fmt.Println("** Redoing last move **")
	notepad.redo()
	fmt.Println("******************************** 5 ****************************")
	fmt.Println("** Deleting first line **")
	notepad.deleteLine(1)
	notepad.display()
	fmt.Println("******************************* 6 ****************************")
	fmt.Println("** Undoing last move **")
	notepad.undo()
	notepad.display()
	fmt.Println("******************************* 7 ***************************")
	fmt.Println("** Undoing last move **")
	notepad.undo()
	notepad.display()
	fmt.Println("**************************** 8 ***********************************")
	fmt.Println("** After deletion of lines 1 to 2 **")
	notepad.deleteRange(1, 2)
	notepad.display()
	fmt.Println("***************************** 9 ****************************")
	fmt.Println("** Undoing last move **")
	notepad.undo()
	notepad.display()
	fmt.Println("***************************** 10 ***************************")
	fmt.Println("** Copying lines 1 to 2 and pasting them on 3rd line **")
	notepad.copyRange(1, 2)
	notepad.paste(3)
	notepad.display()
	fmt.Println("***************************** 11 **************************")
	fmt.Println("** Undoing last move **")
	notepad.undo()
	notepad.display()
	fmt.Println("****************************** 12 ************************")
	fmt.Println("** Redoing last move **")
	notepad.redo()
	notepad.display()
}
