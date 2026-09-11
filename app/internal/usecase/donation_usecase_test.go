package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"donor-service/app/internal/domain"
	"donor-service/app/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ============================================================
// MOCK DONATION REPOSITORY
// ============================================================

type MockDonationRepository struct {
	mock.Mock
}

func (m *MockDonationRepository) Create(
	ctx context.Context,
	donation *entity.Donation,
) error {
	args := m.Called(ctx, donation)
	return args.Error(0)
}

func (m *MockDonationRepository) FindOne(
	ctx context.Context,
	donationID uuid.UUID,
) (*entity.Donation, error) {
	args := m.Called(ctx, donationID)

	var donation *entity.Donation

	if args.Get(0) != nil {
		donation = args.Get(0).(*entity.Donation)
	}

	return donation, args.Error(1)
}

func (m *MockDonationRepository) Completed(
	ctx context.Context,
	donationID uuid.UUID,
	confirmedBy uuid.UUID,
	donationDate time.Time,
) error {
	args := m.Called(
		ctx,
		donationID,
		confirmedBy,
		donationDate,
	)

	return args.Error(0)
}

func (m *MockDonationRepository) GetByRequesterID(
	ctx context.Context,
	requesterID uuid.UUID,
) ([]entity.Donation, error) {
	args := m.Called(
		ctx,
		requesterID,
	)
	var donations []entity.Donation
	return donations, args.Error(0)
}

// ============================================================
// MOCK DONOR MATCH REPOSITORY
// ============================================================

type MockDonorMatchRepository struct {
	mock.Mock
}

func (m *MockDonorMatchRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.DonorMatch, error) {
	args := m.Called(ctx, id)

	var donorMatch *entity.DonorMatch

	if args.Get(0) != nil {
		donorMatch = args.Get(0).(*entity.DonorMatch)
	}

	return donorMatch, args.Error(1)
}

func (m *MockDonorMatchRepository) GetByBloodRequestAndDonor(
	ctx context.Context,
	bloodRequestID uuid.UUID,
	donorID uuid.UUID,
) (*entity.DonorMatch, error) {
	args := m.Called(ctx, bloodRequestID, donorID)

	var donorMatch *entity.DonorMatch

	if args.Get(0) != nil {
		donorMatch = args.Get(0).(*entity.DonorMatch)
	}

	return donorMatch, args.Error(1)
}

func (m *MockDonorMatchRepository) Create(
	ctx context.Context,
	donor *entity.DonorMatch,
) error {
	args := m.Called(ctx, donor)
	return args.Error(0)
}

func (m *MockDonorMatchRepository) CountAccepted(
	ctx context.Context,
	bloodRequestID uuid.UUID,
) (int64, error) {
	args := m.Called(ctx, bloodRequestID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockDonorMatchRepository) UpdateDistance(
	ctx context.Context,
	id uuid.UUID,
	distanceKM float64,
) error {
	args := m.Called(ctx, id, distanceKM)
	return args.Error(0)
}

func (m *MockDonorMatchRepository) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockDonorMatchRepository) GetByDonorID(
	ctx context.Context,
	requesterID uuid.UUID,
) ([]entity.DonorMatch, error) {
	args := m.Called(ctx, requesterID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.DonorMatch), args.Error(1)
}

// ============================================================
// MOCK BLOOD REQUEST REPOSITORY
// ============================================================

type MockBloodRequestRepository struct {
	mock.Mock
}

func (m *MockBloodRequestRepository) Create(
	ctx context.Context,
	req *entity.BloodRequest,
) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockBloodRequestRepository) FindAll(
	ctx context.Context,
	userID uuid.UUID,
) ([]entity.BloodRequest, error) {
	args := m.Called(ctx, userID)

	var bloodRequests []entity.BloodRequest

	if args.Get(0) != nil {
		bloodRequests = args.Get(0).([]entity.BloodRequest)
	}

	return bloodRequests, args.Error(1)
}

func (m *MockBloodRequestRepository) FindOne(
	ctx context.Context,
	userID uuid.UUID,
	bloodRequestID uuid.UUID,
) (*entity.BloodRequest, error) {
	args := m.Called(ctx, userID, bloodRequestID)

	var bloodRequest *entity.BloodRequest

	if args.Get(0) != nil {
		bloodRequest = args.Get(0).(*entity.BloodRequest)
	}

	return bloodRequest, args.Error(1)
}

func (m *MockBloodRequestRepository) Cancel(
	ctx context.Context,
	userID uuid.UUID,
	bloodRequestID uuid.UUID,
) error {
	args := m.Called(ctx, userID, bloodRequestID)
	return args.Error(0)
}

func (m *MockBloodRequestRepository) GetMatches(
	ctx context.Context,
	bloodRequestID uuid.UUID,
) ([]entity.DonorMatch, error) {
	args := m.Called(ctx, bloodRequestID)

	var matches []entity.DonorMatch

	if args.Get(0) != nil {
		matches = args.Get(0).([]entity.DonorMatch)
	}

	return matches, args.Error(1)
}

func (m *MockBloodRequestRepository) GetById(
	ctx context.Context,
	bloodRequestID uuid.UUID,
) (*entity.BloodRequest, error) {
	args := m.Called(ctx, bloodRequestID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.BloodRequest), args.Error(1)
}

func (m *MockBloodRequestRepository) CreateMatches(
	ctx context.Context,
	bloodRequestID uuid.UUID,
	donors []entity.DonorProfile,
) error {
	args := m.Called(ctx, bloodRequestID, donors)
	return args.Error(0)
}

func (m *MockBloodRequestRepository) GetPendingReqs(
	ctx context.Context,
) ([]entity.BloodRequest, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.BloodRequest), args.Error(1)
}
func (m *MockBloodRequestRepository) Complete(ctx context.Context, userID uuid.UUID, bloodRequestID uuid.UUID) error {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return args.Error(1)
	}

	return args.Error(1)
}

// ============================================================
// CREATE DONATION TEST
// ============================================================

func TestDonationUsecase_Create_InvalidInput(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	t.Run("nil request", func(t *testing.T) {
		result, err := uc.Create(
			uuid.New(),
			nil,
		)

		assert.Nil(t, result)
		assert.ErrorIs(
			t,
			err,
			domain.ErrInvalidDonationInput,
		)
	})

	t.Run("empty blood request id", func(t *testing.T) {
		req := &domain.DonationReq{
			BloodRequestID: uuid.Nil,
			DonorMatchID:   uuid.New(),
		}

		result, err := uc.Create(
			uuid.New(),
			req,
		)

		assert.Nil(t, result)
		assert.ErrorIs(
			t,
			err,
			domain.ErrInvalidDonationInput,
		)
	})

	t.Run("empty donor match id", func(t *testing.T) {
		req := &domain.DonationReq{
			BloodRequestID: uuid.New(),
			DonorMatchID:   uuid.Nil,
		}

		result, err := uc.Create(
			uuid.New(),
			req,
		)

		assert.Nil(t, result)
		assert.ErrorIs(
			t,
			err,
			domain.ErrInvalidDonationInput,
		)
	})
}

func TestDonationUsecase_Create_BloodRequestError(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()
	bloodRequestID := uuid.New()
	donorMatchID := uuid.New()

	expectedErr := errors.New("blood request error")

	bloodReqRepo.
		On(
			"FindOne",
			mock.Anything,
			requesterID,
			bloodRequestID,
		).
		Return(nil, expectedErr).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.Create(
		requesterID,
		&domain.DonationReq{
			BloodRequestID: bloodRequestID,
			DonorMatchID:   donorMatchID,
		},
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	bloodReqRepo.AssertExpectations(t)

	donorMatchRepo.AssertNotCalled(
		t,
		"GetByID",
	)

	donationRepo.AssertNotCalled(
		t,
		"Create",
	)
}

func TestDonationUsecase_Create_DonorMatchError(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()
	bloodRequestID := uuid.New()
	donorMatchID := uuid.New()

	expectedErr := errors.New("donor match error")

	bloodReqRepo.
		On(
			"FindOne",
			mock.Anything,
			requesterID,
			bloodRequestID,
		).
		Return(&entity.BloodRequest{}, nil).
		Once()

	donorMatchRepo.
		On(
			"GetByID",
			mock.Anything,
			donorMatchID,
		).
		Return(nil, expectedErr).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.Create(
		requesterID,
		&domain.DonationReq{
			BloodRequestID: bloodRequestID,
			DonorMatchID:   donorMatchID,
		},
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	bloodReqRepo.AssertExpectations(t)
	donorMatchRepo.AssertExpectations(t)

	donationRepo.AssertNotCalled(
		t,
		"Create",
	)
}

func TestDonationUsecase_Create_BloodRequestMismatch(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()
	bloodRequestID := uuid.New()
	donorMatchID := uuid.New()

	differentBloodRequestID := uuid.New()

	bloodReqRepo.
		On(
			"FindOne",
			mock.Anything,
			requesterID,
			bloodRequestID,
		).
		Return(&entity.BloodRequest{}, nil).
		Once()

	donorMatchRepo.
		On(
			"GetByID",
			mock.Anything,
			donorMatchID,
		).
		Return(
			&entity.DonorMatch{
				ID:             donorMatchID,
				BloodRequestID: differentBloodRequestID,
				DonorID:        uuid.New(),
				Status:         "accepted",
			},
			nil,
		).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.Create(
		requesterID,
		&domain.DonationReq{
			BloodRequestID: bloodRequestID,
			DonorMatchID:   donorMatchID,
		},
	)

	assert.Nil(t, result)

	assert.ErrorIs(
		t,
		err,
		domain.ErrInvalidDonationInput,
	)

	donationRepo.AssertNotCalled(
		t,
		"Create",
	)
}

func TestDonationUsecase_Create_MatchNotAccepted(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()
	bloodRequestID := uuid.New()
	donorMatchID := uuid.New()

	bloodReqRepo.
		On(
			"FindOne",
			mock.Anything,
			requesterID,
			bloodRequestID,
		).
		Return(&entity.BloodRequest{}, nil).
		Once()

	donorMatchRepo.
		On(
			"GetByID",
			mock.Anything,
			donorMatchID,
		).
		Return(
			&entity.DonorMatch{
				ID:             donorMatchID,
				BloodRequestID: bloodRequestID,
				DonorID:        uuid.New(),
				Status:         "invited",
			},
			nil,
		).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.Create(
		requesterID,
		&domain.DonationReq{
			BloodRequestID: bloodRequestID,
			DonorMatchID:   donorMatchID,
		},
	)

	assert.Nil(t, result)

	assert.ErrorIs(
		t,
		err,
		domain.ErrDonorMatchNotAccepted,
	)

	donationRepo.AssertNotCalled(
		t,
		"Create",
	)
}

func TestDonationUsecase_Create_RepositoryError(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()
	bloodRequestID := uuid.New()
	donorMatchID := uuid.New()
	donorID := uuid.New()

	expectedErr := errors.New("failed to create donation")

	bloodReqRepo.
		On(
			"FindOne",
			mock.Anything,
			requesterID,
			bloodRequestID,
		).
		Return(&entity.BloodRequest{}, nil).
		Once()

	donorMatchRepo.
		On(
			"GetByID",
			mock.Anything,
			donorMatchID,
		).
		Return(
			&entity.DonorMatch{
				ID:             donorMatchID,
				BloodRequestID: bloodRequestID,
				DonorID:        donorID,
				Status:         "accepted",
			},
			nil,
		).
		Once()

	donationRepo.
		On(
			"Create",
			mock.Anything,
			mock.AnythingOfType("*entity.Donation"),
		).
		Return(expectedErr).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.Create(
		requesterID,
		&domain.DonationReq{
			BloodRequestID: bloodRequestID,
			DonorMatchID:   donorMatchID,
		},
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	donationRepo.AssertExpectations(t)
}

func TestDonationUsecase_Create_Success(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()
	bloodRequestID := uuid.New()
	donorMatchID := uuid.New()
	donorID := uuid.New()

	bloodReqRepo.
		On(
			"FindOne",
			mock.Anything,
			requesterID,
			bloodRequestID,
		).
		Return(&entity.BloodRequest{}, nil).
		Once()

	donorMatchRepo.
		On(
			"GetByID",
			mock.Anything,
			donorMatchID,
		).
		Return(
			&entity.DonorMatch{
				ID:             donorMatchID,
				BloodRequestID: bloodRequestID,
				DonorID:        donorID,
				Status:         "accepted",
			},
			nil,
		).
		Once()

	donationRepo.
		On(
			"Create",
			mock.Anything,
			mock.MatchedBy(
				func(donation *entity.Donation) bool {
					return donation.ID != uuid.Nil &&
						donation.BloodRequestID == bloodRequestID &&
						donation.DonorMatchID == donorMatchID &&
						donation.DonorID == donorID &&
						donation.Status == "pending"
				},
			),
		).
		Return(nil).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.Create(
		requesterID,
		&domain.DonationReq{
			BloodRequestID: bloodRequestID,
			DonorMatchID:   donorMatchID,
		},
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.NotEqual(t, uuid.Nil, result.ID)
	assert.Equal(t, bloodRequestID, result.BloodRequestID)
	assert.Equal(t, donorMatchID, result.DonorMatchID)
	assert.Equal(t, donorID, result.DonorID)
	assert.Equal(t, "pending", result.Status)

	bloodReqRepo.AssertExpectations(t)
	donorMatchRepo.AssertExpectations(t)
	donationRepo.AssertExpectations(t)
}

// ============================================================
// COMPLETED DONATION TEST
// ============================================================

func TestDonationUsecase_Completed_DonationNotFound(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()
	donationID := uuid.New()

	expectedErr := domain.ErrDonationNotFound

	donationRepo.
		On(
			"FindOne",
			mock.Anything,
			donationID,
		).
		Return(nil, expectedErr).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.Completed(
		requesterID,
		donationID,
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	donationRepo.AssertExpectations(t)

	bloodReqRepo.AssertNotCalled(
		t,
		"FindOne",
	)
}

func TestDonationUsecase_Completed_BloodRequestError(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()
	donationID := uuid.New()
	bloodRequestID := uuid.New()

	expectedErr := errors.New("blood request error")

	donation := &entity.Donation{
		ID:             donationID,
		BloodRequestID: bloodRequestID,
		Status:         "pending",
	}

	donationRepo.
		On(
			"FindOne",
			mock.Anything,
			donationID,
		).
		Return(donation, nil).
		Once()

	bloodReqRepo.
		On(
			"FindOne",
			mock.Anything,
			requesterID,
			bloodRequestID,
		).
		Return(nil, expectedErr).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.Completed(
		requesterID,
		donationID,
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	donationRepo.AssertNotCalled(
		t,
		"Completed",
	)
}

func TestDonationUsecase_Completed_InvalidStatus(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()
	donationID := uuid.New()
	bloodRequestID := uuid.New()

	donation := &entity.Donation{
		ID:             donationID,
		BloodRequestID: bloodRequestID,
		Status:         "completed",
	}

	donationRepo.
		On(
			"FindOne",
			mock.Anything,
			donationID,
		).
		Return(donation, nil).
		Once()

	bloodReqRepo.
		On(
			"FindOne",
			mock.Anything,
			requesterID,
			bloodRequestID,
		).
		Return(&entity.BloodRequest{}, nil).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.Completed(
		requesterID,
		donationID,
	)

	assert.Nil(t, result)

	assert.ErrorIs(
		t,
		err,
		domain.ErrInvalidDonationStatus,
	)

	donationRepo.AssertNotCalled(
		t,
		"Completed",
	)
}

func TestDonationUsecase_Completed_RepositoryError(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()
	donationID := uuid.New()
	bloodRequestID := uuid.New()

	expectedErr := errors.New("failed to complete donation")

	donation := &entity.Donation{
		ID:             donationID,
		BloodRequestID: bloodRequestID,
		Status:         "pending",
	}

	donationRepo.
		On(
			"FindOne",
			mock.Anything,
			donationID,
		).
		Return(donation, nil).
		Once()

	bloodReqRepo.
		On(
			"FindOne",
			mock.Anything,
			requesterID,
			bloodRequestID,
		).
		Return(&entity.BloodRequest{}, nil).
		Once()

	donationRepo.
		On(
			"Completed",
			mock.Anything,
			donationID,
			requesterID,
			mock.AnythingOfType("time.Time"),
		).
		Return(expectedErr).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.Completed(
		requesterID,
		donationID,
	)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	donationRepo.AssertExpectations(t)
}

func TestDonationUsecase_Completed_Success(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()
	donationID := uuid.New()
	bloodRequestID := uuid.New()
	donorID := uuid.New()

	donation := &entity.Donation{
		ID:             donationID,
		BloodRequestID: bloodRequestID,
		DonorID:        donorID,
		Status:         "pending",
	}

	donationRepo.
		On(
			"FindOne",
			mock.Anything,
			donationID,
		).
		Return(donation, nil).
		Once()

	bloodReqRepo.
		On(
			"FindOne",
			mock.Anything,
			requesterID,
			bloodRequestID,
		).
		Return(&entity.BloodRequest{}, nil).
		Once()

	donationRepo.
		On(
			"Completed",
			mock.Anything,
			donationID,
			requesterID,
			mock.AnythingOfType("time.Time"),
		).
		Return(nil).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.Completed(
		requesterID,
		donationID,
	)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "completed", result.Status)

	require.NotNil(t, result.DonationDate)
	require.NotNil(t, result.ConfirmedBy)

	assert.Equal(
		t,
		requesterID,
		*result.ConfirmedBy,
	)

	assert.WithinDuration(
		t,
		time.Now(),
		*result.DonationDate,
		time.Second,
	)

	donationRepo.AssertExpectations(t)
	bloodReqRepo.AssertExpectations(t)
}

func (m *MockDonationRepository) FindAllByRequesterID(
	ctx context.Context,
	requesterID uuid.UUID,
) ([]entity.Donation, error) {
	args := m.Called(ctx, requesterID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.Donation), args.Error(1)
}

func TestDonationUsecase_FindAll_Success(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()

	expected := []entity.Donation{
		{
			ID:     uuid.New(),
			Status: "pending",
		},
		{
			ID:     uuid.New(),
			Status: "completed",
		},
	}

	donationRepo.
		On(
			"FindAllByRequesterID",
			mock.Anything,
			requesterID,
		).
		Return(expected, nil).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.FindAll(requesterID)

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	donationRepo.AssertExpectations(t)
}

func TestDonationUsecase_FindAll_RepositoryError(t *testing.T) {
	donationRepo := new(MockDonationRepository)
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)

	requesterID := uuid.New()
	expectedErr := errors.New("failed to get donations")

	donationRepo.
		On(
			"FindAllByRequesterID",
			mock.Anything,
			requesterID,
		).
		Return(nil, expectedErr).
		Once()

	uc := NewDonationUsecase(
		donationRepo,
		donorMatchRepo,
		bloodReqRepo,
	)

	result, err := uc.FindAll(requesterID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)
}
