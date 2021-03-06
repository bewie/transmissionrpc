package transmissionrpc

/*
	Torrent Accessors
	https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L127
*/

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

var validTorrentFields []string

func init() {
	torrentType := reflect.TypeOf(Torrent{})
	for i := 0; i < torrentType.NumField(); i++ {
		validTorrentFields = append(validTorrentFields, torrentType.Field(i).Tag.Get("json"))
	}
}

// TorrentGetAll returns all the known fields for all the torrents.
func (c *Client) TorrentGetAll() (torrents []*Torrent, err error) {
	// Send already validated fields to the low level fx
	return c.torrentGet(validTorrentFields, nil)
}

// TorrentGetAllFor returns all known fields for the given torrent's ids.
func (c *Client) TorrentGetAllFor(ids []int64) (torrents []*Torrent, err error) {
	return c.torrentGet(validTorrentFields, ids)
}

// TorrentGet returns the given of fields (mandatory) for each ids (optionnal).
func (c *Client) TorrentGet(fields []string, ids []int64) (torrents []*Torrent, err error) {
	// Validate fields
	var fieldInvalid bool
	var knownField string
	for _, inputField := range fields {
		fieldInvalid = true
		for _, knownField = range validTorrentFields {
			if inputField == knownField {
				fieldInvalid = false
				break
			}
		}
		if fieldInvalid {
			err = fmt.Errorf("field '%s' is invalid", inputField)
			return
		}
	}
	// Forward to real method
	return c.torrentGet(fields, ids)
}

func (c *Client) torrentGet(fields []string, ids []int64) (torrents []*Torrent, err error) {
	arguments := torrentGetParams{
		Fields: fields,
		IDs:    ids,
	}
	var result torrentGetResults
	if err = c.rpcCall("torrent-get", &arguments, &result); err != nil {
		err = fmt.Errorf("'torrent-get' rpc method failed: %v", err)
		return
	}
	torrents = result.Torrents
	return
}

type torrentGetParams struct {
	Fields []string `json:"fields"`
	IDs    []int64  `json:"ids,omitempty"`
}

type torrentGetResults struct {
	Torrents []*Torrent `json:"torrents"`
}

// Torrent represents all the possible fields of data for a torrent.
// All fields are pointers to detect if the value is nil (field not requested) or default real default value.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L148
type Torrent struct {
	ActivityDate            *time.Time         `json:"activityDate"`
	AddedDate               *time.Time         `json:"addedDate"`
	BandwidthPriority       *int64             `json:"bandwidthPriority"`
	Comment                 *string            `json:"comment"`
	CorruptEver             *int64             `json:"corruptEver"`
	Creator                 *string            `json:"creator"`
	DateCreated             *time.Time         `json:"dateCreated"`
	DesiredAvailable        *int64             `json:"desiredAvailable"`
	DoneDate                *time.Time         `json:"doneDate"`
	DownloadDir             *string            `json:"downloadDir"`
	DownloadedEver          *int64             `json:"downloadedEver"`
	DownloadLimit           *int64             `json:"downloadLimit"`
	DownloadLimited         *bool              `json:"downloadLimited"`
	Error                   *int64             `json:"error"`
	ErrorString             *string            `json:"errorString"`
	Eta                     *int64             `json:"eta"`
	EtaIdle                 *int64             `json:"etaIdle"`
	Files                   []*TorrentFile     `json:"files"`
	FileStats               []*TorrentFileStat `json:"fileStats"`
	HashString              *string            `json:"hashString"`
	HaveUnchecked           *int64             `json:"haveUnchecked"`
	HaveValid               *int64             `json:"haveValid"`
	HonorsSessionLimits     *bool              `json:"honorsSessionLimits"`
	ID                      *int64             `json:"id"`
	IsFinished              *bool              `json:"isFinished"`
	IsPrivate               *bool              `json:"isPrivate"`
	IsStalled               *bool              `json:"isStalled"`
	LeftUntilDone           *int64             `json:"leftUntilDone"`
	MagnetLink              *string            `json:"magnetLink"`
	ManualAnnounceTime      *int64             `json:"manualAnnounceTime"`
	MaxConnectedPeers       *int64             `json:"maxConnectedPeers"`
	MetadataPercentComplete *float64           `json:"metadataPercentComplete"`
	Name                    *string            `json:"name"`
	PeerLimit               *int64             `json:"peer-limit"`
	Peers                   []*Peer            `json:"peers"`
	PeersConnected          *int64             `json:"peersConnected"`
	PeersFrom               *TorrentPeersFrom  `json:"peersFrom"`
	PeersGettingFromUs      *int64             `json:"peersGettingFromUs"`
	PeersSendingToUs        *int64             `json:"peersSendingToUs"`
	PercentDone             *float64           `json:"percentDone"`
	Pieces                  *string            `json:"pieces"` // https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L264
	PieceCount              *int64             `json:"pieceCount"`
	PieceSize               *int64             `json:"pieceSize"`
	Priorities              []int64            `json:"priorities"` // https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L270
	QueuePosition           *int64             `json:"queuePosition"`
	RateDownload            *int64             `json:"rateDownload"` // B/s
	RateUpload              *int64             `json:"rateUpload"`   // B/s
	RecheckProgress         *float64           `json:"recheckProgress"`
	SecondsDownloading      *int64             `json:"secondsDownloading"`
	SecondsSeeding          *time.Duration     `json:"secondsSeeding"`
	SeedIdleLimit           *int64             `json:"seedIdleLimit"`
	SeedIdleMode            *int64             `json:"seedIdleMode"`
	SeedRatioLimit          *float64           `json:"seedRatioLimit"`
	SeedRatioMode           *int64             `json:"seedRatioMode"`
	SizeWhenDone            *int64             `json:"sizeWhenDone"`
	StartDate               *time.Time         `json:"startDate"`
	Status                  *int64             `json:"status"`
	Trackers                []*Tracker         `json:"trackers"`
	TrackerStats            []*TrackerStats    `json:"trackerStats"`
	TotalSize               *int64             `json:"totalSize"`
	TorrentFile             *string            `json:"torrentFile"`
	UploadedEver            *int64             `json:"uploadedEver"`
	UploadLimit             *int64             `json:"uploadLimit"`
	UploadLimited           *bool              `json:"uploadLimited"`
	UploadRatio             *float64           `json:"uploadRatio"`
	Wanted                  []bool             `json:"wanted"`   //https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L310
	WebSeeds                []string           `json:"webseeds"` // https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L314
	WebSeedsSendingToUs     *int64             `json:"webseedsSendingToUs"`
}

// UnmarshalJSON allows to convert timestamps to golang time.Time values.
func (t *Torrent) UnmarshalJSON(data []byte) (err error) {
	// Shadow real type for regular unmarshalling
	type RawTorrent Torrent
	tmp := &struct {
		ActivityDate   *int64  `json:"activityDate"`
		AddedDate      *int64  `json:"addedDate"`
		DateCreated    *int64  `json:"dateCreated"`
		DoneDate       *int64  `json:"doneDate"`
		SecondsSeeding *int64  `json:"secondsSeeding"`
		StartDate      *int64  `json:"startDate"`
		Wanted         []int64 `json:"wanted"` // boolean in number form
		*RawTorrent
	}{
		RawTorrent: (*RawTorrent)(t),
	}
	// Unmarshall (with timestamps as number)
	if err = json.Unmarshal(data, &tmp); err != nil {
		return
	}
	// Create the real time & duration from timsteamps and seconds
	if tmp.ActivityDate != nil {
		ad := time.Unix(*tmp.ActivityDate, 0)
		t.ActivityDate = &ad
	}
	if tmp.AddedDate != nil {
		ad := time.Unix(*tmp.AddedDate, 0)
		t.AddedDate = &ad
	}
	if tmp.DateCreated != nil {
		dc := time.Unix(*tmp.DateCreated, 0)
		t.DateCreated = &dc
	}
	if tmp.DoneDate != nil {
		dd := time.Unix(*tmp.DoneDate, 0)
		t.DoneDate = &dd
	}
	if tmp.SecondsSeeding != nil {
		dur := time.Duration(*tmp.SecondsSeeding) * time.Second
		t.SecondsSeeding = &dur
	}
	if tmp.StartDate != nil {
		st := time.Unix(*tmp.StartDate, 0)
		t.StartDate = &st
	}
	// Boolean slice in decimal form
	if tmp.Wanted != nil {
		t.Wanted = make([]bool, len(tmp.Wanted))
		for index, value := range tmp.Wanted {
			if value == 1 {
				t.Wanted[index] = true
			} else if value != 0 {
				return fmt.Errorf("Can't convert Wanted index %d value '%d' as boolean", index, value)
			}
		}
	}
	return
}

// MarshalJSON allows to convert back golang values to original payload values.
func (t *Torrent) MarshalJSON() (data []byte, err error) {
	// Shadow real type for regular unmarshalling
	type RawTorrent Torrent
	tmp := &struct {
		ActivityDate   *int64  `json:"activityDate"`
		AddedDate      *int64  `json:"addedDate"`
		DateCreated    *int64  `json:"dateCreated"`
		DoneDate       *int64  `json:"doneDate"`
		SecondsSeeding *int64  `json:"secondsSeeding"`
		StartDate      *int64  `json:"startDate"`
		Wanted         []int64 `json:"wanted"` // boolean in number form
		*RawTorrent
	}{
		RawTorrent: (*RawTorrent)(t),
	}
	// Timestamps & Duration
	if t.ActivityDate != nil {
		ad := t.ActivityDate.Unix()
		tmp.ActivityDate = &ad
	}
	if t.AddedDate != nil {
		ad := t.AddedDate.Unix()
		tmp.AddedDate = &ad
	}
	if t.DateCreated != nil {
		dc := t.DateCreated.Unix()
		tmp.DateCreated = &dc
	}
	if t.DoneDate != nil {
		dd := t.DoneDate.Unix()
		tmp.DoneDate = &dd
	}
	if t.SecondsSeeding != nil {
		ss := int64(*t.SecondsSeeding / time.Second)
		tmp.SecondsSeeding = &ss
	}
	if t.StartDate != nil {
		st := t.StartDate.Unix()
		tmp.StartDate = &st
	}
	// Boolean as number
	if t.Wanted != nil {
		tmp.Wanted = make([]int64, len(t.Wanted))
		for index, value := range t.Wanted {
			if value {
				tmp.Wanted[index] = 1
			}
		}
	}
	// Marshall original values within the tmp payload
	return json.Marshal(&tmp)
}

// TorrentFile represent one file from a Torrent.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L221
type TorrentFile struct {
	BytesCompleted int64  `json:"bytesCompleted"`
	Length         int64  `json:"length"`
	Name           string `json:"name"`
}

// TorrentFileStat represents the metadata of a torrent's file.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L227
type TorrentFileStat struct {
	BytesCompleted int64 `json:"bytesCompleted"`
	Wanted         bool  `json:"wanted"`
	Priority       int64 `json:"priority"`
}

// Peer represent a peer metadata of a torrent's peer list.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L235
type Peer struct {
	Address              string  `json:"address"`
	ClientName           string  `json:"clientName"`
	ClientIsChoked       bool    `json:"clientIsChoked"`
	ClientIsint64erested bool    `json:"clientIsint64erested"`
	FlagStr              string  `json:"flagStr"`
	IsDownloadingFrom    bool    `json:"isDownloadingFrom"`
	IsEncrypted          bool    `json:"isEncrypted"`
	IsIncoming           bool    `json:"isIncoming"`
	IsUploadingTo        bool    `json:"isUploadingTo"`
	IsUTP                bool    `json:"isUTP"`
	PeerIsChoked         bool    `json:"peerIsChoked"`
	PeerIsint64erested   bool    `json:"peerIsint64erested"`
	Port                 int64   `json:"port"`
	Progress             float64 `json:"progress"`
	RateToClient         int64   `json:"rateToClient"` // B/s
	RateToPeer           int64   `json:"rateToPeer"`   // B/s
}

// TorrentPeersFrom represents the peers statistics of a torrent.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L254
type TorrentPeersFrom struct {
	FromCache    int64 `json:"fromCache"`
	FromDHT      int64 `json:"fromDht"`
	FromIncoming int64 `json:"fromIncoming"`
	FromLPD      int64 `json:"fromLpd"`
	FromLTEP     int64 `json:"fromLtep"`
	FromPEX      int64 `json:"fromPex"`
	FromTracker  int64 `json:"fromTracker"`
}

// Tracker represent the base data of a torrent's tracker.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L274
type Tracker struct {
	Announce string `json:"announce"`
	ID       int64  `json:"id"`
	Scrape   string `json:"scrape"`
	Tier     int64  `json:"tier"`
}

// TrackerStats represent the extended data of a torrent's tracker.
// https://trac.transmissionbt.com/browser/tags/2.92/extras/rpc-spec.txt?rev=14714#L281
type TrackerStats struct {
	Announce              string    `json:"announce"`
	AnnounceState         int64     `json:"announceState"`
	DownloadCount         int64     `json:"downloadCount"`
	HasAnnounced          bool      `json:"hasAnnounced"`
	HasScraped            bool      `json:"hasScraped"`
	Host                  string    `json:"host"`
	ID                    int64     `json:"id"`
	IsBackup              bool      `json:"isBackup"`
	LastAnnouncePeerCount int64     `json:"lastAnnouncePeerCount"`
	LastAnnounceResult    string    `json:"lastAnnounceResult"`
	LastAnnounceStartTime time.Time `json:"lastAnnounceStartTime"`
	LastAnnounceSucceeded bool      `json:"lastAnnounceSucceeded"`
	LastAnnounceTime      time.Time `json:"lastAnnounceTime"`
	LastAnnounceTimedOut  bool      `json:"lastAnnounceTimedOut"`
	LastScrapeResult      string    `json:"lastScrapeResult"`
	LastScrapeStartTime   time.Time `json:"lastScrapeStartTime"`
	LastScrapeSucceeded   bool      `json:"lastScrapeSucceeded"`
	LastScrapeTime        time.Time `json:"lastScrapeTime"`
	LastScrapeTimedOut    bool      `json:"lastScrapeTimedOut"` // should be boolean but number. Will be converter in UnmarshalJSON
	LeecherCount          int64     `json:"leecherCount"`
	NextAnnounceTime      time.Time `json:"nextAnnounceTime"`
	NextScrapeTime        time.Time `json:"nextScrapeTime"`
	Scrape                string    `json:"scrape"`
	ScrapeState           int64     `json:"scrapeState"`
	SeederCount           int64     `json:"seederCount"`
	Tier                  int64     `json:"tier"`
}

// UnmarshalJSON allows to convert timestamps to golang time.Time values.
func (ts *TrackerStats) UnmarshalJSON(data []byte) (err error) {
	// Shadow real type for regular unmarshalling
	type RawTrackerStats TrackerStats
	tmp := struct {
		LastAnnounceStartTime int64 `json:"lastAnnounceStartTime"`
		LastAnnounceTime      int64 `json:"lastAnnounceTime"`
		LastScrapeStartTime   int64 `json:"lastScrapeStartTime"`
		LastScrapeTime        int64 `json:"lastScrapeTime"`
		LastScrapeTimedOut    int64 `json:"lastScrapeTimedOut"`
		NextAnnounceTime      int64 `json:"nextAnnounceTime"`
		NextScrapeTime        int64 `json:"nextScrapeTime"`
		*RawTrackerStats
	}{
		RawTrackerStats: (*RawTrackerStats)(ts),
	}
	// Unmarshall (with timestamps as number)
	if err = json.Unmarshal(data, &tmp); err != nil {
		return
	}
	// Convert to real boolean
	if tmp.LastScrapeTimedOut == 1 {
		ts.LastScrapeTimedOut = true
	} else if tmp.LastScrapeTimedOut != 0 {
		return fmt.Errorf("can't convert 'lastScrapeTimedOut' value '%v' into boolean", tmp.LastScrapeTimedOut)
	}
	// Create the real time value from the timestamps
	ts.LastAnnounceStartTime = time.Unix(tmp.LastAnnounceStartTime, 0)
	ts.LastAnnounceTime = time.Unix(tmp.LastAnnounceTime, 0)
	ts.LastScrapeStartTime = time.Unix(tmp.LastScrapeStartTime, 0)
	ts.LastScrapeTime = time.Unix(tmp.LastScrapeTime, 0)
	ts.NextAnnounceTime = time.Unix(tmp.NextAnnounceTime, 0)
	ts.NextScrapeTime = time.Unix(tmp.NextScrapeTime, 0)
	return
}

// MarshalJSON allows to convert back golang values to original payload values.
func (ts *TrackerStats) MarshalJSON() (data []byte, err error) {
	// Shadow real type for regular unmarshalling
	type RawTrackerStats TrackerStats
	tmp := struct {
		LastAnnounceStartTime int64 `json:"lastAnnounceStartTime"`
		LastAnnounceTime      int64 `json:"lastAnnounceTime"`
		LastScrapeStartTime   int64 `json:"lastScrapeStartTime"`
		LastScrapeTime        int64 `json:"lastScrapeTime"`
		LastScrapeTimedOut    int64 `json:"lastScrapeTimedOut"`
		NextAnnounceTime      int64 `json:"nextAnnounceTime"`
		NextScrapeTime        int64 `json:"nextScrapeTime"`
		*RawTrackerStats
	}{
		LastAnnounceStartTime: ts.LastAnnounceStartTime.Unix(),
		LastAnnounceTime:      ts.LastAnnounceTime.Unix(),
		LastScrapeStartTime:   ts.LastScrapeStartTime.Unix(),
		LastScrapeTime:        ts.LastScrapeTime.Unix(),
		NextAnnounceTime:      ts.NextAnnounceTime.Unix(),
		NextScrapeTime:        ts.NextScrapeTime.Unix(),
		RawTrackerStats:       (*RawTrackerStats)(ts),
	}
	// Convert real bool to its number form
	if ts.LastScrapeTimedOut {
		tmp.LastScrapeTimedOut = 1
	}
	// MarshalJSON allows to convert back golang values to original payload values
	return json.Marshal(&tmp)
}
