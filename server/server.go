package server

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sample-grpc-server/database"
	"sample-grpc-server/pb"
	"sample-grpc-server/service"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/xerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ pb.BackendServiceServer = (*Server)(nil)

type userID int

const (
	KeyUserID userID = iota
)

type Server struct {
	pb.BackendServiceServer

	db   database.Querier
	hash service.Hasher
	auth service.Auther
}

func NewServer(db database.Querier, hash service.Hasher, auth service.Auther) *Server {
	return &Server{
		db:   db,
		hash: hash,
		auth: auth,
	}
}

func (s *Server) HelloWorld(_ context.Context, _ *emptypb.Empty) (*pb.HelloWorldResponse, error) {
	return &pb.HelloWorldResponse{Message: "Hello, World"}, nil
}

func (s *Server) SignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {

	hash, err := s.hash.CreateHash(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "server error")
	}

	params := database.SignUpParams{
		Email:    req.Email,
		Password: hash,
	}

	dbResp, err := s.db.SignUp(ctx, params)
	if err != nil {
		err = errors.Unwrap(err)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062:
				return nil, status.Error(codes.AlreadyExists, "the email is already registered")
			}
		}
		return nil, status.Error(codes.Internal, "database error")
	}

	token, err := s.auth.CreateAccessToken()
	if err != nil {
		return nil, status.Error(codes.Internal, "server error")
	}

	if err := s.db.CreateSession(ctx, database.CreateSessionParams{
		AccessToken: token,
		UserID:      dbResp.UserID,
	}); err != nil {
		return nil, status.Error(codes.Internal, "server error")
	}

	return &pb.SignUpResponse{AccessToken: token}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	params := database.LoginParams{
		Email: req.Email,
	}

	dbResp, err := s.db.Login(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Errorf(codes.Internal, "database error")
	}

	match, err := s.hash.CompareHash(req.Password, dbResp.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "server error")
	}

	if !match {
		return nil, status.Error(codes.InvalidArgument, "invalid email or password")
	}

	token, err := s.auth.CreateAccessToken()
	if err != nil {
		return nil, status.Error(codes.Internal, "server error")
	}

	if err := s.db.CreateSession(ctx, database.CreateSessionParams{
		AccessToken: token,
		UserID:      dbResp.UserID,
	}); err != nil {
		return nil, status.Error(codes.Internal, "server error")
	}

	return &pb.LoginResponse{AccessToken: token}, nil
}

func (s *Server) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {
	userID := extractUserID(ctx)

	params := database.CreateArticleParams{
		UserID:      userID,
		Title:       req.Title,
		Description: convNullString(req.Description),
		Text:        req.Text,
	}

	dbResp, err := s.db.CreateArticle(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, "server error")
	}

	return &pb.CreateArticleResponse{ArticleId: dbResp.ArticleID}, nil
}

func (s *Server) GetArticles(ctx context.Context, _ *emptypb.Empty) (*pb.GetArticlesResponse, error) {
	userID := extractUserID(ctx)

	params := database.GetArticlesParams{UserID: userID}

	dbResp, err := s.db.GetArticles(ctx, params)
	if err != nil {
		return nil, status.Error(codes.Internal, "server error")
	}

	resp := &pb.GetArticlesResponse{}

	for _, article := range dbResp.Articles {
		resp.Articles = append(resp.Articles, &pb.Article{
			ArticleId:   article.ID,
			Title:       article.Title,
			Description: stringPtr(article.Description),
			Text:        article.Text,
			CreatedAt:   timestampPtr(article.CreatedAt),
		})
	}

	return resp, nil
}

func (s *Server) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleResponse, error) {
	userID := extractUserID(ctx)

	params := database.GetArticleParams{
		ArticleID: req.ArticleId,
		UserID:    userID,
	}

	dbResp, err := s.db.GetArticle(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "article not found")
		}
		return nil, status.Error(codes.Internal, "server error")
	}

	resp := &pb.GetArticleResponse{
		Article: &pb.Article{
			ArticleId:   dbResp.Article.ID,
			Title:       dbResp.Article.Title,
			Description: stringPtr(dbResp.Article.Description),
			Text:        dbResp.Article.Text,
			CreatedAt:   timestampPtr(dbResp.Article.CreatedAt),
		},
	}

	return resp, nil
}

func (s *Server) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*emptypb.Empty, error) {
	userID := extractUserID(ctx)

	params := database.UpdateArticleParams{
		ArticleID:   req.GetArticleId(),
		UserID:      userID,
		Title:       req.GetTitle(),
		Description: convNullString(req.Description),
		Text:        req.GetText(),
	}

	err := s.db.UpdateArticle(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "updating article is aborted")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*emptypb.Empty, error) {
	userID := extractUserID(ctx)

	params := database.DeleteArticleParams{
		ArticleID: req.GetArticleId(),
		UserID:    userID,
	}

	err := s.db.DeleteArticle(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "deleting article is aborted")
	}

	return &emptypb.Empty{}, nil
}

func extractUserID(ctx context.Context) int64 {
	IDStr := ctx.Value(KeyUserID)

	id := IDStr.(int64)
	return id
}

func stringPtr[T string | sql.NullString](t T) *string {
	switch v := any(t).(type) {
	case string:
		return &v
	case sql.NullString:
		if !v.Valid {
			return nil
		}

		return &v.String
	default:
		panic(xerrors.Errorf("failed to convert type: %T", v))
	}
}

func convNullString(sp *string) sql.NullString {
	if sp != nil {
		return sql.NullString{
			String: *sp,
			Valid:  true,
		}
	} else {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
}

func timestampPtr(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
