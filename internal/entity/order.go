package entity

import (
	"errors"
	"time"
)

type Order struct {
	ID         int64
	Price      float64
	Tax        float64
	FinalPrice float64
	CreatedAt  *time.Time
}

func NewOrder(price float64, tax float64) (*Order, error) {
	order := &Order{
		Price: price,
		Tax:   tax,
	}

	err := order.Validate()
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *Order) Validate() error {
	if o.Price <= 0 {
		return errors.New("o preço deve ser meior que zero")
	}
	if o.Tax <= 0 {
		return errors.New("o imposto deve ser maior que zero")
	}

	return nil
}

func (o *Order) CalculateFinalPrice() error {
	o.FinalPrice = o.Price + o.Tax

	err := o.Validate()
	if err != nil {
		return err
	}

	return nil
}
