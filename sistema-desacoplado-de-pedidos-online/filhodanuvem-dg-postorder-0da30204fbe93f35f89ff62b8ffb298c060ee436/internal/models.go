package internal

const MethodCreditCard = "credit"
const MethodDebitCard = "debit"
const MethodPIX = "pix"

type Payment struct {
	Method string
	Value  int
}

type Product struct {
	Category string
}

type Order struct {
	Product        Product
	Payment        Payment
	ShippingLabels []string
}
