package ficc_fut


func (a *FiccFutOrderPositionAdapter) PrepareForRecover(_positionRecord interface{}) {

	positionRecord, ok := _positionRecord.(*PositionRecord)
	if !ok {
		return
	}

	a.PrepareForRecover(positionRecord)
}

func (a *FiccFutOrderPositionAdapter) UpdatePositionBaseDynamically()(dynamicallyUpdate bool) {
	return true
}