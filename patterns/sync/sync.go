// How does this make you feel?
package sync 

import (
	"fmt"
	"sync"
)

func WithLock() {
	var mtx sync.Mutex
	func() {
		defer func() {
			fmt.Println("unlocking")
			mtx.Unlock()
		}()
		mtx.Lock()
		fmt.Println("Hello, 世界")
	}()
	fmt.Println("done")
}

