package logic

import (
	"context"
	"time"

	"github.com/market-place/domain"
	"github.com/market-place/usecase/adapter"
	"github.com/market-place/usecase/repository"
	"github.com/market-place/usecase/usecase_error"
	"golang.org/x/sync/errgroup"
)

type TBuyerUsecase interface {
	Create(ctx context.Context, input adapter.TBuyerCreateInput) (domain.TBuyer, error)
	GetByID(ctx context.Context, transactionID string) (domain.TBuyer, error)
	Fetch(ctx context.Context, cursor string, num int64, search domain.TBuyerSearchOptions) ([]domain.TBuyer, error)
	UploadTransferPhoto(ctx context.Context, fileName string, transactionID string) (domain.TBuyer, error)
	AcceptTransaction(ctx context.Context, transactionID string) (domain.TBuyer, error)
	RejectTransaction(ctx context.Context, input adapter.TbuyerRejectInput, transactionID string) (domain.TBuyer, error)
}

type tbuyerUsecase struct {
	tbuyerRepo     repository.TBuyerRepository
	orderRepo      repository.OrderRepository
	contextTimeout time.Duration
}

func NewTBuyerUsecase(
	tbuyerRepo repository.TBuyerRepository,
	orderRepo repository.OrderRepository,
	contextTimeOut time.Duration,
) TBuyerUsecase {
	return &tbuyerUsecase{
		tbuyerRepo:     tbuyerRepo,
		orderRepo:      orderRepo,
		contextTimeout: contextTimeOut,
	}
}

func (t *tbuyerUsecase) Create(ctx context.Context, input adapter.TBuyerCreateInput) (domain.TBuyer, error) {
	var tbuyer domain.TBuyer

	return tbuyer, nil
}

func (t *tbuyerUsecase) GetByID(ctx context.Context, transactionID string) (domain.TBuyer, error) {
	var tbuyer domain.TBuyer

	return tbuyer, nil
}

func (t *tbuyerUsecase) Fetch(ctx context.Context, cursor string, num int64, search domain.TBuyerSearchOptions) ([]domain.TBuyer, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()
	return t.tbuyerRepo.Fetch(ctx, cursor, num, search)
}

func (t *tbuyerUsecase) UploadTransferPhoto(ctx context.Context, fileName string, transactionID string) (domain.TBuyer, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	tbuyer, err := t.tbuyerRepo.GetByID(ctx, transactionID)
	if err != nil {
		return tbuyer, err
	}

	tbuyer.TransferPhoto = fileName
	tbuyer.PaymentStatus = domain.PEMBAYARAN_SEDANG_DIVERIFIKASI
	return t.tbuyerRepo.UpdateOne(ctx, tbuyer)
}

func (t *tbuyerUsecase) AcceptTransaction(ctx context.Context, transactionID string) (domain.TBuyer, error) {
	credential := ctx.Value("credential")
	if credential == nil {
		return domain.TBuyer{}, usecase_error.ErrNotAuthorization
	}
	userInfo := credential.(domain.Credential)

	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()
	tbuyer, err := t.tbuyerRepo.GetByID(ctx, transactionID)
	if err != nil {
		return tbuyer, err
	}

	tbuyer.PaymentStatus = domain.PEMBAYARAN_SUCCESS
	tbuyer.AdminID = userInfo.UserID

	search := domain.OrderSearchOptions{
		TransactionID: transactionID,
	}
	noCursor := ""
	noNum := int64(0)
	orders, err := t.orderRepo.Fetch(ctx, noCursor, noNum, search)
	if err != nil {
		return tbuyer, err
	}
	for _, order := range orders {
		order.StatusOrder = domain.STATUS_ORDER_SEDANG_DIPROSES
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() error {
			_, err := t.orderRepo.UpdateOne(ctx, order)
			if err != nil {
				return err
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			return tbuyer, err
		}
	}

	return t.tbuyerRepo.UpdateOne(ctx, tbuyer)
}

func (t *tbuyerUsecase) RejectTransaction(ctx context.Context, input adapter.TbuyerRejectInput, transactionID string) (domain.TBuyer, error) {
	var tbuyer domain.TBuyer
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()
	tbuyer, err := t.tbuyerRepo.GetByID(ctx, transactionID)
	if err != nil {
		return tbuyer, err
	}
	tbuyer.Message = input.Message
	tbuyer.PaymentStatus = domain.PEMBAYARAN_GAGAL

	return t.tbuyerRepo.UpdateOne(ctx, tbuyer)
}
