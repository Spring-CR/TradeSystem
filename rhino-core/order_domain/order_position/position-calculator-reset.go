package order_position

func (pc *PositionCalculator) Reset() {
	pc.lock.Lock()
	defer pc.lock.Unlock()

	for k := range pc.positionMap {
		delete(pc.positionMap, k)
	}
}
