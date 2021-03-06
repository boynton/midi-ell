/*
Copyright 2014 Lee Boynton

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package midi

import (
	"github.com/boynton/ell"
	"github.com/rakyll/portmidi"
	"os"
	"sync"
	"time"
)

type Extension struct {
}

func (*Extension) Init() error {
	midiPath := os.Getenv("GOPATH") + "/src/github.com/boynton/midi-ell"
	ell.AddEllDirectory(midiPath)
	midiInit()
	return ell.Load("midi")
}

func (*Extension) Cleanup() {
	midiClose(nil)
}

func (*Extension) String() string {
	return "midi"
}

var inputKey = ell.Intern("input:")
var outputKey = ell.Intern("output:")
var bufsizeKey = ell.Intern("bufsize:")

func midiInit() error {
	ell.DefineFunctionKeyArgs("midi-open", midiOpen, ell.NullType,
		[]*ell.Object{ell.StringType, ell.StringType, ell.NumberType},
		[]*ell.Object{ell.EmptyString, ell.EmptyString, ell.Number(1024)},
		[]*ell.Object{inputKey, outputKey, bufsizeKey})
	ell.DefineFunction("midi-write", midiWrite, ell.NullType, ell.NumberType, ell.NumberType, ell.NumberType, ell.NumberType)
	ell.DefineFunction("midi-close", midiClose, ell.NullType)
	ell.DefineFunction("midi-listen", midiListen, ell.ChannelType)
	ell.DefineFunction("midi-time", midiTime, ell.NumberType)
	return nil
}

var midiOpened = false
var midiInDevice string
var midiOutDevice string
var midiBufsize int64
var midiBaseTime float64

var midiOut *portmidi.Stream
var midiIn *portmidi.Stream
var midiChannel chan portmidi.Event
var midiMutex = &sync.Mutex{}

func findMidiInputDevice(name string) portmidi.DeviceID {
	devcount := portmidi.CountDevices()
	for i := 0; i < devcount; i++ {
		id := portmidi.DeviceID(i)
		info := portmidi.Info(id)
		if info.IsInputAvailable {
			if info.Name == name {
				return id
			}
		}
	}
	return portmidi.DeviceID(-1)
}

func findMidiOutputDevice(name string) (portmidi.DeviceID, string) {
	devcount := portmidi.CountDevices()
	for i := 0; i < devcount; i++ {
		id := portmidi.DeviceID(i)
		info := portmidi.Info(id)
		if info.IsOutputAvailable {
			if info.Name == name {
				return id, info.Name
			}
		}
	}
	id := portmidi.DefaultOutputDeviceID()
	info := portmidi.Info(id)
	return id, info.Name
}

func midiOpen(argv []*ell.Object) (*ell.Object, error) {
	//	defaultInput := "USB Oxygen 8 v2"
	//	defaultOutput := "IAC Driver Bus 1"
	latency := int64(10)
	if !midiOpened {
		err := portmidi.Initialize()
		if err != nil {
			return nil, err
		}
		midiOpened = true
		midiInDevice = ell.StringValue(argv[0])
		midiOutDevice = ell.StringValue(argv[1])
		midiBufsize = ell.Int64Value(argv[2])

		outdev, outname := findMidiOutputDevice(midiOutDevice)
		out, err := portmidi.NewOutputStream(outdev, midiBufsize, latency)
		if err != nil {
			return nil, err
		}
		midiOut = out
		midiOutDevice = outname
		if midiInDevice != "" {
			indev := findMidiInputDevice(midiInDevice)
			if indev >= 0 {
				in, err := portmidi.NewInputStream(indev, midiBufsize)
				if err != nil {
					return nil, err
				}
				midiIn = in
			}
		}
		midiBaseTime = ell.Now()

	}
	result := ell.MakeStruct(4)
	if midiInDevice != "" {
		ell.Put(result, inputKey, ell.String(midiInDevice))
	}
	if midiOutDevice != "" {
		ell.Put(result, outputKey, ell.String(midiOutDevice))
	}
	ell.Put(result, bufsizeKey, ell.Number(float64(midiBufsize)))
	return result, nil
}

func midiAllNotesOff() {
	midiOut.WriteShort(0xB0, 0x7B, 0x00)
}

func midiClose(argv []*ell.Object) (*ell.Object, error) {
	midiMutex.Lock()
	if midiOut != nil {
		midiAllNotesOff()
		midiOut.Close()
		midiOut = nil
	}
	midiMutex.Unlock()
	ell.Sleep(0.5)
	return ell.Null, nil
}

func midiTime(argv []*ell.Object) (*ell.Object, error) {
	ts := (float64(portmidi.Time()) / 1000.0) + midiBaseTime
	return ell.Number(ts), nil
}

// (midi-write (midi-time) 144 60 80) -> middle C note on
// (midi-write 0 128 60 0) -> middle C note off
//note that the timestamp (first argument) can be 0, which means use portmidi.WriteShort is used (using current portmidi.Time()
func midiWrite(argv []*ell.Object) (*ell.Object, error) {
	sec := ell.Float64Value(argv[0])
	status := ell.Int64Value(argv[1])
	data1 := ell.Int64Value(argv[2])
	data2 := ell.Int64Value(argv[3])
	var err error
	midiMutex.Lock()
	if midiOut != nil {
		if sec == 0.0 {
			err = midiOut.WriteShort(status, data1, data2) //this always uses portmidi.Time() as the ts
		} else {
			ts := int64((sec - midiBaseTime) * 1000.0)
			evt := portmidi.Event{
				Timestamp: portmidi.Timestamp(ts),
				Status:    status,
				Data1:     data1,
				Data2:     data2,
			}
			err = midiOut.Write([]portmidi.Event{evt})
		}
	}
	midiMutex.Unlock()
	return ell.Null, err
}

func midiListen(argv []*ell.Object) (*ell.Object, error) {
	ch := ell.Null
	midiMutex.Lock()
	if midiIn != nil {
		ch = ell.Channel(int(midiBufsize), "midi")
		go func(s *portmidi.Stream, ch *ell.Object) {
			for {
				time.Sleep(10 * time.Millisecond)
				events, err := s.Read(1024)
				if err != nil {
					continue
				}
				channel := ell.ChannelValue(ch)
				if channel != nil {
					for _, ev := range events {
						ts := (float64(ev.Timestamp) / 1000) + midiBaseTime
						st := ev.Status
						d1 := ev.Data1
						d2 := ev.Data2
						channel <- ell.List(ell.Number(ts), ell.Number(float64(st)), ell.Number(float64(d1)), ell.Number(float64(d2)))
					}
				}
			}
		}(midiIn, ch)
	}
	midiMutex.Unlock()
	return ch, nil
}
