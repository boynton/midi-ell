# midi-ell
The Ell language combined with real-time MIDI I/O.

## Usage

    $ go get github.com/boynton/midi-ell/...
    $ mell
    ell v0.2 (with midi)
    ? (midi-open input: "USB Oxygen 8 v2") ;; the input is optional, but needed for midi-test3
    = {output: "IAC Driver Bus 1" bufsize: 1024 input: "USB Oxygen 8 v2"}
    ? (midi-test)
    ; should play every note, ascending
    ? (midi-test2)
    ; should play random notes forever
    <ctrl-c>
    *** [interrupt: ] [in random-phrase]
    ? (midi-test3)
    [1445395744.8638463]
    (1445395745.2833228 144 60 62)
    (1445395745.7293227 144 60 0)
    (1445395746.2673228 144 62 77)
    (1445395746.4023228 144 62 0)
    (1445395746.4913228 144 63 79)
    <ctrl-c>
    *** [interrupt: ] [in random-phrase]
    ? <ctrl-d>
    ; exiting will make all notes turn off

## License

Copyright 2015 Lee Boynton

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
                                                                                                                                   
  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
