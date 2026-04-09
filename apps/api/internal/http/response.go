package httpapi

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, err error) {
    c.JSON(http.StatusInternalServerError, gin.H{
        "status": "error",
        "detail": err.Error(),
    })
}

func respondBadRequest(c *gin.Context, detail string) {
    c.JSON(http.StatusBadRequest, gin.H{
        "status": "error",
        "detail": detail,
    })
}

func respondNotFound(c *gin.Context, detail string) {
    c.JSON(http.StatusNotFound, gin.H{
        "status": "error",
        "detail": detail,
    })
}
