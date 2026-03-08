package misc

func ResetChan[T any](ch chan T) {
	for {
		select {
		case <-ch:
			// 清除可能存在的信号
		default:
			return
		}
	}
}
