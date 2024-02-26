package v1

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RoutesV1Init(router *gin.Engine, db *gorm.DB) {

	_ = router.Group("/v1")
	{
	}
}
