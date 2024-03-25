package src

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	pb "github.com/Cprime50/user/profilepb"
	"github.com/Cprime50/user/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedProfileServiceServer
}

func (s *Server) CreateUpdateProfile(ctx context.Context, req *pb.CreateUpdateProfileRequest) (*pb.Profile, error) {
	start := time.Now()
	if err := validateProfile(req.Profile); err != nil {
		log.Printf("CreateProfile error: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "profile validation error: %v", err)
	}

	switch req.Operation {
	case pb.Operation_CREATE:
		existingProfile, _ := getProfileByUserId(req.Profile.UserId)
		if existingProfile != nil {
			log.Printf("CreateProfile error: profile already exists for user ID: %s", req.Profile.UserId)
			return nil, status.Errorf(codes.AlreadyExists, "profile already exists")
		}

		if req.Profile.Username == "" {
			// Generate username if not provided
			username, err := utils.GenerateUsername(req.Profile.Email)
			if err != nil {
				log.Printf("CreateProfile error: generating username failed: %v", err)
				return nil, status.Errorf(codes.Internal, "error generating username: %v", err)
			}
			req.Profile.Username = username
		}

		err := createProfile(req.Profile)
		if err != nil {
			log.Printf("CreateProfile error: creating user profile failed: %v", err)
			return nil, status.Errorf(codes.Internal, "error creating user profile: %v", err)
		}

	case pb.Operation_UPDATE:
		existingProfile, err := getProfileByUserId(req.Profile.UserId)
		if err != nil {
			log.Printf("UpdateProfile error: checking existing profile failed: %v", err)
			return nil, status.Errorf(codes.Internal, "error checking existing profile: %v", err)
		}
		if existingProfile == nil {
			log.Printf("UpdateProfile error: profile not found for user ID: %s", req.Profile.UserId)
			return nil, status.Errorf(codes.NotFound, "profile not found for user ID: %s", req.Profile.UserId)
		}
		if req.Profile.Username != "" {
			existingProfile.Username = req.Profile.Username
		}
		if req.Profile.Bio != "" {
			existingProfile.Bio = req.Profile.Bio
		}
		if req.Profile.Avatar != "" {
			existingProfile.Avatar = req.Profile.Avatar
		}

		err = updateProfile(existingProfile)
		if err != nil {
			log.Printf("UpdateProfile error: updating user profile failed: %v", err)
			return nil, status.Errorf(codes.Internal, "error updating user profile: %v", err)
		}
	default:
		log.Printf("CreateUpdateProfile error: unknown operation: %v", req.Operation)
		return nil, status.Errorf(codes.InvalidArgument, "unknown operation: %v", req.Operation)
	}

	profile, err := getProfileByUserId(req.Profile.UserId)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			log.Printf("CreateUpdateProfile error: profile not found for user ID: %s", req.Profile.UserId)
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		log.Printf("CreateUpdateProfile: failed to get profile: %s", err)
		return nil, status.Errorf(codes.Internal, "failed to get profile: %s", err)
	}
	log.Printf("Successesfully %s"+"D profile", req.Operation)
	slog.Info(fmt.Sprintf("operation: %s, time: %s", req.Operation, time.Since(start)))
	return profile, nil
}

func (s *Server) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.Profile, error) {
	start := time.Now()

	profile, err := getProfileByUserId(req.UserId)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			log.Printf("GetProfile error: profile not found for user ID: %s", req.UserId)
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		log.Printf("GetProfile error: failed to get profile: %s", err)
		return nil, status.Errorf(codes.Internal, "failed to get profile: %s", err)
	}
	log.Printf("GetProfile successful: username=%s", profile.Username)
	slog.Info("GetProfileByUserId", "time", time.Since(start))
	return profile, nil
}

func (s *Server) GetAllProfiles(req *pb.Empty, stream pb.ProfileService_GetAllProfilesServer) error {
	start := time.Now()

	profiles, err := selectProfiles()
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			log.Printf("GetAllProfiles error: profiles not found")
			return status.Errorf(codes.NotFound, err.Error())
		}
		log.Printf("GetAllProfiles error: failed to get profiles: %s", err)
		return status.Errorf(codes.Internal, "failed to get profiles: %s", err)
	}
	// Stream profiles to the client
	for _, profile := range profiles {
		// Send the profile to the client stream
		if err := stream.Send(profile); err != nil {
			log.Printf("GetAllProfiles error: failed to send profiles to client: %s", err)
			return status.Errorf(codes.Internal, "failed to send profiles to client: %s", err)
		}
	}
	log.Printf("GetAllProfiles successful: sent %d profiles", len(profiles))
	slog.Info("GetAllProfiles", "time", time.Since(start))
	return nil
}

func (s *Server) DeleteProfile(ctx context.Context, req *pb.DeleteProfileRequest) (*pb.Empty, error) {
	start := time.Now()

	err := deleteProfileByUserId(req.UserId)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			log.Printf("DeleteProfile error: profile not found for user ID: %s", req.UserId)
			return nil, status.Errorf(codes.NotFound, "profile not found")
		}
		log.Printf("DeleteProfile error: failed to delete profile: %v", err)
		return nil, status.Errorf(codes.Internal, "error deleting profile: %v", err)
	}
	log.Printf("DeleteProfile successful: profile deleted for user ID: %s", req.UserId)
	slog.Info("DeleteProfile", "time", time.Since(start))
	return &pb.Empty{}, nil
}

func (s *Server) UpdateScore(ctx context.Context, req *pb.UpdateScoreRequest) (*pb.Empty, error) {
	start := time.Now()

	err := updateScore(req.UserId, req.Score)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			log.Printf("UpdateScore error: profile not found for user ID: %s", req.UserId)
			return nil, status.Errorf(codes.NotFound, "profile not found")
		}
		log.Printf("UpdateScore error: failed to update score: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update score: %v", err)
	}
	log.Printf("UpdateScore successful: score updated for user ID: %s", req.UserId)
	slog.Info("UpdateScore", "time", time.Since(start))
	return &pb.Empty{}, nil
}
