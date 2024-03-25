package src

import (
	"strings"
	"testing"

	pb "github.com/Cprime50/user/profilepb"
)

func TestProfileValidation(t *testing.T) {

	var profiles1 = []pb.Profile{
		{
			UserId:   "test1",
			Email:    "test1@email.com",
			Username: "Username1",
			Bio:      "test bio 1",
			Avatar:   "testavater1",
			Score:    17,
		},
	}
	// Test case 1: Validate the profile input
	err := validateProfile(&profiles1[0])
	if err != nil {
		t.Errorf("validation error: %v", err)
	}

	// Test case 2: Invalidate inconsistent length of inputs
	profiles1[0].Username = strings.Repeat("a", 101)
	profiles1[0].Bio = strings.Repeat("a", 1001)
	profiles1[0].Avatar = strings.Repeat("a", 1001)
	err = validateProfile(&profiles1[0]) // Corrected: using profiles1 instead of profiles
	containsUsername := strings.Contains(err.Error(), "Username") && strings.Contains(err.Error(), "max")
	containsBio := strings.Contains(err.Error(), "Bio") && strings.Contains(err.Error(), "max")
	containsAvatar := strings.Contains(err.Error(), "Avatar") && strings.Contains(err.Error(), "max")
	if !containsUsername || !containsBio || !containsAvatar {
		//t.Log("email", profiles1[0].Email) // Corrected: using profiles1 instead of profiles
		t.Errorf("validation error: %v", err)
	}

	// Test case 3: Invalidate empty required inputs
	profiles1[0].UserId = ""
	profiles1[0].Email = ""
	err = validateProfile(&profiles1[0]) // Corrected: using profiles1 instead of profiles
	containsUserId := strings.Contains(err.Error(), "Email") && strings.Contains(err.Error(), "required")
	containsEmail := strings.Contains(err.Error(), "UserId") && strings.Contains(err.Error(), "required")
	if containsUserId || containsEmail {
		t.Errorf("validation failed, required fields are empty but validated: %v", err)
	}
}

// func TestScoreValidation(t *testing.T) {
// 	ps := &ProfileService{}
// 	var score = int64(12)
// 	err := ps.validateScore(score)
// 	if err != nil {
// 		t.Errorf("validation error: %v", err)
// 	}

// 	score = -3
// 	err = ps.validateScore(score)
// 	if err == nil {
// 		t.Errorf("validation error: %v", err)
// 	}
// }
