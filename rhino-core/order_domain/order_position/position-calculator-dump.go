package order_position

func (pc *PositionCalculator) Dump() map[string]float64 {
	pc.lock.Lock()
	defer pc.lock.Unlock()

	m := make(map[string]float64)

	for k, v := range pc.positionMap {
		_, quoto := v.GetQuota()
		m[k] = quoto
	}
	return m
}
