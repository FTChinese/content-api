package pkg

// ChannelMap maps channel's name to id.
type ChannelMap map[string]int64

func NewChannelMap(data []ChannelSetting) ChannelMap {
	var m = make(ChannelMap)

	for _, v := range data {
		m[v.KeyName] = v.ID
	}

	return m
}
