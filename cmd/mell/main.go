package main

import (
	"github.com/boynton/ell"
	midi "github.com/boynton/midi-ell"
)

func main() {
	ell.Main(new(midi.Extension))
}
