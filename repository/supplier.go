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
	GetListSupplier(ctx context.Context, filter *ms.SupplierFilter) ([]ms.SupplierList, error)
}

type SupplierDBRepository struct {
	db *sqlx.DB
}

func NewSupplierRepository(db *sqlx.DB) SupplierRepository {
	return &SupplierDBRepository{db: db}
}

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
