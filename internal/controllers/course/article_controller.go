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

func (c *Controller) articlesRoute(rg *gin.RouterGroup) {

	articlesRouter := rg.Group(articlesRoute)
	const contentURL = "/content"

	articlesRouter.POST("/", c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.CreateArticle)
	articlesRouter.PATCH(articleParam.toURL(), c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.UpdateArticle)
	articlesRouter.DELETE(articleParam.toURL(), c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.DeleteArticles)

	articlesRouter.PATCH(articleParam.toURL()+contentURL, c.authController.DeserializeUser(), c.authController.CheckAccessRole(auth.LecturerRole), c.UpdateContent)
	articlesRouter.GET(articleParam.toURL()+contentURL, c.authController.DeserializeUser(), c.GetContent)

}

func (c *Controller) CreateArticle(ctx *gin.Context) {
	var payload *models.CreateArticle

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

	var creatorID string
	res := c.DB.
		Table("courses").
		Select("creator_id").
		Joins("JOIN Chapters ON Chapters.course_id = Courses.id").
		Where("Chapters.id = ?", chapterID).
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
	chapterIDInt, _ := strconv.Atoi(chapterID)

	newArticle := models.Article{
		Name:      payload.Name,
		ChapterID: chapterIDInt,

		CreatedAt: now,
		UpdatedAt: now,
	}

	res = c.DB.Create(&newArticle)
	if res.Error != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusConflict,
			Message:    res.Error.Error(),
		})
		return
	}

	articleResponse := &models.ArticleResponse{
		ID:        newArticle.ID,
		Name:      newArticle.Name,
		ChapterID: newArticle.ChapterID,
		Route:     utils.Latinizer(newArticle.Name),

		CreatedAt: newArticle.CreatedAt,
		UpdatedAt: newArticle.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusCreated,
		Data:       map[string]interface{}{"createdArticle": articleResponse},
	})
}

func (c *Controller) UpdateArticle(ctx *gin.Context) {
	var payload *models.UpdateArticle

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	articleID := ctx.Param(string(articleParam))

	article := models.Article{}
	res := c.DB.
		Preload("Chapter.Course").
		First(&article, "Articles.id = ?", articleID)
	if res.Error != nil && strings.Contains(res.Error.Error(), "record not found") {
		message := fmt.Sprintf("Article with id=%s not found", articleID)
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
	if currentUser.ID != article.Chapter.Course.CreatorID {
		message := "Access denied"
		ctx.JSON(http.StatusForbidden, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusForbidden,
			Message:    message,
		})
		return
	}

	article.Name = payload.Name
	res = c.DB.Updates(&article)
	if res.Error != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusConflict,
			Message:    res.Error.Error(),
		})
		return
	}

	articleResponse := &models.ArticleResponse{
		ID:        article.ID,
		Name:      article.Name,
		ChapterID: article.ChapterID,
		Route:     utils.Latinizer(article.Name),

		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"updatedArticles": articleResponse},
	})
}

func (c *Controller) DeleteArticles(ctx *gin.Context) {
	var articlesResponse []models.ArticleResponse

	params := ctx.Request.URL.Query()
	ids := params["id"]

	for _, id := range ids {

		article := models.Article{}
		res := c.DB.
			Preload("Chapter.Course").
			First(&article, "Articles.id = ?", id)
		if res.Error != nil {
			ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusBadRequest,
				Message:    res.Error.Error(),
			})
			return
		}

		currentUser := ctx.MustGet("currentUser").(models.User)
		if currentUser.ID != article.Chapter.Course.CreatorID {
			message := "Access denied"
			ctx.JSON(http.StatusForbidden, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusForbidden,
				Message:    message,
			})
			return
		}

		res = c.DB.Clauses(clause.Returning{}).Delete(&article)
		if res.Error != nil {
			ctx.JSON(http.StatusConflict, models.HTTPResponse{
				Status:     serv.ErrResponseStatus,
				StatusCode: http.StatusConflict,
				Message:    res.Error.Error(),
			})
			return
		}

		articleResponse := models.ArticleResponse{
			ID:        article.ID,
			Name:      article.Name,
			ChapterID: article.ChapterID,
			Route:     utils.Latinizer(article.Name),
		}
		articlesResponse = append(articlesResponse, articleResponse)
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"deletedArticles": articlesResponse},
	})
}

func (c *Controller) UpdateContent(ctx *gin.Context) {
	var payload *models.UpdateContent

	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	articleID := ctx.Param(string(articleParam))

	article := models.Article{}
	res := c.DB.
		Preload("Chapter.Course").
		First(&article, "Articles.id = ?", articleID)
	if res.Error != nil && strings.Contains(res.Error.Error(), "record not found") {
		message := fmt.Sprintf("Article with id=%s not found", articleID)
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
	if currentUser.ID != article.Chapter.Course.CreatorID {
		message := "Access denied"
		ctx.JSON(http.StatusForbidden, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusForbidden,
			Message:    message,
		})
		return
	}

	article.Content = payload.Content
	res = c.DB.Updates(&article)
	if res.Error != nil {
		ctx.JSON(http.StatusConflict, models.HTTPResponse{
			Status:     serv.ErrResponseStatus,
			StatusCode: http.StatusConflict,
			Message:    res.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Message:    "content updated",
		Data: map[string]interface{}{
			"courseID":  article.Chapter.Course.ID,
			"chapterID": article.Chapter.ID,
			"articleID": article.ID,
		},
	})
}

func (c *Controller) GetContent(ctx *gin.Context) {

	articleID := ctx.Param(string(articleParam))

	content := ""
	res := c.DB.Table("articles").Select("content").Where("id=?", articleID).Scan(&content)
	if res.Error != nil && strings.Contains(res.Error.Error(), "record not found") {
		message := fmt.Sprintf("Article with id=%s not found", articleID)
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

	ctx.JSON(http.StatusOK, models.HTTPResponse{
		Status:     serv.SuccessResponseStatus,
		StatusCode: http.StatusOK,
		Data:       map[string]interface{}{"content": content},
	})
}
