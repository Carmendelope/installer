/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package server

import (
	"fmt"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/grpc-utils/pkg/tools"
	"github.com/nalej/installer/internal/pkg/server/config"
	"github.com/nalej/installer/internal/pkg/server/installer"
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Service struct {
	Configuration config.Config
	Server * tools.GenericGRPCServer
}

// NewService creates a new system model service.
func NewService(conf config.Config) *Service {
	return &Service{
		conf,
		tools.NewGenericGRPCServer(uint32(conf.Port)),
	}
}

// Run the service, launch the REST service handler.
func (s *Service) Run() error {

	s.Configuration.Environment.Resolve()

	s.Configuration.ComponentsPath = utils.ExtendComponentsPath(s.Configuration.ComponentsPath, true)
	vErr := s.Configuration.Validate()
	if vErr != nil {
		log.Error().Str("error", vErr.DebugReport()).Msg("invalid configuration")
		return vErr
	}
	s.Configuration.Print()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}

	installerManager := installer.NewManager(s.Configuration)
	installerHandler := installer.NewHandler(installerManager)

	grpcServer := grpc.NewServer()
	grpc_installer_go.RegisterInstallerServer(grpcServer, installerHandler)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}