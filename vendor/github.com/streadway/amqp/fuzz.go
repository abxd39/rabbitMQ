<<<<<<< HEAD
// +build gofuzz

package amqp

import "bytes"

func Fuzz(data []byte) int {
	r := reader{bytes.NewReader(data)}
	frame, err := r.ReadFrame()
	if err != nil {
		if frame != nil {
			panic("frame is not nil")
		}
		return 0
	}
	return 1
}
=======
// +build gofuzz

package amqp

import "bytes"

func Fuzz(data []byte) int {
	r := reader{bytes.NewReader(data)}
	frame, err := r.ReadFrame()
	if err != nil {
		if frame != nil {
			panic("frame is not nil")
		}
		return 0
	}
	return 1
}
>>>>>>> 7563d76138a21269a66175d7cdad4b86cc127ec9
