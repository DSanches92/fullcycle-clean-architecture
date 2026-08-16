package usecase

import (
	"time"

	"github.com/DSanches92/fullcycle-clean-architecture/internal/entity"
)

type OrderInputDTO struct {
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         int64      `json:"id"`
	Price      float64    `json:"price"`
	Tax        float64    `json:"tax"`
	FinalPrice float64    `json:"final_price"`
	CreatedAt  *time.Time `json:"created_at"`
}

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewCreateOrderUseCase(orderRepository entity.OrderRepositoryInterface) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: orderRepository,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {
	order, err := entity.NewOrder(input.Price, input.Tax)
	if err != nil {
		return OrderOutputDTO{}, err
	}

	err = order.CalculateFinalPrice()
	if err != nil {
		return OrderOutputDTO{}, err
	}

	err = c.OrderRepository.Save(order)
	if err != nil {
		return OrderOutputDTO{}, err
	}

	return OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
		CreatedAt:  order.CreatedAt,
	}, nil
}