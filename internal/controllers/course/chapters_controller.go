package course

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goldlilya1612/diploma-backend/internal/controllers/auth"
	"github.com/goldlilya1612/diploma-backend/internal/models"
	serv "github.com/goldlilya1612/diploma-backend/internal/transport/http"
	"github.com/goldlilya1612/diploma-backend/internal/utils"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (c *Controller) chaptersRoute(rg *gin.RouterGroup) {

	chaptersRouter := rg.Group(chaptersRoute)

	chaptersRouter.POST("/", c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.CreateChapter)
	chaptersRouter.PATCH(chapterParam.toURL(), c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.UpdateChapter)
	chaptersRouter.DELETE(chapterParam.toURL(), c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.DeleteChapters)

	c.articlesRoute(chaptersRouter.Group(chapterParam.toURL()))
}

func (c *Controller) CreateChapter(ctx *gin.Context) {
	var payload *models.CreateChapter

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	courseID := ctx.Param(string(courseParam))

	var creatorID string
	res := c.DB.
		Table("courses").
		Select("creator_id").
		Where("id = ?", courseID).
		Scan(&creatorID)
	if res.Error != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    res.Error.Error(),
		})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	if currentUser.ID.String() != creatorID {
		message := "Access denied"
		ctx.JSON(http.StatusForbidden, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusForbidden,
			Message:    message,
		})
		return
	}

	now := time.Time{}
	courseIDInt, _ := strconv.Atoi(courseID)
	newChapter := models.Chapter{
		Name:     payload.Name,
		CourseID: courseIDInt,

		CreatedAt: now,
		UpdatedAt: now,
	}

	res = c.DB.Create(&newChapter)
	if res.Error != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusConflict,
			Message:    res.Error.Error(),
		})
		return
	}

	chapterResponse := models.ChapterResponse{
		ID:       newChapter.ID,
		Name:     newChapter.Name,
		CourseID: newChapter.CourseID,
		Route:    utils.Latinizer(newChapter.Name),

		CreatedAt: newChapter.CreatedAt,
		UpdatedAt: newChapter.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusCreated,
		Data:       map[string]interface{}{"createdChapter": chapterResponse},
	})
}

func (c *Controller) UpdateChapter(ctx *gin.Context) {
	var payload *models.UpdateChapter

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	chapterID := ctx.Param(string(chapterParam))

	chapter := models.Chapter{}
	res := c.DB.
		Preload("Course").
		First(&chapter, "Chapters.id = ?", chapterID)
	if res.Error != nil && strings.Contains(res.Error.Error(), "record not found") {
		message := fmt.Sprintf("Chapter with id=%s not found", chapterID)
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    message,
		})
		return
	} else if res.Error != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    res.Error.Error(),
		})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	if currentUser.ID != chapter.Course.CreatorID {
		message := "Access denied"
		ctx.JSON(http.StatusForbidden, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusForbidden,
			Message:    message,
		})
		return
	}

	chapter.Name = payload.Name
	res = c.DB.Updates(&chapter)
	if res.Error != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusConflict,
			Message:    res.Error.Error(),
		})
		return
	}

	chapterResponse := models.ChapterResponse{
		ID:       chapter.ID,
		Name:     chapter.Name,
		CourseID: chapter.CourseID,
		Route:    utils.Latinizer(chapter.Name),

		CreatedAt: chapter.CreatedAt,
		UpdatedAt: chapter.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"updatedChapter": chapterResponse},
	})
}

func (c *Controller) DeleteChapters(ctx *gin.Context) {
	var chaptersResponse []models.ChapterResponse

	chapterID := ctx.Param(string(chapterParam))

	chapter := models.Chapter{Articles: []models.Article{}}
	res := c.DB.
		Preload("Articles").
		Joins("Course").
		First(&chapter, "Chapters.id = ?", chapterID)
	if res.Error != nil && strings.Contains(res.Error.Error(), "record not found") {
		message := fmt.Sprintf("Chapter with id=%s not found", chapterID)
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    message,
		})
		return
	} else if res.Error != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    res.Error.Error(),
		})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	if currentUser.ID != chapter.Course.CreatorID {
		message := "Access denied"
		ctx.JSON(http.StatusForbidden, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusForbidden,
			Message:    message,
		})
		return
	}

	res = c.DB.Clauses(clause.Returning{}).Delete(&chapter)
	if res.Error != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusConflict,
			Message:    res.Error.Error(),
		})
		return
	}

	chapterResponse := models.ChapterResponse{
		ID:       chapter.ID,
		Name:     chapter.Name,
		CourseID: chapter.CourseID,
		Route:    utils.Latinizer(chapter.Name),
		Articles: chapter.Articles,
	}
	chaptersResponse = append(chaptersResponse, chapterResponse)

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"deletedChapters": chaptersResponse},
	})
}
