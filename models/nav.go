package models

// ChannelMap maps channel's name to id.
type ChannelMap map[string]int64

func NewChannelMap(data []ChannelSetting) ChannelMap {
	var m = make(ChannelMap)

	for _, v := range data {
		m[v.KeyName] = v.ID
	}

	return m
}

var sections = map[string]ChannelSetting{
	"home": {
		ID:       1,
		ParentID: 0,
		Name:     "首页",
	},
	"home_special": {
		ID:       11,
		ParentID: 1,
		Name:     "特别报导",
	},
	"home_weekly": {
		ID:       12,
		ParentID: 1,
		Name:     "热门文章",
	},
	"home_event": {
		ID:       13,
		ParentID: 1,
		Name:     "会议活动",
	},
	"home_member": {
		ID:       14,
		ParentID: 1,
		Name:     "会员信息中心",
	},
	"home_mba": {
		ID:       15,
		ParentID: 1,
		Name:     "FT商学院",
	},
	"home_ebook": {
		ID:       16,
		ParentID: 1,
		Name:     "FT电子书",
	},
	"home_photo": {
		ID:       17,
		ParentID: 1,
		Name:     "图辑",
	},
	"home_job": {
		ID:       18,
		ParentID: 1,
		Name:     "职业机会",
	},
	"china": {
		ID:       2,
		ParentID: 0,
		Name:     "中国",
	},
	"china_economy": {
		ID:       19,
		ParentID: 2,
		Name:     "中国政经",
	},
	"china_business": {
		ID:       20,
		ParentID: 2,
		Name:     "商业",
	},
	"china_markets": {
		ID:       21,
		ParentID: 2,
		Name:     "金融市场",
	},
	"china_stock": {
		ID:       22,
		ParentID: 2,
		Name:     "股市",
	},
	"china_property": {
		ID:       23,
		ParentID: 2,
		Name:     "房地产",
	},
	"china_culture": {
		ID:       24,
		ParentID: 2,
		Name:     "社会与文化",
	},
	"china_opinion": {
		ID:       25,
		ParentID: 2,
		Name:     "观点",
	},
	"world": {
		ID:       3,
		ParentID: 0,
		Name:     "全球",
	},
	"world_usa": {
		ID:       26,
		ParentID: 3,
		Name:     "美国",
	},
	"world_uk": {
		ID:       27,
		ParentID: 3,
		Name:     "英国",
	},
	"world_asia": {
		ID:       28,
		ParentID: 3,
		Name:     "亚太",
	},
	"world_europe": {
		ID:       29,
		ParentID: 3,
		Name:     "欧洲",
	},
	"world_america": {
		ID:       30,
		ParentID: 3,
		Name:     "美洲",
	},
	"world_africa": {
		ID:       31,
		ParentID: 3,
		Name:     "非洲",
	},
	"economy": {
		ID:       4,
		ParentID: 0,
		Name:     "经济",
	},
	"economy_global": {
		ID:       32,
		ParentID: 4,
		Name:     "全球经济",
	},
	"economy_china": {
		ID:       33,
		ParentID: 4,
		Name:     "中国经济",
	},
	"economy_trade": {
		ID:       34,
		ParentID: 4,
		Name:     "贸易",
	},
	"economy_environment": {
		ID:       35,
		ParentID: 4,
		Name:     "环境",
	},
	"economy_opinions": {
		ID:       36,
		ParentID: 4,
		Name:     "经济评论",
	},
	"markets": {
		ID:       5,
		ParentID: 0,
		Name:     "金融市场",
	},
	"markets_stock": {
		ID:       37,
		ParentID: 5,
		Name:     "股市",
	},
	"markets_forex": {
		ID:       38,
		ParentID: 5,
		Name:     "外汇",
	},
	"markets_bond": {
		ID:       39,
		ParentID: 5,
		Name:     "债市",
	},
	"markets_commodity": {
		ID:       40,
		ParentID: 5,
		Name:     "大宗商品",
	},
	"markets_data": {
		ID:       41,
		ParentID: 5,
		Name:     "金融市场数据",
	},
	"business": {
		ID:       6,
		ParentID: 0,
		Name:     "商业",
	},
	"business_finance": {
		ID:       42,
		ParentID: 6,
		Name:     "金融",
	},
	"business_technology": {
		ID:       43,
		ParentID: 6,
		Name:     "科技",
	},
	"business_auto": {
		ID:       44,
		ParentID: 6,
		Name:     "汽车",
	},
	"business_property": {
		ID:       45,
		ParentID: 6,
		Name:     "房地产",
	},
	"business_agriculture": {
		ID:       46,
		ParentID: 6,
		Name:     "农林",
	},
	"business_energy": {
		ID:       47,
		ParentID: 6,
		Name:     "能源",
	},
	"business_industries": {
		ID:       48,
		ParentID: 6,
		Name:     "工业和采矿",
	},
	"business_airline": {
		ID:       49,
		ParentID: 6,
		Name:     "航空和运输",
	},
	"business_pharma": {
		ID:       50,
		ParentID: 6,
		Name:     "医药",
	},
	"business_entertainment": {
		ID:       51,
		ParentID: 6,
		Name:     "娱乐",
	},
	"business_consumer": {
		ID:       52,
		ParentID: 6,
		Name:     "零售和消费品",
	},
	"business_media": {
		ID:       53,
		ParentID: 6,
		Name:     "传媒和文化",
	},
	"opinion": {
		ID:   7,
		Name: "观点",
	},
	"opinion_lex": {
		ID:       54,
		ParentID: 7,
		Name:     "Lex专栏",
	},
	"opinion_a_list": {
		ID:       55,
		ParentID: 7,
		Name:     "A-List",
	},
	"opinion_columns": {
		ID:       56,
		ParentID: 7,
		Name:     "专栏",
	},
	"opinion_analysis": {
		ID:       57,
		ParentID: 7,
		Name:     "分析",
	},
	"opinion_comment": {
		ID:       58,
		ParentID: 7,
		Name:     "评论",
	},
	"opinion_editorial": {
		ID:       59,
		ParentID: 7,
		Name:     "社评",
	},
	"opinion_book": {
		ID:       60,
		ParentID: 7,
		Name:     "书评",
	},
	"opinion_letter": {
		ID:       61,
		ParentID: 7,
		Name:     "读者有话说",
	},
	"management": {
		ID:   8,
		Name: "管理",
	},
	"management_mba": {
		ID:       62,
		ParentID: 8,
		Name:     "FT商学院",
	},
	"management_career": {
		ID:       63,
		ParentID: 8,
		Name:     "职场",
	},
	"management_leadership": {
		ID:       64,
		ParentID: 8,
		Name:     "领导力",
	},
	"management_wealth": {
		ID:       65,
		ParentID: 8,
		Name:     "财富管理",
	},
	"management_connect": {
		ID:       66,
		ParentID: 8,
		Name:     "商务互联",
	},
	"management_people": {
		ID:       67,
		ParentID: 8,
		Name:     "人物",
	},
	"lifestyle": {
		ID:   9,
		Name: "生活时尚",
	},
	"lifestyle_pursuit": {
		ID:       68,
		ParentID: 9,
		Name:     "乐尚街",
	},
	"lifestyle_taste": {
		ID:       69,
		ParentID: 9,
		Name:     "品味",
	},
	"lifestyle_travel": {
		ID:       70,
		ParentID: 9,
		Name:     "旅行",
	},
	"lifestyle_life": {
		ID:       71,
		ParentID: 9,
		Name:     "生活话题",
	},
	"lifestyle_art": {
		ID:       72,
		ParentID: 9,
		Name:     "艺术与娱乐",
	},
	"lifestyle_spend": {
		ID:       73,
		ParentID: 9,
		Name:     "消费经",
	},
	"stream": {
		ID:   10,
		Name: "视频",
	},
	"stream_politics": {
		ID:       75,
		ParentID: 10,
		Name:     "政经",
	},
	"stream_business": {
		ID:       76,
		ParentID: 10,
		Name:     "产经",
	},
	"stream_finance": {
		ID:       77,
		ParentID: 10,
		Name:     "金融",
	},
	"stream_culture": {
		ID:       78,
		ParentID: 10,
		Name:     "文化",
	},
	"stream_tinted": {
		ID:       79,
		ParentID: 10,
		Name:     "有色眼镜",
	},
}
