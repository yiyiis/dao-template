package jwt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"log/slog"
	"slices"
	"strings"
)

func CheckLogin(gCtx *gin.Context) {
	// 白名单路径直接放行
	if slices.Contains(config.IgnorePath, gCtx.FullPath()) {
		gCtx.Next()
		return
	}

	if !strings.HasPrefix(gCtx.FullPath(), "/api") {
		gCtx.Next()
		return
	}

	token := gCtx.Request.Header.Get("Authorization")

	parse, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Secret), nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		slog.Error("jwt验证失败", "err", err, "jwt", token)
		gCtx.Status(401)
		gCtx.Abort()
		return
	}

	var tc TokenClaims
	split := strings.Split(parse.Raw, ".")
	if len(split) < 2 {
		slog.Error("jwt结构不合法", "jwt", token)
		gCtx.Status(401)
		gCtx.Abort()
		return
	}

	jsonData, err := io.ReadAll(base64.NewDecoder(base64.RawStdEncoding, strings.NewReader(split[1])))
	if err != nil {
		slog.Error("jwt base64解码错误", "err", err)
		gCtx.Status(401)
		gCtx.Abort()
		return
	}

	err = json.Unmarshal(jsonData, &tc)
	if err != nil {
		slog.Error("jwt payload解析错误", "err", err)
		gCtx.Status(401)
		gCtx.Abort()
		return
	}

	ctx := gCtx.Request.Context()
	gCtx.Request = gCtx.Request.WithContext(context.WithValue(ctx, "tokenClaims", tc))
}
