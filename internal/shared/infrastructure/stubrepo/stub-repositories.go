package stubrepo

import (
	"context"
	"time"

	cosmeticdomain "ditza/internal/cosmetic/domain"
	friendshipdomain "ditza/internal/friendship/domain"
	habitcompletiondomain "ditza/internal/habit-completion/domain"
	habitdomain "ditza/internal/habit/domain"
	leaderboarddomain "ditza/internal/leaderboard/domain"
	petdomain "ditza/internal/pet/domain"
	pointtransactiondomain "ditza/internal/point-transaction/domain"
	seasonscoredomain "ditza/internal/season-score/domain"
	seasondomain "ditza/internal/season/domain"
	domainerror "ditza/internal/shared/domain/error"
	valueobject "ditza/internal/shared/domain/value-object"
	usercosmeticdomain "ditza/internal/user-cosmetic/domain"
	userdomain "ditza/internal/user/domain"
)

func notImplementedError() error {
	return domainerror.New("NOT_IMPLEMENTED", "repositorio no implementado", domainerror.ErrNotImplemented)
}

type UnitOfWorkStub struct{}

func (u *UnitOfWorkStub) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type UserRepositoryStub struct{}

func (r *UserRepositoryStub) Create(ctx context.Context, entity *userdomain.User) error {
	return notImplementedError()
}
func (r *UserRepositoryStub) Update(ctx context.Context, entity *userdomain.User) error {
	return notImplementedError()
}
func (r *UserRepositoryStub) FindByID(ctx context.Context, id valueobject.UserID) (*userdomain.User, error) {
	return nil, notImplementedError()
}
func (r *UserRepositoryStub) FindByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	return nil, notImplementedError()
}
func (r *UserRepositoryStub) FindByFriendCode(ctx context.Context, friendCode string) (*userdomain.User, error) {
	return nil, notImplementedError()
}
func (r *UserRepositoryStub) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, notImplementedError()
}

type HabitRepositoryStub struct{}

func (r *HabitRepositoryStub) Create(ctx context.Context, entity *habitdomain.Habit) error {
	return notImplementedError()
}
func (r *HabitRepositoryStub) Update(ctx context.Context, entity *habitdomain.Habit) error {
	return notImplementedError()
}
func (r *HabitRepositoryStub) FindByID(ctx context.Context, id valueobject.HabitID) (*habitdomain.Habit, error) {
	return nil, notImplementedError()
}
func (r *HabitRepositoryStub) FindActiveByUserID(ctx context.Context, userID valueobject.UserID) ([]habitdomain.Habit, error) {
	return nil, notImplementedError()
}
func (r *HabitRepositoryStub) CountActiveByUserID(ctx context.Context, userID valueobject.UserID) (int, error) {
	return 0, notImplementedError()
}

type HabitCompletionRepositoryStub struct{}

func (r *HabitCompletionRepositoryStub) Create(ctx context.Context, entity *habitcompletiondomain.HabitCompletion) error {
	return notImplementedError()
}
func (r *HabitCompletionRepositoryStub) ExistsForHabitOnDate(ctx context.Context, habitID valueobject.HabitID, date time.Time) (bool, error) {
	return false, notImplementedError()
}
func (r *HabitCompletionRepositoryStub) CountByUserOnDate(ctx context.Context, userID valueobject.UserID, date time.Time) (int, error) {
	return 0, notImplementedError()
}
func (r *HabitCompletionRepositoryStub) FindByUserIDAndDateRange(ctx context.Context, userID valueobject.UserID, from, to time.Time) ([]habitcompletiondomain.HabitCompletion, error) {
	return nil, notImplementedError()
}

type PetRepositoryStub struct{}

func (r *PetRepositoryStub) Create(ctx context.Context, entity *petdomain.Pet) error {
	return notImplementedError()
}
func (r *PetRepositoryStub) Update(ctx context.Context, entity *petdomain.Pet) error {
	return notImplementedError()
}
func (r *PetRepositoryStub) FindByUserID(ctx context.Context, userID valueobject.UserID) (*petdomain.Pet, error) {
	return nil, notImplementedError()
}

type CosmeticRepositoryStub struct{}

func (r *CosmeticRepositoryStub) Create(ctx context.Context, entity *cosmeticdomain.Cosmetic) error {
	return notImplementedError()
}
func (r *CosmeticRepositoryStub) Update(ctx context.Context, entity *cosmeticdomain.Cosmetic) error {
	return notImplementedError()
}
func (r *CosmeticRepositoryStub) FindByID(ctx context.Context, id valueobject.CosmeticID) (*cosmeticdomain.Cosmetic, error) {
	return nil, notImplementedError()
}
func (r *CosmeticRepositoryStub) FindAllActive(ctx context.Context) ([]cosmeticdomain.Cosmetic, error) {
	return nil, notImplementedError()
}

type UserCosmeticRepositoryStub struct{}

func (r *UserCosmeticRepositoryStub) Create(ctx context.Context, entity *usercosmeticdomain.UserCosmetic) error {
	return notImplementedError()
}
func (r *UserCosmeticRepositoryStub) Exists(ctx context.Context, userID valueobject.UserID, cosmeticID valueobject.CosmeticID) (bool, error) {
	return false, notImplementedError()
}
func (r *UserCosmeticRepositoryStub) FindByUserID(ctx context.Context, userID valueobject.UserID) ([]usercosmeticdomain.UserCosmetic, error) {
	return nil, notImplementedError()
}

type FriendshipRepositoryStub struct{}

func (r *FriendshipRepositoryStub) Create(ctx context.Context, entity *friendshipdomain.Friendship) error {
	return notImplementedError()
}
func (r *FriendshipRepositoryStub) Update(ctx context.Context, entity *friendshipdomain.Friendship) error {
	return notImplementedError()
}
func (r *FriendshipRepositoryStub) FindByID(ctx context.Context, id valueobject.FriendshipID) (*friendshipdomain.Friendship, error) {
	return nil, notImplementedError()
}
func (r *FriendshipRepositoryStub) FindBetweenUsers(ctx context.Context, userA, userB valueobject.UserID) (*friendshipdomain.Friendship, error) {
	return nil, notImplementedError()
}
func (r *FriendshipRepositoryStub) FindAcceptedByUserID(ctx context.Context, userID valueobject.UserID) ([]friendshipdomain.Friendship, error) {
	return nil, notImplementedError()
}
func (r *FriendshipRepositoryStub) FindPendingByAddresseeID(ctx context.Context, addresseeID valueobject.UserID) ([]friendshipdomain.Friendship, error) {
	return nil, notImplementedError()
}

type SeasonRepositoryStub struct{}

func (r *SeasonRepositoryStub) Create(ctx context.Context, entity *seasondomain.Season) error {
	return notImplementedError()
}
func (r *SeasonRepositoryStub) Update(ctx context.Context, entity *seasondomain.Season) error {
	return notImplementedError()
}
func (r *SeasonRepositoryStub) FindByID(ctx context.Context, id valueobject.SeasonID) (*seasondomain.Season, error) {
	return nil, notImplementedError()
}
func (r *SeasonRepositoryStub) FindActive(ctx context.Context) (*seasondomain.Season, error) {
	return nil, notImplementedError()
}
func (r *SeasonRepositoryStub) DeactivateAll(ctx context.Context) error {
	return notImplementedError()
}

type SeasonScoreRepositoryStub struct{}

func (r *SeasonScoreRepositoryStub) Create(ctx context.Context, entity *seasonscoredomain.SeasonScore) error {
	return notImplementedError()
}
func (r *SeasonScoreRepositoryStub) Update(ctx context.Context, entity *seasonscoredomain.SeasonScore) error {
	return notImplementedError()
}
func (r *SeasonScoreRepositoryStub) FindByUserAndSeason(ctx context.Context, userID valueobject.UserID, seasonID valueobject.SeasonID) (*seasonscoredomain.SeasonScore, error) {
	return nil, notImplementedError()
}
func (r *SeasonScoreRepositoryStub) FindByUserIDsAndSeason(ctx context.Context, userIDs []valueobject.UserID, seasonID valueobject.SeasonID) ([]seasonscoredomain.SeasonScore, error) {
	return nil, notImplementedError()
}
func (r *SeasonScoreRepositoryStub) ResetBySeasonID(ctx context.Context, seasonID valueobject.SeasonID) error {
	return notImplementedError()
}

type PointTransactionRepositoryStub struct{}

func (r *PointTransactionRepositoryStub) Create(ctx context.Context, entity *pointtransactiondomain.PointTransaction) error {
	return notImplementedError()
}
func (r *PointTransactionRepositoryStub) FindByUserID(ctx context.Context, userID valueobject.UserID, limit int) ([]pointtransactiondomain.PointTransaction, error) {
	return nil, notImplementedError()
}

type LeaderboardRepositoryStub struct{}

func (r *LeaderboardRepositoryStub) GetFriendRanking(ctx context.Context, userID valueobject.UserID, seasonID valueobject.SeasonID) ([]leaderboarddomain.FriendEntry, error) {
	return nil, notImplementedError()
}
