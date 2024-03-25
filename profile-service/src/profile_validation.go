package src

import (
	"fmt"
	"regexp"

	pb "github.com/Cprime50/user/profilepb"
	"github.com/Cprime50/user/utils"
)

func validateProfile(in *pb.Profile) error {
	rules := map[string]string{
		"UserId":   "required",
		"Email":    "required",
		"Username": "max=100",
		"Bio":      "max=1000",
		"Avatar":   "max=1000",
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(in.Email) {
		return fmt.Errorf("invalid email format")
	}
	err := utils.ValidateStruct[pb.Profile](rules, pb.Profile{}, in)
	if err != nil {
		return fmt.Errorf("validateProfile error: %w", err)
	}
	return nil
}

// func (ps *ProfileService) validateScore(score int64) error {
// 	if score < 0 {
// 		return fmt.Errorf("must be positive int64")
// 	}
// 	return nil
// }
