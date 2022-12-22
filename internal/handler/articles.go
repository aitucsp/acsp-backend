package handler

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"acsp/internal/apperror"
	"acsp/internal/dto"
	"acsp/internal/logs"
	"acsp/internal/validation"
)

// @Summary Create an article
// @Security ApiKeyAuth
// @Tags articles
// @Description Method for creating an article
// @ID project-id
// @Accept  json
// @Produce  json
// @Param input body model.Article true "project information"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/articles [post]
func (h *Handler) createArticle(c *fiber.Ctx) error {
	log.Println("Creating an article... ")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	input := dto.CreateArticle{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrBadInputBody,
		})
	}

	// validate := validator.New()
	// if err := validate.Struct(input); err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
	// 		"errors":  true,
	// 		"message": validation.ValidatorErrors(err),
	// 	})
	// }

	id, err := h.services.Articles.Create(c.UserContext(), userId, input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"errors":    false,
		"message":   nil,
		"projectId": id,
	})
}

// @Summary Get all articles by user id
// @Security ApiKeyAuth
// @Tags articles
// @Description Get all articles of user
// @ID get-all-articles-by-user-id
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/articles [get]
func (h *Handler) getAllProjects(c *fiber.Ctx) error {
	logs.Log().Info("Getting all articles... ")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	articles, err := h.services.Articles.GetAll(c.UserContext(), userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"errors":   false,
		"message":  nil,
		"count":    len(articles),
		"user_id":  c.GetRespHeader(userCtx, ""),
		"articles": articles,
	})
}

// @Summary Update an article by id
// @Security ApiKeyAuth
// @Tags articles
// @Description Update data of article
// @ID update-title-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/title/:id [put]
func (h *Handler) updateArticle(c *fiber.Ctx) error {
	log.Println("Updating a project... ")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	var input dto.UpdateArticle
	projectId := c.Params("id", "")
	if projectId == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrBodyParsed,
		})
	}

	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": validation.ValidatorErrors(err),
		})
	}

	result, err := h.services.Articles.Update(c.UserContext(), userId, input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"result": result})
}

// @Summary Delete an article by id
// @Security ApiKeyAuth
// @Tags articles
// @Description Delete an article from database
// @ID delete-article-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Failure default {object} map[string]interface{}
// @Router /api/v1/article/:id [delete]
func (h *Handler) deleteArticle(c *fiber.Ctx) error {
	logs.Log().Info("Deleting an article")

	userId, err := getUserId(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	id := c.Params("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors":  true,
			"message": apperror.ErrParameterNotFound,
		})
	}

	result, err := h.services.Articles.Delete(c.UserContext(), userId, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"errors":  true,
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusNoContent).JSON(fiber.Map{"result": result})
}
