package script

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"gitlab.shoplazza.site/shoplaza/cobra/config"
	"gitlab.shoplazza.site/shoplaza/cobra/models"
	"time"
)

func UpdateTagAndApptag() error {
	//先获取旧版tag数据

	tx := config.DB.Begin()
	oldTags, _ := models.GetTags(tx)
	//执行更新脚本后插入新的数据
	if err := insertTags(tx); err != nil {
		log.Error("Insert New Tags Fail")
		tx.Rollback()
		return err
	}
	newTags, _ := models.GetTags(tx)

	var updateMap = map[uint64]uint64{}
	var exists bool
	var deletedIds []uint64
	for _, oldTag := range oldTags {
		exists = false
		for _, newTag := range newTags {
			if oldTag.NameEn == newTag.NameEn && newTag.NameEn != "Other" {
				updateMap[oldTag.Id] = newTag.Id
				exists = true
				break
			}
		}
		if !exists {
			deletedIds = append(deletedIds, oldTag.Id)
		}
	}
	//删除appTags中的多余数据
	if err := models.DeleteAppTag(tx, deletedIds); err != nil {
		log.Error("Delete Extra AppTags Fail")
		tx.Rollback()
		return err
	}
	//更新并插入一二三级分类对应的tagId
	if err := models.UpdateAppTag(tx, updateMap); err != nil {
		log.Error("Update AppTags Fail")
		tx.Rollback()
		return err
	}
	var newAppTags []*models.AppTag
	infos, err := models.GetAppTagsParent(tx)
	if err != nil {
		log.Error("Get AppTags Fail")
		tx.Rollback()
		return err
	}
	for _, info := range infos {
		if info.TagId != 0 {
			newAppTags = append(newAppTags, &models.AppTag{
				ApplicationId: info.ApplicationId,
				TagId:         info.TagId,
				CreatedAt:     time.Now().UTC(),
				UpdatedAt:     time.Now().UTC(),
			})
		}
		if info.ParentId != 0 {
			newAppTags = append(newAppTags, &models.AppTag{
				ApplicationId: info.ApplicationId,
				TagId:         info.ParentId,
				CreatedAt:     time.Now().UTC(),
				UpdatedAt:     time.Now().UTC(),
			})
		}
	}
	if err := models.CreateAppTags(tx, newAppTags); err != nil {
		log.Error("Insert AppTags Fail")
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func insertTags(tx *gorm.DB) error {

	return tx.Exec("INSERT INTO tags (id, tag_level, parent_id, name_zh, name_en, created_at, updated_at)\nVALUES ('1', '1', '0', '寻找产品', 'Finding products', now(), now()),\n       ('2', '1', '0', '销售产品', 'Selling products', now(), now()),\n       ('3', '1', '0', '订单和运输', 'Orders and shipping', now(), now()),\n       ('4', '1', '0', '商店设计', 'Store design', now(), now()),\n       ('5', '1', '0', '营销和转换', 'Marketing and conversion', now(), now()),\n       ('6', '1', '0', '商店管理', 'Store management', now(), now()),\n       ('7', '2', '1', '产品采购', 'Product sourcing', now(), now()),\n       ('8', '2', '2', '销售方法', 'Selling methods', now(), now()),\n       ('9', '2', '2', '产品展示', 'Product display', now(), now()),\n       ('10', '2', '2', '产品变量', 'Product variants', now(), now()),\n       ('11', '2', '2', '购买方式', 'Purchase options', now(), now()),\n       ('12', '2', '2', '礼物', 'Gifts', now(), now()),\n       ('13', '2', '2', '数字产品', 'Digital products', now(), now()),\n       ('14', '2', '3', '订单管理', 'Managing orders', now(), now()),\n       ('15', '2', '3', '履行订单', 'Fulfilling orders', now(), now()),\n       ('16', '2', '3', '库存管理', 'Managing inventory', now(), now()),\n       ('17', '2', '4', '店铺页面', 'Store pages', now(), now()),\n       ('18', '2', '4', '导航与搜索', 'Navigation and search', now(), now()),\n       ('19', '2', '4', '图片和媒体', 'Images and media', now(), now()),\n       ('20', '2', '4', '通知', 'Notifications', now(), now()),\n       ('21', '2', '4', '店铺提示', 'Store alerts', now(), now()),\n       ('22', '2', '4', '国际化', 'Internationalization', now(), now()),\n       ('23', '2', '4', '社交证明', 'Social proof', now(), now()),\n       ('24', '2', '5', '搜索引擎优化', 'Search engine optimization', now(), now()),\n       ('25', '2', '5', '广告', 'Advertising', now(), now()),\n       ('26', '2', '5', '电子邮件营销', 'Email marketing', now(), now()),\n       ('27', '2', '5', '直接营销', 'Direct marketing', now(), now()),\n       ('28', '2', '5', '促销', 'Promotions', now(), now()),\n       ('29', '2', '5', '向上销售和交叉销售', 'Upselling and cross-selling', now(), now()),\n       ('30', '2', '5', '购物车调整', 'Cart modification', now(), now()),\n       ('31', '2', '5', '顾客账户', 'Customer accounts', now(), now()),\n       ('32', '2', '6', '支持服务', 'Support', now(), now()),\n       ('33', '2', '6', '店铺数据', 'Store data', now(), now()),\n       ('34', '2', '6', '安全与隐私', 'Privacy and security', now(), now()),\n       ('35', '2', '6', '财务管理', 'Finances', now(), now()),\n       ('36', '2', '6', '分析报告', 'Analytics', now(), now()),\n       ('37', '3', '7', '代发货', 'Dropshipping', now(), now()),\n       ('38', '3', '7', '寻找供应商', 'Finding suppliers', now(), now()),\n       ('39', '3', '7', '按需印刷（POD）', 'Print on demand (POD)', now(), now()),\n       ('40', '3', '7', '批发采购', 'Buying wholesale', now(), now()),\n       ('41', '3', '7', '其他', 'Other', now(), now()),\n       ('42', '3', '8', '社交媒体', 'Social media', now(), now()),\n       ('43', '3', '8', '在线购物', 'Live shopping', now(), now()),\n       ('44', '3', '8', '定制店面', 'Custom storefronts', now(), now()),\n       ('45', '3', '8', '其他', 'Other', now(), now()),\n       ('46', '3', '9', '3D/AR/VR', '3D/AR/VR', now(), now()),\n       ('47', '3', '9', '其他', 'Other', now(), now()),\n       ('48', '3', '10', '产品定制选项', 'Product options', now(), now()),\n       ('49', '3', '10', '色卡', 'Color swatches', now(), now()),\n       ('50', '3', '10', '自定义文件上传', 'Custom file upload', now(), now()),\n       ('51', '3', '10', '其他', 'Other', now(), now()),\n       ('52', '3', '11', '订阅', 'Subscriptions', now(), now()),\n       ('53', '3', '11', '预定', 'Pre-orders', now(), now()),\n       ('54', '3', '11', '订单限制', 'Order limits', now(), now()),\n       ('55', '3', '11', '其他', 'Other', now(), now()),\n       ('56', '3', '12', '礼品卡', 'Gift cards', now(), now()),\n       ('57', '3', '12', '其他', 'Other', now(), now()),\n       ('58', '3', '13', '数字产品下载', 'Digital downloads', now(), now()),\n       ('59', '3', '13', 'NFTs', 'NFTs', now(), now()),\n       ('60', '3', '13', '其他', 'Other', now(), now()),\n       ('61', '3', '14', '订单标记器', 'Order tagger', now(), now()),\n       ('62', '3', '14', '订单同步', 'Order sync', now(), now()),\n       ('63', '3', '14', '订单编辑', 'Order editing', now(), now()),\n       ('64', '3', '14', '商家订单通知', 'Merchant order notifications', now(), now()),\n       ('65', '3', '14', '采购订单', 'Purchase orders', now(), now()),\n       ('66', '3', '14', '发票和收据', 'Invoices and receipts', now(), now()),\n       ('67', '3', '14', '订单和运输报告', 'Order and shipping reports', now(), now()),\n       ('68', '3', '14', '其他', 'Other', now(), now()),\n       ('69', '3', '15', '订单扫描器', 'Order scanner', now(), now()),\n       ('70', '3', '15', '仓库管理', 'Warehouse management', now(), now()),\n       ('71', '3', '15', '外包履行', 'Outsourced fulfillment', now(), now()),\n       ('72', '3', '15', 'SKU和条形码', 'SKU and barcodes', now(), now()),\n       ('73', '3', '15', '包装单据', 'Packing slips', now(), now()),\n       ('74', '3', '15', '运输标签', 'Shipping labels', now(), now()),\n       ('75', '3', '15', '运费计算器', 'Shipping rate calculator', now(), now()),\n       ('76', '3', '15', '包装', 'Packaging', now(), now()),\n       ('77', '3', '15', '其他', 'Other', now(), now()),\n       ('78', '3', '16', '库存同步', 'Inventory sync', now(), now()),\n       ('79', '3', '16', '库存提醒', 'Stock alerts', now(), now()),\n       ('80', '3', '16', '库存优化', 'Inventory optimization', now(), now()),\n       ('81', '3', '16', '库存跟踪', 'Inventory tracking', now(), now()),\n       ('82', '3', '16', '其他', 'Other', now(), now()),\n       ('83', '3', '17', '页面制作工具', 'Page builder', now(), now()),\n       ('84', '3', '17', '预先启动', 'Pre-launch', now(), now()),\n       ('85', '3', '17', '其他', 'Other', now(), now()),\n       ('86', '3', '18', '导航和过滤器', 'Navigation and filters', now(), now()),\n       ('87', '3', '18', '搜索', 'Search', now(), now()),\n       ('88', '3', '18', '其他', 'Other', now(), now()),\n       ('89', '3', '19', '幻灯片', 'Image slider', now(), now()),\n       ('90', '3', '19', '图片库', 'Image galleries', now(), now()),\n       ('91', '3', '19', '图片编辑', 'Image editor', now(), now()),\n       ('92', '3', '19', '视频编辑', 'Video editor', now(), now()),\n       ('93', '3', '19', '音频播放器', 'Audio player', now(), now()),\n       ('94', '3', '19', '内容管理器', 'Content manager', now(), now()),\n       ('95', '3', '19', '其他', 'Other', now(), now()),\n       ('96', '3', '20', '产品名片', 'Product badges', now(), now()),\n       ('97', '3', '20', '弹出式窗口', 'Popups', now(), now()),\n       ('98', '3', '20', '横幅广告', 'Banners', now(), now()),\n       ('99', '3', '20', '通知', 'Notifications', now(), now()),\n       ('100', '3', '20', '其他', 'Other', now(), now()),\n       ('101', '3', '21', '倒数计时器', 'Countdown timer', now(), now()),\n       ('102', '3', '21', '库存计数器', 'Stock counter', now(), now()),\n       ('103', '3', '21', '有库存提示', 'Back in stock alert', now(), now()),\n       ('104', '3', '21', '价格变化提示', 'Price change alert', now(), now()),\n       ('105', '3', '21', '其他', 'Other', now(), now()),\n       ('106', '3', '22', '货币', 'Currency', now(), now()),\n       ('107', '3', '22', '语言与翻译', 'Language and translation', now(), now()),\n       ('108', '3', '22', '其他', 'Other', now(), now()),\n       ('109', '3', '23', '社交证明', 'Social proof', now(), now()),\n       ('110', '3', '23', '信任徽章', 'Trust badges', now(), now()),\n       ('111', '3', '23', '商品评价', 'Product reviews', now(), now()),\n       ('112', '3', '23', '其他', 'Other', now(), now()),\n       ('113', '3', '24', '搜索引擎优化', 'SEO', now(), now()),\n       ('114', '3', '24', '其他', 'Other', now(), now()),\n       ('115', '3', '25', '广告', 'Advertising', now(), now()),\n       ('116', '3', '25', '重定向广告', 'Retargeting ads', now(), now()),\n       ('117', '3', '25', '社交媒体广告', 'Social media ads', now(), now()),\n       ('118', '3', '25', '联盟计划', 'Affiliate programs', now(), now()),\n       ('119', '3', '25', '货源', 'Product feeds', now(), now()),\n       ('120', '3', '25', '其他', 'Other', now(), now()),\n       ('121', '3', '26', '电子邮件营销', 'Email marketing', now(), now()),\n       ('122', '3', '26', '营销活动管理', 'Campaign management', now(), now()),\n       ('123', '3', '26', '其他', 'Other', now(), now()),\n       ('124', '3', '27', '推送通知', 'Push notifications', now(), now()),\n       ('125', '3', '27', '短信营销', 'SMS marketing', now(), now()),\n       ('126', '3', '27', '其他', 'Other', now(), now()),\n       ('127', '3', '28', '折扣', 'Discounts', now(), now()),\n       ('128', '3', '28', '赠品', 'Gift with purchase', now(), now()),\n       ('129', '3', '28', '买一送一', 'Buy one, get one (BOGO)', now(), now()),\n       ('130', '3', '28', '其他', 'Other', now(), now()),\n       ('131', '3', '29', '向上销售和交叉销售', 'Upselling and cross-selling', now(), now()),\n       ('132', '3', '29', '产品捆绑销售', 'Product bundles', now(), now()),\n       ('133', '3', '29', '近期浏览记录', 'Recently viewed', now(), now()),\n       ('134', '3', '29', '推荐商品', 'Recommended products', now(), now()),\n       ('135', '3', '29', '其他', 'Other', now(), now()),\n       ('136', '3', '30', '购物车调整', 'Cart modification', now(), now()),\n       ('137', '3', '30', '一键结账', 'One-click checkout', now(), now()),\n       ('138', '3', '30', '加入购物车', 'Add to cart', now(), now()),\n       ('139', '3', '30', '其他', 'Other', now(), now()),\n       ('140', '3', '31', '账户与登录', 'Accounts and login', now(), now()),\n       ('141', '3', '31', '忠诚度与奖励', 'Loyalty and rewards', now(), now()),\n       ('142', '3', '31', '许愿单', 'Wishlists', now(), now()),\n       ('143', '3', '31', '其他', 'Other', now(), now()),\n       ('144', '3', '32', '聊天', 'Chat', now(), now()),\n       ('145', '3', '32', '联系方式表单', 'Contact form', now(), now()),\n       ('146', '3', '32', '常见问题与回答', 'FAQ', now(), now()),\n       ('147', '3', '32', '顾客订单追踪', 'Customer order tracking', now(), now()),\n       ('148', '3', '32', '其他', 'Other', now(), now()),\n       ('149', '3', '33', '店铺数据导入', 'Store data importers', now(), now()),\n       ('150', '3', '33', '备份', 'Backup', now(), now()),\n       ('151', '3', '33', '其他', 'Other', now(), now()),\n       ('152', '3', '34', 'IP屏蔽器', 'IP blocker', now(), now()),\n       ('153', '3', '34', '法律', 'Legal', now(), now()),\n       ('154', '3', '34', '隐私', 'Privacy', now(), now()),\n       ('155', '3', '34', '安全', 'Security', now(), now()),\n       ('156', '3', '34', '其他', 'Other', now(), now()),\n       ('157', '3', '35', '会计', 'Accounting', now(), now()),\n       ('158', '3', '35', '税费', 'Taxes', now(), now()),\n       ('159', '3', '35', '财务报告', 'Financial reports', now(), now()),\n       ('160', '3', '35', '其他', 'Other', now(), now()),\n       ('161', '3', '36', '销售分析', 'Sales analytics', now(), now()),\n       ('162', '3', '36', '顾客分析', 'Customer analytics', now(), now()),\n       ('163', '3', '36', '店铺活动', 'Store activity', now(), now()),\n       ('164', '3', '36', '市场分析', 'Marketing analytics', now(), now()),\n       ('165', '3', '36', '其他', 'Other', now(), now()) ON DUPLICATE KEY\nUPDATE tag_level =\nvalues (tag_level), parent_id =\nvalues (parent_id), name_zh =\nvalues (name_zh), name_en =\nvalues (name_en), created_at =\nvalues (created_at), updated_at =\nvalues (updated_at);\n").Error
}
