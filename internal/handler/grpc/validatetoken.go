package grpc

import (
	"context"

	pb "github.com/mhasnanr/ewallet-ums/cmd/tokenvalidation"
	"github.com/mhasnanr/ewallet-ums/constants"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type JWTManager interface {
	ValidateToken(ctx context.Context, token string) (*helpers.ClaimToken, error)
}

type TokenValidation struct {
	pb.UnimplementedTokenValidationServer
	jwtManager JWTManager
}

func NewTokenValidationHandler(jwtManager JWTManager) *TokenValidation {
	return &TokenValidation{
		jwtManager: jwtManager,
	}
}

func (h *TokenValidation) ValidateToken(ctx context.Context, request *pb.TokenRequest) (*pb.TokenResponse, error) {
	token := request.GetToken()
	if token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	claims, err := h.jwtManager.ValidateToken(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return &pb.TokenResponse{
		Message: constants.ValidToken,
		Data: &pb.UserData{
			UserId:   int64(claims.UserID),
			Username: claims.Username,
			FullName: claims.Fullname,
			Email:    claims.Email,
		},
	}, nil
}
