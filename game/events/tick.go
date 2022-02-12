package events

func SyncTick(gametick uint64) bool {
	return (gametick%60)%20 == 0
}
