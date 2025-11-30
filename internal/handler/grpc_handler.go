package handler

import (
	"context"
	"github.com/spitfy/urlshortener/internal/auth"
	pb "github.com/spitfy/urlshortener/pkg/shortener"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedShortenerServiceServer
	service ServiceShortener
	auth    *auth.Manager
}

func newGRPC(service ServiceShortener, auth *auth.Manager) *server {
	return &server{
		service: service,
		auth:    auth,
	}
}

func (s *server) ShortenURL(ctx context.Context, req *pb.URLShortenRequest) (*pb.URLShortenResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authorization required")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token missing")
	}

	token := authHeader[0]
	userID, err := s.auth.ParseUserID(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	shortURL, err := s.service.Add(ctx, req.GetUrl(), userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.URLShortenResponse{Result: shortURL}, nil
}

func (s *server) ExpandURL(ctx context.Context, req *pb.URLExpandRequest) (*pb.URLExpandResponse, error) {
	originalURL, err := s.service.GetByHash(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "URL not found")
	}

	return &pb.URLExpandResponse{Result: originalURL.Link}, nil
}

func (s *server) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*pb.UserURLsResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authorization required")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token missing")
	}

	token := authHeader[0]
	userID, err := s.auth.ParseUserID(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	urls, err := s.service.GetByUserID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbURLs []*pb.URLData
	for _, u := range urls {
		pbURLs = append(pbURLs, &pb.URLData{
			ShortUrl:    u.ShortURL,
			OriginalUrl: u.OriginalURL,
		})
	}

	return &pb.UserURLsResponse{Url: pbURLs}, nil
}
