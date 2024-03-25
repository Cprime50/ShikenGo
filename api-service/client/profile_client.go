package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/Cprime50/api-service/middleware"
	profilepb "github.com/Cprime50/api-service/pb"
	"github.com/Cprime50/api-service/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	// import profile pb here
)

type ProfileClient struct {
	Client profilepb.ProfileServiceClient
}

var (
	_               = utils.LoadEnv()
	ENV             = utils.MustHaveEnv("ENV")
	PROFILE_SVC_URL = utils.MustHaveEnv("PROFILE_SVC_URL")
	GRPC_PORT       = utils.MustHaveEnv("GRPC_PORT")
	CERT_PATH       = utils.MustHaveEnv("CERT_PATH")
	KEY_PATH        = utils.MustHaveEnv("KEY_PATH")
)

type Profile struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}

func InitProfileServiceClient(c *context.Context) (profilepb.ProfileServiceClient, error) {
	if ENV == "production" {
		certificate, err := tls.LoadX509KeyPair(CERT_PATH, KEY_PATH)
		if err != nil {
			slog.Error("Error loading TLS certificate", "tls.LoadX509KeyPair \n", err)
			return nil, err
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{certificate},
		}
		creds := credentials.NewTLS(tlsConfig)
		conn, err := grpc.DialContext(*c, PROFILE_SVC_URL, []grpc.DialOption{
			grpc.WithTransportCredentials(creds),
			grpc.WithBlock(),
		}...)
		if err != nil {
			return nil, fmt.Errorf("connection to profile gRPC service failed: %v", err)
		}
		return profilepb.NewProfileServiceClient(conn), nil
	} else {
		// Non-production environment, use insecure connection
		conn, err := grpc.DialContext(*c, PROFILE_SVC_URL, []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		}...)

		if err != nil {
			return nil, fmt.Errorf("connection to profile gRPC service failed: %v", err)
		}
		return profilepb.NewProfileServiceClient(conn), nil
	}
}

func CreateUpdateProfile(c *gin.Context, ctx context.Context, method string) (*profilepb.Profile, error) {
	var profile *Profile
	client, err := InitProfileServiceClient(&ctx)
	if err != nil {
		return nil, err
	}

	userValue, exists := c.Get("user")
	if !exists {
		log.Println("User not found in context")
		c.AbortWithStatus(http.StatusUnauthorized)
		return nil, err
	}

	user, ok := userValue.(*middleware.User)
	if !ok || user == nil {
		log.Println("Invalid user data in context")
		c.AbortWithStatus(http.StatusUnauthorized)
		return nil, err
	}

	if err := c.BindJSON(&profile); err != nil {
		log.Print("error binding data for createUpdateProfileRequest \n: Invalid Json format")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Json format"})
		return nil, err
	}

	op := profilepb.Operation_CREATE
	if method == http.MethodPut {
		op = profilepb.Operation_UPDATE
	}

	req := &profilepb.CreateUpdateProfileRequest{
		Operation: op,
		Profile: &profilepb.Profile{
			UserId:   user.UserID,
			Email:    user.Email,
			Username: profile.Username,
			Bio:      profile.Bio,
			Avatar:   profile.Avatar,
		},
	}
	response, err := client.CreateUpdateProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetProfile(ctx context.Context, userID string) (*profilepb.Profile, error) {
	client, err := InitProfileServiceClient(&ctx)
	if err != nil {
		return nil, err
	}

	req := &profilepb.GetProfileRequest{
		UserId: userID,
	}

	res, err := client.GetProfile(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetAllProfiles(ctx context.Context) ([]*profilepb.Profile, error) {
	client, err := InitProfileServiceClient(&ctx)
	if err != nil {
		return nil, err
	}

	req := &profilepb.Empty{}

	stream, err := client.GetAllProfiles(ctx, req)
	if err != nil {
		return nil, err
	}

	// Read profiles from stream
	var profiles []*profilepb.Profile
	for {
		profile, err := stream.Recv()
		if err != nil {
			break
		}
		profiles = append(profiles, profile)
	}

	return profiles, nil
}

func DeleteProfile(ctx context.Context, userID string) error {
	client, err := InitProfileServiceClient(&ctx)
	if err != nil {
		return err
	}

	req := &profilepb.DeleteProfileRequest{
		UserId: userID,
	}

	_, err = client.DeleteProfile(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

// UpdateScore updates the score for the given user ID.
func UpdateScore(ctx context.Context, userID string, score int64) error {
	client, err := InitProfileServiceClient(&ctx)
	if err != nil {
		return err
	}

	req := &profilepb.UpdateScoreRequest{
		UserId: userID,
		Score:  score,
	}

	_, err = client.UpdateScore(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
