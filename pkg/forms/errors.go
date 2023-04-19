package forms

type errors map[string][]string

func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) []string {
	if es, ok := e[field]; ok {
		return es
	}
	return nil
}
