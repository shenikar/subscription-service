package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shenikar/subscription-service/internal/dto"
	"github.com/shenikar/subscription-service/internal/logger"
	"github.com/shenikar/subscription-service/internal/mapper"
	"github.com/shenikar/subscription-service/internal/service"
	"github.com/sirupsen/logrus"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// Create godoc
// @Summary Создать подписку
// @Description Создать новую запись о подписке пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.CreateSubscriptionRequest true "Данные подписки"
// @Success 201 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	log := logger.GetLogger()
	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("Create: invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		log.WithError(err).Error("Create: failed to create subscription")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
		return
	}

	log.WithFields(logrus.Fields{
		"id":          sub.ID,
		"serviceName": sub.ServiceName,
		"userID":      sub.UserID,
	}).Info("Subscription created")

	c.JSON(http.StatusCreated, mapper.ToResponseDTO(sub))
}

// GetByID godoc
// @Summary Получить подписку по ID
// @Description Получить запись подписки по ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	log := logger.GetLogger()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.WithError(err).Warn("GetByID: invalid id param")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		log.WithField("id", id).Warn("GetByID: subscription not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	log.WithField("id", id).Info("GetByID: subscription fetched")

	c.JSON(http.StatusOK, mapper.ToResponseDTO(*sub))
}

// GetAll godoc
// @Summary Получить все подписки
// @Description Получить список всех подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} dto.SubscriptionResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAll(c *gin.Context) {
	log := logger.GetLogger()

	subs, err := h.service.GelAll(c.Request.Context())
	if err != nil {
		log.WithError(err).Error("List: failed to list subscriptions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list subscriptions"})
		return
	}

	var res []dto.SubscriptionResponse
	for _, sub := range subs {
		res = append(res, mapper.ToResponseDTO(sub))
	}

	log.WithField("count", len(res)).Info("List: subscriptions listed")

	c.JSON(http.StatusOK, res)
}

// Update godoc
// @Summary Обновить подписку
// @Description Обновить запись подписки по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Param subscription body dto.UpdateSubscriptionRequest true "Обновленные данные подписки"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	log := logger.GetLogger()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.WithError(err).Warn("Update: invalid id param")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("Update: invalid request payload")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		log.WithFields(logrus.Fields{
			"id":  id,
			"err": err,
		}).Error("Update: failed to update subscription")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subscription"})
		return
	}

	log.WithField("id", id).Info("Update: subscription updated")
	c.JSON(http.StatusOK, mapper.ToResponseDTO(sub))
}

// Delete godoc
// @Summary Удалить подписку
// @Description Удалить подписку по ID
// @Tags subscriptions
// @Param id path int true "ID подписки"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	log := logger.GetLogger()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.WithError(err).Warn("Delete: invalid id param")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		log.WithError(err).Error("Delete: failed to delete subscription")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete subscription"})
		return
	}

	log.WithField("id", id).Info("Delete: subscription deleted")
	c.Status(http.StatusNoContent)
}

// TotalPrice godoc
// @Summary Получить суммарную стоимость подписок
// @Description Подсчитывает общую стоимость подписок за период с фильтрацией по user_id и service_name
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "UUID пользователя"
// @Param service_name query string false "Название сервиса"
// @Param from query string false "Дата начала периода (dd-MM-YYYY)"
// @Param to query string false "Дата конца периода (dd-MM-YYYY)"
// @Success 200 {object} map[string]int
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) TotalPrice(c *gin.Context) {
	log := logger.GetLogger()

	var filter dto.TotalPriceFilterDTO
	if err := c.ShouldBindQuery(&filter); err != nil {
		log.WithError(err).Warn("TotalPrice: invalid query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	total, err := h.service.TotalPrice(c.Request.Context(), filter)
	if err != nil {
		log.WithError(err).Error("TotalPrice: failed to calculate total price")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate total price"})
		return
	}

	log.WithFields(logrus.Fields{
		"user_id":      filter.UserID,
		"service_name": filter.ServiceName,
		"from":         filter.FromDate,
		"to":           filter.ToDate,
		"total":        total,
	}).Info("TotalPrice: total calculated")

	c.JSON(http.StatusOK, gin.H{"total": total})
}
