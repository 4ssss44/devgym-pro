package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_free_shipping(t *testing.T) {
	p := PostOrder{}
	p.AddRule(MinValue{1000}, AddShippingLabel{"frete-gratis"})
	order := Order{
		Payment: Payment{
			Method: "Boleto",
			Value:  1001,
		},
		Product: Product{
			Category: "eletrodoméstico",
		},
	}
	newOrder := p.Execute(order)

	labels := []string{"frete-gratis"}
	assert.Equal(t, newOrder.ShippingLabels, labels)

	order.ShippingLabels = labels
	assert.Equal(t, order, newOrder)
}

func Test_no_rule_same_order(t *testing.T) {
	p := PostOrder{} // no rules
	order := Order{
		Payment: Payment{
			Method: "Boleto",
		},
		Product: Product{
			Category: "eletrodoméstico",
		},
	}
	newOrder := p.Execute(order)

	assert.Equal(t, order, newOrder)
}

func Test_MultipleLabels(t *testing.T) {
	p := PostOrder{}
	p.AddRule(MinValue{1000}, AddShippingLabel{"frete-gratis"})
	p.AddRule(ExpectedCategory{"eletrodoméstico"}, AddShippingLabel{"frágil"})
	order := Order{
		Payment: Payment{
			Method: "Boleto",
			Value:  1001,
		},
		Product: Product{
			Category: "eletrodoméstico",
		},
	}
	newOrder := p.Execute(order)

	labels := []string{"frete-gratis", "frágil"}
	assert.Equal(t, newOrder.ShippingLabels, labels)

	order.ShippingLabels = labels
	assert.Equal(t, order, newOrder)
}

func Test_update_order_total_value(t *testing.T) {
	p := PostOrder{}
	p.AddRule(ExpectedPaymentMethod{"Boleto"}, AddDiscountPercent{10})
	order := Order{
		Payment: Payment{
			Method: "Boleto",
			Value:  100,
		},
		Product: Product{
			Category: "eletrodoméstico",
		},
	}
	newOrder := p.Execute(order)

	assert.Equal(t, newOrder.Payment.Value, 90)

	newOrder.Payment.Value = 100
	assert.Equal(t, order, newOrder)
}
