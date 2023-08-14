package permession

import (
	"context"
	"time"
	"transfer/api/types/schema"
	permessionModel "transfer/internal/app/model/permession"

	"github.com/google/wire"
)

var ServiceSet = wire.NewSet(wire.Struct(new(Service), "*"))

type Service struct {
	PermessionRepo *permessionModel.PermessionRepo
}

func (s *Service) Create(c context.Context, userPermession *permessionModel.Permession) error {
	if err := s.PermessionRepo.Create(userPermession); err != nil {
		return err
	}
	return nil
}

func (s *Service) Info(c context.Context, userPermession *permessionModel.Permession) *[]schema.FilePermession {
	return s.PermessionRepo.Info(userPermession)

}

func (s *Service) CheckFile(c context.Context, userPermession *permessionModel.Permession) bool {
	filePermession := s.PermessionRepo.Info(userPermession)
	if filePermession != nil {
		for _, permession := range *filePermession {
			return time.Now().Before(permession.Expiredate)
		}
	}
	return false
}

func (s *Service) Update(c context.Context, userPermession *permessionModel.Permession) bool {
	return s.PermessionRepo.Update(userPermession)
}
