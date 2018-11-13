package middleware

import (
	"github.com/gin-gonic/gin"
	"monsterblog/controller"
)

func SiteMessageInject() gin.HandlerFunc {
	return func(c *gin.Context) {
        host:=c.Request.Host
		//globalSiteInfo:=&controller.GlobalSiteInfo{
        	//RecommendBlogs:controller.RecommendBlogGetterMapInstance.Get(host).Get(),
        	//TopNRecommendBlogs:controller.RecommendBlogGetterMapInstance.Get(host).GetTopN(6),
        	//RecommendSeries:controller.RecommendSeriesGetterMapInstance.Get(host).Get(),
		//	TopNRecommendSeries:controller.RecommendSeriesGetterMapInstance.Get(host).GetTopN(6),
        	//MainCategories:controller.MainCategoryGetterMapInstance.Get(host).Get(),
        	//TopNMainCategories:controller.MainCategoryGetterMapInstance.Get(host).GetTopN(4),
        	//SpecialRecommendations:controller.SpecialRecommendationGetterMapInstance.Get(host).Get(),
        	//TopNSpecialRecommendations:controller.SpecialRecommendationGetterMapInstance.Get(host).GetTopN(3),
        	//CommonInfo:controller.CommonInfoGetterInstance.Get(),
		//}
		//c.Set("globalSiteInfo",globalSiteInfo)
		if _,ok:=controller.GlobalSiteInfoMapInstance[host];ok{
			info:=controller.GlobalSiteInfoMapInstance[host]
			info["CommonInfo"]=controller.CommonInfoGetterInstance.Get()
			c.Set("globalSiteInfo",controller.GlobalSiteInfoMapInstance[host])
		}else{
			c.Set("globalSiteInfo",nil)
		}
		c.Next()
	}
}