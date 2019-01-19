package tftp

import (
	"sync"
	"os"
)

type ConnInfo struct {
	FileInfo      *os.File
	Options       Options
	Mode          string
	DataBlockSize int
	BlockId       uint16 // 从1起索引
}

func NewConnInfo(fs *os.File, options Options, mode string, dataBlockSize int) *ConnInfo {
	return &ConnInfo{fs, options, mode, dataBlockSize, 1}
}

type ConnInfos struct {
	fmap sync.Map
}

func NewConnInfos() *ConnInfos {
	return &ConnInfos{
		sync.Map{},
	}
}

func (f *ConnInfos) IsExist(key string) bool {
	_, ok := f.fmap.Load(key)
	return ok
}

func (f *ConnInfos) Get(key string) *ConnInfo {
	if v, ok := f.fmap.Load(key); ok {
		if fv, fok := v.(*ConnInfo); fok {
			if fv.FileInfo != nil {
				return fv
			}
		}
	}
	return nil
}

func (f *ConnInfos) Set(key string, fs *ConnInfo) {
	if fs == nil || fs.FileInfo == nil {
		return
	}
	f.fmap.Store(key, fs)
}

func (f *ConnInfos) Del(key string) bool {
	defer f.fmap.Delete(key)
	if v, ok := f.fmap.Load(key); ok {
		if fv, fok := v.(*ConnInfo); fok {
			if fv.FileInfo != nil {
				fv.FileInfo.Close()
			}
		}
		return true
	}
	return false
}
