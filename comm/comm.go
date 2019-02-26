package comm

import (
	"errors"
	"git.zhugefang.com/gocore/zgo.git/config"
	"sync"
)

//getCurrentLabel 着重判断输入的label与zgo engine 在内存中的用户态的label
func GetCurrentLabel(label []string, mu sync.RWMutex, cm map[string][]config.ConnDetail) (string, error) {
	mu.RLock()
	defer mu.RUnlock()

	lcl := len(cm)
	if lcl == 0 {
		return "", errors.New("invalid label in zgo engine or engine not start.")
	}
	if len(label) == 0 { //用户没有选择
		if lcl >= 1 {
			//自动返回默认的第一个
			l := ""
			for k, _ := range cm {
				l = k
				break
			}
			return l, nil
		} else {
			return "", errors.New("invalid label in zgo engine.")
		}
	} else if len(label) > 1 {
		return "", errors.New("you are choose must be one label or defalut zero.")
	} else {
		if _, ok := cm[label[0]]; ok {
			return label[0], nil
		} else {
			return "", errors.New("invalid label for u input.")
		}
	}
}
