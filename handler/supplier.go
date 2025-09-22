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
