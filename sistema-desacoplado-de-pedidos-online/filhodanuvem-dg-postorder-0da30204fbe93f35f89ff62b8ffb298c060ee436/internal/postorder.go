package internal

type PostOrder struct {
	rules []struct {
		rule   Rule
		action Action
	}
}

func (p *PostOrder) AddRule(r Rule, a Action) {
	p.rules = append(p.rules, struct {
		rule   Rule
		action Action
	}{r, a})
}

func (p *PostOrder) Execute(order Order) Order {
	for _, config := range p.rules {
		if config.rule.Satisfy(order) {
			order = config.action.Execute(order)
		}
	}
	return order
}
