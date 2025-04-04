package internal

type Rule interface {
	Satisfy(order Order) bool
}

type MinValue struct {
	Value int
}

func (f MinValue) Satisfy(order Order) bool {
	return order.Payment.Value >= f.Value
}

type ExpectedCategory struct {
	Category string
}

func (f ExpectedCategory) Satisfy(order Order) bool {
	return order.Product.Category == f.Category
}

type ExpectedPaymentMethod struct {
	PaymentMethod string
}

func (f ExpectedPaymentMethod) Satisfy(order Order) bool {
	return order.Payment.Method == f.PaymentMethod
}
