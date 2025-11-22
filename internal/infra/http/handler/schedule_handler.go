package handler

import (
	"github.com/labstack/echo/v4"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/infra/http/response"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/pkg/validation"
	"go.bankkrud.com/bankkrud/backend/krudapp/internal/usecase/schedule"
)

type ScheduleHandler struct {
	va *validation.Validator
	uc *schedule.Usecase
}

func NewScheduleHandler(va *validation.Validator, uc *schedule.Usecase) *ScheduleHandler {
	return &ScheduleHandler{
		va: va,
		uc: uc,
	}
}

// GetSchedules swaggo annotation.
//
//	@Summary		Get schedules
//	@Description	Get list of schedules for the authenticated user
//	@Tags			schedule
//	@Accept			json
//	@Produce		json
//	@Param			Authorization		header		string	true	"Authorization token"
//	@Param			transaction_type	query		string	false	"Transaction type"
//	@Param			status				query		string	false	"Schedule status"
//	@Param			limit				query		integer	false	"Limit number of schedules to return"
//	@Param			startID				query		integer	false	"Start id for pagination"
//	@Param			uuid				path		string	true	"Transaction UUID"
//	@Success		200					{object}	response.Response
//	@Failure		400					{object}	response.Response
//	@Failure		404					{object}	response.Response
//	@Failure		500					{object}	response.Response
//	@Router			/schedules [get]
func (h *ScheduleHandler) GetSchedules(ctx echo.Context) error {
	req := new(schedule.GetSchedulesRequest)
	err := ctx.Bind(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.GetSchedules(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	code, res := response.Success(resp.Schedules)
	return ctx.JSON(code, res.SetPagination(&response.Pagination{
		Limit:   req.Limit,
		StartID: req.StartID,
		NextID:  resp.NextID,
	}))
}

// GetSchedule swaggo annotation.
//
//	@Summary		Get schedule by UUID
//	@Description	Get schedule by UUID
//	@Tags			schedule
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Authorization token"
//	@Param			uuid			path		string	true	"Schedule UUID"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/schedules/{uuid} [get]
func (h *ScheduleHandler) GetSchedule(ctx echo.Context) error {
	req := new(schedule.GetScheduleByUUIDRequest)
	err := ctx.Bind(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.GetScheduleByUUID(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}

// CreateSchedule swaggo annotation.
//
//	@Summary		Create a new schedule
//	@Description	Create a new schedule for the authenticated user
//	@Tags			schedule
//	@Accept			json
//	@Produce		json
//	@Param			Authorization			header		string							true	"Authorization token"
//	@Param			CreateScheduleRequest	body		schedule.CreateScheduleRequest	true	"Schedule to create"
//	@Success		200						{object}	response.Response
//	@Failure		400						{object}	response.Response
//	@Failure		500						{object}	response.Response
//	@Router			/schedules [post]
func (h *ScheduleHandler) CreateSchedule(ctx echo.Context) error {
	req := new(schedule.CreateScheduleRequest)
	err := ctx.Bind(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.CreateSchedule(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}

// UpdateScheduleStatus swaggo annotation.
//
//	@Summary		Update a schedule status
//	@Description	Activate or deactivate a schedule
//	@Tags			schedule
//	@Accept			json
//	@Produce		json
//	@Param			Authorization				header		string									true	"Authorization token"
//	@Param			uuid						path		string									true	"UUID of schedule to update"
//	@Param			UpdateScheduleStatusRequest	body		schedule.UpdateScheduleStatusRequest	true	"Status body"
//	@Success		200							{object}	response.Response
//	@Failure		400							{object}	response.Response
//	@Failure		500							{object}	response.Response
//	@Router			/schedules/{uuid} [patch]
func (h *ScheduleHandler) UpdateScheduleStatus(ctx echo.Context) error {
	req := new(schedule.UpdateScheduleStatusRequest)
	err := ctx.Bind(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.UpdateScheduleStatus(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}

// DeleteSchedule swaggo annotation.
//
//	@Summary		Delete schedule by UUID
//	@Description	Delete a schedule by its UUID
//	@Tags			schedule
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Authorization token"
//	@Param			uuid			path		string	true	"Schedule UUID"
//	@Success		200				{object}	response.Response
//	@Failure		400				{object}	response.Response
//	@Failure		404				{object}	response.Response
//	@Failure		500				{object}	response.Response
//	@Router			/schedules/{uuid} [delete]
func (h *ScheduleHandler) DeleteSchedule(ctx echo.Context) error {
	req := new(schedule.DeleteScheduleRequest)
	err := ctx.Bind(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	err = h.va.Validate(req)
	if err != nil {
		return ctx.JSON(response.BadRequest(err))
	}
	resp, err := h.uc.DeleteSchedule(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(response.Error(err))
	}
	return ctx.JSON(response.Success(resp))
}
