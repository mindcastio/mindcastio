package search

import (
	"strings"
	"gopkg.in/mgo.v2/bson"
	"strconv"

	"github.com/franela/goreq"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/services/datastore"
	"github.com/mindcastio/mindcastio/backend/services/logger"
	"github.com/mindcastio/mindcastio/backend/services/metrics"
	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/util"
)

type (
	PodcastSearchMetadata struct {
		Uid         string `json:"uid"`
		Title       string `json:"title"`
		Subtitle    string `json:"subtitle"`
		Description string `json:"description"`
		Published   int64  `json:"published"`
		Language    string `json:"language"`
		OwnerName   string `json:"owner_name"`
		OwnerEmail  string `json:"owner_email"`
		Tags        string `json:"tags"`
	}

	EpisodeSearchMetadata struct {
		Uid         string `json:"uid"`
		Title       string `json:"title"`
		Link        string `json:"link"`
		Description string `json:"description"`
		Published   int64  `json:"published"`
		Author      string `json:"author"`
		PodcastUid  string `json:"puid"`
	}

)

func SchedulePodcastIndexing() {

	logger.Log("schedule_podcast_indexing")

	// search for podcasts that are candidates for indexing
	notIndexed := podcastSearchNotIndexed(backend.DEFAULT_INDEX_UPDATE_BATCH, backend.SEARCH_REVISION)
	count := len(notIndexed)

	logger.Log("schedule_podcast_indexing.scheduling", strconv.FormatInt((int64)(count), 10))

	if count > 0 {
		ds := datastore.GetDataStore()
		defer ds.Close()

		podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

		for i := 0; i < count; i++ {
			err := podcastAddToSearchIndex(&notIndexed[i])
			if err != nil {
				logger.Error("schedule_podcast_indexing.error.1", err, notIndexed[i].Uid)
				metrics.Error("schedule_podcast_indexing.error.1", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}

			// update the metadata
			notIndexed[i].IndexVersion = backend.SEARCH_REVISION
			notIndexed[i].Updated = util.Timestamp()
			err = podcast_metadata.Update(bson.M{"uid": notIndexed[i].Uid}, &notIndexed[i])
			if err != nil {
				logger.Error("schedule_podcast_indexing.error.2", err, notIndexed[i].Uid)
				metrics.Error("schedule_podcast_indexing.error.2", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}
		}
		metrics.Count("indexer.podcasts.scheduled", count)
	}

	logger.Log("schedule_podcast_indexing.done")
}

func ScheduleEpisodeIndexing() {

	logger.Log("schedule_episode_indexing")

	// search for podcasts that are candidates for indexing
	notIndexed := episodesSearchNotIndexed(backend.DEFAULT_INDEX_UPDATE_BATCH, backend.SEARCH_REVISION)
	count := len(notIndexed)

	logger.Log("schedule_episode_indexing.scheduling", strconv.FormatInt((int64)(count), 10))

	if count > 0 {
		ds := datastore.GetDataStore()
		defer ds.Close()

		episodes_metadata := ds.Collection(datastore.EPISODES_COL)

		for i := 0; i < count; i++ {
			err := episodeAddToSearchIndex(&notIndexed[i])
			if err != nil {
				logger.Error("schedule_episode_indexing.error.1", err, notIndexed[i].Uid)
				metrics.Error("schedule_episode_indexing.error.1", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}

			// update the metadata
			notIndexed[i].IndexVersion = backend.SEARCH_REVISION
			notIndexed[i].Updated = util.Timestamp()
			err = episodes_metadata.Update(bson.M{"uid": notIndexed[i].Uid}, &notIndexed[i])
			if err != nil {
				logger.Error("schedule_episode_indexing.error.2", err, notIndexed[i].Uid)
				metrics.Error("schedule_episode_indexing.error.2", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}
		}
		metrics.Count("indexer.episodes.scheduled", count)
	}

	logger.Log("schedule_episode_indexing.done")
}

func podcastAddToSearchIndex(podcast *backend.PodcastMetadata) error {

	uri := strings.Join([]string{environment.GetEnvironment().SearchServiceUrl(), "/search/podcast/", podcast.Uid}, "")

	payload := PodcastSearchMetadata{
		podcast.Uid,
		podcast.Title,
		podcast.Subtitle,
		podcast.Description,
		podcast.Published,
		podcast.Language,
		podcast.OwnerName,
		podcast.OwnerEmail,
		podcast.Tags,
	}

	// post the payload to elasticsearch
	res, err := goreq.Request{
		Method:      "PUT",
		Uri:         uri,
		ContentType: "application/json",
		Body:        payload,
	}.Do()

	if res != nil {
		res.Body.Close()
	}
	return err
}

func episodeAddToSearchIndex(episode *backend.EpisodeMetadata) error {

	uri := strings.Join([]string{environment.GetEnvironment().SearchServiceUrl(), "/search/episode/", episode.Uid}, "")

	payload := EpisodeSearchMetadata{
		episode.Uid,
		episode.Title,
		episode.Url,
		episode.Description,
		episode.Published,
		episode.Author,
		episode.PodcastUid,
	}

	// post the payload to elasticsearch
	res, err := goreq.Request{
		Method:      "PUT",
		Uri:         uri,
		ContentType: "application/json",
		Body:        payload,
	}.Do()

	if res != nil {
		res.Body.Close()
	}
	return err
}

func podcastSearchNotIndexed(limit int, version int) []backend.PodcastMetadata {

	ds := datastore.GetDataStore()
	defer ds.Close()

	podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

	results := []backend.PodcastMetadata{}
	q := bson.M{"indexversion": bson.M{"$lt": version}}

	if limit <= 0 {
		// return all
		podcast_metadata.Find(q).All(&results)
	} else {
		// with a limit
		podcast_metadata.Find(q).Limit(limit).All(&results)
	}

	return results
}

func episodesSearchNotIndexed(limit int, version int) []backend.EpisodeMetadata {

	ds := datastore.GetDataStore()
	defer ds.Close()

	episodes_metadata := ds.Collection(datastore.EPISODES_COL)

	results := []backend.EpisodeMetadata{}
	q := bson.M{"indexversion": bson.M{"$lt": version}}

	if limit <= 0 {
		// return all
		episodes_metadata.Find(q).All(&results)
	} else {
		// with a limit
		episodes_metadata.Find(q).Limit(limit).All(&results)
	}

	return results
}
