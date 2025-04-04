package internal

type Action interface {
	Execute(order Order) Order
}

type AddShippingLabel struct {
	Label string
}

func (f AddShippingLabel) Execute(order Order) Order {
	order.ShippingLabels = append(order.ShippingLabels, f.Label)
	return order
}

type AddDiscountPercent struct {
	Value int
}

func (f AddDiscountPercent) Execute(order Order) Order {
	order.Payment.Value = order.Payment.Value - (order.Payment.Value * f.Value / 100)
	return order
}
