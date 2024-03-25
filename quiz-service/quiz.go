package main

import (
	"context"
	"fmt"
	"math/rand"

	db "github.com/Cprime50/quiz/db"
	pb "github.com/Cprime50/quiz/quizpb"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *server) GetQuiz(ctx context.Context, req *pb.GetQuizRequest) (*pb.GetQuizResponse, error) {
	userScore := req.Score
	// Get quizzes based on user score
	rows, err := db.Db.Query("SELECT id, japanese FROM quiz WHERE id > ? LIMIT 20", userScore)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	var quizzes []*pb.Quiz
	for rows.Next() {,
		var quiz pb.Quiz
		if err := rows.Scan(&quiz.Id, &quiz.Japanese, &quiz.Pronounce, &quiz.English); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		quizzes = append(quizzes, &quiz)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}
	// Shuffle options for each quiz
	for _, quiz := range quizzes {
		shuffleOptions(quiz)
	}

	return &pb.GetQuizResponse{Quizes: quizzes}, nil
}

func shuffleOptions(quiz *pb.Quiz) {
	// Randomly shuffle the options (English translations)
	options := []string{quiz.English, "Option 1", "Option 2"}
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})
	quiz.English = options[0]
}

func (s *server) Answer(ctx context.Context, req *pb.AnswerRequest) (*pb.AnswerResponse, error) {
	// Check user authentication

	// Calculate score
	score := int64(0)
	for _, ans := range req.Answers {
		// Check if the answer is correct
		var correctAnswer string
		err := db.Db.QueryRow("SELECT english FROM quiz WHERE id = ?", ans.Id).Scan(&correctAnswer)
		if err != nil {
			return nil, fmt.Errorf("error querying database: %v", err)
		}
		if ans.UserAnswer == correctAnswer {
			score++
		}
	}

	// Check if user can proceed to the next quiz
	nextAllowed := false
	if score >= 16 {
		nextAllowed = true
		// Update user score in Firebase Custom Claims
	}

	return &pb.AnswerResponse{Score: score, NextAllowed: nextAllowed}, nil
}

func (s *server) CreateQuiz(ctx context.Context, req *pb.CreateQuizRequest) (*emptypb.Empty, error) {

	// Insert quizzes into the database
	stmt, err := db.Db.Prepare("INSERT INTO quiz (japanese, pronounce, english) VALUES (?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, quiz := range req.Quizes {
		_, err := stmt.Exec(quiz.Japanese, quiz.Pronounce, quiz.English)
		if err != nil {
			return nil, fmt.Errorf("error inserting data: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *server) UpdateQuiz(ctx context.Context, req *pb.UpdateQuizRequest) (*emptypb.Empty, error) {
	// Update quizzes in the database
	stmt, err := db.Db.Prepare("UPDATE quiz SET japanese=?, pronounce=?, english=? WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, quiz := range req.Quizes {
		_, err := stmt.Exec(quiz.Japanese, quiz.Pronounce, quiz.English, quiz.Id)
		if err != nil {
			return nil, fmt.Errorf("error updating data: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *server) DeleteQuiz(ctx context.Context, req *pb.DeleteQuizRequest) (*emptypb.Empty, error) {

	// Delete quizzes from the database
	stmt, err := db.Db.Prepare("DELETE FROM quiz WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	for _, id := range req.QuizIds {
		_, err := stmt.Exec(id)
		if err != nil {
			return nil, fmt.Errorf("error deleting data: %v", err)
		}
	}

	return &emptypb.Empty{}, nil
}

func (s *server) GetAllQuizzes(ctx context.Context, req *pb.GetAllQuizzesRequest) (*pb.GetAllQuizzesResponse, error) {
	// Retrieve all quizzes from the database
	rows, err := db.Db.Query("SELECT id, japanese, pronounce, english FROM quiz")
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	var quizzes []*pb.Quiz
	for rows.Next() {
		var quiz pb.Quiz
		if err := rows.Scan(&quiz.Id, &quiz.Japanese, &quiz.Pronounce, &quiz.English); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		quizzes = append(quizzes, &quiz)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return &pb.GetAllQuizzesResponse{Quizes: quizzes}, nil
}
