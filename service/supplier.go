package service

import (
	"context"
	"encoding/json"
	"errors"
	ms "supplier-be/model/supplier"
	"supplier-be/model/utils"
	repository "supplier-be/repository"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// SupplierService is the interface that provides supplier-related methods.
type Supplier interface {
	CreateSupplier(ctx context.Context, request *ms.SupplierRequest, traceId string) (*ms.Supplier, error)
	GetListSupplier(ctx context.Context, filter *ms.SupplierFilter, traceId string) (*ms.SupplierListResponse, error)
}

// SupplierServiceImpl is the implementation of SupplierService.
type SupplierService struct {
	supplierRepository repository.SupplierRepository
	supplierDB         *sqlx.DB
	log                *logrus.Logger
}

// NewSupplierService creates a new SupplierService instance.
func NewSupplierService(supplierRepository repository.SupplierRepository, supplierDB *sqlx.DB, log *logrus.Logger) Supplier {
	return &SupplierService{supplierRepository: supplierRepository, supplierDB: supplierDB, log: log}
}

// CreateSupplier creates a new supplier..
func (s *SupplierService) CreateSupplier(ctx context.Context, request *ms.SupplierRequest, traceId string) (*ms.Supplier, error) {

	// simpan request json ke field address JSONB
	stringObj, err := json.Marshal(request)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"TraceID": traceId,
			"Request": request,
		}).Error("CreateSupplier - Failed to marshal request to JSON: ", err)
		return nil, err
	}

	for _, contact := range request.Contacts {
		if !utils.IsPhoneNumber(contact.PhoneNumber) {
			s.log.WithFields(logrus.Fields{
				"TraceID": traceId,
				"Request": request,
			}).Error("CreateSupplier - Invalid phone number format for contact: ", contact.PhoneNumber)
			return nil, errors.New("invalid phone number format for contact: " + contact.PhoneNumber)
		}

		if !utils.IsPhoneNumber(contact.MobileNumber) {
			s.log.WithFields(logrus.Fields{
				"TraceID": traceId,
				"Request": request,
			}).Error("CreateSupplier - Invalid phone number format for contact: ", contact.MobileNumber)
			return nil, errors.New("invalid phone number format for contact: " + contact.MobileNumber)
		}

		if !utils.IsValidEmail(contact.Email) {
			s.log.WithFields(logrus.Fields{
				"TraceID": traceId,
				"Request": request,
			}).Error("CreateSupplier - Invalid email format for contact: ", contact.Email)
			return nil, errors.New("invalid email format for contact: " + contact.Email)
		}
	}

	// insert supplier utama
	id, err := s.supplierRepository.CreateSupplier(ctx, &ms.Supplier{
		SupplierName: request.SupplierName,
		Logo:         request.Logo,
		Nickname:     request.Nickname,
		Status:       ms.StatusSupplier[1],
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Address:      string(stringObj),
	})
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"TraceID": traceId,
			"Request": request,
		}).Error("CreateSupplier - Failed to create supplier: ", err)
		return nil, err
	}

	s.log.Printf("CreateSupplier - Supplier created with ID: %d", id)

	// proses multiple address/contact/group di background
	go func(request *ms.SupplierRequest, id int64, traceId string) {
		g := new(errgroup.Group) // no WithContext, biar gak auto cancel

		for _, addr := range request.Address {
			addr := addr
			g.Go(func() error {
				_, err := s.supplierRepository.CreateAddress(context.Background(), &ms.Address{
					SupplierID: id,
					Name:       addr.Name,
					Address:    addr.Address,
					IsMain:     addr.IsMain,
				})
				if err != nil {
					s.log.WithFields(logrus.Fields{
						"TraceID": traceId,
						"Request": request,
					}).Error("Failed to create address: ", err)
				}
				return nil // jangan propagate error
			})
		}

		for _, contact := range request.Contacts {
			contact := contact
			g.Go(func() error {
				_, err := s.supplierRepository.CreateContacts(context.Background(), &ms.Contacts{
					SupplierID:   id,
					Name:         contact.Name,
					JobPosition:  contact.JobPosition,
					Email:        contact.Email,
					PhoneNumber:  utils.FormatPhoneNumber(contact.PhoneNumber),
					MobileNumber: utils.FormatPhoneNumber(contact.MobileNumber),
					IsMain:       contact.IsMain,
				})
				if err != nil {
					s.log.WithFields(logrus.Fields{
						"TraceID": traceId,
						"Request": request,
					}).Error("Failed to create contact: ", err)
				}
				return nil
			})
		}

		for _, group := range request.Groups {
			group := group
			g.Go(func() error {
				_, err := s.supplierRepository.CreateGroups(context.Background(), &ms.Groups{
					SupplierID: id,
					GroupName:  group.GroupName,
					Value:      group.Value,
					IsActive:   group.IsActive,
				})
				if err != nil {
					s.log.WithFields(logrus.Fields{
						"TraceID": traceId,
						"Request": request,
					}).Error("Failed to create groups: ", err)
				}
				return nil
			})
		}

		if err := g.Wait(); err != nil {
			s.log.WithFields(logrus.Fields{
				"TraceID": traceId,
				"Request": request,
			}).Error("Background process error details: ", err)
		}
	}(request, id, traceId)

	return &ms.Supplier{
		SupplierID:   id,
		SupplierName: request.SupplierName,
		Logo:         request.Logo,
		Nickname:     request.Nickname,
		Status:       ms.StatusSupplier[1],
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Address:      string(stringObj),
	}, nil
}

func (s *SupplierService) GetListSupplier(ctx context.Context, filter *ms.SupplierFilter, traceId string) (*ms.SupplierListResponse, error) {

	dataSuppliers, err := s.supplierRepository.GetListSupplier(ctx, filter)
	if err != nil {
		s.log.WithFields(logrus.Fields{
			"TraceID": traceId,
			"Request": filter,
		}).Error("GetListSupplier - Failed to get supplier: ", err)
		return nil, err
	}

	return &ms.SupplierListResponse{
		TotalData: int64(len(dataSuppliers)),
		Data:      dataSuppliers,
	}, nil
}
