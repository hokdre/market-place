package logic

import (
	"context"
	"sync"
	"time"

	guuid "github.com/google/uuid"
	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/helper"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
)

type OrderUsecase interface {
	Create(ctx context.Context, input adapter.OrderCreateInput) ([]domain.Order, error)
	GetByID(ctx context.Context, orderID string) (domain.Order, error)
	Fetch(ctx context.Context, cursor string, num int64, search adapter.OrderSearchOptions) ([]domain.Order, error)
	RejectOrder(ctx context.Context, orderId string) (domain.Order, error)
	FinishOrder(ctx context.Context, orderID string) (domain.Order, error)
	InputResiNumber(ctx context.Context, updateInput adapter.OrderResiInput, orderID string) (domain.Order, error)
	UploadShippingPhoto(ctx context.Context, fileName string, orderID string) (domain.Order, error)
	AjukanPaketSampai(ctx context.Context, orderID string) (domain.Order, error)
	EstimasiPendapatan(ctx context.Context, startDay string, endDay string) ([]map[string]interface{}, error)
	OrderSummary(ctx context.Context, startDay string, endDay string) (map[string]int64, error)
}

type orderUsecase struct {
	orderRepo      repository.OrderRepository
	productRepo    repository.ProductRepository
	customerRepo   repository.CustomerRepository
	merchantRepo   repository.MerchantRepository
	cartRepo       repository.CartRepository
	tBuyerRepo     repository.TBuyerRepository
	contextTimeout time.Duration
}

func NewOrderUsecase(
	orderRepo repository.OrderRepository,
	productRepo repository.ProductRepository,
	customerRepo repository.CustomerRepository,
	merchantRepo repository.MerchantRepository,
	cartRepo repository.CartRepository,
	tBuyerRepo repository.TBuyerRepository,
	contextTimeout time.Duration,
) OrderUsecase {
	return &orderUsecase{
		orderRepo:      orderRepo,
		productRepo:    productRepo,
		customerRepo:   customerRepo,
		merchantRepo:   merchantRepo,
		cartRepo:       cartRepo,
		tBuyerRepo:     tBuyerRepo,
		contextTimeout: contextTimeout,
	}
}

func (c *orderUsecase) validate(value interface{}) error {
	if entityErr := helper.NewValidationEntity().Validate(value); entityErr != nil {
		return entityErr
	}

	return nil
}

func (o *orderUsecase) Create(ctx context.Context, input adapter.OrderCreateInput) ([]domain.Order, error) {
	credential := ctx.Value("credential")
	if credential == nil {
		return nil, usecase_error.ErrNotAuthorization
	}
	userInfo := credential.(domain.Credential)
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()

	//fetch user and cart
	var customer domain.Customer
	var cart domain.Cart
	var errUserCart error
	var wgUserandCart sync.WaitGroup
	wgUserandCart.Add(2)
	go func() {
		defer wgUserandCart.Done()
		c, err := o.customerRepo.GetByID(ctx, userInfo.UserID)
		customer = c
		errUserCart = err
	}()
	go func() {
		defer wgUserandCart.Done()
		c, err := o.cartRepo.GetByID(ctx, userInfo.CartID)
		cart = c
		errUserCart = err
	}()
	wgUserandCart.Wait()
	if errUserCart != nil {
		return []domain.Order{}, errUserCart
	}

	// Invoice
	transaction := domain.TBuyer{}
	transaction.ID = guuid.New().String()
	transaction.CustomerID = customer.ID
	transaction.CreatedAt = time.Now().Truncate(time.Millisecond)
	transaction.UpdatedAt = time.Now().Truncate(time.Millisecond)
	transaction.PaymentStatus = domain.PEMBAYARAN_MENUNGGU_VERIFIKASI

	// convert input data to order object
	pOrder := func(ctx context.Context, input adapter.OrderCreateInput) chan domain.Order {
		chOrder := make(chan domain.Order)
		go func() {
			defer close(chOrder)
			for _, orderData := range input.Orders {
				select {
				case <-ctx.Done():
					return
				default:
					order := domain.Order{}
					order.ID = guuid.New().String()
					order.StatusOrder = domain.STATUS_ORDER_MENUNGGU_PEMBAYARAN
					order.Customer = customer.DenomalizationCustomer()
					order.Merchant = domain.DenormalizationMerchant{
						ID: orderData.MerchantID,
					}
					order.Shipping = domain.ShippingProvider{
						ID: orderData.ShippingID,
					}
					order.ServiceName = orderData.ServiceName
					order.TransactionsID = transaction.ID
					order.ReceiverName = orderData.ReceiverName
					order.ReceiverPhone = orderData.ReceiverPhone
					order.ReceiverAddress = domain.Address{
						ID:     guuid.New().String(),
						City:   orderData.ReceiverAddress.City,
						Street: orderData.ReceiverAddress.Street,
						Number: orderData.ReceiverAddress.Number,
					}
					order.ShippingCost = orderData.ShippingCost
					transaction.TotalTransfer += order.ShippingCost

					order.CreatedAt = time.Now().Truncate(time.Millisecond)
					order.UpdatedAt = time.Now().Truncate(time.Millisecond)

					order.OrderItems = []domain.OrderItems{}
					for _, product := range orderData.Products {
						item := domain.OrderItems{}
						item.Product.ID = product.ProductID
						item.Quantity = product.Quantity
						item.BuyerNote = product.BuyerNote
						item.Colors = product.Colors
						item.Sizes = product.Sizes

						transaction.TotalTransfer += item.Price * item.Quantity
						order.OrderItems = append(order.OrderItems, item)
					}
					chOrder <- order
				}
			}
		}()
		return chOrder
	}

	//populate order object merchant data
	type ResultOrder struct {
		Order domain.Order
		Err   error
	}
	stageFetchMerchant := func(ctx context.Context, chOrders chan domain.Order) chan ResultOrder {
		chMerchant := make(chan ResultOrder)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			for order := range chOrders {
				wg.Add(1)
				go func(order domain.Order) {
					defer wg.Done()
					select {
					case <-ctx.Done():
						return
					default:
						merchant, err := o.merchantRepo.GetByID(ctx, order.Merchant.ID)
						//check shipping if merchant provide that
						isMerchantProvideShippingProvider := false
						for _, shipping := range merchant.Shippings {
							if shipping.ID == order.Shipping.ID {
								isMerchantProvideShippingProvider = true
								order.Shipping = shipping
							}
						}
						if !isMerchantProvideShippingProvider {
							err = usecase_error.ErrBadEntityInput{
								usecase_error.ErrEntityField{
									Field:   "ShippingID",
									Message: "Shipping is not provided by merchant",
								},
							}
						}

						select {
						case <-ctx.Done():
							return
						default:
							order.Merchant = merchant.DenomarlizationData()
							chMerchant <- ResultOrder{
								Order: order,
								Err:   err,
							}
						}
					}
				}(order)
			}
		}()

		go func() {
			wg.Wait()
			close(chMerchant)
		}()

		return chMerchant
	}

	//populate order object produk data
	type ResultProduct struct {
		Product domain.Product
		Err     error
	}
	stageFetchProduct := func(ctx context.Context, chMerchant chan ResultOrder) chan ResultOrder {
		chProducerProduct := make(chan ResultOrder)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			for resultMerchant := range chMerchant {
				wg.Add(1)
				go func(resultMerchant ResultOrder) {
					defer wg.Done()
					if resultMerchant.Err == nil {
						// process fetch products if fetch merchant before success
						items := resultMerchant.Order.OrderItems
						order := resultMerchant.Order

						resultProduct := make(chan ResultProduct, len(items))

						var wgFetchProduct sync.WaitGroup
						wgFetchProduct.Add(len(items))

						for _, item := range items {
							go func(item domain.OrderItems) {
								defer wgFetchProduct.Done()

								product, err := o.productRepo.GetByID(ctx, item.Product.ID)

								//validate merchant sell the product
								if err == nil {
									merchantID := order.Merchant.ID
									productMerchantID := product.Merchant.ID
									if merchantID != productMerchantID {
										err = usecase_error.ErrNotFound
									}
								}

								select {
								case <-ctx.Done():
									return
								default:
									resultProduct <- ResultProduct{
										Product: product,
										Err:     err,
									}
								}
							}(item)
						}

						wgFetchProduct.Wait()
						for i := 0; i < len(items); i++ {
							resProduct := <-resultProduct
							product := resProduct.Product
							err := resProduct.Err
							//pass the error
							if err != nil {
								select {
								case <-ctx.Done():
									return
								default:
									resultMerchant.Err = err
									chProducerProduct <- resultMerchant
									return
								}
							}

							//assingn product to order
							for i, item := range order.OrderItems {
								if item.Product.ID == product.ID {
									order.OrderItems[i].Product = product.DenormalizationData()
									break
								}
							}
						}
					}

					select {
					case <-ctx.Done():
						return
					default:
						chProducerProduct <- resultMerchant
					}
				}(resultMerchant)
			}
		}()

		go func() {
			wg.Wait()
			close(chProducerProduct)
		}()
		return chProducerProduct
	}

	//save order to database
	type ResultCreateOrder struct {
		Order domain.Order
		Err   error
	}
	stageCreateOrder := func(ctx context.Context, chProduct chan ResultOrder) chan ResultOrder {
		chCreateOrder := make(chan ResultOrder)
		var wgCreateOrder sync.WaitGroup
		wgCreateOrder.Add(len(input.Orders))
		go func() {
			for result := range chProduct {
				// skip procces if process before error
				if result.Err != nil {
					wgCreateOrder.Done()
					select {
					case <-ctx.Done():
						break
					default:
						chCreateOrder <- result
					}
				} else {
					for _, item := range result.Order.OrderItems {
						productID := item.Product.ID

						//remove oredered item from cart
						index := -1
						for i, item := range cart.Items {
							if item.Product.ID == productID {
								index = i
								break
							}
						}
						if index != -1 {
							newItems := []domain.Item{}
							copy(newItems, cart.Items)

							firstItem := 0
							lastItem := len(cart.Items) - 1
							if index == firstItem && index == lastItem {
								newItems = []domain.Item{}
							} else if index == firstItem {
								newItems = append(newItems, newItems[index+1:]...)
							} else if index == lastItem {
								newItems = append(newItems, newItems[:index]...)
							} else {
								newItems = append(newItems[:index], newItems[index+1:]...)
							}

							cart.Items = newItems
						}
					}

					go func(result ResultOrder) {
						defer wgCreateOrder.Done()
						order, err := o.orderRepo.Create(ctx, result.Order)
						select {
						case <-ctx.Done():
							return
						default:
							chCreateOrder <- ResultOrder{
								Order: order,
								Err:   err,
							}
						}
					}(result)
				}

			}
		}()
		go func() {
			wgCreateOrder.Wait()
			close(chCreateOrder)
		}()
		return chCreateOrder
	}

	cOrder := pOrder(ctx, input)
	cMerchant := stageFetchMerchant(ctx, cOrder)
	cProduct := stageFetchProduct(ctx, cMerchant)
	cCreateOrder := stageCreateOrder(ctx, cProduct)
	orders := []domain.Order{}
	for result := range cCreateOrder {
		if result.Err != nil {
			cancel()
			return orders, result.Err
		}

		orders = append(orders, result.Order)
	}

	//update cart and save transaction
	var wgCartAndTransaction sync.WaitGroup
	var errCartAndTransaction error
	wgCartAndTransaction.Add(2)
	go func() {
		defer wgCartAndTransaction.Done()
		_, errCartAndTransaction = o.cartRepo.UpdateOne(ctx, cart)
	}()
	go func() {
		defer wgCartAndTransaction.Done()
		_, errCartAndTransaction = o.tBuyerRepo.Create(ctx, transaction)
	}()
	wgCartAndTransaction.Wait()
	if errCartAndTransaction != nil {
		return orders, errCartAndTransaction
	}

	return orders, nil
}

func (o *orderUsecase) GetByID(ctx context.Context, orderID string) (domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()

	return o.orderRepo.GetByID(ctx, orderID)
}

func (o *orderUsecase) Fetch(ctx context.Context, cursor string, num int64, options adapter.OrderSearchOptions) ([]domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()

	search := domain.OrderSearchOptions{
		CustomerID:    options.CustomerID,
		MerchantID:    options.MerchantID,
		TransactionID: options.TransactionID,
		Status:        options.Status,
	}

	return o.orderRepo.Fetch(ctx, cursor, num, search)
}

func (o *orderUsecase) RejectOrder(ctx context.Context, orderID string) (domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()
	order, err := o.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return order, err
	}
	order.StatusOrder = domain.STATUS_ORDER_DI_CANCEL
	if err := o.validate(order); err != nil {
		return order, err
	}

	return o.orderRepo.UpdateOne(ctx, order)
}

func (o *orderUsecase) FinishOrder(ctx context.Context, orderID string) (domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()
	order, err := o.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return order, err
	}
	order.StatusOrder = domain.STATUS_ORDER_SELESAI
	if err := o.validate(order); err != nil {
		return order, err
	}

	return o.orderRepo.UpdateOne(ctx, order)
}

func (o *orderUsecase) InputResiNumber(ctx context.Context, input adapter.OrderResiInput, orderID string) (domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()
	order, err := o.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return order, err
	}

	order.ResiNumber = input.ResiNumber
	order.StatusOrder = domain.STATUS_ORDER_SEDANG_DIKIRIM
	if err := o.validate(order); err != nil {
		return order, err
	}

	return o.orderRepo.UpdateOne(ctx, order)
}

func (o *orderUsecase) UploadShippingPhoto(ctx context.Context, fileName string, orderID string) (domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()
	order, err := o.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return order, err
	}

	order.ShippingPhoto = fileName
	if err := o.validate(order); err != nil {
		return order, err
	}

	return o.orderRepo.UpdateOne(ctx, order)
}

func (o *orderUsecase) AjukanPaketSampai(ctx context.Context, orderID string) (domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()
	order, err := o.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return order, err
	}

	order.Delivered = true

	return o.orderRepo.UpdateOne(ctx, order)
}

func (o *orderUsecase) EstimasiPendapatan(ctx context.Context, startDay string, endDay string) ([]map[string]interface{}, error) {
	credential := ctx.Value("credential")
	if credential == nil {
		return nil, usecase_error.ErrNotAuthorization
	}
	userInfo := credential.(domain.Credential)
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()

	return o.orderRepo.EstimasiPendapatan(ctx, userInfo.MerchantID, startDay, endDay)
}

func (o *orderUsecase) OrderSummary(ctx context.Context, startDay string, endDay string) (map[string]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, o.contextTimeout)
	defer cancel()
	credential := ctx.Value("credential")
	if credential == nil {
		return nil, usecase_error.ErrNotAuthorization
	}
	userInfo := credential.(domain.Credential)

	return o.orderRepo.OrderSummary(ctx, userInfo.MerchantID, startDay, endDay)
}
