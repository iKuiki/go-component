package language

func SafeDeleteMapItem(m *map[string]interface{}, field string) bool {
	if _, ok := (*m)[field]; ok {
		delete(*m, field)
		return true
	}
	return false
}
