package apiwarp

import (
	"context"
	"dao-template/pkg/errors"
	"dao-template/pkg/validate"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
)

type FilePathData struct {
	FileType string // 文件类型
	Path     string // 文件路径
}

type FileBytesData struct {
	FileType string // 文件类型
	Data     []byte // 文件数据
}

// Controller 泛型控制器包装器：将 (ctx, *In) (*Out, error) 自动包装为标准 gin.HandlerFunc
// 并统一处理参数绑定校验、异常捕获堆栈日志、以及规范的 JSON 响应格式
func Controller[In any, Out any](fc func(ctx context.Context, in *In) (*Out, error)) gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		var in In
		err := gCtx.ShouldBind(&in)
		if err != nil {
			errMap := make(map[string]string)

			var validatorErr validator.ValidationErrors
			if errors.As(err, &validatorErr) {
				for _, fieldErr := range validatorErr {
					errMap[fieldErr.Field()] = fieldErr.Translate(validate.GetTranslator())
				}
			}

			gCtx.JSON(200, map[string]any{
				"code": 1,
				"msg":  "参数错误",
				"data": errMap,
			})
			return
		}

		// 执行业务函数
		out, err := fc(gCtx.Request.Context(), &in)
		if err != nil {
			slog.Default().With("url", gCtx.FullPath()).ErrorContext(gCtx.Request.Context(), fmt.Sprintf("\n%+v\n", err))

			var msgErr *errors.MsgErr
			msg := "系统异常"

			if errors.As(err, &msgErr) {
				msg = msgErr.Error()
			}

			gCtx.JSON(200, map[string]any{
				"code": 1,
				"msg":  msg,
			})
			return
		}

		var outAny any = out

		// 多返回值类型支持（支持文件路径流、字节流等特殊响应）
		switch outData := outAny.(type) {
		case *FilePathData:
			gCtx.Writer.Header().Set("Content-Type", outData.FileType)
			gCtx.File(outData.Path)
		case FilePathData:
			gCtx.Writer.Header().Set("Content-Type", outData.FileType)
			gCtx.File(outData.Path)
		case *FileBytesData:
			gCtx.Data(200, outData.FileType, outData.Data)
		case FileBytesData:
			gCtx.Data(200, outData.FileType, outData.Data)
		default:
			// 默认标准 JSON 响应包装
			gCtx.JSON(200, map[string]any{
				"code": 0,
				"msg":  "",
				"data": out,
			})
		}
	}
}
