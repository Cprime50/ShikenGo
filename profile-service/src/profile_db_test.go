package src

import (
	"log"
	"testing"

	"github.com/Cprime50/user/db"
	pb "github.com/Cprime50/user/profilepb"
)

var profiles = []pb.Profile{
	{
		UserId:   "test1",
		Email:    "test1@email.com",
		Username: "Username1",
		Bio:      "test bio 1",
		Avatar:   "testavater1",
		Score:    17,
	},
	{
		UserId:   "test2",
		Email:    "test2@email.com",
		Username: "Username2",
		Bio:      "test bio 2",
		Avatar:   "testavater2",
		Score:    20,
	},
	{
		UserId:   "test3",
		Email:    "test3@email.com",
		Username: "Username3",
		Bio:      "test bio 3",
		Avatar:   "testavater3",
		Score:    30,
	},
}

func clearProfiles() {
	_, err := db.Db.Exec("delete from profiles")
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetProfileByUserId(t *testing.T) {
	clearProfiles()
	// Test case 1: Select a profile by id
	profile := &profiles[0]
	err := createProfile(profile)
	if err != nil {
		t.Errorf("createProfile error: %v", err)
	}
	gottenProfile, err := getProfileByUserId(profile.UserId)
	if err != nil {
		t.Errorf("getProfileByUserId error: %v", err)
	}
	if gottenProfile.UserId != profile.UserId {
		t.Errorf("getProfileByUserId error: not equal")
	}

	// Test case 2: Select a profile by id that does not exist
	_, err = getProfileByUserId("not_exist")
	if err == nil {
		t.Errorf("getProfileByUserId error: %v", err)
	}
}

func TestInsertProfile(t *testing.T) {
	clearProfiles()
	profile := &profiles[0]
	// Test case 1: Insert a valid profile
	err := createProfile(profile)
	if err != nil {
		t.Errorf("createProfile error: %v", err)
	}
	gottenProfile, _ := getProfileByUserId(profile.UserId)
	equal := gottenProfile.Username == profile.Username &&
		gottenProfile.Bio == profile.Bio &&
		gottenProfile.Avatar == profile.Avatar &&
		gottenProfile.Email == profile.Email

	if !equal {
		t.Errorf("createProfile error: not equal")
	}
	if gottenProfile.Id == "" {
		t.Errorf("createProfile error: id field not auto generated")
	}

	// Test case 2: Insert a second valid profile
	err = createProfile(&profiles[1])
	if err != nil {
		t.Errorf("createProfile error: %v", err)
	}
	gottenProfile, _ = getProfileByUserId(profiles[1].UserId)
	if gottenProfile.Username != profiles[1].Username {
		t.Errorf("createProfile error: not equal")
	}

	// Test case 3: Insert a profile that already exist
	err = createProfile(&profiles[1])
	if err == nil {
		t.Errorf("creating duplicate profile error: %v", err)
	}
}

func TestUpdateProfile(t *testing.T) {
	clearProfiles()
	// Test case 1: Update a valid profile
	err := createProfile(&profiles[0])
	if err != nil {
		t.Errorf("createProfile error: %v", err)
	}
	profile, _ := getProfileByUserId(profiles[0].UserId)
	newProfile := pb.Profile{
		Id:       profile.Id,
		Email:    profile.Email,
		UserId:   profile.UserId,
		Username: "New username",
		Bio:      "New bio",
		Avatar:   "New avatar",
	}
	err = updateProfile(&newProfile)
	if err != nil {
		t.Errorf("updateProfile error: %v", err)
	}
	profile, _ = getProfileByUserId(profiles[0].UserId)
	if profile.Username != newProfile.Username {
		t.Errorf("updateProfile error: not equal")
	}

	// Test case 2: Update a profile that does not exist
	newProfile = pb.Profile{
		Id: "not_exist",
	}
	err = updateProfile(&newProfile)
	if err == nil {
		t.Errorf("updateProfile error: %v", err)
	}
}

func TestSelectProfiles(t *testing.T) {
	clearProfiles()

	// Insert profiles
	_ = createProfile(&profiles[0])
	_ = createProfile(&profiles[1])
	_ = createProfile(&profiles[2])

	// Get profiles
	gottenProfiles, err := selectProfiles()
	if err != nil {
		t.Errorf("Error selecting profiles: %v", err)
		return
	}

	// Check if the number of retrieved profiles matches the expected number
	expectedProfiles := len(profiles)
	if len(gottenProfiles) != expectedProfiles {
		t.Errorf("Expected %d profiles, got %d", expectedProfiles, len(gottenProfiles))
	}

	// Check if the retrieved profiles match the expected profiles
	for i := 0; i < len(gottenProfiles); i++ {
		if gottenProfiles[i].UserId != profiles[i].UserId {
			t.Errorf("Expected profile %d to have UserId %s, got %s", i, profiles[i].UserId, gottenProfiles[i].UserId)
		}
	}
}
