package main

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/flaneur4dev/good-limiter/internal/server/grpc/pb"
)

type client struct {
	conn       *grpc.ClientConn
	grpcClient pb.RateLimiterClient
}

func newClient(target string) (*client, error) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &client{
		conn:       conn,
		grpcClient: pb.NewRateLimiterClient(conn),
	}, nil
}

func (r *client) AddNet(ctx context.Context, ipNet, list string) string {
	if !validateIPNet(ipNet) || !validateList(list) {
		return invalidArgs
	}

	res, err := r.grpcClient.AddNet(ctx, &pb.AddRequest{SubNet: ipNet, List: list})
	if err != nil {
		return err.Error()
	}

	return res.GetMessage()
}

func (r *client) DeleteNet(ctx context.Context, ipNet, list string) string {
	if !validateIPNet(ipNet) || !validateList(list) {
		return invalidArgs
	}

	res, err := r.grpcClient.DeleteNet(ctx, &pb.DeleteRequest{SubNet: ipNet, List: list})
	if err != nil {
		return err.Error()
	}

	return res.GetMessage()
}

func (r *client) DropBucket(ctx context.Context, login, ip string) string {
	if !validateIP(ip) || len(strings.TrimSpace(login)) == 0 {
		return invalidArgs
	}

	res, err := r.grpcClient.DropBucket(ctx, &pb.DropRequest{Login: login, Ip: ip})
	if err != nil {
		return err.Error()
	}

	return res.GetMessage()
}

func (r *client) Close() {
	r.conn.Close()
}

func validateList(list string) bool {
	return list == "black" || list == "white"
}

func validateIPNet(ip string) bool {
	_, ipNet, err := net.ParseCIDR(ip)
	if err != nil {
		return false
	}

	return ip == ipNet.String()
}

func validateIP(ip string) bool {
	return ip == net.ParseIP(ip).String()
}
