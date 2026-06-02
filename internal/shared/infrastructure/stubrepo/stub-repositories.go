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
	"ditza/internal/shared/infrastructure/logger"
	"ditza/internal/shared/infrastructure/monitoring"
	usercosmeticdomain "ditza/internal/user-cosmetic/domain"
	userdomain "ditza/internal/user/domain"
)

func notImplementedError() error {
	return domainerror.New("NOT_IMPLEMENTED", "repositorio no implementado", domainerror.ErrNotImplemented)
}

func logRepo(model, operation string, attrs map[string]any) error {
	monitoring.Repository(model, operation, attrs)
	return notImplementedError()
}

type UnitOfWorkStub struct{}

func (u *UnitOfWorkStub) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	logger.App().Info("unit_of_work_transaction", "operation", "within_transaction")
	return fn(ctx)
}

type UserRepositoryStub struct{}

func (r *UserRepositoryStub) Create(ctx context.Context, entity *userdomain.User) error {
	return logRepo(logger.ModelUser, "create", map[string]any{"user_id": entity.ID, "email": entity.Email})
}
func (r *UserRepositoryStub) Update(ctx context.Context, entity *userdomain.User) error {
	return logRepo(logger.ModelUser, "update", map[string]any{"user_id": entity.ID})
}
func (r *UserRepositoryStub) FindByID(ctx context.Context, id valueobject.UserID) (*userdomain.User, error) {
	return nil, logRepo(logger.ModelUser, "find_by_id", map[string]any{"user_id": id})
}
func (r *UserRepositoryStub) FindByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	return nil, logRepo(logger.ModelUser, "find_by_email", map[string]any{"email": email})
}
func (r *UserRepositoryStub) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, logRepo(logger.ModelUser, "exists_by_email", map[string]any{"email": email})
}

type HabitRepositoryStub struct{}

func (r *HabitRepositoryStub) Create(ctx context.Context, entity *habitdomain.Habit) error {
	return logRepo(logger.ModelHabit, "create", map[string]any{"habit_id": entity.ID, "user_id": entity.UserID})
}
func (r *HabitRepositoryStub) Update(ctx context.Context, entity *habitdomain.Habit) error {
	return logRepo(logger.ModelHabit, "update", map[string]any{"habit_id": entity.ID})
}
func (r *HabitRepositoryStub) FindByID(ctx context.Context, id valueobject.HabitID) (*habitdomain.Habit, error) {
	return nil, logRepo(logger.ModelHabit, "find_by_id", map[string]any{"habit_id": id})
}
func (r *HabitRepositoryStub) FindActiveByUserID(ctx context.Context, userID valueobject.UserID) ([]habitdomain.Habit, error) {
	return nil, logRepo(logger.ModelHabit, "find_active_by_user_id", map[string]any{"user_id": userID})
}
func (r *HabitRepositoryStub) CountActiveByUserID(ctx context.Context, userID valueobject.UserID) (int, error) {
	return 0, logRepo(logger.ModelHabit, "count_active_by_user_id", map[string]any{"user_id": userID})
}

type HabitCompletionRepositoryStub struct{}

func (r *HabitCompletionRepositoryStub) Create(ctx context.Context, entity *habitcompletiondomain.HabitCompletion) error {
	return logRepo(logger.ModelHabitCompletion, "create", map[string]any{"completion_id": entity.ID, "habit_id": entity.HabitID})
}
func (r *HabitCompletionRepositoryStub) ExistsForHabitOnDate(ctx context.Context, habitID valueobject.HabitID, date time.Time) (bool, error) {
	return false, logRepo(logger.ModelHabitCompletion, "exists_for_habit_on_date", map[string]any{"habit_id": habitID, "date": date.Format("2006-01-02")})
}
func (r *HabitCompletionRepositoryStub) CountByUserOnDate(ctx context.Context, userID valueobject.UserID, date time.Time) (int, error) {
	return 0, logRepo(logger.ModelHabitCompletion, "count_by_user_on_date", map[string]any{"user_id": userID, "date": date.Format("2006-01-02")})
}
func (r *HabitCompletionRepositoryStub) FindByUserIDAndDateRange(ctx context.Context, userID valueobject.UserID, from, to time.Time) ([]habitcompletiondomain.HabitCompletion, error) {
	return nil, logRepo(logger.ModelHabitCompletion, "find_by_user_id_and_date_range", map[string]any{"user_id": userID, "from": from.Format("2006-01-02"), "to": to.Format("2006-01-02")})
}

type PetRepositoryStub struct{}

func (r *PetRepositoryStub) Create(ctx context.Context, entity *petdomain.Pet) error {
	return logRepo(logger.ModelPet, "create", map[string]any{"user_id": entity.UserID})
}
func (r *PetRepositoryStub) Update(ctx context.Context, entity *petdomain.Pet) error {
	return logRepo(logger.ModelPet, "update", map[string]any{"user_id": entity.UserID})
}
func (r *PetRepositoryStub) FindByUserID(ctx context.Context, userID valueobject.UserID) (*petdomain.Pet, error) {
	return nil, logRepo(logger.ModelPet, "find_by_user_id", map[string]any{"user_id": userID})
}

type CosmeticRepositoryStub struct{}

func (r *CosmeticRepositoryStub) Create(ctx context.Context, entity *cosmeticdomain.Cosmetic) error {
	return logRepo(logger.ModelCosmetic, "create", map[string]any{"cosmetic_id": entity.ID})
}
func (r *CosmeticRepositoryStub) Update(ctx context.Context, entity *cosmeticdomain.Cosmetic) error {
	return logRepo(logger.ModelCosmetic, "update", map[string]any{"cosmetic_id": entity.ID})
}
func (r *CosmeticRepositoryStub) FindByID(ctx context.Context, id valueobject.CosmeticID) (*cosmeticdomain.Cosmetic, error) {
	return nil, logRepo(logger.ModelCosmetic, "find_by_id", map[string]any{"cosmetic_id": id})
}
func (r *CosmeticRepositoryStub) FindAllActive(ctx context.Context) ([]cosmeticdomain.Cosmetic, error) {
	return nil, logRepo(logger.ModelCosmetic, "find_all_active", nil)
}

type UserCosmeticRepositoryStub struct{}

func (r *UserCosmeticRepositoryStub) Create(ctx context.Context, entity *usercosmeticdomain.UserCosmetic) error {
	return logRepo(logger.ModelUserCosmetic, "create", map[string]any{"user_id": entity.UserID, "cosmetic_id": entity.CosmeticID})
}
func (r *UserCosmeticRepositoryStub) Exists(ctx context.Context, userID valueobject.UserID, cosmeticID valueobject.CosmeticID) (bool, error) {
	return false, logRepo(logger.ModelUserCosmetic, "exists", map[string]any{"user_id": userID, "cosmetic_id": cosmeticID})
}
func (r *UserCosmeticRepositoryStub) FindByUserID(ctx context.Context, userID valueobject.UserID) ([]usercosmeticdomain.UserCosmetic, error) {
	return nil, logRepo(logger.ModelUserCosmetic, "find_by_user_id", map[string]any{"user_id": userID})
}

type FriendshipRepositoryStub struct{}

func (r *FriendshipRepositoryStub) Create(ctx context.Context, entity *friendshipdomain.Friendship) error {
	return logRepo(logger.ModelFriendship, "create", map[string]any{"friendship_id": entity.ID})
}
func (r *FriendshipRepositoryStub) Update(ctx context.Context, entity *friendshipdomain.Friendship) error {
	return logRepo(logger.ModelFriendship, "update", map[string]any{"friendship_id": entity.ID})
}
func (r *FriendshipRepositoryStub) FindByID(ctx context.Context, id valueobject.FriendshipID) (*friendshipdomain.Friendship, error) {
	return nil, logRepo(logger.ModelFriendship, "find_by_id", map[string]any{"friendship_id": id})
}
func (r *FriendshipRepositoryStub) FindBetweenUsers(ctx context.Context, userA, userB valueobject.UserID) (*friendshipdomain.Friendship, error) {
	return nil, logRepo(logger.ModelFriendship, "find_between_users", map[string]any{"user_a": userA, "user_b": userB})
}
func (r *FriendshipRepositoryStub) FindAcceptedByUserID(ctx context.Context, userID valueobject.UserID) ([]friendshipdomain.Friendship, error) {
	return nil, logRepo(logger.ModelFriendship, "find_accepted_by_user_id", map[string]any{"user_id": userID})
}
func (r *FriendshipRepositoryStub) FindPendingByAddresseeID(ctx context.Context, addresseeID valueobject.UserID) ([]friendshipdomain.Friendship, error) {
	return nil, logRepo(logger.ModelFriendship, "find_pending_by_addressee_id", map[string]any{"addressee_id": addresseeID})
}

type SeasonRepositoryStub struct{}

func (r *SeasonRepositoryStub) Create(ctx context.Context, entity *seasondomain.Season) error {
	return logRepo(logger.ModelSeason, "create", map[string]any{"season_id": entity.ID})
}
func (r *SeasonRepositoryStub) Update(ctx context.Context, entity *seasondomain.Season) error {
	return logRepo(logger.ModelSeason, "update", map[string]any{"season_id": entity.ID})
}
func (r *SeasonRepositoryStub) FindByID(ctx context.Context, id valueobject.SeasonID) (*seasondomain.Season, error) {
	return nil, logRepo(logger.ModelSeason, "find_by_id", map[string]any{"season_id": id})
}
func (r *SeasonRepositoryStub) FindActive(ctx context.Context) (*seasondomain.Season, error) {
	return nil, logRepo(logger.ModelSeason, "find_active", nil)
}
func (r *SeasonRepositoryStub) DeactivateAll(ctx context.Context) error {
	return logRepo(logger.ModelSeason, "deactivate_all", nil)
}

type SeasonScoreRepositoryStub struct{}

func (r *SeasonScoreRepositoryStub) Create(ctx context.Context, entity *seasonscoredomain.SeasonScore) error {
	return logRepo(logger.ModelSeasonScore, "create", map[string]any{"user_id": entity.UserID, "season_id": entity.SeasonID})
}
func (r *SeasonScoreRepositoryStub) Update(ctx context.Context, entity *seasonscoredomain.SeasonScore) error {
	return logRepo(logger.ModelSeasonScore, "update", map[string]any{"user_id": entity.UserID, "season_id": entity.SeasonID})
}
func (r *SeasonScoreRepositoryStub) FindByUserAndSeason(ctx context.Context, userID valueobject.UserID, seasonID valueobject.SeasonID) (*seasonscoredomain.SeasonScore, error) {
	return nil, logRepo(logger.ModelSeasonScore, "find_by_user_and_season", map[string]any{"user_id": userID, "season_id": seasonID})
}
func (r *SeasonScoreRepositoryStub) FindByUserIDsAndSeason(ctx context.Context, userIDs []valueobject.UserID, seasonID valueobject.SeasonID) ([]seasonscoredomain.SeasonScore, error) {
	return nil, logRepo(logger.ModelSeasonScore, "find_by_user_ids_and_season", map[string]any{"user_count": len(userIDs), "season_id": seasonID})
}
func (r *SeasonScoreRepositoryStub) ResetBySeasonID(ctx context.Context, seasonID valueobject.SeasonID) error {
	return logRepo(logger.ModelSeasonScore, "reset_by_season_id", map[string]any{"season_id": seasonID})
}

type PointTransactionRepositoryStub struct{}

func (r *PointTransactionRepositoryStub) Create(ctx context.Context, entity *pointtransactiondomain.PointTransaction) error {
	return logRepo(logger.ModelPointTransaction, "create", map[string]any{"transaction_id": entity.ID, "user_id": entity.UserID})
}
func (r *PointTransactionRepositoryStub) FindByUserID(ctx context.Context, userID valueobject.UserID, limit int) ([]pointtransactiondomain.PointTransaction, error) {
	return nil, logRepo(logger.ModelPointTransaction, "find_by_user_id", map[string]any{"user_id": userID, "limit": limit})
}

type LeaderboardRepositoryStub struct{}

func (r *LeaderboardRepositoryStub) GetFriendRanking(ctx context.Context, userID valueobject.UserID, seasonID valueobject.SeasonID) ([]leaderboarddomain.FriendEntry, error) {
	return nil, logRepo(logger.ModelLeaderboard, "get_friend_ranking", map[string]any{"user_id": userID, "season_id": seasonID})
}
