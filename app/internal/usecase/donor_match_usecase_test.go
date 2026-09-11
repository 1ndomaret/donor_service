package usecase

import (
	"context"
	"errors"
	"testing"

	"donor-service/app/internal/domain"
	"donor-service/app/internal/dto"
	"donor-service/app/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

/*
Unit Test Donor Match Usecase

Yang dites:

1. Invite
   - Blood request gagal ditemukan / repository error
   - Requester bukan pemilik blood request
   - Donor match gagal ditemukan
   - Donor match sudah accepted / declined
   - Donor sudah pernah di-invite
   - Gagal update status menjadi invited
   - Berhasil invite donor

2. Accept
   - Donor match tidak ditemukan
   - Donor yang login bukan pemilik donor match
   - Status donor match bukan invited
   - Gagal mengambil blood request
   - Gagal menghitung accepted donor
   - Jumlah accepted donor sudah mencapai quantity
   - Gagal mengambil donor profile dari user-service
   - Gagal mengambil route dari Geoapify
   - Gagal update distance
   - Gagal update status menjadi accepted
   - Gagal auto-create donation
   - Berhasil accept:
     * quantity masih tersedia
     * route Geoapify berhasil
     * distance meter dikonversi ke KM
     * donor match menjadi accepted
     * donation otomatis dibuat dengan status pending

3. Decline
   - Donor match tidak ditemukan
   - Donor yang login bukan pemilik donor match
   - Status donor match bukan invited
   - Gagal update status menjadi declined
   - Berhasil decline

4.GetByDonorID
   - Repository error
   - Berhasil mengambil donor matches milik requester

Catatan:
MockDonationRepository, MockDonorMatchRepository, dan MockBloodRequestRepository
direuse dari donation_usecase_test.go karena berada di package usecase yang sama.
*/

// ============================================================
// MOCK USER SERVICE HTTP REPOSITORY
// ============================================================

type MockUserServiceHttpRepo struct {
	mock.Mock
}

func (m *MockUserServiceHttpRepo) SearchMatches(
	ctx context.Context,
	token string,
	req *domain.SearchMatchesRequest,
) ([]entity.DonorProfile, error) {
	args := m.Called(ctx, token, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]entity.DonorProfile), args.Error(1)
}

func (m *MockUserServiceHttpRepo) GetDonorProfile(
	ctx context.Context,
	donorID uuid.UUID,
) (*entity.DonorProfile, error) {
	args := m.Called(ctx, donorID)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.DonorProfile), args.Error(1)
}

// ============================================================
// MOCK GEOAPIFY REPOSITORY
// ============================================================

type MockGeoapifyRepository struct {
	mock.Mock
}

func (m *MockGeoapifyRepository) GetGeoapifyHospitals(
	ctx context.Context,
	city string,
) ([]domain.GeoapifyHospital, error) {
	args := m.Called(ctx, city)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]domain.GeoapifyHospital), args.Error(1)
}

func (m *MockGeoapifyRepository) GetGeoapifyRoute(
	ctx context.Context,
	req *dto.GeoapifyRoutingRequest,
) (*domain.GeoapifyRoute, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.GeoapifyRoute), args.Error(1)
}

// ============================================================
// TEST HELPER
// ============================================================

func newDonorMatchUsecaseForTest() (
	domain.DonorMatchUsecase,
	*MockDonorMatchRepository,
	*MockBloodRequestRepository,
	*MockUserServiceHttpRepo,
	*MockGeoapifyRepository,
	*MockDonationRepository,
) {
	donorMatchRepo := new(MockDonorMatchRepository)
	bloodReqRepo := new(MockBloodRequestRepository)
	userServiceRepo := new(MockUserServiceHttpRepo)
	geoapifyRepo := new(MockGeoapifyRepository)
	donationRepo := new(MockDonationRepository)

	uc := NewDonorMatchUsecase(
		donorMatchRepo,
		bloodReqRepo,
		userServiceRepo,
		geoapifyRepo,
		donationRepo,
	)

	return uc, donorMatchRepo, bloodReqRepo, userServiceRepo, geoapifyRepo, donationRepo
}

// ============================================================
// INVITE TEST
// ============================================================

// Test ketika blood request tidak berhasil diambil.
// Proses harus langsung berhenti dan donor match tidak boleh diakses.
func TestDonorMatchUsecase_Invite_BloodRequestError(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, _, _, _ := newDonorMatchUsecaseForTest()

	requesterID := uuid.New()
	bloodRequestID := uuid.New()
	donorID := uuid.New()
	expectedErr := errors.New("blood request error")

	bloodReqRepo.
		On("FindOne", mock.Anything, requesterID, bloodRequestID).
		Return(nil, expectedErr).
		Once()

	result, err := uc.Invite(requesterID, bloodRequestID, donorID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	bloodReqRepo.AssertExpectations(t)
	donorMatchRepo.AssertNotCalled(t, "GetByBloodRequestAndDonor")
}

// Test requester tidak boleh invite donor ke blood request milik requester lain.
func TestDonorMatchUsecase_Invite_Forbidden(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, _, _, _ := newDonorMatchUsecaseForTest()

	requesterID := uuid.New()
	otherRequesterID := uuid.New()
	bloodRequestID := uuid.New()
	donorID := uuid.New()

	bloodReqRepo.
		On("FindOne", mock.Anything, requesterID, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:          bloodRequestID,
			RequesterID: otherRequesterID,
		}, nil).
		Once()

	result, err := uc.Invite(requesterID, bloodRequestID, donorID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrForbidden)

	donorMatchRepo.AssertNotCalled(t, "GetByBloodRequestAndDonor")
}

// Test ketika donor match gagal diambil dari repository.
func TestDonorMatchUsecase_Invite_DonorMatchError(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, _, _, _ := newDonorMatchUsecaseForTest()

	requesterID := uuid.New()
	bloodRequestID := uuid.New()
	donorID := uuid.New()
	expectedErr := errors.New("donor match error")

	bloodReqRepo.
		On("FindOne", mock.Anything, requesterID, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:          bloodRequestID,
			RequesterID: requesterID,
		}, nil).
		Once()

	donorMatchRepo.
		On("GetByBloodRequestAndDonor", mock.Anything, bloodRequestID, donorID).
		Return(nil, expectedErr).
		Once()

	result, err := uc.Invite(requesterID, bloodRequestID, donorID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)
}

// Test donor match yang sudah accepted atau declined tidak boleh di-invite ulang.
func TestDonorMatchUsecase_Invite_InvalidStatus(t *testing.T) {
	statuses := []string{"accepted", "declined"}

	for _, status := range statuses {
		t.Run(status, func(t *testing.T) {
			uc, donorMatchRepo, bloodReqRepo, _, _, _ := newDonorMatchUsecaseForTest()

			requesterID := uuid.New()
			bloodRequestID := uuid.New()
			donorID := uuid.New()
			matchID := uuid.New()

			bloodReqRepo.
				On("FindOne", mock.Anything, requesterID, bloodRequestID).
				Return(&entity.BloodRequest{
					ID:          bloodRequestID,
					RequesterID: requesterID,
				}, nil).
				Once()

			donorMatchRepo.
				On("GetByBloodRequestAndDonor", mock.Anything, bloodRequestID, donorID).
				Return(&entity.DonorMatch{
					ID:             matchID,
					BloodRequestID: bloodRequestID,
					DonorID:        donorID,
					Status:         status,
				}, nil).
				Once()

			result, err := uc.Invite(requesterID, bloodRequestID, donorID)

			assert.Nil(t, result)
			assert.ErrorIs(t, err, domain.ErrInvalidMatchStatus)
			donorMatchRepo.AssertNotCalled(t, "UpdateStatus")
		})
	}
}

// Test donor yang sebelumnya sudah berstatus invited cukup dikembalikan.
// Repository tidak perlu update status lagi.
func TestDonorMatchUsecase_Invite_AlreadyInvited(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, _, _, _ := newDonorMatchUsecaseForTest()

	requesterID := uuid.New()
	bloodRequestID := uuid.New()
	donorID := uuid.New()
	matchID := uuid.New()

	existingMatch := &entity.DonorMatch{
		ID:             matchID,
		BloodRequestID: bloodRequestID,
		DonorID:        donorID,
		Status:         "invited",
	}

	bloodReqRepo.
		On("FindOne", mock.Anything, requesterID, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:          bloodRequestID,
			RequesterID: requesterID,
		}, nil).
		Once()

	donorMatchRepo.
		On("GetByBloodRequestAndDonor", mock.Anything, bloodRequestID, donorID).
		Return(existingMatch, nil).
		Once()

	result, err := uc.Invite(requesterID, bloodRequestID, donorID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, existingMatch, result)
	assert.Equal(t, "invited", result.Status)

	donorMatchRepo.AssertNotCalled(t, "UpdateStatus")
}

// Test error repository ketika status donor match gagal diubah menjadi invited.
func TestDonorMatchUsecase_Invite_UpdateStatusError(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, _, _, _ := newDonorMatchUsecaseForTest()

	requesterID := uuid.New()
	bloodRequestID := uuid.New()
	donorID := uuid.New()
	matchID := uuid.New()
	expectedErr := errors.New("failed to update status")

	bloodReqRepo.
		On("FindOne", mock.Anything, requesterID, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:          bloodRequestID,
			RequesterID: requesterID,
		}, nil).
		Once()

	donorMatchRepo.
		On("GetByBloodRequestAndDonor", mock.Anything, bloodRequestID, donorID).
		Return(&entity.DonorMatch{
			ID:             matchID,
			BloodRequestID: bloodRequestID,
			DonorID:        donorID,
			Status:         "pending",
		}, nil).
		Once()

	donorMatchRepo.
		On("UpdateStatus", mock.Anything, matchID, "invited").
		Return(expectedErr).
		Once()

	result, err := uc.Invite(requesterID, bloodRequestID, donorID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)
}

// Test requester berhasil mengubah donor match dari pending menjadi invited.
func TestDonorMatchUsecase_Invite_Success(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, _, _, _ := newDonorMatchUsecaseForTest()

	requesterID := uuid.New()
	bloodRequestID := uuid.New()
	donorID := uuid.New()
	matchID := uuid.New()

	donorMatch := &entity.DonorMatch{
		ID:             matchID,
		BloodRequestID: bloodRequestID,
		DonorID:        donorID,
		Status:         "pending",
	}

	bloodReqRepo.
		On("FindOne", mock.Anything, requesterID, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:          bloodRequestID,
			RequesterID: requesterID,
		}, nil).
		Once()

	donorMatchRepo.
		On("GetByBloodRequestAndDonor", mock.Anything, bloodRequestID, donorID).
		Return(donorMatch, nil).
		Once()

	donorMatchRepo.
		On("UpdateStatus", mock.Anything, matchID, "invited").
		Return(nil).
		Once()

	result, err := uc.Invite(requesterID, bloodRequestID, donorID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "invited", result.Status)

	donorMatchRepo.AssertExpectations(t)
}

// ============================================================
// ACCEPT TEST
// ============================================================

// Test ketika donor match tidak ditemukan / repository mengembalikan error.
func TestDonorMatchUsecase_Accept_DonorMatchError(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, _, _, donationRepo := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()
	expectedErr := domain.ErrDonorMatchNotFound

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(nil, expectedErr).
		Once()

	result, err := uc.Accept(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	bloodReqRepo.AssertNotCalled(t, "GetById")
	donationRepo.AssertNotCalled(t, "Create")
}

// Test hanya donor pemilik donor match yang boleh accept.
func TestDonorMatchUsecase_Accept_Forbidden(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, _, _, donationRepo := newDonorMatchUsecaseForTest()

	loggedInDonorID := uuid.New()
	ownerDonorID := uuid.New()
	matchID := uuid.New()

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:      matchID,
			DonorID: ownerDonorID,
			Status:  "invited",
		}, nil).
		Once()

	result, err := uc.Accept(loggedInDonorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrForbidden)

	bloodReqRepo.AssertNotCalled(t, "GetById")
	donationRepo.AssertNotCalled(t, "Create")
}

// Test hanya donor match berstatus invited yang boleh di-accept.
func TestDonorMatchUsecase_Accept_InvalidStatus(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, _, _, donationRepo := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:      matchID,
			DonorID: donorID,
			Status:  "accepted",
		}, nil).
		Once()

	result, err := uc.Accept(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrInvalidMatchStatus)

	bloodReqRepo.AssertNotCalled(t, "GetById")
	donationRepo.AssertNotCalled(t, "Create")
}

// Test ketika blood request gagal diambil.
func TestDonorMatchUsecase_Accept_BloodRequestError(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, _, _, donationRepo := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()
	bloodRequestID := uuid.New()
	expectedErr := errors.New("blood request error")

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:             matchID,
			DonorID:        donorID,
			BloodRequestID: bloodRequestID,
			Status:         "invited",
		}, nil).
		Once()

	bloodReqRepo.
		On("GetById", mock.Anything, bloodRequestID).
		Return(nil, expectedErr).
		Once()

	result, err := uc.Accept(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	donorMatchRepo.AssertNotCalled(t, "CountAccepted")
	donationRepo.AssertNotCalled(t, "Create")
}

// Test error ketika repository gagal menghitung jumlah donor accepted.
func TestDonorMatchUsecase_Accept_CountAcceptedError(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, _, _, donationRepo := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()
	bloodRequestID := uuid.New()
	expectedErr := errors.New("count accepted error")

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:             matchID,
			DonorID:        donorID,
			BloodRequestID: bloodRequestID,
			Status:         "invited",
		}, nil).
		Once()

	bloodReqRepo.
		On("GetById", mock.Anything, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:       bloodRequestID,
			Quantity: 2,
		}, nil).
		Once()

	donorMatchRepo.
		On("CountAccepted", mock.Anything, bloodRequestID).
		Return(int64(0), expectedErr).
		Once()

	result, err := uc.Accept(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	donationRepo.AssertNotCalled(t, "Create")
}

// Test quantity = 2 dan acceptedCount = 2.
// Donor berikutnya harus ditolak dengan ErrBloodRequestFulfilled.
// Update distance, update status, dan create donation tidak boleh dijalankan.
func TestDonorMatchUsecase_Accept_BloodRequestFulfilled(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, userServiceRepo, geoapifyRepo, donationRepo := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()
	bloodRequestID := uuid.New()

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:             matchID,
			DonorID:        donorID,
			BloodRequestID: bloodRequestID,
			Status:         "invited",
		}, nil).
		Once()

	bloodReqRepo.
		On("GetById", mock.Anything, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:       bloodRequestID,
			Quantity: 2,
		}, nil).
		Once()

	donorMatchRepo.
		On("CountAccepted", mock.Anything, bloodRequestID).
		Return(int64(2), nil).
		Once()

	result, err := uc.Accept(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrBloodRequestFulfilled)

	userServiceRepo.AssertNotCalled(t, "GetDonorProfile")
	geoapifyRepo.AssertNotCalled(t, "GetGeoapifyRoute")
	donorMatchRepo.AssertNotCalled(t, "UpdateDistance")
	donorMatchRepo.AssertNotCalled(t, "UpdateStatus")
	donationRepo.AssertNotCalled(t, "Create")
}

// Test error ketika donor profile gagal diambil dari user-service.
func TestDonorMatchUsecase_Accept_DonorProfileError(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, userServiceRepo, _, donationRepo := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()
	bloodRequestID := uuid.New()
	expectedErr := errors.New("donor profile error")

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:             matchID,
			DonorID:        donorID,
			BloodRequestID: bloodRequestID,
			Status:         "invited",
		}, nil).
		Once()

	bloodReqRepo.
		On("GetById", mock.Anything, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:       bloodRequestID,
			Quantity: 2,
		}, nil).
		Once()

	donorMatchRepo.
		On("CountAccepted", mock.Anything, bloodRequestID).
		Return(int64(1), nil).
		Once()

	userServiceRepo.
		On("GetDonorProfile", mock.Anything, donorID).
		Return(nil, expectedErr).
		Once()

	result, err := uc.Accept(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	donationRepo.AssertNotCalled(t, "Create")
}

// Test error dari Geoapify harus diteruskan oleh usecase.
func TestDonorMatchUsecase_Accept_GeoapifyError(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, userServiceRepo, geoapifyRepo, donationRepo := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()
	bloodRequestID := uuid.New()
	expectedErr := errors.New("geoapify error")

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:             matchID,
			DonorID:        donorID,
			BloodRequestID: bloodRequestID,
			Status:         "invited",
		}, nil).
		Once()

	bloodReqRepo.
		On("GetById", mock.Anything, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:        bloodRequestID,
			Quantity:  2,
			Latitude:  -6.2,
			Longitude: 106.8,
		}, nil).
		Once()

	userServiceRepo.
		On("GetDonorProfile", mock.Anything, donorID).
		Return(&entity.DonorProfile{
			ID:        donorID,
			Latitude:  -6.3,
			Longitude: 106.9,
		}, nil).
		Once()

	donorMatchRepo.
		On("CountAccepted", mock.Anything, bloodRequestID).
		Return(int64(1), nil).
		Once()

	geoapifyRepo.
		On(
			"GetGeoapifyRoute",
			mock.Anything,
			mock.AnythingOfType("*dto.GeoapifyRoutingRequest"),
		).
		Return(nil, expectedErr).
		Once()

	result, err := uc.Accept(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	donorMatchRepo.AssertNotCalled(t, "UpdateDistance")
	donationRepo.AssertNotCalled(t, "Create")
}

// Test error repository ketika distance gagal disimpan.
func TestDonorMatchUsecase_Accept_UpdateDistanceError(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, userServiceRepo, geoapifyRepo, donationRepo := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()
	bloodRequestID := uuid.New()
	expectedErr := errors.New("update distance error")

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:             matchID,
			DonorID:        donorID,
			BloodRequestID: bloodRequestID,
			Status:         "invited",
		}, nil).
		Once()

	bloodReqRepo.
		On("GetById", mock.Anything, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:        bloodRequestID,
			Quantity:  2,
			Latitude:  -6.2,
			Longitude: 106.8,
		}, nil).
		Once()

	donorMatchRepo.
		On("CountAccepted", mock.Anything, bloodRequestID).
		Return(int64(1), nil).
		Once()

	userServiceRepo.
		On("GetDonorProfile", mock.Anything, donorID).
		Return(&entity.DonorProfile{
			ID:        donorID,
			Latitude:  -6.3,
			Longitude: 106.9,
		}, nil).
		Once()

	geoapifyRepo.
		On("GetGeoapifyRoute", mock.Anything, mock.Anything).
		Return(&domain.GeoapifyRoute{
			Distance:      10000,
			DistanceUnits: "meters",
		}, nil).
		Once()

	donorMatchRepo.
		On("UpdateDistance", mock.Anything, matchID, float64(10)).
		Return(expectedErr).
		Once()

	result, err := uc.Accept(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	donorMatchRepo.AssertNotCalled(t, "UpdateStatus")
	donationRepo.AssertNotCalled(t, "Create")
}

// Test error repository ketika status donor match gagal menjadi accepted.
func TestDonorMatchUsecase_Accept_UpdateStatusError(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, userServiceRepo, geoapifyRepo, donationRepo := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()
	bloodRequestID := uuid.New()
	expectedErr := errors.New("update status error")

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:             matchID,
			DonorID:        donorID,
			BloodRequestID: bloodRequestID,
			Status:         "invited",
		}, nil).
		Once()

	bloodReqRepo.
		On("GetById", mock.Anything, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:        bloodRequestID,
			Quantity:  2,
			Latitude:  -6.2,
			Longitude: 106.8,
		}, nil).
		Once()

	donorMatchRepo.
		On("CountAccepted", mock.Anything, bloodRequestID).
		Return(int64(1), nil).
		Once()

	userServiceRepo.
		On("GetDonorProfile", mock.Anything, donorID).
		Return(&entity.DonorProfile{
			ID:        donorID,
			Latitude:  -6.3,
			Longitude: 106.9,
		}, nil).
		Once()

	geoapifyRepo.
		On("GetGeoapifyRoute", mock.Anything, mock.Anything).
		Return(&domain.GeoapifyRoute{
			Distance:      10000,
			DistanceUnits: "meters",
		}, nil).
		Once()

	donorMatchRepo.
		On("UpdateDistance", mock.Anything, matchID, float64(10)).
		Return(nil).
		Once()

	donorMatchRepo.
		On("UpdateStatus", mock.Anything, matchID, "accepted").
		Return(expectedErr).
		Once()

	result, err := uc.Accept(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	donationRepo.AssertNotCalled(t, "Create")
}

// Test error saat auto-create donation.
// Pada implementasi sekarang, update distance dan status terjadi sebelum create donation.
func TestDonorMatchUsecase_Accept_DonationCreateError(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, userServiceRepo, geoapifyRepo, donationRepo := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()
	bloodRequestID := uuid.New()
	expectedErr := errors.New("create donation error")

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:             matchID,
			DonorID:        donorID,
			BloodRequestID: bloodRequestID,
			Status:         "invited",
		}, nil).
		Once()

	bloodReqRepo.
		On("GetById", mock.Anything, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:        bloodRequestID,
			Quantity:  2,
			Latitude:  -6.2,
			Longitude: 106.8,
		}, nil).
		Once()

	donorMatchRepo.
		On("CountAccepted", mock.Anything, bloodRequestID).
		Return(int64(1), nil).
		Once()

	userServiceRepo.
		On("GetDonorProfile", mock.Anything, donorID).
		Return(&entity.DonorProfile{
			ID:        donorID,
			Latitude:  -6.3,
			Longitude: 106.9,
		}, nil).
		Once()

	geoapifyRepo.
		On("GetGeoapifyRoute", mock.Anything, mock.Anything).
		Return(&domain.GeoapifyRoute{
			Distance:      10000,
			DistanceUnits: "meters",
		}, nil).
		Once()

	donorMatchRepo.
		On("UpdateDistance", mock.Anything, matchID, float64(10)).
		Return(nil).
		Once()

	donorMatchRepo.
		On("UpdateStatus", mock.Anything, matchID, "accepted").
		Return(nil).
		Once()

	donationRepo.
		On(
			"Create",
			mock.Anything,
			mock.MatchedBy(func(donation *entity.Donation) bool {
				return donation.BloodRequestID == bloodRequestID &&
					donation.DonorID == donorID &&
					donation.DonorMatchID == matchID &&
					donation.Status == "pending"
			}),
		).
		Return(expectedErr).
		Once()

	result, err := uc.Accept(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	donationRepo.AssertExpectations(t)
}

// Test quantity = 2 dan acceptedCount = 1.
// Donor kedua masih boleh accept.
// Test juga memastikan:
// - Geoapify route dipanggil
// - distance 12500 meter diubah menjadi 12.5 KM
// - status menjadi accepted
// - donation pending otomatis dibuat.
func TestDonorMatchUsecase_Accept_Success(t *testing.T) {
	uc, donorMatchRepo, bloodReqRepo, userServiceRepo, geoapifyRepo, donationRepo := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()
	bloodRequestID := uuid.New()

	donorMatch := &entity.DonorMatch{
		ID:             matchID,
		DonorID:        donorID,
		BloodRequestID: bloodRequestID,
		Status:         "invited",
	}

	bloodReqRepo.
		On("GetById", mock.Anything, bloodRequestID).
		Return(&entity.BloodRequest{
			ID:        bloodRequestID,
			Quantity:  2,
			Latitude:  -6.2,
			Longitude: 106.8,
		}, nil).
		Once()

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(donorMatch, nil).
		Once()

	donorMatchRepo.
		On("CountAccepted", mock.Anything, bloodRequestID).
		Return(int64(1), nil).
		Once()

	userServiceRepo.
		On("GetDonorProfile", mock.Anything, donorID).
		Return(&entity.DonorProfile{
			ID:        donorID,
			Latitude:  -6.3,
			Longitude: 106.9,
		}, nil).
		Once()

	geoapifyRepo.
		On(
			"GetGeoapifyRoute",
			mock.Anything,
			mock.MatchedBy(func(req *dto.GeoapifyRoutingRequest) bool {
				return req.OriginLat == -6.3 &&
					req.OriginLon == 106.9 &&
					req.DestinationLat == -6.2 &&
					req.DestinationLon == 106.8
			}),
		).
		Return(&domain.GeoapifyRoute{
			Distance:      12500,
			DistanceUnits: "meters",
			Time:          900,
			Mode:          "drive",
		}, nil).
		Once()

	donorMatchRepo.
		On("UpdateDistance", mock.Anything, matchID, float64(12.5)).
		Return(nil).
		Once()

	donorMatchRepo.
		On("UpdateStatus", mock.Anything, matchID, "accepted").
		Return(nil).
		Once()

	donationRepo.
		On(
			"Create",
			mock.Anything,
			mock.MatchedBy(func(donation *entity.Donation) bool {
				return donation.ID != uuid.Nil &&
					donation.BloodRequestID == bloodRequestID &&
					donation.DonorID == donorID &&
					donation.DonorMatchID == matchID &&
					donation.Status == "pending"
			}),
		).
		Return(nil).
		Once()

	result, err := uc.Accept(donorID, matchID)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "accepted", result.Status)
	assert.Equal(t, float64(12.5), result.DistanceKM)

	donorMatchRepo.AssertExpectations(t)
	bloodReqRepo.AssertExpectations(t)
	userServiceRepo.AssertExpectations(t)
	geoapifyRepo.AssertExpectations(t)
	donationRepo.AssertExpectations(t)
}

// ============================================================
// DECLINE TEST
// ============================================================

// Test ketika donor match tidak ditemukan.
func TestDonorMatchUsecase_Decline_DonorMatchError(t *testing.T) {
	uc, donorMatchRepo, _, _, _, _ := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()
	expectedErr := domain.ErrDonorMatchNotFound

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(nil, expectedErr).
		Once()

	result, err := uc.Decline(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)
}

// Test hanya donor pemilik donor match yang boleh decline invitation.
func TestDonorMatchUsecase_Decline_Forbidden(t *testing.T) {
	uc, donorMatchRepo, _, _, _, _ := newDonorMatchUsecaseForTest()

	loggedInDonorID := uuid.New()
	ownerDonorID := uuid.New()
	matchID := uuid.New()

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:      matchID,
			DonorID: ownerDonorID,
			Status:  "invited",
		}, nil).
		Once()

	result, err := uc.Decline(loggedInDonorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrForbidden)

	donorMatchRepo.AssertNotCalled(t, "UpdateStatus")
}

// Test hanya status invited yang boleh di-decline.
func TestDonorMatchUsecase_Decline_InvalidStatus(t *testing.T) {
	uc, donorMatchRepo, _, _, _, _ := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:      matchID,
			DonorID: donorID,
			Status:  "accepted",
		}, nil).
		Once()

	result, err := uc.Decline(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrInvalidMatchStatus)

	donorMatchRepo.AssertNotCalled(t, "UpdateStatus")
}

// Test error repository ketika status gagal diubah menjadi declined.
func TestDonorMatchUsecase_Decline_UpdateStatusError(t *testing.T) {
	uc, donorMatchRepo, _, _, _, _ := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()
	expectedErr := errors.New("update status error")

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(&entity.DonorMatch{
			ID:      matchID,
			DonorID: donorID,
			Status:  "invited",
		}, nil).
		Once()

	donorMatchRepo.
		On("UpdateStatus", mock.Anything, matchID, "declined").
		Return(expectedErr).
		Once()

	result, err := uc.Decline(donorID, matchID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)
}

// Test invitation berhasil di-decline dan status donor match berubah menjadi declined.
func TestDonorMatchUsecase_Decline_Success(t *testing.T) {
	uc, donorMatchRepo, _, _, _, _ := newDonorMatchUsecaseForTest()

	donorID := uuid.New()
	matchID := uuid.New()

	donorMatch := &entity.DonorMatch{
		ID:      matchID,
		DonorID: donorID,
		Status:  "invited",
	}

	donorMatchRepo.
		On("GetByID", mock.Anything, matchID).
		Return(donorMatch, nil).
		Once()

	donorMatchRepo.
		On("UpdateStatus", mock.Anything, matchID, "declined").
		Return(nil).
		Once()

	result, err := uc.Decline(donorID, matchID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "declined", result.Status)

	donorMatchRepo.AssertExpectations(t)
}

// ============================================================
// GET BY REQUESTER ID TEST
// ============================================================

// Test repository error ketika mengambil donor match berdasarkan requester.
func TestDonorMatchUsecase_GetByDonorID_RepositoryError(t *testing.T) {
	uc, donorMatchRepo, _, _, _, _ := newDonorMatchUsecaseForTest()

	requesterID := uuid.New()
	expectedErr := errors.New("get donor matches error")

	donorMatchRepo.
		On("GetByDonorID", mock.Anything, requesterID).
		Return(nil, expectedErr).
		Once()

	result, err := uc.GetByDonorID(requesterID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	donorMatchRepo.AssertExpectations(t)
}

// Test berhasil mengambil semua donor match milik requester.
func TestDonorMatchUsecase_GetByDonorID_Success(t *testing.T) {
	uc, donorMatchRepo, _, _, _, _ := newDonorMatchUsecaseForTest()

	requesterID := uuid.New()

	expected := []entity.DonorMatch{
		{
			ID:     uuid.New(),
			Status: "invited",
		},
		{
			ID:     uuid.New(),
			Status: "accepted",
		},
	}

	donorMatchRepo.
		On("GetByDonorID", mock.Anything, requesterID).
		Return(expected, nil).
		Once()

	result, err := uc.GetByDonorID(requesterID)

	require.NoError(t, err)
	assert.Equal(t, expected, result)

	donorMatchRepo.AssertExpectations(t)
}
