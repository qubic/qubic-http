package rpc

import (
	"context"
	"encoding/base64"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	qubic "github.com/qubic/go-node-connector"
	"github.com/qubic/qubic-http/protobuff"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"google.golang.org/protobuf/encoding/protojson"

	"log"
	"net"
	"net/http"
)

var _ protobuff.QubicLiveServiceServer = &Server{}

type Server struct {
	protobuff.UnimplementedQubicLiveServiceServer
	listenAddrGRPC string
	listenAddrHTTP string
	qPool          *qubic.Pool
}

func NewServer(listenAddrGRPC, listenAddrHTTP string, qPool *qubic.Pool) *Server {
	return &Server{
		listenAddrGRPC: listenAddrGRPC,
		listenAddrHTTP: listenAddrHTTP,
		qPool:          qPool,
	}
}

func (s *Server) GetBalance(ctx context.Context, req *protobuff.GetBalanceRequest) (*protobuff.GetBalanceResponse, error) {
	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection :%v", err)
	}

	identityInfo, err := client.GetIdentity(ctx, req.Id)
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "getting identity info from node %v", err)
	}

	s.qPool.Put(client)

	balance := protobuff.Balance{
		Id:                         req.Id,
		Balance:                    identityInfo.AddressData.IncomingAmount - identityInfo.AddressData.OutgoingAmount,
		ValidForTick:               identityInfo.Tick,
		LatestIncomingTransferTick: identityInfo.AddressData.LatestIncomingTransferTick,
		LatestOutgoingTransferTick: identityInfo.AddressData.LatestOutgoingTransferTick,
	}
	return &protobuff.GetBalanceResponse{Balance: &balance}, nil
}

func (s *Server) GetTickInfo(ctx context.Context, _ *emptypb.Empty) (*protobuff.GetTickInfoResponse, error) {
	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection %v", err)
	}

	tickInfo, err := client.GetTickInfo(ctx)
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "getting tick info from node %v", err)
	}

	s.qPool.Put(client)
	return &protobuff.GetTickInfoResponse{TickInfo: &protobuff.TickInfo{
		Tick:        tickInfo.Tick,
		Duration:    uint32(tickInfo.TickDuration),
		Epoch:       uint32(tickInfo.Epoch),
		InitialTick: tickInfo.InitialTick,
	}}, nil
}

func (s *Server) GetBlockHeight(ctx context.Context, _ *emptypb.Empty) (*protobuff.GetBlockHeightResponse, error) {
	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection %v", err)
	}

	tickInfo, err := client.GetTickInfo(ctx)
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "getting tick info from node %v", err)
	}

	s.qPool.Put(client)
	return &protobuff.GetBlockHeightResponse{BlockHeight: &protobuff.TickInfo{
		Tick:        tickInfo.Tick,
		Duration:    uint32(tickInfo.TickDuration),
		Epoch:       uint32(tickInfo.Epoch),
		InitialTick: tickInfo.InitialTick,
	}}, nil
}

func (s *Server) BroadcastTransaction(ctx context.Context, req *protobuff.BroadcastTransactionRequest) (*protobuff.BroadcastTransactionResponse, error) {
	decodedTx, err := base64.StdEncoding.DecodeString(req.EncodedTransaction)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &protobuff.BroadcastTransactionResponse{PeersBroadcasted: int32(broadcastTxToMultiple(ctx, s.qPool, decodedTx)), EncodedTransaction: req.EncodedTransaction}, nil
}

func broadcastTxToMultiple(ctx context.Context, pool *qubic.Pool, decodedTx []byte) int {
	nrSuccess := 0
	for i := 0; i < 3; i++ {
		func() {
			client, err := pool.Get()
			if err != nil {
				return
			}

			err = client.SendRawTransaction(ctx, decodedTx)
			if err != nil {
				pool.Close(client)
				return
			}
			pool.Put(client)
			nrSuccess++
		}()
	}

	return nrSuccess
}

func (s *Server) Start() error {
	srv := grpc.NewServer(
		grpc.MaxRecvMsgSize(600*1024*1024),
		grpc.MaxSendMsgSize(600*1024*1024),
	)
	protobuff.RegisterQubicLiveServiceServer(srv, s)
	reflection.Register(srv)

	lis, err := net.Listen("tcp", s.listenAddrGRPC)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()

	if s.listenAddrHTTP != "" {
		go func() {
			mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{EmitDefaultValues: true, EmitUnpopulated: false},
			}))
			opts := []grpc.DialOption{
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithDefaultCallOptions(
					grpc.MaxCallRecvMsgSize(600*1024*1024),
					grpc.MaxCallSendMsgSize(600*1024*1024),
				),
			}

			if err := protobuff.RegisterQubicLiveServiceHandlerFromEndpoint(
				context.Background(),
				mux,
				s.listenAddrGRPC,
				opts,
			); err != nil {
				panic(err)
			}

			if err := http.ListenAndServe(s.listenAddrHTTP, mux); err != nil {
				panic(err)
			}
		}()
	}

	return nil
}
