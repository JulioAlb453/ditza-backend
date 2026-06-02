package httpserver

import (
	"net/http"

	cosmetichttp "ditza/internal/cosmetic/infrastructure/http"
	friendshiphttp "ditza/internal/friendship/infrastructure/http"
	habitcompletionhttp "ditza/internal/habit-completion/infrastructure/http"
	habithttp "ditza/internal/habit/infrastructure/http"
	leaderboardhttp "ditza/internal/leaderboard/infrastructure/http"
	pethttp "ditza/internal/pet/infrastructure/http"
	seasonhttp "ditza/internal/season/infrastructure/http"
	usercosmetichttp "ditza/internal/user-cosmetic/infrastructure/http"
	userhttp "ditza/internal/user/infrastructure/http"
)

type Controllers struct {
	User            *userhttp.Controller
	Habit           *habithttp.Controller
	HabitCompletion *habitcompletionhttp.Controller
	Cosmetic        *cosmetichttp.Controller
	UserCosmetic    *usercosmetichttp.Controller
	Friendship      *friendshiphttp.Controller
	Leaderboard     *leaderboardhttp.Controller
	Season          *seasonhttp.Controller
	Pet             *pethttp.Controller
}

func NewRouter(controllers Controllers) *http.ServeMux {
	mux := http.NewServeMux()

	// Auth/User
	mux.HandleFunc("POST /auth/register", controllers.User.Register)
	mux.HandleFunc("POST /auth/login", controllers.User.Login)
	mux.HandleFunc("GET /me", controllers.User.GetMe)
	mux.HandleFunc("GET /users/search", controllers.User.Search)

	// Pet
	mux.HandleFunc("GET /pet", controllers.Pet.Get)
	mux.HandleFunc("PATCH /pet/equip", controllers.Pet.Equip)

	// Habits
	mux.HandleFunc("GET /habits", controllers.Habit.List)
	mux.HandleFunc("POST /habits", controllers.Habit.Create)
	mux.HandleFunc("DELETE /habits/{id}", controllers.Habit.Deactivate)
	mux.HandleFunc("PATCH /habits/{id}/complete", controllers.HabitCompletion.Complete)

	// Shop / Cosmetics
	mux.HandleFunc("GET /shop/items", controllers.Cosmetic.ListActive)
	mux.HandleFunc("POST /shop/buy", controllers.UserCosmetic.Buy)
	mux.HandleFunc("GET /shop/inventory", controllers.UserCosmetic.ListInventory)

	// Friendships
	mux.HandleFunc("POST /friends/request", controllers.Friendship.SendRequest)
	mux.HandleFunc("PATCH /friends/{id}/accept", controllers.Friendship.Accept)
	mux.HandleFunc("PATCH /friends/{id}/reject", controllers.Friendship.Reject)
	mux.HandleFunc("GET /friends", controllers.Friendship.ListFriends)
	mux.HandleFunc("GET /friends/pending", controllers.Friendship.ListPending)

	// Ranking / Season
	mux.HandleFunc("GET /leaderboard/friends", controllers.Leaderboard.GetFriendRanking)
	mux.HandleFunc("GET /seasons/current", controllers.Season.GetActive)

	return mux
}
