package admin

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	formutil "gitlab.shoplazza.site/xiabing/goat.git/form"
	"gitlab.shoplazza.site/xiabing/goat.git/i18n"
	"gitlab.shoplazza.site/xiabing/goat.git/middlewares"
)

type GetInstalledAppsForm struct {
	Keyword string `form:"keyword" json:"keyword"`
}

type InstalledAppView struct {
	Id            uint64     `json:"id"`
	AppUri        string     `json:"app_uri"`
	Embbed        bool       `json:"embbed"`
	Icon          string     `json:"icon"`
	Link          string     `json:"link"`
	OauthDancable bool       `json:"oauth_dancable"`
	RedirectUri   string     `json:"redirect_uri"`
	Scopes        []string   `json:"scopes"`
	Status        string     `json:"status"`
	Uid           string     `json:"uid"`
	Subtitle      string     `json:"subtitle"`
	Title         string     `json:"title"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	Owner         string     `json:"owner"`
	Locale        string     `json:"locale"`
	InstalledAt   *time.Time `json:"installed_at"`
}

func (form *GetInstalledAppsForm) Do(c *gin.Context) (interface{}, error) {
	appctx, ok := c.MustGet("context").(*middlewares.Context)
	if !ok {
		return nil, formutil.NewInternalStateError(i18n.T("no_context", appctx.Locale))
	}

	db := config.DB
	count, err := models.GetInstalledAppsCount(db, appctx.StoreID)
	if err != nil {
		log.WithError(err).Errorln("failed to get installed application by store id")
		return nil, formutil.NewInternalStateError(i18n.T("installed_app_count_error", appctx.Locale))
	}

	installedApps, err := models.GetInstalledApps(db, strconv.Itoa(int(appctx.StoreID)), []uint64{})
	if err != nil {
		return nil, err
	}

	appIds := []uint64{}
	for _, v := range installedApps {
		appIds = append(appIds, v.Id)
	}

	appLocales, err := models.GetAppLocalesByAppIds(db, appIds)
	if err != nil {
		return nil, err
	}

	appLocalsMap := models.ResolveApplocale(appLocales, appctx.Locale)

	title := ""
	subTitle := ""
	icon := ""
	installedAppViews := []InstalledAppView{}
	for _, v := range installedApps {
		keyWord := strings.ToLower(form.Keyword)
		appName := strings.ToLower(v.Name)
		var appLocaleTitle string
		locale := appctx.Locale
		if appLocale, ok := appLocalsMap[v.Id]; ok {
			appLocaleTitle = strings.ToLower(appLocale.Name)
			locale = appLocale.Locale
		}

		// 如果app状态为已上架，显示listing信息；否则，显示app信息
		if v.Status == models.APPSTATUS_PUBLISHED && (keyWord == "" || strings.Contains(appLocaleTitle, keyWord)) {
			if appLocale, ok := appLocalsMap[v.Id]; ok {
				title = appLocale.Name
				subTitle = appLocale.Subtitle
				icon = appLocale.Icon
			}
		} else if v.Status != models.APPSTATUS_PUBLISHED && (keyWord == "" || strings.Contains(appName, keyWord)) {
			title = v.Name
			subTitle = ""
			json.Unmarshal([]byte(v.Icon), &icon)
		} else {
			continue
		}

		installedAppView := InstalledAppView{}
		installedAppView.Id = v.Id
		installedAppView.AppUri = v.AppUri
		installedAppView.Embbed = v.Embbed
		installedAppView.Link = v.Link
		installedAppView.OauthDancable = v.OauthDancable
		installedAppView.RedirectUri = v.RedirectUri
		installedAppView.Status = v.Status
		installedAppView.Uid = v.Uid
		installedAppView.InstalledAt = v.InstalledAt
		installedAppView.Scopes = strings.Split(v.Scopes, " ")
		installedAppView.Title = title
		installedAppView.Icon = icon
		installedAppView.Subtitle = subTitle
		installedAppView.Name = v.Name
		installedAppView.Email = v.Email
		installedAppView.Locale = locale

		owner := v.Owner
		if v.PartnerId == 0 {
			owner = ""
		}
		installedAppView.Owner = owner

		installedAppViews = append(installedAppViews, installedAppView)
	}

	return map[string]interface{}{"count": count, "apps": installedAppViews}, nil
}
