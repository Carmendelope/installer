/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package server

import (
	"fmt"
	"github.com/nalej/grpc-installer-go"
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
}

// NewService creates a new system model service.
func NewService(conf config.Config) *Service {
	return &Service{
		conf,
	}
}

// Run the service, launch the REST service handler.
func (s *Service) Run() error {
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
