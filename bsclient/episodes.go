package bsclient

import (
	"net/url"
	"strconv"
)

type episodeItem struct {
	Episode *Episode      `json:"episode"`
	Errors  []interface{} `json:"errors"`
}

// EpisodesList returns a slice of unseen episodes ordered by shows.
func (bs *BetaSeries) EpisodesList(showID, userID int) ([]Show, error) {
	usedAPI := "/episodes/list"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	if showID >= 0 {
		q.Set("showId", strconv.Itoa(showID))
	}
	if userID >= 0 {
		q.Set("userId", strconv.Itoa(userID))
	}
	u.RawQuery = q.Encode()

	return bs.doGetShows(u, usedAPI)
}

func (bs *BetaSeries) episodeUpdate(method string, id int) (*Episode, error) {
	usedAPI := "/episodes/downloaded"
	u, err := url.Parse(bs.baseURL + usedAPI)
	if err != nil {
		return nil, errURLParsing
	}
	q := u.Query()
	q.Set("id", strconv.Itoa(id))
	u.RawQuery = q.Encode()

	resp, err := bs.do(method, u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	episode := &episodeItem{}
	err = bs.decode(episode, resp, usedAPI, u.RawQuery)
	if err != nil {
		return nil, err
	}

	return episode.Episode, nil
}

// EpisodeDownloaded marks the episode with the given 'id' as downloaded.
func (bs *BetaSeries) EpisodeDownloaded(id int) (*Episode, error) {
	return bs.episodeUpdate("POST", id)
}

// EpisodeNotDownloaded marks the episode with the given 'id' as downloaded.
func (bs *BetaSeries) EpisodeNotDownloaded(id int) (*Episode, error) {
	return bs.episodeUpdate("DELETE", id)
}
