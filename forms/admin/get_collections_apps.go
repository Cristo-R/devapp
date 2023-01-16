package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"

	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"gitlab.shoplazza.site/shoplaza/cobra/pkg/errs"
	"gitlab.shoplazza.site/shoplaza/cobra/service"
)

type GetCollectionsAppsForm struct {
	Handle string `form:"handle" json:"handle" binding:"required"`

	applicationRepo models.ApplicationRepo
}

func NewGetCollectionsAppsForm(appRepo models.ApplicationRepo) *GetCollectionsAppsForm {
	return &GetCollectionsAppsForm{
		applicationRepo: appRepo,
	}
}

type CollectionView struct {
	CollectionId          uint64            `json:"collection_id"`
	CollectionName        string            `json:"collection_name"`
	CollectionDescription string            `json:"collection_description"`
	AppsCount             int               `json:"apps_count"`
	Apps                  []*models.AppView `json:"apps"`
}

func (form *GetCollectionsAppsForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	// 根据handle获取一个page
	page, err := models.GetPageByHandle(config.DB, form.Handle)
	if err != nil {
		return nil, formutil.NewInternalStateError(i18n.T("get_page_error", appctx.Locale))
	}
	if page == nil {
		return nil, formutil.NewNotFoundError(i18n.T("page_no_existed", appctx.Locale))
	}

	// 根据page id获取多个collection id
	pageCollects, err := models.GetPageCollectsByPageId(config.DB, page.ID)
	if err != nil {
		return nil, formutil.NewInternalStateError(i18n.T("get_collection_locale_error", appctx.Locale))
	}
	var region string

	user, err := service.GetShoplazzaAccount(appctx.LoginUserID)
	if err != nil {
		region = ""
	}
	region = user.Region
	res := []CollectionView{}

	// 分类下所有应用ID
	var allCollectAppIds []uint64
	// 遍历collection id
	for _, v := range pageCollects {
		// 根据collection id获取apps count
		count, err := models.GetAppsCountByCollectionIdAndRegion(config.DB, v.CollectionId, region)
		if err != nil {
			return nil, formutil.NewInternalStateError(i18n.T("app_count_error", appctx.Locale))
		}

		collectionLocale, err := models.GetCollectionLocale(config.DB, v.CollectionId, appctx.Locale)
		if err != nil {
			return nil, formutil.NewInternalStateError(i18n.T("collection_error", appctx.Locale))
		}
		if collectionLocale == nil {
			continue
		}

		appIds, err := models.GetAppIdsByCollectionIdAndRegion(config.DB, v.CollectionId, 0, 4, user.Region)
		if err != nil {
			return nil, formutil.NewInternalStateError(i18n.T("apps_error", appctx.Locale))
		}

		allCollectAppIds = append(allCollectAppIds, appIds...)

		appLocales, err := models.GetAppLocalesByAppIds(config.DB, appIds)
		if err != nil {
			return nil, err
		}

		appsTags, err := models.GetAppsTagsByIds(config.DB, appIds)
		if err != nil {
			return nil, err
		}

		appsTagsMap := models.TagAppViewSliceGroupByAppid(appsTags)
		appViews := []*models.AppView{}
		appLocalsMap := models.ResolveApplocale(appLocales, appctx.Locale)
		review := models.Reviews{}
		appReviews := []models.AppReviewsWithOverview{}

		appReviews, err = review.GetRatingHistsByStatusAndAppId(config.DB, appIds,
			[]string{models.REVIEWS_STATUS_ACCEPTED,
				models.REVIEWS_STATUS_CREATED,
				models.REVIEWS_STATUS_UPDATED},
			appctx.StoreID)
		if err != nil {
			return nil, formutil.NewInternalStateError(i18n.T("review_error", appctx.Locale))
		}

		reviewMap := models.ReviewsliceConverToMap(appReviews)
		for _, v := range appIds {
			if appLocal, ok := appLocalsMap[v]; ok {
				tmp := models.AppView{
					Id:            appLocal.ApplicationId,
					ListingName:   appLocal.Name,
					ListingIcon:   appLocal.Icon,
					Subtitle:      appLocal.Subtitle,
					OverallRating: 0,
					ReviewTotal:   0,
				}
				if review, ok := reviewMap[v]; ok {
					tmp.OverallRating = review.OverallRating
					tmp.ReviewTotal = review.Total
				}
				if appTags, ok := appsTagsMap[v]; ok {
					tmp.Tags = models.TrimAppTagViews(appTags)
				}

				appViews = append(appViews, &tmp)
			}
		}

		res = append(res, CollectionView{
			CollectionId:          v.CollectionId,
			CollectionName:        collectionLocale.Name,
			CollectionDescription: collectionLocale.Description,
			AppsCount:             count,
			Apps:                  appViews,
		})
	}
	applications, err := form.applicationRepo.QueryApplications(c, models.AqWithIds(allCollectAppIds))
	if err != nil && !errors.Is(err, errs.ParamsError) {
		return nil, err
	}
	appMap := models.TransToApplicationMap(applications)
	for _, cv := range res {
		models.SetAppViewsHandle(cv.Apps, appMap)
		models.SetAppInstallBase(config.DB, appctx.StoreID, allCollectAppIds, cv.Apps, appMap, appctx.Locale)
	}
	return map[string]interface{}{"collections": res}, nil
}
