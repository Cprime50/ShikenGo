package src

import (
	pb "github.com/Cprime50/quiz/quizpb"
)

type Server struct {
	pb.UnimplementedQuizServiceServer
}
