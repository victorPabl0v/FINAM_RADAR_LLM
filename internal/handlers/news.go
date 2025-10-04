package handlers

import (
	"FINAM_RADAR/internal/models"
	parser "FINAM_RADAR/internal/parsers"
	"FINAM_RADAR/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type NewsHandler struct {
	service service.INewsService
}

func NewNewsHandler(service service.INewsService) *NewsHandler {
	return &NewsHandler{
		service: service,
	}
}

// AddNews godoc
// @Summary      Add news
// @Description  Добавляет одну или несколько новостей с черновиком, источниками и сущностями
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        request body models.NewsRequest true "Список новостей"
// @Success      200 {object} service.AddNewsResponse "Успешное добавление"
// @Failure      400 {object} service.AddNewsResponse "Неверный JSON"
// @Failure      500 {object} service.AddNewsResponse "Ошибка сервера"
// @Router       /news [post]
func (h *NewsHandler) AddNews(c echo.Context) error {
	var req models.NewsRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid JSON"})
	}

	ctx := c.Request().Context()

	resp, err := h.service.AddNews(req, ctx)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, resp)
	}

	return c.JSON(http.StatusOK, resp)
}

// GetNews godoc
// @Summary      Получить новости с сайтов
// @Description  Парсит новости с указанных сайтов (например, TASS) и возвращает заголовки, ссылки и текст статей
// @Tags         news
// @Accept       json
// @Produce      json
// @Param        site  query  string  true  "Название сайта (например: tass)"
// @Param        limit query  int     false "Количество новостей (по умолчанию 10)"
// @Success      200 {object} service.GetNewsFromSitesResponse
// @Failure      400 {object} map[string]string "Некорректный запрос или ошибка парсера"
// @Router       /news [get]
func (h *NewsHandler) GetNews(c echo.Context) error {
	var req GetNewsReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid JSON"})
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	var news service.GetNewsFromSitesResponse

	if req.Site == "tass" {
		items, err := parser.ParseEconomyWithArticlesFromTass(req.Limit)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		for _, item := range items {
			news.News = append(news.News, service.News{Title: item.Title, Text: item.Text, URL: item.URL})
		}
	} else {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	return c.JSON(http.StatusOK, news)
}
