package services

type ReqTriggerDiscord struct {
	Type          int64     `json:"type"`
	GuildID       string    `json:"guild_id"`
	ChannelID     string    `json:"channel_id"`
	ApplicationId string    `json:"application_id"`
	SessionId     string    `json:"session_id"`
	Data          DSCommand `json:"data"`
}

type DSCommand struct {
	Version            string               `json:"version"`
	Id                 string               `json:"id"`
	Name               string               `json:"name"`
	Type               int64                `json:"type"`
	Options            []DSOption           `json:"options"`
	ApplicationCommand DSApplicationCommand `json:"application_command"`
	Attachments        []interface{}        `json:"attachments"`
}

type DSOption struct {
	Type  int64  `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DSApplicationCommand struct {
	Id                       string            `json:"id"`
	ApplicationId            string            `json:"application_id"`
	Version                  string            `json:"version"`
	DefaultPermission        bool              `json:"default_permission"`
	DefaultMemberPermissions map[string]int    `json:"default_member_permissions"`
	Type                     int64             `json:"type"`
	Nsfw                     bool              `json:"nsfw"`
	Name                     string            `json:"name"`
	Description              string            `json:"description"`
	DmPermission             bool              `json:"dm_permission"`
	Options                  []DSCommandOption `json:"options"`
}

type DSCommandOption struct {
	Type        int64  `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type ReqUpscaleDiscord struct {
	Type          int64       `json:"type"`
	GuildId       string      `json:"guild_id"`
	ChannelId     string      `json:"channel_id"`
	MessageFlags  int64       `json:"message_flags"`
	MessageId     string      `json:"message_id"`
	ApplicationId string      `json:"application_id"`
	SessionId     string      `json:"session_id"`
	Data          UpscaleData `json:"data"`
}

type UpscaleData struct {
	ComponentType int64  `json:"component_type"`
	CustomId      string `json:"custom_id"`
}

type ReqVariationDiscord = ReqUpscaleDiscord

type ReqResetDiscord = ReqUpscaleDiscord
