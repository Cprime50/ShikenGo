package src

import (
	"context"
	"testing"

	pb "github.com/Cprime50/user/profilepb"
	"google.golang.org/grpc"
)

func TestCreateUpdateProfile(t *testing.T) {

	clearProfiles()
	s := &Server{}

	// Test case 1: CREATE operation
	reqCreate := &pb.CreateUpdateProfileRequest{
		Operation: pb.Operation_CREATE,
		Profile:   &profiles[0],
	}
	profile, err := s.CreateUpdateProfile(context.Background(), reqCreate)
	if err != nil {
		t.Errorf("CreateUpdateProfile() error = %v", err)
		return
	}
	if profile == nil {
		t.Error("CreateUpdateProfile() response is nil")
		return
	}

	equal := profiles[0].Username == profile.Username &&
		profiles[0].Bio == profile.Bio &&
		profiles[0].Avatar == profile.Avatar &&
		profiles[0].Email == profile.Email

	if !equal {
		t.Errorf("createUpdateProfile() error: not equal")
	}
	if profile.Id == "" {
		t.Errorf("createUpdateProfile() error: id field not auto generated")
	}
	if profile.Score != 0 {
		t.Errorf("createUpdateProfile() error: score should be 0")
	}

	// Test case 2: Create operation (without username)
	reqCreate1 := &pb.CreateUpdateProfileRequest{
		Operation: pb.Operation_CREATE,
		Profile: &pb.Profile{
			UserId:   "test4",
			Email:    "test4@email.com",
			Username: "",
			Bio:      "test bio 4",
			Avatar:   "testavater4",
			Score:    30,
		},
	}
	profile1, err := s.CreateUpdateProfile(context.Background(), reqCreate1)
	if err != nil {
		t.Errorf("CreateUpdateProfile() error = %v", err)
		return
	}
	if profile1 == nil {
		t.Error("CreateUpdateProfile() response is nil")
		return
	}
	if profile1.Username == "" {
		t.Errorf("createUpdateProfile() error: failed to generate username")
		return
	}

	// Test case 3: UPDATE operation
	newProfile := pb.Profile{
		Id:       profile.Id,
		Email:    profile.Email,
		UserId:   profile.UserId,
		Username: "New username",
		Bio:      "New bio",
		Avatar:   "New avatar",
	}
	reqUpdate := &pb.CreateUpdateProfileRequest{
		Operation: pb.Operation_UPDATE,
		Profile:   &newProfile,
	}
	profileUpdate, err := s.CreateUpdateProfile(context.Background(), reqUpdate)
	if err != nil {
		t.Errorf("CreateUpdateProfile() error = %v", err)
		return
	}
	if profileUpdate == nil {
		t.Error("CreateUpdateProfile() response is nil")
		return
	}
	//Check data is equal

	//Check update cannot change score
}

func TestGetProfile(t *testing.T) {
	clearProfiles()
	s := &Server{}

	reqCreate := &pb.CreateUpdateProfileRequest{
		Operation: pb.Operation_CREATE,
		Profile:   &profiles[0],
	}

	_, err := s.CreateUpdateProfile(context.Background(), reqCreate)
	if err != nil {
		t.Errorf("CreateUpdateProfile() error = %v", err)
		return
	}
	req := &pb.GetProfileRequest{
		UserId: "test1",
	}
	profile, err := s.GetProfile(context.Background(), req)
	if err != nil {
		t.Errorf("GetProfile() error = %v", err)
		return
	}
	if profile == nil {
		t.Error("GetProfile() failed to return a profile")
		return
	}

	// check is  profile data is equal
}

// Mock for testing stream gRPC
type mockProfileService_GetAllProfilesServer struct {
	grpc.ServerStream
	Results []*pb.Profile
}

func (_m *mockProfileService_GetAllProfilesServer) Send(p *pb.Profile) error {
	_m.Results = append(_m.Results, p)
	return nil
}

func TestGetAllProfiles(t *testing.T) {
	clearProfiles()
	s := &Server{}

	mock := &mockProfileService_GetAllProfilesServer{}

	reqCreate := &pb.CreateUpdateProfileRequest{
		Operation: pb.Operation_CREATE,
		Profile:   &profiles[0],
	}
	_, _ = s.CreateUpdateProfile(context.Background(), reqCreate)

	reqCreate1 := &pb.CreateUpdateProfileRequest{
		Operation: pb.Operation_CREATE,
		Profile:   &profiles[1],
	}
	_, _ = s.CreateUpdateProfile(context.Background(), reqCreate1)

	reqCreate2 := &pb.CreateUpdateProfileRequest{
		Operation: pb.Operation_CREATE,
		Profile:   &profiles[2],
	}
	_, _ = s.CreateUpdateProfile(context.Background(), reqCreate2)

	err := s.GetAllProfiles(&pb.Empty{}, mock)
	if err != nil {
		t.Fatalf("GetAllProfiles returned error: %v", err)
	}

	// Verify the results in the mock
	if len(mock.Results) != 3 {
		t.Errorf("Expected 3 profile, got %d", len(mock.Results))
	}
	// Check if results isequal
	for i, result := range mock.Results {
		if result.UserId != profiles[i].UserId {
			t.Errorf("Profile at index %d: expected user ID %s, got %s", i, profiles[i].UserId, result.UserId)
		}
	}
}

func TestDeleteProfile(t *testing.T) {
	clearProfiles()
	s := &Server{}
	reqCreate := &pb.CreateUpdateProfileRequest{
		Operation: pb.Operation_CREATE,
		Profile:   &profiles[0],
	}
	_, err := s.CreateUpdateProfile(context.Background(), reqCreate)
	if err != nil {
		t.Errorf("CreateUpdateProfile() error = %v", err)
		return
	}

	req := &pb.DeleteProfileRequest{
		UserId: "test1",
	}
	resp, err := s.DeleteProfile(context.Background(), req)
	if err != nil {
		t.Errorf("DeleteProfile() error = %v", err)
		return
	}
	if resp == nil {
		t.Error("DeleteProfile() response is nil")
		return
	}

	reqProfile := &pb.GetProfileRequest{
		UserId: "test1",
	}
	profile, err := s.GetProfile(context.Background(), reqProfile)
	if err == nil {
		t.Errorf("failed to delete profile %v %s", err, profile)
		return
	}
}

func TestUpdateScore(t *testing.T) {
	clearProfiles()
	s := &Server{}
	reqCreate := &pb.CreateUpdateProfileRequest{
		Operation: pb.Operation_CREATE,
		Profile:   &profiles[0],
	}
	_, err := s.CreateUpdateProfile(context.Background(), reqCreate)
	if err != nil {
		t.Errorf("CreateUpdateProfile() error = %v", err)
		return
	}
	req := &pb.UpdateScoreRequest{
		UserId: "test1",
		Score:  100,
	}
	resp, err := s.UpdateScore(context.Background(), req)
	if err != nil {
		t.Errorf("UpdateScore() error = %v", err)
		return
	}
	if resp == nil {
		t.Error("UpdateScore() response is nil")
		return
	}
}
