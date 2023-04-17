package server

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"sample-grpc-server/database/model"
	"testing"
	"time"

	"sample-grpc-server/database"
	mock_database "sample-grpc-server/database/mock"
	"sample-grpc-server/pb"
	"sample-grpc-server/service"
	mock_service "sample-grpc-server/service/mock"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServer_HelloWorld(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expect := &pb.HelloWorldResponse{Message: "Hello, World"}

	resp, err := callHello(&emptypb.Empty{}, nil, nil, nil)

	if err != nil {
		t.Errorf("err should be nil: %v", err)
	}

	if !reflect.DeepEqual(resp, expect) {
		t.Errorf("Expect: %v, Got: %v", expect, resp)
	}
}

func TestServer_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.SignUpRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	t.Run("リクエスト成功", func(t *testing.T) {
		hash := mock_service.NewMockHasher(ctrl)
		hash.EXPECT().CreateHash(gomock.Any()).Return("hash", nil)

		db := mock_database.NewMockQuerier(ctrl)
		db.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&database.SignUpResult{UserID: 1}, nil)
		db.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(nil)

		auth := mock_service.NewMockAuther(ctrl)
		auth.EXPECT().CreateAccessToken().Return("access_token", nil)

		expect := &pb.SignUpResponse{AccessToken: "access_token"}

		got, err := callSignUp(req, db, hash, auth)

		if err != nil {
			t.Errorf("err should be nil: %v", err)
		}

		if got == nil {
			t.Error("response should not be nil")
		}

		if !reflect.DeepEqual(got, expect) {
			t.Errorf("Expect: %v, Got: %v", expect, got)
		}
	})

	t.Run("ハッシュ化エラー", func(t *testing.T) {
		hash := mock_service.NewMockHasher(ctrl)
		hash.EXPECT().CreateHash(gomock.Any()).Return("", errors.New("some error"))

		_, err := callSignUp(req, nil, hash, nil)

		if err == nil {
			t.Error("err should not be nil")
		}

		if s, ok := status.FromError(err); ok {
			if s.Code() != codes.Internal {
				t.Errorf("Expect: %v, Got: %v", codes.Internal, s.Code())
			}
		}
	})

	t.Run("データベースエラー", func(t *testing.T) {
		t.Run("データベースエラー", func(t *testing.T) {
			hash := mock_service.NewMockHasher(ctrl)
			hash.EXPECT().CreateHash(gomock.Any()).Return("hash", nil)

			db := mock_database.NewMockQuerier(ctrl)
			db.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))

			_, err := callSignUp(req, db, hash, nil)

			if err == nil {
				t.Errorf("err should not be nil: %v", err)
			}

			if s, ok := status.FromError(err); ok {
				if s.Code() != codes.Internal {
					t.Errorf("Expect: %v, Got: %v", codes.Internal, s.Code())
				}
			}
		})

		t.Run("セッション作成エラー", func(t *testing.T) {
			hash := mock_service.NewMockHasher(ctrl)
			hash.EXPECT().CreateHash(gomock.Any()).Return("hash", nil)

			db := mock_database.NewMockQuerier(ctrl)
			db.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&database.SignUpResult{UserID: 1}, nil)
			db.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(errors.New("some error"))

			auth := mock_service.NewMockAuther(ctrl)
			auth.EXPECT().CreateAccessToken().Return("access_token", nil)

			_, err := callSignUp(req, db, hash, auth)

			if err == nil {
				t.Errorf("err should not be nil: %v", err)
			}

			if s, ok := status.FromError(err); ok {
				if s.Code() != codes.Internal {
					t.Errorf("Expect: %v, Got: %v", codes.Internal, s.Code())
				}
			}
		})
	})

	t.Run("アクセストークン生成エラー", func(t *testing.T) {
		hash := mock_service.NewMockHasher(ctrl)
		hash.EXPECT().CreateHash(gomock.Any()).Return("hash", nil)

		db := mock_database.NewMockQuerier(ctrl)
		db.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(&database.SignUpResult{UserID: 1}, nil)

		auth := mock_service.NewMockAuther(ctrl)
		auth.EXPECT().CreateAccessToken().Return("", errors.New("some error"))

		_, err := callSignUp(req, db, hash, auth)

		if err == nil {
			t.Error("err should not be nil")
		}

		if s, ok := status.FromError(err); ok {
			if s.Code() != codes.Internal {
				t.Errorf("Expect: %v, Got: %v", codes.Internal, s.Code())
			}
		}
	})
}

func TestServer_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	t.Run("リクエスト成功", func(t *testing.T) {
		db := mock_database.NewMockQuerier(ctrl)
		db.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&database.LoginResult{
			UserID:   1,
			Password: "password",
		}, nil)

		hash := mock_service.NewMockHasher(ctrl)
		hash.EXPECT().CompareHash(gomock.Any(), gomock.Any()).Return(true, nil)

		auth := mock_service.NewMockAuther(ctrl)
		auth.EXPECT().CreateAccessToken().Return("access_token", nil)

		db.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(nil)

		expect := &pb.LoginResponse{
			AccessToken: "access_token",
		}

		resp, err := callLogin(req, db, hash, auth)

		if err != nil {
			t.Errorf("err should be nil: %v", err)
		}

		if resp == nil {
			t.Error("response should not be nil")
		}

		if !reflect.DeepEqual(resp, expect) {
			t.Errorf("Expect: %v, Got: %v", expect, resp)
		}
	})

	t.Run("データベースエラー", func(t *testing.T) {
		t.Run("emailまたはpasswordが不一致", func(t *testing.T) {
			db := mock_database.NewMockQuerier(ctrl)
			db.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows)

			_, err := callLogin(req, db, nil, nil)

			if err == nil {
				t.Error("err should not be nil")
			}

			if s, ok := status.FromError(err); ok {
				if s.Code() != codes.InvalidArgument {
					t.Errorf("Expect: %v, Got: %v", codes.InvalidArgument, s.Code())
				}
			}
		})

		t.Run("サーバーエラー", func(t *testing.T) {
			db := mock_database.NewMockQuerier(ctrl)
			db.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))

			_, err := callLogin(req, db, nil, nil)

			if err == nil {
				t.Error("err should not be nil")
			}

			if s, ok := status.FromError(err); ok {
				if s.Code() != codes.Internal {
					t.Errorf("Expect: %v, Got: %v", codes.Internal, s.Code())
				}
			}
		})

		t.Run("セッション作成エラー", func(t *testing.T) {
			db := mock_database.NewMockQuerier(ctrl)
			db.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&database.LoginResult{
				UserID:   1,
				Password: "password",
			}, nil)

			hash := mock_service.NewMockHasher(ctrl)
			hash.EXPECT().CompareHash(gomock.Any(), gomock.Any()).Return(true, nil)

			auth := mock_service.NewMockAuther(ctrl)
			auth.EXPECT().CreateAccessToken().Return("access_token", nil)

			db.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Return(errors.New("some error"))

			_, err := callLogin(req, db, hash, auth)

			if err == nil {
				t.Errorf("err should not be nil")
			}

			if s, ok := status.FromError(err); ok {
				if s.Code() != codes.Internal {
					t.Errorf("Expect: %v, Got: %v", codes.Internal, s.Code())
				}
			}
		})
	})

	t.Run("ハッシュ", func(t *testing.T) {
		t.Run("比較処理失敗", func(t *testing.T) {
			db := mock_database.NewMockQuerier(ctrl)
			db.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&database.LoginResult{
				UserID:   1,
				Password: "password",
			}, nil)

			hash := mock_service.NewMockHasher(ctrl)
			hash.EXPECT().CompareHash(gomock.Any(), gomock.Any()).Return(false, errors.New("some error"))

			_, err := callLogin(req, db, hash, nil)

			if err == nil {
				t.Error("err should not be nil")
			}

			if s, ok := status.FromError(err); ok {
				if s.Code() != codes.Internal {
					t.Errorf("Expect: %v, Got: %v", codes.Internal, s.Code())
				}
			}
		})

		t.Run("不一致の場合", func(t *testing.T) {
			db := mock_database.NewMockQuerier(ctrl)
			db.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&database.LoginResult{
				UserID:   1,
				Password: "password",
			}, nil)

			hash := mock_service.NewMockHasher(ctrl)
			hash.EXPECT().CompareHash(gomock.Any(), gomock.Any()).Return(false, nil)

			_, err := callLogin(req, db, hash, nil)

			if err == nil {
				t.Error("err should not be nil")
			}

			if s, ok := status.FromError(err); ok {
				if s.Code() != codes.InvalidArgument {
					t.Errorf("Expect: %v, Got: %v", codes.InvalidArgument, s.Code())
				}
			}
		})
	})

	t.Run("アクセストークン生成エラー", func(t *testing.T) {
		db := mock_database.NewMockQuerier(ctrl)
		db.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&database.LoginResult{
			UserID:   1,
			Password: "password",
		}, nil)

		hash := mock_service.NewMockHasher(ctrl)
		hash.EXPECT().CompareHash(gomock.Any(), gomock.Any()).Return(true, nil)

		auth := mock_service.NewMockAuther(ctrl)
		auth.EXPECT().CreateAccessToken().Return("", errors.New("some error"))

		_, err := callLogin(req, db, hash, auth)

		if err == nil {
			t.Error("err should not be nil")
		}

		if s, ok := status.FromError(err); ok {
			if s.Code() != codes.Internal {
				t.Errorf("Expect: %v, Got: %v", codes.Internal, s.Code())
			}
		}
	})
}

func TestServer_CreateArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.CreateArticleRequest{
		Title:       "test_title",
		Description: stringPtr("test_description"),
		Text:        "test_text",
	}

	t.Run("リクエスト成功", func(t *testing.T) {
		db := mock_database.NewMockQuerier(ctrl)
		db.EXPECT().CreateArticle(gomock.Any(), gomock.Any()).Return(&database.CreateArticleResult{ArticleID: 1}, nil)

		expect := &pb.CreateArticleResponse{
			ArticleId: 1,
		}

		got, err := callCreateArticle(req, db, nil, nil)

		if err != nil {
			t.Error(err)
		}
		if got == nil {
			t.Error("response should not be nil")
		}

		if !reflect.DeepEqual(got, expect) {
			t.Errorf("Expect: %v, Got: %v", expect, got)
		}
	})

	t.Run("データベースエラー", func(t *testing.T) {
		db := mock_database.NewMockQuerier(ctrl)
		db.EXPECT().CreateArticle(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))

		got, err := callCreateArticle(req, db, nil, nil)

		if got != nil {
			t.Error("response should be nil")
		}
		if err == nil {
			t.Error("err should not be nil")
		}
		if s, ok := status.FromError(err); ok {
			if s.Code() != codes.Internal {
				t.Errorf("Expect: %v, Got: %v", codes.Internal, s.Code())
			}
		}

	})
}

func TestServer_GetArticles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &emptypb.Empty{}

	t.Run("リクエスト成功", func(t *testing.T) {
		t.Run("multi Articles", func(t *testing.T) {
			mockCreatedAt := time.Now()
			db := mock_database.NewMockQuerier(ctrl)
			db.EXPECT().GetArticles(gomock.Any(), gomock.Any()).Return(&database.GetArticlesResult{
				Articles: []model.Article{
					{
						ID:        1,
						UserID:    1,
						Title:     "test_title_1",
						Text:      "test_text_1",
						CreatedAt: mockCreatedAt,
					},
					{
						ID:        2,
						UserID:    1,
						Title:     "test_title_2",
						Text:      "test_text_2",
						CreatedAt: mockCreatedAt,
					},
				},
			}, nil)

			expect := &pb.GetArticlesResponse{
				Articles: []*pb.Article{
					{
						ArticleId: 1,
						Title:     "test_title_1",
						Text:      "test_text_1",
						CreatedAt: timestampPtr(mockCreatedAt),
					},
					{
						ArticleId: 2,
						Title:     "test_title_2",
						Text:      "test_text_2",
						CreatedAt: timestampPtr(mockCreatedAt),
					},
				}}

			got, err := callGetArticles(req, db, nil, nil)

			if err != nil {
				t.Error("err should be nil")
			}
			if got == nil {
				t.Error("response should not be nil")
			}
			if !reflect.DeepEqual(got, expect) {
				t.Errorf("Expect: %v, Got: %v", expect, got)
			}
		})

		t.Run("one Article", func(t *testing.T) {
			mockCreatedAt := time.Now()
			db := mock_database.NewMockQuerier(ctrl)
			db.EXPECT().GetArticles(gomock.Any(), gomock.Any()).Return(&database.GetArticlesResult{
				Articles: []model.Article{
					{
						ID:        1,
						UserID:    1,
						Title:     "test_title_1",
						Text:      "test_text_1",
						CreatedAt: mockCreatedAt,
					},
				},
			}, nil)

			expect := &pb.GetArticlesResponse{
				Articles: []*pb.Article{
					{
						ArticleId: 1,
						Title:     "test_title_1",
						Text:      "test_text_1",
						CreatedAt: timestampPtr(mockCreatedAt),
					},
				}}

			got, err := callGetArticles(req, db, nil, nil)

			if err != nil {
				t.Error("err should be nil")
			}
			if got == nil {
				t.Error("response should not be nil")
			}
			if !reflect.DeepEqual(got, expect) {
				t.Errorf("Expect: %v, Got: %v", expect, got)
			}
		})

		t.Run("no Article", func(t *testing.T) {
			db := mock_database.NewMockQuerier(ctrl)
			db.EXPECT().GetArticles(gomock.Any(), gomock.Any()).Return(&database.GetArticlesResult{
				Articles: []model.Article{},
			}, nil)

			expect := &pb.GetArticlesResponse{}

			got, err := callGetArticles(req, db, nil, nil)

			if err != nil {
				t.Error("err should be nil")
			}
			if got == nil {
				t.Error("response should not be nil")
			}
			if !reflect.DeepEqual(got, expect) {
				t.Errorf("Expect: %v, Got: %v", expect, got)
			}
		})
	})

	t.Run("データベースエラー", func(t *testing.T) {
		db := mock_database.NewMockQuerier(ctrl)
		db.EXPECT().GetArticles(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))

		got, err := callGetArticles(req, db, nil, nil)

		if got != nil {
			t.Error("response should be nil")
		}

		if err == nil {
			t.Error("err should not be nil")
		}

		if s, ok := status.FromError(err); ok {
			if s.Code() != codes.Internal {
				t.Errorf("Expect: %v, Got: %v,", codes.Internal, s.Code())
			}
		}
	})
}

func TestServer_GetArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCreatedAt := time.Now()
	req := &pb.GetArticleRequest{
		ArticleId: 1,
	}

	t.Run("リクエスト成功", func(t *testing.T) {

		expect := &pb.GetArticleResponse{
			Article: &pb.Article{
				ArticleId: 1,
				Title:     "test_title",
				Text:      "test_text",
				CreatedAt: timestampPtr(mockCreatedAt),
			},
		}

		db := mock_database.NewMockQuerier(ctrl)
		db.EXPECT().GetArticle(gomock.Any(), gomock.Any()).Return(&database.GetArticleResult{
			Article: model.Article{
				ID:        1,
				UserID:    1,
				Title:     "test_title",
				Text:      "test_text",
				CreatedAt: mockCreatedAt,
			},
		}, nil)

		got, err := callGetArticle(req, db, nil, nil)

		if err != nil {
			t.Error("should be nil")
		}

		if got == nil {
			t.Error("should not be nil")
		}

		if !reflect.DeepEqual(got, expect) {
			t.Errorf("Expect: %v, Got: %v", expect, got)
		}
	})

	t.Run("データベースエラー", func(t *testing.T) {
		t.Run("not found", func(t *testing.T) {
			db := mock_database.NewMockQuerier(ctrl)
			db.EXPECT().GetArticle(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows)

			resp, err := callGetArticle(req, db, nil, nil)

			if resp != nil {
				t.Error("should not be nil")
			}

			if err == nil {
				t.Error("should not be nil")
			}

			if s, ok := status.FromError(err); ok {
				if s.Code() != codes.NotFound {
					t.Errorf("Expect: %v, Got: %v", codes.NotFound, s.Code())
				}
			}
		})

		t.Run("internal", func(t *testing.T) {
			db := mock_database.NewMockQuerier(ctrl)
			db.EXPECT().GetArticle(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))

			resp, err := callGetArticle(req, db, nil, nil)

			if resp != nil {
				t.Error("should not be nil")
			}

			if err == nil {
				t.Error("should not be nil")
			}

			if s, ok := status.FromError(err); ok {
				if s.Code() != codes.Internal {
					t.Errorf("Expect: %v, Got: %v", codes.Internal, s.Code())
				}
			}
		})
	})
}

func TestServer_UpdateArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.UpdateArticleRequest{
		ArticleId: 1,
		Title:     "test_title",
		Text:      "test_text",
	}

	t.Run("リクエスト成功", func(t *testing.T) {
		db := mock_database.NewMockQuerier(ctrl)
		db.EXPECT().UpdateArticle(gomock.Any(), gomock.Any()).Return(nil)

		_, err := callUpdateArticle(req, db, nil, nil)

		if err != nil {
			t.Error("err should be nil")
		}
	})

	t.Run("データベースエラー", func(t *testing.T) {
		db := mock_database.NewMockQuerier(ctrl)
		db.EXPECT().UpdateArticle(gomock.Any(), gomock.Any()).Return(errors.New("some error"))

		_, err := callUpdateArticle(req, db, nil, nil)

		if err == nil {
			t.Error("err should not be nil")
		}

		if s, ok := status.FromError(err); ok {
			if s.Code() != codes.Aborted {
				t.Errorf("Expect: %v, Got: %v", codes.Aborted, s.Code())
			}
		}
	})
}

func TestServer_DeleteArticle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &pb.DeleteArticleRequest{
		ArticleId: 1,
	}

	t.Run("リクエスト成功", func(t *testing.T) {
		db := mock_database.NewMockQuerier(ctrl)
		db.EXPECT().DeleteArticle(gomock.Any(), gomock.Any()).Return(nil)

		_, err := callDeleteArticle(req, db, nil, nil)

		if err != nil {
			t.Errorf("err should be nil: %v", err)
		}
	})

	t.Run("データベースエラー", func(t *testing.T) {
		db := mock_database.NewMockQuerier(ctrl)
		db.EXPECT().DeleteArticle(gomock.Any(), gomock.Any()).Return(errors.New("some error"))

		_, err := callDeleteArticle(req, db, nil, nil)

		if err == nil {
			t.Error("err should not be nil")
		}

		if s, ok := status.FromError(err); ok {
			if s.Code() != codes.Aborted {
				t.Errorf("Expect: %v, Got: %v", codes.Aborted, s.Code())
			}
		}
	})
}

func Test_extractUserID(t *testing.T) {
	f := func(ctx context.Context) context.Context {
		ctx = context.WithValue(ctx, KeyUserID, int64(1))
		return ctx
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "正常系",
			args: args{
				ctx: f(context.Background()),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractUserID(tt.args.ctx); got != tt.want {
				t.Errorf("extractUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringPtr(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		got := stringPtr("test")

		if got == nil {
			t.Error("should not be nil")
			return
		}

		if *got != "test" {
			t.Errorf("Expect: %v, Got: %v", "test", *got)
		}
	})

	t.Run("sql.NullString", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			got := stringPtr(sql.NullString{String: "test", Valid: true})

			if got == nil {
				t.Error("should not be nil")
				return
			}

			if *got != "test" {
				t.Errorf("Expect: %v, Got: %v", "test", *got)
			}
		})

		t.Run("invalid", func(t *testing.T) {
			got := stringPtr(sql.NullString{String: "", Valid: false})

			if got != nil {
				t.Errorf("should be nil: %v", *got)
			}
		})
	})
}

func TestConvNullString(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		expect := sql.NullString{
			String: "test",
			Valid:  true,
		}

		ps := stringPtr("test")
		got := convNullString(ps)

		if !reflect.DeepEqual(expect, got) {
			t.Errorf("Expect: %v, Got: %v", expect, got)
		}
	})

	t.Run("not nil", func(t *testing.T) {
		expect := sql.NullString{
			String: "",
			Valid:  false,
		}

		got := convNullString(nil)

		if !reflect.DeepEqual(expect, got) {
			t.Errorf("Expect: %v, Got: %v", expect, got)
		}
	})
}

func TestTimestampPtr(t *testing.T) {
	loc := time.FixedZone("Local", 9*60*60)

	tm, err := time.Parse(time.DateTime, "2023-03-29 09:00:00")
	tm = tm.In(loc)
	if err != nil {
		t.Fatalf("failed to parse time: %v", err)
	}

	expect := &timestamppb.Timestamp{
		Seconds: 1680080400,
	}

	got := timestampPtr(tm)

	if !reflect.DeepEqual(expect, got) {
		t.Errorf("Expect: %v, Got: %v", expect, got)
	}
}

func callHello(req *emptypb.Empty, db database.Querier, hash service.Hasher, auth service.Auther) (*pb.HelloWorldResponse, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyUserID, int64(1))
	s := NewServer(db, hash, auth)

	return s.HelloWorld(ctx, req)
}

func callSignUp(req *pb.SignUpRequest, db database.Querier, hash service.Hasher, auth service.Auther) (*pb.SignUpResponse, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyUserID, int64(1))
	s := NewServer(db, hash, auth)

	return s.SignUp(ctx, req)
}

func callLogin(req *pb.LoginRequest, db database.Querier, hash service.Hasher, auth service.Auther) (*pb.LoginResponse, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyUserID, int64(1))
	s := NewServer(db, hash, auth)

	return s.Login(ctx, req)
}

func callCreateArticle(req *pb.CreateArticleRequest, db database.Querier, hash service.Hasher, auth service.Auther) (*pb.CreateArticleResponse, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyUserID, int64(1))
	s := NewServer(db, hash, auth)

	return s.CreateArticle(ctx, req)
}

func callGetArticles(req *emptypb.Empty, db database.Querier, hash service.Hasher, auth service.Auther) (*pb.GetArticlesResponse, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyUserID, int64(1))
	s := NewServer(db, hash, auth)

	return s.GetArticles(ctx, req)
}

func callGetArticle(req *pb.GetArticleRequest, db database.Querier, hash service.Hasher, auth service.Auther) (*pb.GetArticleResponse, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyUserID, int64(1))
	s := NewServer(db, hash, auth)

	return s.GetArticle(ctx, req)
}

func callUpdateArticle(req *pb.UpdateArticleRequest, db database.Querier, hash service.Hasher, auth service.Auther) (*emptypb.Empty, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyUserID, int64(1))
	s := NewServer(db, hash, auth)

	return s.UpdateArticle(ctx, req)
}

func callDeleteArticle(req *pb.DeleteArticleRequest, db database.Querier, hash service.Hasher, auth service.Auther) (*emptypb.Empty, error) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KeyUserID, int64(1))
	s := NewServer(db, hash, auth)

	return s.DeleteArticle(ctx, req)
}
