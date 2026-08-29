package ficc_fut


func (m *MarginExposureManager) PrepareForRecover(positionRecord *PositionRecord) {
	
	delete(m.orderListsMap, positionRecord.Key)

	for _, symbolQtyAggregratorItem := range m.symbolQtyAggregrator.entry {
		for key, p := range symbolQtyAggregratorItem.entry {
			if p.Key == positionRecord.Key {
				delete(symbolQtyAggregratorItem.entry, key)
			}
		}
		symbolQtyAggregratorItem.Refresh()
	}

	for _, accountAmountAggregratorItem := range m.accountAmountAggregrator.entry {
		for key, p := range accountAmountAggregratorItem.entry {
			if p.Key == positionRecord.Key {
				delete(accountAmountAggregratorItem.entry, key)
			}
		}
		accountAmountAggregratorItem.Refresh()
	}
}
