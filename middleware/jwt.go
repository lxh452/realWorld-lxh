package middleware

// 身份认证中间件
import (
	"realWorld/model/resp"
	"realWorld/utils"

	"github.com/gin-gonic/gin"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取token，并验证token是否为空
		token := utils.GetToken(c)
		if token == "" {
			resp.FailWithMessage("未登录或非法访问", c)
			// c.JSON(403, "令牌不正确")
			c.Abort()
			return
		}
		// 2. 解析token
		j := utils.NewJwt()
		_, err := j.ParseToken(token)
		if err != nil {
			resp.FailWithMessage(err.Error(), c)
			c.Abort()
		}
	}
}
