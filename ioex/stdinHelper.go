package ioex

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/season-studio/go-utils/misc"

	"golang.org/x/term"
)

var (
	ctrlChan           = make(chan byte)
	inputChan chan any = nil
	locker    sync.Mutex
)

func getByteRaw(reader *bufio.Reader) {
	var retByte byte = 0
	state, err := term.MakeRaw(int(os.Stdin.Fd()))
	defer func() {
		if err == nil {
			term.Restore(int(os.Stdin.Fd()), state)
		}
		inputChan <- retByte
	}()
	for {
		retByte, err = reader.ReadByte()
		if err != nil {
			continue
		}
		return
	}
}

func getByte(reader *bufio.Reader) {
	for {
		b, err := reader.ReadByte()
		if err != nil {
			continue
		}
		inputChan <- b
		return
	}
}

func getStr(reader *bufio.Reader, buf *misc.ByteBuffer) {
	for {
		b, err := reader.ReadByte()
		if err != nil {
			continue
		}
		if b == '\n' {
			str := strings.Trim(string(buf.Bytes()), "\r\n")
			inputChan <- str
			buf.Reset()
			return
		} else {
			buf.Write(b)
		}
	}
}

func StartupStdinInteractive() {
	locker.Lock()
	defer locker.Unlock()

	if inputChan != nil {
		return
	}

	reader := bufio.NewReader(os.Stdin)
	inputChan = make(chan any)
	buf := misc.CreateByteBuffer(256, 256)

	go func() {
		for {
			ctrl, ok := <-ctrlChan
			if !ok {
				return
			}
			switch ctrl {
			case 0:
				getByteRaw(reader)
			case 1:
				getByte(reader)
			default:
				getStr(reader, buf)
			}
		}
	}()
}

func ReadStdinByte() (byte, error) {
	// locker.Lock()
	// defer locker.Unlock()

	if inputChan == nil {
		return 0, fmt.Errorf("stdin no ready")
	}
	ctrlChan <- 0
	v := <-inputChan
	switch s := v.(type) {
	case byte:
		return s, nil
	case string:
		if len(s) > 0 {
			return s[0], nil
		} else {
			return 0, nil
		}
	default:
		return 0, fmt.Errorf("input is not byte")
	}
}

func ReadStdin() (string, error) {
	// locker.Lock()
	// defer locker.Unlock()

	if inputChan == nil {
		return "", fmt.Errorf("stdin no ready")
	}
	ctrlChan <- 255
	v := <-inputChan
	if s, ok := v.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("input is not string")
}
