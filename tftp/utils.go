package tftp

import (
	"math"
	"os"
	"time"
)

func uint16ToBytes(v uint16) (b []byte) {
	b = make([]byte, 2)
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	b[0] = byte(v >> 8)
	b[1] = byte(v)
	return b
}

func bytesToUint16(b []byte) uint16 {
	_ = b[1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[1]) | uint16(b[0])<<8
}

func bytesFill(bs []byte, size int) ([]byte) {
	if size <= len(bs) {
		return bs
	}
	nbs := make([]byte, size)
	copy(nbs[:len(bs)], bs)
	return nbs
}

// 切割数据
func splitDataSegment(data []byte, size int) ([][]byte, int) {
	min := func(v1, v2 int) int {
		return int(math.Min(float64(v1), float64(v2)))
	}
	result := [][]byte{}
	segment := 0
	l := len(data)
	for i := 0; i <= l; i += size {
		minB := min(i, l)
		maxB := min(i+size, l)
		item := make([]byte, maxB-minB)
		item = data[minB:maxB]
		result = append(result, item)
		segment += 1
	}
	return result, segment
}

func saveFileData(fileName string, data []byte) error {
	os.Remove(fileName)
	fs, err := os.Create(fileName)
	defer fs.Close()
	if err == nil {
		_, err = fs.Write(data)
		return err
	}
	return err
}

func retryFunc(f func() error, time int) error {
	var err error
	for i := 0; i < time; time++ {
		err = f()
		if err == nil {
			return nil
		}
	}
	return err
}

func waitRecvTimeout(f func(), timeout time.Duration) {
	wg := make(chan struct{}, )
	go func() {
		f()
		wg <- struct{}{}
	}()
	ticker := time.NewTicker(timeout)
	select {
	case <-ticker.C:
		return
	case <-wg:
		return
	}
}

func fileIsExist(fileName string) bool {
	fs, err := os.Open(fileName)
	defer fs.Close()
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}