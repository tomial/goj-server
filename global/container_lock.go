package global

import "sync"

var GoConLock *sync.RWMutex
var CConLock *sync.RWMutex

func init() {
	GoConLock = new(sync.RWMutex)
	CConLock = new(sync.RWMutex)
}
