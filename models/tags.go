package models

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type AppTag struct {
	Id            uint64    `gorm:"PRIMARY_KEY" json:"id"`
	ApplicationId uint64    `json:"application_id"`
	TagId         uint64    `json:"tag_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (at AppTag) TableName() string {
	return "app_tags"
}

type Tag struct {
	Id        uint64    `gorm:"PRIMARY_KEY" json:"id"`
	NameZh    string    `json:"name_zh"`
	TagLevel  string    `json:"tag_level"`
	ParentId  uint64    `json:"parent_id"`
	NameEn    string    `json:"name_en"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TagView struct {
	Id     uint64    `json:"id"`
	NameZh string    `json:"name_zh"`
	NameEn string    `json:"name_en"`
	Tags   []TagView `json:"tags"`
}

type TagInfo struct {
	ApplicationId uint64 `json:"application_id"`
	TagId         uint64 `json:"tag_id"`
	NameZh        string `json:"name_zh"`
	NameEn        string `json:"name_en"`
	TagLevel      string `json:"tag_level"`
	ParentId      uint64 `json:"parent_id"`
	Num           int    `json:"num"`
}

func GetAppTags(db *gorm.DB, id uint64) ([]AppTag, error) {
	appTags := []AppTag{}
	if err := db.Where("application_id = ?", id).Find(&appTags).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return appTags, nil
}

func GetAppTagsByIds(db *gorm.DB, ids []uint64) ([]TagView, error) {
	tagViews := []TagView{}
	if err := db.Table("app_tags at").
		Select("at.application_id as id, t.name_zh, t.name_en").
		Joins("left join tags t on at.tag_id = t.id").
		Where("t.parent_id = 0 ").
		Where("application_id in (?)", ids).Find(&tagViews).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return tagViews, nil
}

func GetTag(db *gorm.DB, id uint64) (*Tag, error) {
	tag := Tag{}
	if err := db.Where("id = ?", id).Find(&tag).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &tag, nil
}

func GetTags(db *gorm.DB) ([]Tag, error) {
	tags := []Tag{}
	if err := db.Find(&tags).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return tags, nil
}

func GetTagsByParentId(db *gorm.DB, parentId []uint64) ([]Tag, error) {
	tags := []Tag{}
	if err := db.Where("parent_id in (?)", parentId).Order("parent_id, id").Find(&tags).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return tags, nil
}

func GetTagsByLevel(db *gorm.DB, level string) ([]Tag, error) {
	tags := []Tag{}
	if err := db.Where("tag_level = ?", level).Find(&tags).Order("id").Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return tags, nil

}

func ConvertTagsToTagViews(tags []Tag) ([]TagView, error) {

	tagsView := make([]TagView, 0)
	for _, tagsValue := range tags {
		//ParentId 为0说明数据 为没有父级 可以作为一级数据。
		if 0 == tagsValue.ParentId {
			t := make([]TagView, 0)
			tagsView = append(tagsView, TagView{
				Id:     tagsValue.Id,
				NameZh: tagsValue.NameZh,
				NameEn: tagsValue.NameEn,
				Tags:   t,
			})
			continue
		}

		//for tagsViewKey, tagsViewValue := range tagsView {
		//	// 如果传入数据tags当中一条父级id 等于 传出结构的id那么说明已经存在可以叠加
		//	if tagsValue.ParentId == tagsViewValue.Id {
		//		tagsView[tagsViewKey].Tags = append(tagsView[tagsViewKey].Tags, TagView{
		//			Id:     tagsValue.Id,
		//			NameZh: tagsValue.NameZh,
		//			NameEn: tagsValue.NameEn,
		//		})
		//		continue
		//	}
		//
		//}

		tagsView = handleTagMapByParentId(tagsValue, tagsView)

	}
	return tagsView, nil

}
func CreateAppTags(db *gorm.DB, tags []*AppTag) error {

	var buffer bytes.Buffer
	sql := "INSERT INTO app_tags " +
		"(application_id, tag_id, created_at, updated_at) " +
		"VALUES"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}

	for i, tag := range tags {
		if i == len(tags)-1 {
			buffer.WriteString(fmt.Sprintf("('%d','%d',now(),now()) ",
				tag.ApplicationId, tag.TagId))
			buffer.WriteString(" ON duplicate KEY UPDATE application_id = values(application_id);")
		} else {
			buffer.WriteString(fmt.Sprintf("('%d','%d',now(),now()),",
				tag.ApplicationId, tag.TagId))
		}
	}

	return db.Exec(buffer.String()).Error
}
func DeleteAppTag(db *gorm.DB, ids []uint64) error {
	if err := db.Where("tag_id in (?)", ids).Delete(AppTag{}).Error; err != nil {
		return err
	}

	return nil
}

type TagInfos struct {
	ApplicationId uint64 `json:"application_id"`
	TagId         uint64 `json:"tag_id"`
	ParentId      uint64 `json:"parent_id"`
}

func GetAppTagsParent(db *gorm.DB) (taginfos []*TagInfos, err error) {
	if err := db.Table("app_tags a").
		Select("a.application_id, b.parent_id tag_id, c.parent_id ").
		Joins("inner join tags b on a.tag_id = b.id ").
		Joins("inner join tags c on b.parent_id = c.id ").
		Find(&taginfos).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return
}
func UpdateAppTag(db *gorm.DB, updateMap map[uint64]uint64) error {
	var buffer bytes.Buffer

	sql := "update app_tags set tag_id = case tag_id "
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for oldV, newV := range updateMap {
		buffer.WriteString(fmt.Sprintf("when %d then %d\n", oldV, newV))
	}
	keys := make([]uint64, 0, len(updateMap))
	for k := range updateMap {
		keys = append(keys, k)
	}
	buffer.WriteString("end\n;")
	return db.Where(" tag_id in (?)", keys).Exec(buffer.String()).Error
}

func handleTagMapByParentId(tag Tag, tagViews []TagView) []TagView {

	for key, tagView := range tagViews {
		if tagView.Id == 0 {
			continue
		}

		if tag.ParentId == tagView.Id {
			tagViews[key].Tags = append(tagViews[key].Tags, TagView{
				Id:     tag.Id,
				NameZh: tag.NameZh,
				NameEn: tag.NameEn,
			})
		} else {
			tagViews[key].Tags = handleTagMapByParentId(tag, tagViews[key].Tags)
		}

	}

	return tagViews
}

func DeleteAppTags(db *gorm.DB, id uint64) error {
	if err := db.Where("application_id = ?", id).Delete(AppTag{}).Error; err != nil {
		return err
	}

	return nil
}

func CreateAppTag(db *gorm.DB, appTag AppTag) error {
	if err := db.Create(&appTag).Error; err != nil {
		return err
	}

	return nil
}

// CheckAppTagsIsValid check multiple app tagid is exist
func CheckAppTagsIsValid(db *gorm.DB, appIds []uint64) (bool, error) {
	var count int
	if err := db.Table("tags").Where("id IN (?)", appIds).Count(&count).Error; err != nil {
		return false, err
	}
	return len(appIds) == count, nil
}

func GetAppCountGroubyTag(db *gorm.DB, appIds []uint64) ([]TagInfo, error) {
	var tagInfo []TagInfo

	if err := db.Raw(`select tt.id as tag_id, tt.name_zh, tt.name_en, tt.tag_level, tt.parent_id, tmp.num from tags tt left join (
		SELECT t.id as tag_id, t.name_zh, t.name_en, t.tag_level, t.parent_id, count(*) as num FROM app_tags at 
		left join tags t on t.id = at.tag_id WHERE (at.application_id in (?)) GROUP BY t.id, t.name_zh, t.name_en, t.tag_level, t.parent_id
	) tmp on tmp.tag_id = tt.id`, appIds).Find(&tagInfo).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return tagInfo, nil
}

func GetPublishedAppByTagId(db *gorm.DB, tagIds ...uint64) ([]uint64, error) {
	var appIds []uint64
	var tempDB *gorm.DB

	if len(tagIds) > 0 {
		tempDB = db.Raw(`select oa.id from oauth_applications oa 
		where status = ? and listing = 1 and private_app = 0 and 
		id in (select distinct(application_id) from app_tags where tag_id in (?))`, "published", tagIds)
	} else {
		tempDB = db.Raw(`select oa.id from oauth_applications oa 
		where status = ? and listing = 1 and private_app = 0 and 
		id in (select distinct(application_id) from app_tags)`, "published")
	}

	if err := tempDB.Order("oa.published_at desc").Pluck("id", &appIds).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return appIds, nil
}

func GerAppIdsOrderbyInstalledCount(db *gorm.DB, appIds []uint64) ([]uint64, error) {
	var ids []uint64

	if err := db.Table("install_tracks it").Select("it.application_id as id").
		Where("it.application_id in (?)", appIds).Group("it.application_id").
		Order("count(*) desc").Pluck("id", &ids).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	idsMap := make(map[uint64]struct{})
	for _, v := range ids {
		idsMap[v] = struct{}{}
	}

	for _, v := range appIds {
		if _, ok := idsMap[v]; !ok {
			ids = append(ids, v)
		}
	}

	return ids, nil
}

type TagAppView struct {
	AppId    uint64 `json:"app_id"`
	TagId    uint64 `json:"tag_id"`
	NameZh   string `json:"name_zh"`
	NameEn   string `json:"name_en"`
	TagLevel string `json:"tag_level"`
	ParentId uint64 `json:"parent_id"`
}

func GetAppsTagsByIds(db *gorm.DB, ids []uint64) ([]TagAppView, error) {
	tagAppViews := []TagAppView{}
	if err := db.Table("app_tags at").
		Select("at.application_id as app_id,at.tag_id, t.name_zh, t.name_en, t.tag_level, t.parent_id").
		Joins("left join tags t on at.tag_id = t.id").
		Where("at.application_id in (?)", ids).Find(&tagAppViews).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return tagAppViews, nil
}

func TagAppViewSliceGroupByAppid(tagAppViews []TagAppView) map[uint64][]TagAppView {
	TagAppMap := make(map[uint64][]TagAppView)
	for _, v := range tagAppViews {
		if _, ok := TagAppMap[v.AppId]; !ok {
			TagAppMap[v.AppId] = make([]TagAppView, 0)
			TagAppMap[v.AppId] = append(TagAppMap[v.AppId], v)
			continue
		}

		TagAppMap[v.AppId] = append(TagAppMap[v.AppId], v)
	}
	return TagAppMap
}

func TrimAppTagViews(tav []TagAppView) []TagAppView {
	newTav := make([]TagAppView, 0)
	for i := 0; i < len(tav); i++ {

		if tav[i].TagLevel != "1" && tav[i].NameZh != "其他" {
			//普通标签保存
			newTav = append(newTav, tav[i])
		} else if tav[i].TagLevel != "1" {
			//其他标签保存父级
			parent := getAppTagViewsParent(tav, tav[i].ParentId)

			if parent != nil {
				newTav = append(newTav, *parent)
				// 保留parent
			}
		}
	}
	return newTav
}
func getAppTagViewsParent(tavs []TagAppView, parentId uint64) *TagAppView {
	for i := 0; i < len(tavs); i++ {
		if tavs[i].TagId == parentId {
			return &tavs[i]
		}
	}
	return nil
}

func GetPublishedAppByTagIdAndAppRegion(db *gorm.DB, region string, tagIds ...uint64) ([]uint64, error) {
	var appIds []uint64
	var tempDB *gorm.DB

	if len(tagIds) > 0 {
		tempDB = db.Raw(`select oa.id from oauth_applications oa 
		where status = ? and listing = 1 and private_app = 0 and 
		id in (select distinct(application_id) from app_tags where tag_id in (?))
		and (regions LIKE ? OR regions LIKE '%All%' ) and category not in (?)`, "published", tagIds, "%"+region+"%", NotNeedListingAPPCategoryPayments)
	} else {
		tempDB = db.Raw(`select oa.id from oauth_applications oa 
		where status = ? and listing = 1 and private_app = 0 and 
		id in (select distinct(application_id) from app_tags)
		and (regions LIKE  ?  OR regions LIKE '%All%' ) and category not in (?)`, "published", "%"+region+"%", NotNeedListingAPPCategoryPayments)
	}

	if err := tempDB.Order("oa.published_at desc").Pluck("id", &appIds).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return appIds, nil
}
