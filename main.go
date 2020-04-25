package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
	"youtube-integrations/pkgs"
)
//var api_key string = os.Getenv("API_KEY")
//var channel_id string = os.Getenv("CHANNEL_ID")
//var part string = "contentDetails,snippet"
//var listOfVideos []string = make([]string,0,2000)
type logs struct {
	Timestamp int64
}

var AllPlaylist []pkgs.PlaylistItems = make([]pkgs.PlaylistItems, 0, 100)
var AllVideo []pkgs.VideoItems = make([]pkgs.VideoItems, 0, 2000)
var VideoWithNoPlaylist []pkgs.VideoItems = make([]pkgs.VideoItems,0,2000)
var username string = os.Getenv("USERNAME")
var password string = os.Getenv("PASSWORD")
var PlaylistVideoMap map[string] []pkgs.VideoItems
func _handlePlaylistVideo(r *gin.Context){
	r.JSON(200, map[string] interface{} {
		"playlists": AllPlaylist,
		"videos": PlaylistVideoMap[""],
		"status": "success",
	})
}
func _handleParticularPlaylist(r *gin.Context){
	pID := r.Query("playlistId")
	if len(pID) == 0 {
		r.JSON (400, map[string] string {
			"error" : "playlistId missing.",
		})
		return
	}
	resp := PlaylistVideoMap[pID]
	if resp == nil {
		r.JSON(400, map[string] string {
			"error": "PlaylistId not in the database.",
		})
		return
	}
	r.JSON(200, map[string] interface{} {
		"playlistId": pID,
		"videos" : resp,
	})
}

func getData() {
	AllPlaylist = make([]pkgs.PlaylistItems, 0, 100)
	AllVideo = make([]pkgs.VideoItems, 0, 2000)
	VideoWithNoPlaylist = make([]pkgs.VideoItems,0,2000)

	ctx := context.TODO()
	opt := options.Client().ApplyURI(
		"mongodb+srv://admin:Q4QSTxbfKlgHWrC4@main-2t9o8.mongodb.net/test?retryWrites=true&w=majority",
	)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		log.Fatal(err)
	}
	playlistsCollection := client.Database("Main").Collection("Playists")
	videosCollection := client.Database("Main").Collection("Videos")

	findOptions := options.Find()
	cur, err := playlistsCollection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(ctx) {
		var elem pkgs.PlaylistItems
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		AllPlaylist = append(AllPlaylist, elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(ctx)

	cur, err = videosCollection.Find(ctx, bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(ctx) {
		var elem pkgs.VideoItems
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		AllVideo = append(AllVideo, elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(ctx)
	PlaylistVideoMap = make(map[string] []pkgs.VideoItems)
	for i, _ := range AllVideo {
		if PlaylistVideoMap[AllVideo[i].PlaylistId] == nil {
			PlaylistVideoMap[AllVideo[i].PlaylistId] = make([]pkgs.VideoItems,0,2000)
			PlaylistVideoMap[AllVideo[i].PlaylistId] = append(PlaylistVideoMap[AllVideo[i].PlaylistId],AllVideo[i])
		} else {
			PlaylistVideoMap[AllVideo[i].PlaylistId] = append(PlaylistVideoMap[AllVideo[i].PlaylistId], AllVideo[i])
		}
	}
	client.Disconnect(ctx)
}

func _handelUpdate(r *gin.Context){
	log.Println("Update in progress...")
	pkgs.Update()
	getData()
	r.JSON(200, map[string] string{
		"status": "Done",
		"message": "Update was done sucessfully",
	})
}
func routineUpdate() {
	for true {
		ctx := context.TODO()
		opt := options.Client().ApplyURI(
			"mongodb+srv://"+username+":"+password+"@main-2t9o8.mongodb.net/test?retryWrites=true&w=majority",
		)
		client, err := mongo.Connect(ctx, opt)
		if err != nil {
			log.Fatal(err)
		}
		updateCollection := client.Database("Logs").Collection("UPDATED")
		l := options.Find()
		l.SetSort(bson.D{{"timestamp", -1}})
		var tStamp logs
		cur ,_ := updateCollection.Find(ctx,bson.D{{}},l)
		cur.Next(ctx)
		cur.Decode(&tStamp)
		if (time.Now().Unix()-tStamp.Timestamp) > 21600 {
			pkgs.Update()
		}
		client.Disconnect(ctx)
		time.Sleep(15*60*time.Second)
	}
}

func main() {
	getData()
	go routineUpdate()
	r := gin.Default()
	r.Use(cors.Default())
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/getPlaylistVideos",_handlePlaylistVideo)
			v1.GET("/getPlaylist",_handleParticularPlaylist)
		}
		api.GET("/update",_handelUpdate)
	}
	r.Run()

}