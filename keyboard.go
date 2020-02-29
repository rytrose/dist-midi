package util

import (
	"fmt"
	"os"
	"sync"

	"github.com/eiannone/keyboard"
)

// KeyboardReader reads keyboard input and maintains key press toggle state.
type KeyboardReader struct {
	state     map[rune]bool
	functions map[rune]func(bool, bool)
	moot      sync.Mutex
	reading   bool
	stopChan  chan struct{}
}

// NewKeyboardReader is a KeyboardReader factory.
func NewKeyboardReader() *KeyboardReader {
	return &KeyboardReader{
		state:     make(map[rune]bool),
		functions: make(map[rune]func(bool, bool)),
		stopChan:  make(chan struct{}),
	}
}

// Read begins reading from standard input, updating state of key press toggle.
func (kr *KeyboardReader) Read() {
	err := keyboard.Open()
	Must(err)

	kr.reading = true
	go func(k *KeyboardReader) {
		defer keyboard.Close()
		for {
			r, key, err := keyboard.GetKey()
			Must(err)
			if key == keyboard.KeyCtrlC {
				fmt.Println("Goodbye!")
				os.Exit(0)
			}

			k.moot.Lock()
			ks := k.state[r]
			f, ok := k.functions[r]
			if ok {
				go f(ks, !ks)
			}
			k.state[r] = !ks
			k.moot.Unlock()

			select {
			case <-k.stopChan:
				return
			default:
			}
		}
	}(kr)
}

// Close stops reading from standard input.
func (kr *KeyboardReader) Close() {
	if !kr.reading {
		return
	}
	kr.reading = false
	kr.state = make(map[rune]bool)
	kr.functions = make(map[rune]func(bool, bool))
	fmt.Println("Keyboard reading will stop on next key press.")
	kr.stopChan <- struct{}{}
}

// GetState gets the current toggle state of a rune since keyboard input reading has started.
func (kr *KeyboardReader) GetState(r rune) bool {
	if !kr.reading {
		fmt.Println("[WARN] Not reading keyboard input! Call Read() to read keyboard input.")
		return false
	}
	kr.moot.Lock()
	defer kr.moot.Unlock()
	s, ok := kr.state[r]
	if ok {
		return s
	}
	return ok
}

// Register sets up a function to be called on rune state toggle.
func (kr *KeyboardReader) Register(r rune, f func(bool, bool)) {
	kr.moot.Lock()
	defer kr.moot.Unlock()
	kr.functions[r] = f
}
