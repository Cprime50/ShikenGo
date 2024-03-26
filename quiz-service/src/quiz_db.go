package src

import (
	"errors"

	pb "github.com/Cprime50/quiz/quizpb"
)

var (
	ErrProfileNotFound           = errors.New("not found")
	ErrDuplicateEntry            = errors.New("duplicate entry")
	ErrForeignKeyViolation       = errors.New("foreign key violation")
	ErrUniqueConstraintViolation = errors.New("unique constraint violation")
)

func getQuizByUserId(userId string) (*pb.Quiz, error) {

}

func createQuiz(q *pb.Quiz) error {

}

func updateQuiz(q *pb.Quiz) error {

}

func selectLeaderBoard() ([]*pb.LeaderBoard, error){

}

