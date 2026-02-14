package handlers

import (
	"log"
	"net/http"
	"strconv"

	authorizationService "github.com/darksuei/suei-intelligence/internal/application/authorization"
	datasourceService "github.com/darksuei/suei-intelligence/internal/application/datasource"
	"github.com/darksuei/suei-intelligence/internal/config"
	authorizationDomain "github.com/darksuei/suei-intelligence/internal/domain/authorization"
	datasourceDomain "github.com/darksuei/suei-intelligence/internal/domain/datasource"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/etl"
	"github.com/darksuei/suei-intelligence/internal/infrastructure/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

func SupportedDatasources(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"datasources": datasourceDomain.SupportedDatasources,
	})
	return
}

func NewDatasource(c *gin.Context) {
	var req struct {
		SourceType string `json:"sourceType" binding:"required"`
		Connection interface{} `json:"connection" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid request: Missing required fields.",
		})
		return
	}

	key := c.Param("key") // assumes route is like /projects/:key
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project key is required",
		})
		return
	}

	// Authorization
	allow, err := authorizationService.EnforceRoles(utils.GetUserRolesFromContext(c), "org", authorizationDomain.Organization, "write")

	if err != nil || !allow {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
		})
		return
	}

	// Validate datasource
	// - sourceType must be supported
	// - connection params must be complete and valid

	var databaseCfg config.DatabaseConfig
	if err := envconfig.Process("", &databaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	var _datasource *datasourceDomain.Datasource
	var datasourceID uint

	createdByEmail, err := utils.GetUserEmailFromContext(c)

	if err != nil || createdByEmail == nil{
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get account",
		})
		return
	}

	// Create datasource
	_datasource, err = datasourceService.NewDatasource(key, *createdByEmail, &databaseCfg)

	datasourceID = _datasource.ID

	if err != nil {
		log.Printf("Error creating datasource: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = etl.GetInstance().CreateSourceConnection(strconv.FormatUint(uint64(datasourceID), 10), req.Connection)

	if err != nil {
		// Rollback CREATED Datasource
		err = datasourceService.HardDeleteDatasource(uint(datasourceID), key, &databaseCfg)

		log.Printf("Error creating datasource: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"datasource": _datasource,
	})
	return
}

func RetrieveDatasources(c *gin.Context) {
	key := c.Param("key") // assumes route is like /projects/:key
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project key is required",
		})
		return
	}

	var databaseCfg config.DatabaseConfig
	if err := envconfig.Process("", &databaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// Authorization
	allow, err := authorizationService.EnforceRoles(utils.GetUserRolesFromContext(c), "org", authorizationDomain.Organization, "read")

	if err != nil || !allow {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
		})
		return
	}

	// Retrieve datasources
	_datasources, err := datasourceService.RetrieveDatasources(key, &databaseCfg)

	if err != nil {
		log.Printf("Error retrieving datasources: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"datasources": _datasources,
	})
}

func DeleteDatasource(c *gin.Context) {
	key := c.Param("key") // assumes route is like /projects/:key
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project key is required",
		})
		return
	}

	idParam := c.Param("id") // /projects/:id
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Project id is required",
		})
		return
	}

	datasourceID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid datasource id",
		})
		return
	}

	var databaseCfg config.DatabaseConfig
	if err := envconfig.Process("", &databaseCfg); err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// Authorization
	allow, err := authorizationService.EnforceRoles(utils.GetUserRolesFromContext(c), "org", authorizationDomain.Organization, "write")

	if err != nil || !allow {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
		})
		return
	}

	// Delete datasource
	err = datasourceService.SoftDeleteDatasource(uint(datasourceID), key, &databaseCfg)

	if err != nil {
		log.Printf("Error deleting datasource: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}