package rpc

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/qubic/go-node-connector/types"
	"github.com/qubic/go-schnorrq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"math"
	"net"

	qubic "github.com/qubic/go-node-connector"
	"github.com/qubic/qubic-http/protobuff"
	"log"
	"net/http"
)

var _ protobuff.QubicLiveServiceServer = &Server{}

type Server struct {
	protobuff.UnimplementedQubicLiveServiceServer
	logger          *log.Logger
	listenAddrGRPC  string
	listenAddrHTTP  string
	qPool           *qubic.Pool
	maxTickFetchUrl string
}

func NewServer(listenAddrGRPC, listenAddrHTTP string, logger *log.Logger, qPool *qubic.Pool, maxTickFetchUrl string) *Server {
	return &Server{
		listenAddrGRPC:  listenAddrGRPC,
		listenAddrHTTP:  listenAddrHTTP,
		logger:          logger,
		qPool:           qPool,
		maxTickFetchUrl: maxTickFetchUrl,
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
		IncomingAmount:             identityInfo.AddressData.IncomingAmount,
		OutgoingAmount:             identityInfo.AddressData.OutgoingAmount,
		NumberOfIncomingTransfers:  identityInfo.AddressData.NumberOfIncomingTransfers,
		NumberOfOutgoingTransfers:  identityInfo.AddressData.NumberOfOutgoingTransfers,
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

func (s *Server) QuerySmartContract(ctx context.Context, req *protobuff.QuerySmartContractRequest) (*protobuff.QuerySmartContractResponse, error) {
	reqData, err := base64.StdEncoding.DecodeString(req.RequestData)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "failed to decode from base64 the request data: %s", req.RequestData)
	}

	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection %v", err)
	}

	scData, err := client.QuerySmartContract(ctx, qubic.RequestContractFunction{
		ContractIndex: req.ContractIndex,
		InputType:     uint16(req.InputType),
		InputSize:     uint16(req.InputSize),
	}, reqData)
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "query smart contract %v", err)
	}

	s.qPool.Put(client)

	return &protobuff.QuerySmartContractResponse{ResponseData: base64.StdEncoding.EncodeToString(scData.Data)}, nil
}

type maxTickResponse struct {
	MaxTick uint32 `json:"max_tick"`
}

func fetchMaxTick(ctx context.Context, maxTickFetchUrl string) (uint32, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, maxTickFetchUrl, nil)
	if err != nil {
		return 0, errors.Wrap(err, "creating new request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "performing request")
	}
	defer res.Body.Close()

	var resp maxTickResponse
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, errors.Wrap(err, "reading response body")
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return 0, errors.Wrap(err, "unmarshalling response")
	}

	tick := resp.MaxTick

	if tick == 0 {
		return 0, errors.New("Fetched max tick is 0.")
	}

	return tick, nil
}

func (s *Server) BroadcastTransaction(ctx context.Context, req *protobuff.BroadcastTransactionRequest) (*protobuff.BroadcastTransactionResponse, error) {
	decodedTx, err := base64.StdEncoding.DecodeString(req.EncodedTransaction)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	reader := bytes.NewReader(decodedTx)

	var transaction types.Transaction
	err = transaction.UnmarshallBinary(reader)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	digest, err := transaction.GetUnsignedDigest()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = schnorrq.Verify(transaction.SourcePublicKey, digest, transaction.Signature)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	maxTick, err := fetchMaxTick(ctx, s.maxTickFetchUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	offsetTick := int32(transaction.Tick) - int32(maxTick)
	if offsetTick <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Target tick: %d for the transaction should be greater than max tick: %d", transaction.Tick, maxTick)
	}

	transactionId, err := transaction.ID()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var sourceID types.Identity
	sourceID, err = sourceID.FromPubKey(transaction.SourcePublicKey, false)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "getting source ID").Error())
	}

	var destID types.Identity
	destID, err = destID.FromPubKey(transaction.DestinationPublicKey, false)
	if err != nil {
		return nil, status.Error(codes.Internal, errors.Wrap(err, "getting dest ID").Error())
	}

	peersBroadcasted := broadcastTxToMultiple(ctx, s.qPool, decodedTx)
	s.logger.Printf("Tx ID: %s | Source: %s | Dest: %s | Target tick: %d | Max tick: %d | Offset tick: %d | Peers broadcasted: %d\n", transactionId, sourceID, destID, transaction.Tick, maxTick, offsetTick, peersBroadcasted)
	if peersBroadcasted == 0 {
		return nil, status.Error(codes.Internal, "tx wasn't broadcast to any peers, please retry")
	}

	return &protobuff.BroadcastTransactionResponse{
		PeersBroadcasted:   int32(peersBroadcasted),
		EncodedTransaction: req.EncodedTransaction,
		TransactionId:      transactionId,
	}, nil
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

func int8ArrayToString(array []int8) string {
	runes := make([]rune, 0)

	for _, char := range array {
		if char == 0 {
			continue
		}

		runes = append(runes, rune(char))
	}
	return string(runes)
}

func int8ArrayToInt32Array(array []int8) []int32 {
	ints := make([]int32, 0)

	for _, smallInt := range array {
		ints = append(ints, int32(smallInt))
	}
	return ints
}

func (s *Server) GetIssuedAssets(ctx context.Context, req *protobuff.IssuedAssetsRequest) (*protobuff.IssuedAssetsResponse, error) {
	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection :%v", err)
	}

	assets, err := client.GetIssuedAssets(ctx, req.Identity)
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "getting issued assets from node %v", err)
	}

	s.qPool.Put(client)

	issuedAssets := make([]*protobuff.IssuedAsset, 0)

	for _, asset := range assets {

		iAsset := asset.Data
		var iAssetIdentity types.Identity
		iAssetIdentity, err = iAssetIdentity.FromPubKey(iAsset.PublicKey, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get identity for issued asset public key")
		}

		data := protobuff.IssuedAssetData{
			IssuerIdentity:        iAssetIdentity.String(),
			Type:                  uint32(iAsset.Type),
			Name:                  int8ArrayToString(iAsset.Name[:]),
			NumberOfDecimalPlaces: int32(iAsset.NumberOfDecimalPlaces),
			UnitOfMeasurement:     int8ArrayToInt32Array(iAsset.UnitOfMeasurement[:]),
		}

		info := protobuff.AssetInfo{
			Tick:          asset.Info.Tick,
			UniverseIndex: asset.Info.UniverseIndex,
		}

		issuedAsset := protobuff.IssuedAsset{
			Data: &data,
			Info: &info,
		}

		issuedAssets = append(issuedAssets, &issuedAsset)
	}

	return &protobuff.IssuedAssetsResponse{IssuedAssets: issuedAssets}, nil
}

func (s *Server) GetOwnedAssets(ctx context.Context, req *protobuff.OwnedAssetsRequest) (*protobuff.OwnedAssetsResponse, error) {
	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection :%v", err)
	}

	assets, err := client.GetOwnedAssets(ctx, req.Identity)
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "getting owned assets from node %v", err)
	}

	s.qPool.Put(client)

	ownedAssets := make([]*protobuff.OwnedAsset, 0)

	for _, asset := range assets {

		iAsset := asset.Data.IssuedAsset

		var iAssetIdentity types.Identity
		iAssetIdentity, err = iAssetIdentity.FromPubKey(iAsset.PublicKey, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get identity for issued asset public key")
		}

		issuedAsset := protobuff.IssuedAssetData{
			IssuerIdentity:        iAssetIdentity.String(),
			Type:                  uint32(iAsset.Type),
			Name:                  int8ArrayToString(iAsset.Name[:]),
			NumberOfDecimalPlaces: int32(iAsset.NumberOfDecimalPlaces),
			UnitOfMeasurement:     int8ArrayToInt32Array(iAsset.UnitOfMeasurement[:]),
		}

		var oAssetIdentity types.Identity
		oAssetIdentity, err = oAssetIdentity.FromPubKey(asset.Data.PublicKey, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get identity for owned asset public key")
		}

		data := protobuff.OwnedAssetData{
			OwnerIdentity:         oAssetIdentity.String(),
			Type:                  uint32(asset.Data.Type),
			Padding:               int32(asset.Data.Padding[0]),
			ManagingContractIndex: uint32(asset.Data.ManagingContractIndex),
			IssuanceIndex:         asset.Data.IssuanceIndex,
			NumberOfUnits:         asset.Data.NumberOfUnits,
			IssuedAsset:           &issuedAsset,
		}

		info := protobuff.AssetInfo{
			Tick:          asset.Info.Tick,
			UniverseIndex: asset.Info.UniverseIndex,
		}

		ownedAsset := protobuff.OwnedAsset{
			Data: &data,
			Info: &info,
		}

		ownedAssets = append(ownedAssets, &ownedAsset)

	}

	return &protobuff.OwnedAssetsResponse{OwnedAssets: ownedAssets}, nil
}

func (s *Server) GetPossessedAssets(ctx context.Context, req *protobuff.PossessedAssetsRequest) (*protobuff.PossessedAssetsResponse, error) {
	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection :%v", err)
	}

	assets, err := client.GetPossessedAssets(ctx, req.Identity)
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "getting possessed assets from node %v", err)
	}

	s.qPool.Put(client)

	possessedAssets := make([]*protobuff.PossessedAsset, 0)

	for _, asset := range assets {

		oAsset := asset.Data.OwnedAsset
		var oAssetIdentity types.Identity
		oAssetIdentity, err = oAssetIdentity.FromPubKey(oAsset.PublicKey, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get identity for owned asset public key")
		}

		iAsset := oAsset.IssuedAsset
		var iAssetIdentity types.Identity
		iAssetIdentity, err = iAssetIdentity.FromPubKey(iAsset.PublicKey, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get identity for issued asset public key")
		}

		issuedAsset := protobuff.IssuedAssetData{
			IssuerIdentity:        iAssetIdentity.String(),
			Type:                  uint32(iAsset.Type),
			Name:                  int8ArrayToString(iAsset.Name[:]),
			NumberOfDecimalPlaces: int32(iAsset.NumberOfDecimalPlaces),
			UnitOfMeasurement:     int8ArrayToInt32Array(iAsset.UnitOfMeasurement[:]),
		}

		ownedAsset := protobuff.OwnedAssetData{
			OwnerIdentity:         oAssetIdentity.String(),
			Type:                  uint32(asset.Data.Type),
			Padding:               int32(asset.Data.Padding[0]),
			ManagingContractIndex: uint32(asset.Data.ManagingContractIndex),
			IssuanceIndex:         asset.Data.IssuanceIndex,
			NumberOfUnits:         asset.Data.NumberOfUnits,
			IssuedAsset:           &issuedAsset,
		}

		var pAssetIdentity types.Identity
		pAssetIdentity, err = pAssetIdentity.FromPubKey(asset.Data.PublicKey, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get identity for possessed asset public key")
		}

		data := protobuff.PossessedAssetData{
			PossessorIdentity:     pAssetIdentity.String(),
			Type:                  uint32(asset.Data.Type),
			Padding:               int32(asset.Data.Padding[0]),
			ManagingContractIndex: uint32(asset.Data.ManagingContractIndex),
			IssuanceIndex:         asset.Data.IssuanceIndex,
			NumberOfUnits:         asset.Data.NumberOfUnits,
			OwnedAsset:            &ownedAsset,
		}

		info := protobuff.AssetInfo{
			Tick:          asset.Info.Tick,
			UniverseIndex: asset.Info.UniverseIndex,
		}

		possessedAsset := protobuff.PossessedAsset{
			Data: &data,
			Info: &info,
		}

		possessedAssets = append(possessedAssets, &possessedAsset)

	}

	return &protobuff.PossessedAssetsResponse{PossessedAssets: possessedAssets}, nil
}

const assetIssuanceType = 1
const assetOwnershipType = 2
const assetPossessionType = 3

func (s *Server) GetIssuedAssetByUniverseIndex(ctx context.Context, request *protobuff.GetByUniverseIndexRequest) (*protobuff.AssetIssuance, error) {
	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection: %v", err)
	}

	assets, err := client.GetAssetIssuancesByUniverseIndex(ctx, request.Index)
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "getting asset issuances: %v", err)
	}
	s.qPool.Put(client)

	// don't return empty or wrong type
	if len(assets) > 0 && assets[0].Asset != (types.AssetIssuanceData{}) && assets[0].Asset.Type == assetIssuanceType {

		converted, err := convertAssetIssuance(assets[0])
		if err != nil {
			return nil, errors.Wrap(err, "converting response")
		}
		return converted, nil

	} else {
		log.Printf("asset issuance with universe index [%d] not found", request.Index)
		return nil, status.Error(codes.NotFound, "asset not found")
	}

}

func (s *Server) GetOwnedAssetByUniverseIndex(ctx context.Context, request *protobuff.GetByUniverseIndexRequest) (*protobuff.AssetOwnership, error) {
	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection: %v", err)
	}

	assets, err := client.GetAssetOwnershipsByUniverseIndex(ctx, request.Index)
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "getting asset ownerships: %v", err)
	}
	s.qPool.Put(client)

	// don't return empty or wrong type
	if len(assets) > 0 && assets[0].Asset != (types.AssetOwnershipData{}) && assets[0].Asset.Type == assetOwnershipType {

		converted, err := convertAssetOwnership(assets[0])
		if err != nil {
			return nil, errors.Wrap(err, "converting response")
		}
		return converted, nil

	} else {
		log.Printf("asset ownership with universe index [%d] not found", request.Index)
		return nil, status.Error(codes.NotFound, "asset not found")
	}

}

func (s *Server) GetPossessedAssetByUniverseIndex(ctx context.Context, request *protobuff.GetByUniverseIndexRequest) (*protobuff.AssetPossession, error) {
	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection: %v", err)
	}

	assets, err := client.GetAssetPossessionsByUniverseIndex(ctx, request.Index)
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "getting asset possessions: %v", err)
	}
	s.qPool.Put(client)

	// don't return empty or wrong type
	if len(assets) > 0 && assets[0].Asset != (types.AssetPossessionData{}) && assets[0].Asset.Type == assetPossessionType {

		converted, err := convertAssetPossession(assets[0])
		if err != nil {
			return nil, errors.Wrap(err, "converting response")
		}
		return converted, nil

	} else {
		log.Printf("asset possession with universe index [%d] not found", request.Index)
		return nil, status.Error(codes.NotFound, "asset not found")
	}
}

func (s *Server) GetIssuedAssetsByFilter(ctx context.Context, request *protobuff.GetIssuedAssetsByFilterRequest) (*protobuff.AssetIssuances, error) {

	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection: %v", err)
	}

	assets, err := client.GetAssetIssuancesByFilter(ctx, request.GetIssuerIdentity(), request.GetAssetName())
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "getting asset issuances: %v", err)
	}
	s.qPool.Put(client)

	assetIssuances := make([]*protobuff.AssetIssuance, 0)

	for _, asset := range assets {
		assetIssuance, err := convertAssetIssuance(asset)
		if err != nil {
			return nil, errors.Wrap(err, "converting response")
		}
		assetIssuances = append(assetIssuances, assetIssuance)
	}

	return &protobuff.AssetIssuances{Assets: assetIssuances}, nil
}

func (s *Server) GetOwnedAssetsByFilter(ctx context.Context, request *protobuff.GetOwnedAssetsByFilterRequest) (*protobuff.AssetOwnerships, error) {

	if request.GetOwnershipManagingContract() > math.MaxUint16 {
		return nil, status.Errorf(codes.InvalidArgument, "ownership managing contract value out of range")
	}
	contract := uint16(request.GetOwnershipManagingContract())

	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection: %v", err)
	}

	assets, err := client.GetAssetOwnershipsByFilter(ctx, request.GetIssuerIdentity(), request.GetAssetName(), request.GetOwnerIdentity(), contract)
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "getting asset ownerships: %v", err)
	}
	s.qPool.Put(client)

	assetOwnerships := make([]*protobuff.AssetOwnership, 0)

	for _, asset := range assets {
		assetOwnership, err := convertAssetOwnership(asset)
		if err != nil {
			return nil, errors.Wrap(err, "converting response")
		}

		assetOwnerships = append(assetOwnerships, assetOwnership)
	}

	return &protobuff.AssetOwnerships{Assets: assetOwnerships}, nil

}

func (s *Server) GetPossessedAssetsByFilter(ctx context.Context, request *protobuff.GetPossessedAssetsByFilterRequest) (*protobuff.AssetPossessions, error) {
	if request.GetOwnershipManagingContract() > math.MaxUint16 {
		return nil, status.Errorf(codes.InvalidArgument, "ownership managing contract value out of range")
	}
	ownerContract := uint16(request.GetOwnershipManagingContract())

	if request.GetPossessionManagingContract() > math.MaxUint16 {
		return nil, status.Errorf(codes.InvalidArgument, "possession managing contract value out of range")
	}
	possessorContract := uint16(request.GetPossessionManagingContract())

	client, err := s.qPool.Get()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting pool connection: %v", err)
	}

	assets, err := client.GetAssetPossessionsByFilter(ctx, request.GetIssuerIdentity(), request.GetAssetName(), request.GetOwnerIdentity(), request.GetPossessorIdentity(), ownerContract, possessorContract)
	if err != nil {
		s.qPool.Close(client)
		return nil, status.Errorf(codes.Internal, "getting asset possessions: %v", err)
	}
	s.qPool.Put(client)

	assetPossessions := make([]*protobuff.AssetPossession, 0)

	for _, asset := range assets {
		assetPossession, err := convertAssetPossession(asset)
		if err != nil {
			return nil, errors.Wrap(err, "converting response")
		}
		assetPossessions = append(assetPossessions, assetPossession)
	}

	return &protobuff.AssetPossessions{Assets: assetPossessions}, nil
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
		return errors.Wrap(err, "listening gRPC")
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
