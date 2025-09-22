# 📘Supplier API Documentation

---

### 📌 Overview

Project ini adalah implementasi **Supplier Management API** sesuai wireframe yang diberikan, menggunakan:

- **Backend**: Golang + Fiber
- **Database**: PostgreSQL
- **Logger**: Logrus (dengan TraceID)
- **API Contract**: JSON RESTful API

---
---
---
---

### 4. DB Architecture/ERD

diagram erd 

![Screenshot 2025-09-22 at 23.28.18.png](attachment:dba023d5-9ca7-44c9-b930-b093d8b457bd:Screenshot_2025-09-22_at_23.28.18.png)


### 5. Contoh Hasil Pengerjaan

### 🚀 Endpoint: Create Supplier

**POST** `/supplier`

Request Body:

```json
{
    "supplier_name": "PT. Indofood Tbk",
    "logo":"http://indofood-logo.co.id",
    "nickname":"indofood",
    "address": [{
        "name":"indofood HO",
        "address":"address ho",
        "is_main": true
    },{
        "name":"indofood Cabang JKT 1",
        "address":"address JKT 1",
        "is_main": false
    },{
        "name":"indofood Cabang JKT 2",
        "address":"address JKT 2",
        "is_main": false
    }
    ],
    "contacts":[{
        "name":"Jhon Doe",
        "job_position":"Manage Ops",
        "email":"jhondoe@gmail.com",
        "phone_number":"62221210442",
        "mobile_number":"62258124781442",
        "is_main": true
    },{
        "name":"Jhon Doe 2",
        "job_position":"Spv Ops",
        "email":"jhondoe2@gmail.com",
        "phone_number":"6221210442",
        "mobile_number":"6258124781442",
        "is_main": false
    }],
    "groups":[{
        "group_name":"Industry",
        "value":"Food",
        "active":"true"
    }]
}
```

Response (201 Created):

```json
{
    "code": 201,
    "message": {
        "en": "Suppliers created successfully",
        "id": "Supplier berhasil dibuat"
    },
    "data": {
        "supplier_id": 9,
        "supplier_name": "PT. Indofood Tbk",
        "logo": "http://indofood-logo.co.id",
        "nickname": "indofood",
        "status": "Draft",
        "created_at": "2025-09-22T15:59:48.606394Z",
        "updated_at": "2025-09-22T15:59:48.606394Z",
        "address": "{\"supplier_name\":\"PT. Indofood Tbk\",\"logo\":\"http://indofood-logo.co.id\",\"nickname\":\"indofood\",\"address\":[{\"address_id\":0,\"supplier_id\":0,\"name\":\"indofood HO\",\"address\":\"address ho\",\"is_main\":true},{\"address_id\":0,\"supplier_id\":0,\"name\":\"indofood Cabang JKT 1\",\"address\":\"address JKT 1\",\"is_main\":false},{\"address_id\":0,\"supplier_id\":0,\"name\":\"indofood Cabang JKT 2\",\"address\":\"address JKT 2\",\"is_main\":false}],\"contacts\":[{\"contact_id\":0,\"supplier_id\":0,\"name\":\"Jhon Doe\",\"job_position\":\"Manage Ops\",\"email\":\"jhondoe@gmail.com\",\"phone_number\":\"62221210442\",\"mobile_number\":\"62258124781442\",\"is_main\":true},{\"contact_id\":0,\"supplier_id\":0,\"name\":\"Jhon Doe 2\",\"job_position\":\"Spv Ops\",\"email\":\"jhondoe2@gmail.com\",\"phone_number\":\"6221210442\",\"mobile_number\":\"6258124781442\",\"is_main\":false}],\"groups\":[{\"group_id\":0,\"supplier_id\":0,\"group_name\":\"Industry\",\"value\":\"Food\",\"is_active\":false}],\"materials\":null,\"others\":null}"
    }
}
```
---
### 6. Route Setup (Code)

```go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	handler "supplier-be/handler"
)

// SetupRoutes sets up the routes for the application.
func SetupRoutes(app *fiber.App, supplierHandler *handler.SupplierHandler, log *logrus.Logger) {
	app.Post("/supplier", supplierHandler.CreateSupplierHandler)
	app.Post("/suppliers", supplierHandler.GetListSuppliers)
	app.Get("/", func(c *fiber.Ctx) error {
		response := map[string]string{
			"message": "Health Check OK",
		}
		return c.JSON(response)
	})

}

```

---

### 7. Handler Setup (Code)

```go
package handler

import (
	ms "supplier-be/model/supplier"
	utils "supplier-be/model/utils"
	ss "supplier-be/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// SupplierHandler represents the HTTP handler for supplier-related operations.
type SupplierHandler struct {
	supplierService ss.Supplier
	log             *logrus.Logger
}

// NewSupplierHandler creates a new SupplierHandler instance.
func NewSupplierHandler(supplierService ss.Supplier, logger *logrus.Logger) *SupplierHandler {
	return &SupplierHandler{supplierService: supplierService, log: logger}
}
```

---

### 7.1 Handler CreateSupplier Setup (Code)

```go
// CreateSupplierHandler handles the "create supplier" HTTP request.
func (h *SupplierHandler) CreateSupplierHandler(c *fiber.Ctx) error {
	traceId := utils.GenerateTraceID()

	supplier := new(ms.SupplierRequest)
	if err := c.BodyParser(supplier); err != nil {
		h.log.WithField("Request : ", supplier).Infoln("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, map[string]string{
			"en": "Invalid request payload",
			"id": "Muatan permintaan tidak valid"},
		))
	}

	h.log.WithFields(logrus.Fields{
		"TraceID": traceId,
		"Request": supplier,
	}).Info("CreateSupplierHandler - start")

	createdSupplier, err := h.supplierService.CreateSupplier(c.Context(), supplier, traceId)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"TraceID": traceId,
			"Request": supplier,
		}).Error("CreateSupplierHandler - Error details: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(fiber.StatusInternalServerError, map[string]string{
			"en": "Failed to create supplier",
			"id": "Gagal membuat supplier",
		},
		))
	}

	h.log.WithFields(logrus.Fields{
		"TraceID":  traceId,
		"Request":  supplier,
		"Response": createdSupplier,
	}).Info("CreateSupplierHandler - Success - end")

	return c.Status(fiber.StatusCreated).JSON(utils.SuccessResponse(fiber.StatusCreated, map[string]string{
		"en": "Suppliers created successfully",
		"id": "Supplier berhasil dibuat",
	},
		createdSupplier))
}
```
---

### 7.1.1 Service CreateSupplier Setup (Code)

```go
// CreateSupplier creates a new supplier.
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
					MobileNumber: contact.MobileNumber,
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
```
---
### 7.1.2 Repository CreateSupplier Setup (Code)

```go
package repository

import (
	"context"
	"database/sql"
	"strings"
	ms "supplier-be/model/supplier"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// Repository interface with CRUD operations
type SupplierRepository interface {
	CreateSupplier(ctx context.Context, request *ms.Supplier) (int64, error)
	CreateAddress(ctx context.Context, request *ms.Address) (int64, error)
	CreateContacts(ctx context.Context, request *ms.Contacts) (int64, error)
	CreateGroups(ctx context.Context, request *ms.Groups) (int64, error)
	GetListSuppliers(ctx context.Context, filter *ms.SupplierFilter) ([]ms.SupplierList, error)
}

type SupplierDBRepository struct {
	db *sqlx.DB
}

func NewSupplierRepository(db *sqlx.DB) SupplierRepository {
	return &SupplierDBRepository{db: db}
}
```

```go
// CreateSupplier creates a new supplier in the database
func (r *SupplierDBRepository) CreateSupplier(ctx context.Context, request *ms.Supplier) (int64, error) {
	var (
		param = make([]interface{}, 0)
		id    int64
	)
	param = append(param, request.SupplierName)
	param = append(param, request.Logo)
	param = append(param, request.Nickname)
	param = append(param, request.Status)
	param = append(param, request.CreatedAt)
	param = append(param, request.UpdatedAt)
	param = append(param, request.Address)

	query := `INSERT INTO public.suppliers (supplier_id, supplier_name, nickname, logo, status, created_at, updated_at, address) 
			VALUES(nextval('suppliers_supplier_id_seq'::regclass), ?, ?, ?, ?, ?, ?, ?) RETURNING supplier_id`
	query, args, err := sqlx.In(query, param...)
	if err != nil {
		return 0, err
	}

	query = r.db.Rebind(query)

	res := r.db.QueryRowContext(ctx, query, args...)

	err = res.Err()
	if err != nil {
		return 0, err
	}

	err = res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// CreateAddress creates a new address in the database
func (r *SupplierDBRepository) CreateAddress(ctx context.Context, request *ms.Address) (int64, error) {
	var (
		param = make([]interface{}, 0)
		id    int64
	)
	param = append(param, request.SupplierID)
	param = append(param, request.Name)
	param = append(param, request.Address)
	param = append(param, request.IsMain)

	query := `INSERT INTO public.supplier_addresses
			(address_id, supplier_id, "name", address, is_main)
			VALUES(nextval('supplier_addresses_address_id_seq'::regclass), ?, ?, ?, ?) RETURNING address_id`
	query, args, err := sqlx.In(query, param...)
	if err != nil {
		return 0, err
	}

	query = r.db.Rebind(query)

	res := r.db.QueryRowContext(ctx, query, args...)

	err = res.Err()
	if err != nil {
		return 0, err
	}

	err = res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *SupplierDBRepository) CreateContacts(ctx context.Context, request *ms.Contacts) (int64, error) {
	var (
		param = make([]interface{}, 0)
		id    int64
	)
	param = append(param, request.SupplierID)
	param = append(param, request.Name)
	param = append(param, request.JobPosition)
	param = append(param, request.Email)
	param = append(param, request.PhoneNumber)
	param = append(param, request.MobileNumber)
	param = append(param, request.IsMain)

	query := `INSERT INTO public.supplier_contacts
(contact_id, supplier_id, "name", job_position, email, phone, mobile, is_main)
VALUES(nextval('supplier_contacts_contact_id_seq'::regclass), ?, ?, ?, ?, ?, ?, ?) RETURNING contact_id`
	query, args, err := sqlx.In(query, param...)
	if err != nil {
		return 0, err
	}

	query = r.db.Rebind(query)

	res := r.db.QueryRowContext(ctx, query, args...)

	err = res.Err()
	if err != nil {
		return 0, err
	}

	err = res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *SupplierDBRepository) CreateGroups(ctx context.Context, request *ms.Groups) (int64, error) {
	var (
		param = make([]interface{}, 0)
		id    int64
	)
	param = append(param, request.SupplierID)
	param = append(param, request.GroupName)
	param = append(param, request.Value)
	param = append(param, request.IsActive)
	query := `INSERT INTO public.supplier_groups
(group_id, supplier_id, group_name, value, active)
VALUES(nextval('supplier_groups_group_id_seq'::regclass), ?, ?, ?, ?) RETURNING group_id`
	query, args, err := sqlx.In(query, param...)
	if err != nil {
		return 0, err
	}

	query = r.db.Rebind(query)

	res := r.db.QueryRowContext(ctx, query, args...)

	err = res.Err()
	if err != nil {
		return 0, err
	}

	err = res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
```
---
### 📋 Endpoint: Get List Supplier

**POST** `/suppliers`

Request Body:

```json
{
    "offset": 1, //opsional but default by code 1
    "limit": 10, //opsional but default by code 10
    "supplier_name":"", //opsional
    "status":"" //opsional
}

```

Response (200 OK):

```json
{
    "code": 200,
    "message": {
        "en": "Suppliers retrieved successfully",
        "id": "Suppplier berhasil diambil"
    },
    "data": {
        "total_data": 3,
        "data": [
            {
                "supplier_name": "PT. Indofood Tbk",
                "logo": "",
                "address": "address ho",
                "contact": "Jhon Doe 2",
                "status": "Draft"
            },
            {
                "supplier_name": "PT. Telkom Tbk",
                "logo": "",
                "address": "address ho",
                "contact": "Jhon Doe",
                "status": "Draft"
            },
            {
                "supplier_name": "PT. ByteByteGo Tbk",
                "logo": "",
                "address": "address ho",
                "contact": "Jhon Doe",
                "status": "Draft"
            }
        ]
    }
}

```

---

### 7.2 Handler GetListSupplier Setup (Code)

```go
func (h *SupplierHandler) GetListSupplier(c *fiber.Ctx) error {

	traceId := utils.GenerateTraceID()

	filter := new(ms.SupplierFilter)
	if err := c.BodyParser(filter); err != nil {
		h.log.WithField("Request : ", filter).Infoln("Error parsing request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, map[string]string{
			"en": "Invalid request payload",
			"id": "Muatan permintaan tidak valid"},
		))
	}

	h.log.WithFields(logrus.Fields{
		"TraceID": traceId,
		"Request": filter,
	}).Info("GetListSuppliers - start")

	suppliers, err := h.supplierService.GetListSupplier(c.Context(), filter, traceId)
	if err != nil {
		h.log.WithFields(logrus.Fields{
			"TraceID": traceId,
			"Request": filter,
		}).Error("GetListSuppliers - Error details: ", err)
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(fiber.StatusNotFound, map[string]string{
			"en": "Suppliers not found",
			"id": "Supplier tidak ditemukan",
		}))
	}

	h.log.WithFields(logrus.Fields{
		"TraceID":  traceId,
		"Request":  filter,
		"Response": suppliers,
	}).Info("GetListSuppliers - Success - end")

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(fiber.StatusOK, map[string]string{
		"en": "Suppliers retrieved successfully",
		"id": "Suppplier berhasil diambil",
	}, suppliers))
}

```
---
### 7.2.1 Service GetListSupplier Setup (Code)

```go
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

```
---
### 7.2.1 Repository GetListSupplier Setup (Code)

```go
func (r *SupplierDBRepository) GetListSupplier(ctx context.Context, filter *ms.SupplierFilter) ([]ms.SupplierList, error) {

	var result []ms.SupplierList

	q := `SELECT 
	s.supplier_name,
	sa.address,
	sc.name,
	s.status
	FROM public.suppliers s
	LEFT JOIN public.supplier_addresses sa ON sa.supplier_id = s.supplier_id and sa.is_main = true
	LEFT JOIN public.supplier_contacts sc ON sc.supplier_id = s.supplier_id and sc.is_main = true`

	queryStatement, args2 := buildQueryStatementGetListSuppliers(q, filter)

	query, args, err := sqlx.In(queryStatement, args2...)
	if err != nil {
		return result, err
	}

	query = r.db.Rebind(query)
	err = r.db.Select(&result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return result, err
	}

	return result, nil

}

func buildQueryStatementGetListSuppliers(baseQuery string, filter *ms.SupplierFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filter.SupplierName.Valid && filter.SupplierName.String != "" {
		conditions = append(conditions, "s.supplier_name ILIKE '%' || ? || '%'")
		args = append(args, filter.SupplierName.String)
	}

	if filter.Status.Valid && filter.Status.String != "" {
		conditions = append(conditions, "s.status = ?")
		args = append(args, filter.Status.String)
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
	}

	if filter.Offset.Valid && filter.Limit.Valid && filter.Offset.Int64 != 0 && filter.Limit.Int64 != 0 {
		baseQuery += " OFFSET ((? - 1) * ?) ROWS FETCH NEXT ? ROWS ONLY"
		args = append(args, filter.Offset.Int64, filter.Limit.Int64, filter.Limit.Int64)
	} else {
		baseQuery += " OFFSET ((1 - 1) * 10) ROWS FETCH NEXT 10 ROWS ONLY"
	}

	return baseQuery, args
}

```
---