package pkgs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var api_key string = os.Getenv("API_KEY")
var channel_id string = os.Getenv("CHANNEL_ID")
var username string = os.Getenv("USERNAME")
var password string = os.Getenv("PASSWORD")

var part string = "contentDetails,snippet"
var listOfVideos []string = make([]string, 0, 2000)

type VideoItems struct {
	VideoId          string      `json:"videoId"`
	VideoPublishedAt string      `json:"videoPublishedAt"`
	PlaylistId       string      `json:"playlistId"`
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	Thumbnails       interface{} `json:"thumbnails"`
}
type logs struct {
	Timestamp int64
}
type PlaylistItems struct {
	PlaylistId          string      `json:"playlistId"`
	PlaylistPublishedAt string      `json:"playlistPubishedAt"`
	Title               string      `json:"title"`
	Description         string      `json:"description"`
	Thumbnails          interface{} `json:"thumbnails"`
}

var PlaylistItemsList []PlaylistItems
var allVideos []VideoItems = make([]VideoItems, 0, 2000)

func validateVideo(id string) bool {
	for _, val := range listOfVideos {
		if id == val {
			return true
		}
	}
	return false
}
func getPlaylistItems(playlist_id string) {
	var url string = fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=%s&key=%s&playlistId=%s&maxResults=50", part, api_key, playlist_id)
	log.Printf("Visting : %s\n", url)
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	var body_map map[string]interface{}
	var nextPageToken string
	err = json.Unmarshal(body, &body_map)
	if body_map["nextPageToken"] != nil {
		nextPageToken = body_map["nextPageToken"].(string)
	}
	items := body_map["items"].([]interface{})
	for _, value := range items {
		mainBody := value.(map[string]interface{})
		contentDetails := mainBody["contentDetails"].(map[string]interface{})
		videoId := contentDetails["videoId"].(string)
		videoPubDateTime := contentDetails["videoPublishedAt"].(string)
		snippet := mainBody["snippet"].(map[string]interface{})
		thumbnails := snippet["thumbnails"]
		title := snippet["title"].(string)
		description := snippet["description"].(string)
		listOfVideos = append(listOfVideos, videoId)
		allVideos = append(allVideos, VideoItems{VideoId: videoId, VideoPublishedAt: videoPubDateTime, Title: title, Description: description, Thumbnails: thumbnails, PlaylistId: playlist_id})
	}
	for len(nextPageToken) > 0 {
		var body_map map[string]interface{}
		nextPageURL := url + "&pageToken=" + nextPageToken
		log.Printf("Visting : %s\n", nextPageURL)
		response, _ = http.Get(nextPageURL)
		body, _ = ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, &body_map)
		nextPageToken = ""
		if body_map["nextPageToken"] != nil {
			nextPageToken = body_map["nextPageToken"].(string)
		}
		items = body_map["items"].([]interface{})
		for _, value := range items {
			mainBody := value.(map[string]interface{})
			contentDetails := mainBody["contentDetails"].(map[string]interface{})
			videoId := contentDetails["videoId"].(string)
			videoPubDateTime := contentDetails["videoPublishedAt"].(string)
			snippet := mainBody["snippet"].(map[string]interface{})
			thumbnails := snippet["thumbnails"]
			title := snippet["title"].(string)
			description := snippet["description"].(string)
			listOfVideos = append(listOfVideos, videoId)
			allVideos = append(allVideos, VideoItems{VideoId: videoId, VideoPublishedAt: videoPubDateTime, Title: title, Description: description, Thumbnails: thumbnails, PlaylistId: playlist_id})
		}
	}
}
func getAllVideos() {
	var url string = fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?channelId=%s&key=%s&part=snippet&maxResults=50&type=video", channel_id, api_key)
	log.Printf("Visting : %s\n", url)
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	var body_map map[string]interface{}
	var nextPageToken string
	err = json.Unmarshal(body, &body_map)
	if body_map["nextPageToken"] != nil {
		nextPageToken = body_map["nextPageToken"].(string)
	}
	items := body_map["items"].([]interface{})
	for _, value := range items {
		mainBody := value.(map[string]interface{})
		id := mainBody["id"].(map[string]interface{})
		videoId := id["videoId"].(string)
		if validateVideo(videoId) {
			continue
		}
		snippet := mainBody["snippet"].(map[string]interface{})
		thumbnails := snippet["thumbnails"]
		videoPubDateTime := snippet["publishedAt"].(string)
		title := snippet["title"].(string)
		description := snippet["description"].(string)
		allVideos = append(allVideos, VideoItems{VideoId: videoId, VideoPublishedAt: videoPubDateTime, Title: title, Description: description, Thumbnails: thumbnails})
	}
	for len(nextPageToken) > 0 {
		var body_map map[string]interface{}
		nextPageURL := url + "&pageToken=" + nextPageToken
		response, _ = http.Get(nextPageURL)
		log.Printf("Visting : %s\n", nextPageURL)
		//fmt.Println(nextPageURL)
		body, _ = ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, &body_map)
		nextPageToken = ""
		if body_map["nextPageToken"] != nil {
			nextPageToken = body_map["nextPageToken"].(string)
		}
		items = body_map["items"].([]interface{})
		for _, value := range items {
			mainBody := value.(map[string]interface{})
			id := mainBody["id"].(map[string]interface{})
			videoId := id["videoId"].(string)
			if validateVideo(videoId) {
				continue
			}
			snippet := mainBody["snippet"].(map[string]interface{})
			thumbnails := snippet["thumbnails"]
			videoPubDateTime := snippet["publishedAt"].(string)
			title := snippet["title"].(string)
			description := snippet["description"].(string)
			allVideos = append(allVideos, VideoItems{VideoId: videoId, VideoPublishedAt: videoPubDateTime, Title: title, Description: description, Thumbnails: thumbnails})
		}
	}
}
func getAllPlaylists() []PlaylistItems {
	var url string = fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?channelId=%s&key=%s&part=snippet&maxResults=50&type=playlist", channel_id, api_key)
	log.Printf("Visting : %s\n", url)
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	var body_map map[string]interface{}
	var nextPageToken string
	err = json.Unmarshal(body, &body_map)
	pageInfo := body_map["pageInfo"].(map[string]interface{})
	totalResults := int64(pageInfo["totalResults"].(float64))
	_items := make([]PlaylistItems, 0, totalResults)
	if body_map["nextPageToken"] != nil {
		nextPageToken = body_map["nextPageToken"].(string)
	}
	items := body_map["items"].([]interface{})
	for _, value := range items {
		mainBody := value.(map[string]interface{})
		id := mainBody["id"].(map[string]interface{})
		playlistId := id["playlistId"].(string)
		snippet := mainBody["snippet"].(map[string]interface{})
		thumbnails := snippet["thumbnails"]
		playlistPubDateTime := snippet["publishedAt"].(string)
		title := snippet["title"].(string)
		description := snippet["description"].(string)
		_items = append(_items, PlaylistItems{PlaylistId: playlistId, PlaylistPublishedAt: playlistPubDateTime, Title: title, Description: description, Thumbnails: thumbnails})
	}
	for len(nextPageToken) > 0 {
		var body_map map[string]interface{}
		nextPageURL := url + "&pageToken=" + nextPageToken
		response, _ = http.Get(nextPageURL)
		log.Printf("Visting : %s\n", nextPageURL)
		body, _ = ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, &body_map)
		nextPageToken = ""
		if body_map["nextPageToken"] != nil {
			nextPageToken = body_map["nextPageToken"].(string)
		}
		items = body_map["items"].([]interface{})
		for _, value := range items {
			mainBody := value.(map[string]interface{})
			id := mainBody["id"].(map[string]interface{})
			playlistId := id["playlistId"].(string)
			snippet := mainBody["snippet"].(map[string]interface{})
			thumbnails := snippet["thumbnails"]
			playlistPubDateTime := snippet["publishedAt"].(string)
			title := snippet["title"].(string)
			description := snippet["description"].(string)
			_items = append(_items, PlaylistItems{PlaylistId: playlistId, PlaylistPublishedAt: playlistPubDateTime, Title: title, Description: description, Thumbnails: thumbnails})
		}
	}
	return _items
}
func makeMap() {
	for _, value := range PlaylistItemsList {
		getPlaylistItems(value.PlaylistId)
	}
}
func Update() {
	PlaylistItemsList = getAllPlaylists()
	makeMap()
	getAllVideos()
	ctx := context.TODO()
	opt := options.Client().ApplyURI(
		"mongodb+srv://" + username + ":" + password + "@<Reacted>",
	)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatal(err)
	}
	playlistsCollection := client.Database("Main").Collection("Playists")
	videosCollection := client.Database("Main").Collection("Videos")
	updateCollection := client.Database("Logs").Collection("UPDATED")

	updateCollection.InsertOne(ctx, logs{
		Timestamp: time.Now().Unix(),
	})

	deleteResult, err := playlistsCollection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleted %v documents in the Playlists collection\n", deleteResult.DeletedCount)

	deleteResult, err = videosCollection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleted %v documents in the Videos collection\n", deleteResult.DeletedCount)
	valuesPlaylist := make([]interface{}, len(PlaylistItemsList))
	for i, value := range PlaylistItemsList {
		valuesPlaylist[i] = value
	}
	playlistsCollection.InsertMany(ctx, valuesPlaylist)
	log.Printf("Videos Update...")
	valuesVideos := make([]interface{}, len(allVideos))
	for i, value := range allVideos {
		valuesVideos[i] = value
	}
	videosCollection.InsertMany(ctx, valuesVideos)
	client.Disconnect(ctx)
}
