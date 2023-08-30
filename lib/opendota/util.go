package opendota

func GetTotalField(field string, totals []TotalField) (TotalField, bool) {
	for i := range totals {
		if field == totals[i].FieldName {
			return totals[i], true
		}
	}
	return TotalField{}, false
}
