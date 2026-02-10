package admin

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/narendhupati/bluejay-cms/db/sqlc"
	"github.com/narendhupati/bluejay-cms/internal/services"
)

type AdminContactHandler struct {
	queries *sqlc.Queries
	logger  *slog.Logger
	cache   *services.Cache
}

func NewAdminContactHandler(queries *sqlc.Queries, logger *slog.Logger, cache *services.Cache) *AdminContactHandler {
	return &AdminContactHandler{
		queries: queries,
		logger:  logger,
		cache:   cache,
	}
}

// ListSubmissions displays all contact form submissions with pagination and filtering
func (h *AdminContactHandler) ListSubmissions(c echo.Context) error {
	page := int64(1)
	if v := c.QueryParam("page"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil && parsed > 0 {
			page = parsed
		}
	}

	perPage := int64(25)
	offset := (page - 1) * perPage

	status := c.QueryParam("status")
	submissionType := c.QueryParam("type")
	search := c.QueryParam("search")

	ctx := c.Request().Context()

	// Collect counts for tabs
	totalAll, _ := h.queries.CountContactSubmissions(ctx)
	totalContact, _ := h.queries.CountContactSubmissionsByType(ctx, "contact")
	totalRFQ, _ := h.queries.CountContactSubmissionsByType(ctx, "rfq")
	newCount, _ := h.queries.CountContactSubmissionsByStatus(ctx, "new")

	var submissions []listSubmissionRow
	var totalCount int64

	if search != "" {
		searchParam := sql.NullString{String: search, Valid: true}
		items, err := h.queries.SearchContactSubmissions(ctx, sqlc.SearchContactSubmissionsParams{
			Column1: searchParam,
			Column2: searchParam,
			Column3: searchParam,
			Column4: searchParam,
			Limit:   perPage,
			Offset:  offset,
		})
		if err != nil {
			h.logger.Error("Failed to search contact submissions", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to load submissions")
		}
		for _, item := range items {
			submissions = append(submissions, listSubmissionRow{
				ID: item.ID, Name: item.Name, Email: item.Email, Phone: item.Phone,
				Company: item.Company, InquiryType: item.InquiryType, Status: item.Status,
				SubmissionType: item.SubmissionType, CreatedAt: item.CreatedAt,
			})
		}
		cnt, _ := h.queries.CountContactSubmissionsSearch(ctx, sqlc.CountContactSubmissionsSearchParams{
			Column1: searchParam, Column2: searchParam, Column3: searchParam, Column4: searchParam,
		})
		totalCount = cnt
	} else if status != "" && submissionType != "" {
		items, err := h.queries.ListContactSubmissionsByStatusAndType(ctx, sqlc.ListContactSubmissionsByStatusAndTypeParams{
			Status: status, SubmissionType: submissionType, Limit: perPage, Offset: offset,
		})
		if err != nil {
			h.logger.Error("Failed to list contact submissions", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to load submissions")
		}
		for _, item := range items {
			submissions = append(submissions, listSubmissionRow{
				ID: item.ID, Name: item.Name, Email: item.Email, Phone: item.Phone,
				Company: item.Company, InquiryType: item.InquiryType, Status: item.Status,
				SubmissionType: item.SubmissionType, CreatedAt: item.CreatedAt,
			})
		}
		cnt, _ := h.queries.CountContactSubmissionsByStatusAndType(ctx, sqlc.CountContactSubmissionsByStatusAndTypeParams{
			Status: status, SubmissionType: submissionType,
		})
		totalCount = cnt
	} else if status != "" {
		items, err := h.queries.ListContactSubmissionsByStatus(ctx, sqlc.ListContactSubmissionsByStatusParams{
			Status: status, Limit: perPage, Offset: offset,
		})
		if err != nil {
			h.logger.Error("Failed to list contact submissions", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to load submissions")
		}
		for _, item := range items {
			submissions = append(submissions, listSubmissionRow{
				ID: item.ID, Name: item.Name, Email: item.Email, Phone: item.Phone,
				Company: item.Company, InquiryType: item.InquiryType, Status: item.Status,
				SubmissionType: item.SubmissionType, CreatedAt: item.CreatedAt,
			})
		}
		cnt, _ := h.queries.CountContactSubmissionsByStatus(ctx, status)
		totalCount = cnt
	} else if submissionType != "" {
		items, err := h.queries.ListContactSubmissionsByType(ctx, sqlc.ListContactSubmissionsByTypeParams{
			SubmissionType: submissionType, Limit: perPage, Offset: offset,
		})
		if err != nil {
			h.logger.Error("Failed to list contact submissions", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to load submissions")
		}
		for _, item := range items {
			submissions = append(submissions, listSubmissionRow{
				ID: item.ID, Name: item.Name, Email: item.Email, Phone: item.Phone,
				Company: item.Company, InquiryType: item.InquiryType, Status: item.Status,
				SubmissionType: item.SubmissionType, CreatedAt: item.CreatedAt,
			})
		}
		cnt, _ := h.queries.CountContactSubmissionsByType(ctx, submissionType)
		totalCount = cnt
	} else {
		items, err := h.queries.ListContactSubmissions(ctx, sqlc.ListContactSubmissionsParams{
			Limit: perPage, Offset: offset,
		})
		if err != nil {
			h.logger.Error("Failed to list contact submissions", "error", err)
			return c.String(http.StatusInternalServerError, "Failed to load submissions")
		}
		for _, item := range items {
			submissions = append(submissions, listSubmissionRow{
				ID: item.ID, Name: item.Name, Email: item.Email, Phone: item.Phone,
				Company: item.Company, InquiryType: item.InquiryType, Status: item.Status,
				SubmissionType: item.SubmissionType, CreatedAt: item.CreatedAt,
			})
		}
		totalCount = totalAll
	}

	totalPages := int64(math.Ceil(float64(totalCount) / float64(perPage)))

	return c.Render(http.StatusOK, "admin/pages/contact_submissions_list.html", map[string]interface{}{
		"Title":        "Contact Submissions",
		"Submissions":  submissions,
		"Page":         page,
		"TotalPages":   totalPages,
		"TotalCount":   totalCount,
		"Status":       status,
		"Type":         submissionType,
		"Search":       search,
		"TotalAll":     totalAll,
		"TotalContact": totalContact,
		"TotalRFQ":     totalRFQ,
		"NewCount":     newCount,
	})
}

type listSubmissionRow struct {
	ID             int64
	Name           string
	Email          string
	Phone          string
	Company        string
	InquiryType    sql.NullString
	Status         string
	SubmissionType string
	CreatedAt      interface{}
}

// ViewSubmission displays a single contact submission
func (h *AdminContactHandler) ViewSubmission(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid submission ID")
	}

	ctx := c.Request().Context()

	submission, err := h.queries.GetContactSubmissionByID(ctx, id)
	if err != nil {
		h.logger.Error("Failed to get contact submission", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load submission")
	}

	// Get prev/next IDs for navigation
	prevID := int64(0)
	nextID := int64(0)
	if pid, err := h.queries.GetPreviousSubmissionID(ctx, id); err == nil {
		prevID = pid
	}
	if nid, err := h.queries.GetNextSubmissionID(ctx, id); err == nil {
		nextID = nid
	}

	return c.Render(http.StatusOK, "admin/pages/contact_submission_detail.html", map[string]interface{}{
		"Title":      "Contact Submission",
		"Submission": submission,
		"PrevID":     prevID,
		"NextID":     nextID,
	})
}

// UpdateSubmissionStatus updates the status of a contact submission
func (h *AdminContactHandler) UpdateSubmissionStatus(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid submission ID")
	}

	status := c.FormValue("status")
	notes := c.FormValue("notes")

	err = h.queries.UpdateContactSubmissionStatus(c.Request().Context(), sqlc.UpdateContactSubmissionStatusParams{
		Status: status,
		Notes:  sql.NullString{String: notes, Valid: notes != ""},
		ID:     id,
	})
	if err != nil {
		h.logger.Error("Failed to update submission status", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to update status")
	}

	logActivity(c, "updated", "contact_submission", id, "", "Updated Contact Submission #%d status", id)
	return c.Redirect(http.StatusSeeOther, "/admin/contact/submissions/"+c.Param("id"))
}

// BulkMarkRead marks all new submissions as read
func (h *AdminContactHandler) BulkMarkRead(c echo.Context) error {
	err := h.queries.BulkMarkContactSubmissionsRead(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to bulk mark submissions as read", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to update")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/contact/submissions")
}

// DeleteSubmission deletes a contact submission
func (h *AdminContactHandler) DeleteSubmission(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid submission ID")
	}

	err = h.queries.DeleteContactSubmission(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete contact submission", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to delete submission")
	}

	logActivity(c, "deleted", "contact_submission", id, "", "Deleted Contact Submission #%d", id)
	return c.NoContent(http.StatusOK)
}

// ListOffices displays all office locations
func (h *AdminContactHandler) ListOffices(c echo.Context) error {
	offices, err := h.queries.ListAllOfficeLocations(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to list office locations", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load offices")
	}

	return c.Render(http.StatusOK, "admin/pages/office_locations_list.html", map[string]interface{}{
		"Title":   "Office Locations",
		"Offices": offices,
	})
}

// NewOffice displays the form for creating a new office
func (h *AdminContactHandler) NewOffice(c echo.Context) error {
	return c.Render(http.StatusOK, "admin/pages/office_locations_form.html", map[string]interface{}{
		"Title":      "New Office Location",
		"FormAction": "/admin/contact/offices",
		"Item":       nil,
		"IsNew":      true,
	})
}

// CreateOffice handles office location creation
func (h *AdminContactHandler) CreateOffice(c echo.Context) error {
	name := c.FormValue("name")
	addressLine1 := c.FormValue("address_line1")
	addressLine2 := c.FormValue("address_line2")
	city := c.FormValue("city")
	state := c.FormValue("state")
	postalCode := c.FormValue("postal_code")
	country := c.FormValue("country")
	phone := c.FormValue("phone")
	email := c.FormValue("email")
	mapURL := c.FormValue("map_url")
	_ = mapURL // map_url not in DB yet, placeholder

	isPrimary := int64(0)
	if v := c.FormValue("is_primary"); v == "on" || v == "1" {
		isPrimary = 1
	}

	isActive := int64(0)
	if v := c.FormValue("is_active"); v == "on" || v == "1" {
		isActive = 1
	}

	displayOrder := int64(0)
	if v := c.FormValue("display_order"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			displayOrder = parsed
		}
	}

	// If setting as primary, unset all others first
	if isPrimary == 1 {
		err := h.queries.UnsetPrimaryOfficeLocations(c.Request().Context())
		if err != nil {
			h.logger.Error("Failed to unset primary office locations", "error", err)
		}
	}

	params := sqlc.CreateOfficeLocationParams{
		Name:         name,
		AddressLine1: addressLine1,
		AddressLine2: sql.NullString{String: addressLine2, Valid: addressLine2 != ""},
		City:         city,
		State:        state,
		PostalCode:   postalCode,
		Country:      country,
		Phone:        sql.NullString{String: phone, Valid: phone != ""},
		Email:        sql.NullString{String: email, Valid: email != ""},
		IsPrimary:    isPrimary,
		IsActive:     isActive,
		DisplayOrder: displayOrder,
	}

	_, err := h.queries.CreateOfficeLocation(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to create office location", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to create office")
	}

	h.cache.DeleteByPrefix("page:contact")
	logActivity(c, "created", "office", 0, c.FormValue("name"), "Created Office '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/contact/offices")
}

// EditOffice displays the form for editing an office location
func (h *AdminContactHandler) EditOffice(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid office ID")
	}

	office, err := h.queries.GetOfficeLocationByID(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to get office location", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to load office")
	}

	return c.Render(http.StatusOK, "admin/pages/office_locations_form.html", map[string]interface{}{
		"Title":      "Edit Office Location",
		"FormAction": fmt.Sprintf("/admin/contact/offices/%d", id),
		"Item":       office,
		"IsNew":      false,
	})
}

// UpdateOffice handles office location updates
func (h *AdminContactHandler) UpdateOffice(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid office ID")
	}

	name := c.FormValue("name")
	addressLine1 := c.FormValue("address_line1")
	addressLine2 := c.FormValue("address_line2")
	city := c.FormValue("city")
	state := c.FormValue("state")
	postalCode := c.FormValue("postal_code")
	country := c.FormValue("country")
	phone := c.FormValue("phone")
	email := c.FormValue("email")

	isPrimary := int64(0)
	if v := c.FormValue("is_primary"); v == "on" || v == "1" {
		isPrimary = 1
	}

	isActive := int64(0)
	if v := c.FormValue("is_active"); v == "on" || v == "1" {
		isActive = 1
	}

	displayOrder := int64(0)
	if v := c.FormValue("display_order"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			displayOrder = parsed
		}
	}

	// If setting as primary, unset all others first
	if isPrimary == 1 {
		err := h.queries.UnsetPrimaryOfficeLocations(c.Request().Context())
		if err != nil {
			h.logger.Error("Failed to unset primary office locations", "error", err)
		}
	}

	params := sqlc.UpdateOfficeLocationParams{
		Name:         name,
		AddressLine1: addressLine1,
		AddressLine2: sql.NullString{String: addressLine2, Valid: addressLine2 != ""},
		City:         city,
		State:        state,
		PostalCode:   postalCode,
		Country:      country,
		Phone:        sql.NullString{String: phone, Valid: phone != ""},
		Email:        sql.NullString{String: email, Valid: email != ""},
		IsPrimary:    isPrimary,
		IsActive:     isActive,
		DisplayOrder: displayOrder,
		ID:           id,
	}

	err = h.queries.UpdateOfficeLocation(c.Request().Context(), params)
	if err != nil {
		h.logger.Error("Failed to update office location", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to update office")
	}

	h.cache.DeleteByPrefix("page:contact")
	logActivity(c, "updated", "office", id, c.FormValue("name"), "Updated Office '%s'", c.FormValue("name"))
	return c.Redirect(http.StatusSeeOther, "/admin/contact/offices")
}

// DeleteOffice handles office location deletion
func (h *AdminContactHandler) DeleteOffice(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid office ID")
	}

	err = h.queries.DeleteOfficeLocation(c.Request().Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete office location", "error", err)
		return c.String(http.StatusInternalServerError, "Failed to delete office")
	}

	h.cache.DeleteByPrefix("page:contact")
	logActivity(c, "deleted", "office", id, "", "Deleted Office #%d", id)
	return c.NoContent(http.StatusOK)
}
