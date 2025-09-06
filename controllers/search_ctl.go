package controllers

import (
	"CloudMusic/controllers/Respond"
	"CloudMusic/utils"
	"github.com/gin-gonic/gin"
	"github.com/meilisearch/meilisearch-go"
	"net/http"
	"strconv"
)

// SearchHandler 搜索（关键词 + 分页 + 高亮）GET /api/search
func SearchHandler(c *gin.Context) {
	index := utils.GetIndex()
	q := c.DefaultQuery("q", "")
	if len(q) == 0 || q == "" {
		Respond.Resp.Fail(c, http.StatusBadRequest, "搜索内容不能为空")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	res, err := index.Search(q, &meilisearch.SearchRequest{
		Limit:                 int64(size),
		Offset:                int64((page - 1) * size),
		AttributesToHighlight: []string{"*"},  // 高亮信息 = * 搜索值的所有字符
		AttributesToRetrieve:  []string{"id"}, // 留空 = 返回全部原始字段 写入什么返回什么
		AttributesToCrop:      nil,            // 摘要信息 = 不要 - 搜索值 分页信息等等
	})
	if err != nil {
		Respond.Resp.Fail(c, http.StatusInternalServerError, "服务器错误")
		return
	}
	Respond.Resp.Success(c, "查询成功", res)
}
