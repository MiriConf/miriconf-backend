package model

// User model info
// @Description User account information
type User struct {
	// ID of user account
	ID int `json:"id" example:"10" format:"int64"`
	// Fullname string of user
	Fullname string `json:"fullname" example:"John Smith" format:"string"`
	// Denotes what teams the user is a member of
	MemberOfTeams []int `json:"member_of_teams" example:"2, 5, 7" format:"int64"`
	// Declares the user account role the user has
	// within the dashboard when logged in
	Role string `json:"role" example:"admin" format:"string"`
	// Determines whether or not the user is allowed to log into the dashboard
	LoginAllowed bool `json:"login_allowed" example:"No" format:"bool"`
	// Creation timestamp
	CreatedAt string `json:"created_at" example:"08-23-2022" format:"string"`
}
