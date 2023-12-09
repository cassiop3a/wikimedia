package events

// RevisionCreateEvent corresponds to the JSON event objects returned by the recent changes stream (see:
// https://schema.wikimedia.org/repositories/primary/jsonschema/mediawiki/revision/create/2.0.0.yaml).
type RevisionCreateEvent struct {
	Meta              Meta      `json:"meta"`
	Database          string    `json:"database"`
	PageID            int       `json:"page_id"`
	PageTitle         string    `json:"page_title"`
	PageNamespace     int       `json:"page_namespace"`
	RevID             int       `json:"rev_id"`
	RevTimestamp      string    `json:"rev_timestamp"`
	RevSHA1           string    `json:"rev_sha1"`
	RevMinorEdit      bool      `json:"rev_minor_edit"`
	RevLen            int       `json:"rev_len"`
	RevContentModel   string    `json:"rev_content_model"`
	RevContentFormat  string    `json:"rev_content_format"`
	Performer         Performer `json:"performer"`
	PageIsRedirect    bool      `json:"page_is_redirect"`
	Comment           string    `json:"comment"`
	ParsedComment     string    `json:"parsedcomment"`
	RevParentID       int       `json:"rev_parent_id"`
	Dt                string    `json:"dt"`
	RevSlots          RevSlots  `json:"rev_slots"`
	RevContentChanged bool      `json:"rev_content_changed"`
	EventData
}

type Performer struct {
	UserText           string   `json:"user_text"`
	UserGroups         []string `json:"user_groups"`
	UserIsBot          bool     `json:"user_is_bot"`
	UserID             int      `json:"user_id"`
	UserRegistrationDt string   `json:"user_registration_dt"`
	UserEditCount      int      `json:"user_edit_count"`
}

type RevSlots struct {
	Main Slot `json:"main"`
}

type Slot struct {
	RevSlotContentModel string `json:"rev_slot_content_model"`
	RevSlotSha1         string `json:"rev_slot_sha1"`
	RevSlotSize         int    `json:"rev_slot_size"`
	RevSlotOriginRevID  int    `json:"rev_slot_origin_rev_id"`
}
