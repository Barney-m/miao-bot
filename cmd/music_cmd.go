package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"miao-bot/services"
	"miao-bot/utils"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

type YoutubeService struct {
	maxDurationInSeconds int
	fileDirectory        string
}

func ExecuteMusic(discord *discordgo.Session, message *discordgo.MessageCreate) {
	cmdBlock := strings.Split(message.Content, " ")

	// TODO: Implements Pause, Skip, Queue
	switch cmdBlock[1] {
	case "play", "p":
		HandleYoutubeCommand(discord, message, strings.Join(cmdBlock[2:], ""), NewYoutubeService(9999))
	case "pause":
	case "skip":
	default:
	}
}

// Initialize Service
func NewYoutubeService(maxDurationInSeconds int) *YoutubeService {
	_ = os.Mkdir("tmpAudio", os.ModePerm)
	return &YoutubeService{
		maxDurationInSeconds: maxDurationInSeconds,
		fileDirectory:        "tmpAudio",
	}
}

// Search and Download Youtube Audio
func (y *YoutubeService) SearchAndDownload(query string) (*services.Media, error) {
	timeout := make(chan bool, 1)
	result := make(chan searchAndDownloadResult, 1)
	go func() {
		time.Sleep(60 * time.Second)
		timeout <- true
	}()
	go func() {
		result <- y.doSearchAndDownload(query)
	}()
	select {
	case <-timeout:
		return nil, errors.New("timed out")
	case result := <-result:
		return result.Media, result.Error
	}
}

func (y *YoutubeService) doSearchAndDownload(query string) searchAndDownloadResult {
	start := time.Now()
	youtubeDownloader, err := exec.LookPath("yt-dlp") // Make sure there is yt-dlp installed in local machine
	if err != nil {
		return searchAndDownloadResult{Error: errors.New("yt-dlp not found in path")}
	} else {
		args := []string{
			fmt.Sprintf("ytsearch10:%s", strings.ReplaceAll(query, "\"", "")),
			"--extract-audio",
			"--audio-format", "opus",
			"--no-playlist",
			"--match-filter", fmt.Sprintf("duration < %d & !is_live", y.maxDurationInSeconds),
			"--max-downloads", "1",
			"--output", fmt.Sprintf("%s/%d-%%(id)s.opus", y.fileDirectory, start.Unix()),
			"--quiet",
			"--print-json",
			"--ignore-errors", // Ignores unavailable videos
			"--no-color",
			"--no-check-formats",
		} // Set yt-dtp Arguments
		log.Printf("yt-dlp %s", strings.Join(args, " "))
		cmd := exec.Command(youtubeDownloader, args...)
		if data, err := cmd.Output(); err != nil && err.Error() != "exit status 101" {
			return searchAndDownloadResult{Error: fmt.Errorf("failed to search and download audio: %s\n%s", err.Error(), string(data))}
		} else {
			videoMetadata := videoMetadata{}
			err = json.Unmarshal(data, &videoMetadata)
			if err != nil {
				return searchAndDownloadResult{Error: fmt.Errorf("failed to unmarshal video metadata: %w", err)}
			}
			return searchAndDownloadResult{
				Media: services.NewMedia(
					videoMetadata.Title,
					strings.Split(videoMetadata.Filename, ".")[0]+".opus", // Replace actual file name with .webm to .opus which is audio file
					videoMetadata.Uploader,
					fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoMetadata.ID),
					videoMetadata.Thumbnail,
					int(videoMetadata.Duration),
				),
			}
		}
	}
}

type searchAndDownloadResult struct {
	Media *services.Media
	Error error
}

type videoMetadata struct {
	ID                   string      `json:"id"`
	Title                string      `json:"title"`
	Thumbnail            string      `json:"thumbnail"`
	Description          string      `json:"description"`
	Uploader             string      `json:"uploader"`
	UploaderID           string      `json:"uploader_id"`
	UploaderURL          string      `json:"uploader_url"`
	ChannelID            string      `json:"channel_id"`
	ChannelURL           string      `json:"channel_url"`
	Duration             int         `json:"duration"`
	ViewCount            int         `json:"view_count"`
	AverageRating        interface{} `json:"average_rating"`
	AgeLimit             int         `json:"age_limit"`
	WebpageURL           string      `json:"webpage_url"`
	Categories           []string    `json:"categories"`
	Tags                 []string    `json:"tags"`
	PlayableInEmbed      bool        `json:"playable_in_embed"`
	LiveStatus           interface{} `json:"live_status"`
	ReleaseTimestamp     interface{} `json:"release_timestamp"`
	CommentCount         interface{} `json:"comment_count"`
	LikeCount            int         `json:"like_count"`
	Channel              string      `json:"channel"`
	ChannelFollowerCount int         `json:"channel_follower_count"`
	UploadDate           string      `json:"upload_date"`
	Availability         string      `json:"availability"`
	OriginalURL          string      `json:"original_url"`
	WebpageURLBasename   string      `json:"webpage_url_basename"`
	WebpageURLDomain     string      `json:"webpage_url_domain"`
	Extractor            string      `json:"extractor"`
	ExtractorKey         string      `json:"extractor_key"`
	PlaylistCount        int         `json:"playlist_count"`
	Playlist             string      `json:"playlist"`
	PlaylistID           string      `json:"playlist_id"`
	PlaylistTitle        string      `json:"playlist_title"`
	PlaylistUploader     interface{} `json:"playlist_uploader"`
	PlaylistUploaderID   interface{} `json:"playlist_uploader_id"`
	NEntries             int         `json:"n_entries"`
	PlaylistIndex        int         `json:"playlist_index"`
	LastPlaylistIndex    int         `json:"__last_playlist_index"`
	PlaylistAutonumber   int         `json:"playlist_autonumber"`
	DisplayID            string      `json:"display_id"`
	Fulltitle            string      `json:"fulltitle"`
	DurationString       string      `json:"duration_string"`
	RequestedSubtitles   interface{} `json:"requested_subtitles"`
	Asr                  int         `json:"asr"`
	Filesize             int         `json:"filesize"`
	FormatID             string      `json:"format_id"`
	FormatNote           string      `json:"format_note"`
	SourcePreference     int         `json:"source_preference"`
	Fps                  interface{} `json:"fps"`
	AudioChannels        int         `json:"audio_channels"`
	Height               interface{} `json:"height"`
	Quality              float64     `json:"quality"`
	HasDrm               bool        `json:"has_drm"`
	Tbr                  float64     `json:"tbr"`
	URL                  string      `json:"url"`
	Width                interface{} `json:"width"`
	Language             string      `json:"language"`
	LanguagePreference   int         `json:"language_preference"`
	Preference           interface{} `json:"preference"`
	Ext                  string      `json:"ext"`
	Vcodec               string      `json:"vcodec"`
	Acodec               string      `json:"acodec"`
	DynamicRange         interface{} `json:"dynamic_range"`
	Abr                  float64     `json:"abr"`
	Filename             string      `json:"filename"`
}

// To handle overall youtube play audio command
func HandleYoutubeCommand(session *discordgo.Session, message *discordgo.MessageCreate, query string, ytService *YoutubeService) {
	guildsMutex.Lock()
	activeGuild := guilds[message.GuildID]
	guildsMutex.Unlock()

	if activeGuild != nil {
		if activeGuild.IsMediaQueueFull() {
			_, _ = session.ChannelMessageSend(message.ChannelID, "The queue is full!")
			return
		}
	} else {
		activeGuild = services.NewActiveGuild(GetGuildNameByID(session, message.GuildID))
		guildsMutex.Lock()
		guilds[message.GuildID] = activeGuild
		guildsMutex.Unlock()
	}
	// Find the voice channel the user is in
	voiceChannelId, err := GetVoiceChannelWhereMessageAuthorIs(session, message)
	if err != nil {
		log.Printf("[%s] Failed to find voice channel where message author is located: %s", activeGuild.Name, err.Error())
		_ = session.MessageReactionAdd(message.ChannelID, message.ID, "❌")
		_, _ = session.ChannelMessageSend(message.ChannelID, err.Error())
		return
	} else {
		log.Printf("[%s] Found user %s in voice channel %s", activeGuild.Name, message.Author.Username, voiceChannelId)
		_ = session.MessageReactionAdd(message.ChannelID, message.ID, "✅")
	}
	log.Printf("[%s] Searching for \"%s\"", activeGuild.Name, query)
	sessionMessage, _ := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf(":mag: Searching for `%s`...", query))
	media, err := ytService.SearchAndDownload(query)
	if err != nil {
		log.Printf("[%s] Unable to find video for query \"%s\": %s", activeGuild.Name, query, err.Error())
		_, _ = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Unable to find video for query `%s`: %s", query, err.Error()))
		return
	}
	log.Printf("[%s] Successfully searched for and extracted audio from video with title \"%s\" to \"%s\"", activeGuild.Name, media.Title, media.FilePath)
	sessionMessage, _ = session.ChannelMessageEdit(sessionMessage.ChannelID, sessionMessage.ID, fmt.Sprintf(":white_check_mark: Found matching video titled `%s`!", media.Title))
	go func(session *discordgo.Session, message *discordgo.Message) {
		time.Sleep(5 * time.Second)
		_ = session.ChannelMessageDelete(sessionMessage.ChannelID, sessionMessage.ID)
	}(session, sessionMessage)
	// Add song to guild queue
	createNewWorker := false
	if !activeGuild.IsStreaming() {
		log.Printf("[%s] Preparing for streaming", activeGuild.Name)
		activeGuild.PrepareForStreaming(viper.GetInt("music.max-queue-size"))
		// If the channel was nil, it means that there was no worker
		createNewWorker = true
	}
	activeGuild.EnqueueMedia(media)
	log.Printf("[%s] Added media with title \"%s\" to queue at position %d", activeGuild.Name, media.Title, activeGuild.MediaQueueSize())
	_, _ = session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
		URL:         media.URL,
		Title:       media.Title,
		Description: fmt.Sprintf("Position in queue: %d", activeGuild.MediaQueueSize()),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: media.Thumbnail,
		},
	})
	if createNewWorker {
		log.Printf("[%s] Starting worker", activeGuild.Name)
		go func() {
			err = utils.NewMusicWorker(session, activeGuild, message.GuildID, voiceChannelId)
			if err != nil {
				log.Printf("[%s] Failed to start worker: %s", activeGuild.Name, err.Error())
				_, _ = session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("❌ Unable to start voice worker: %s", err.Error()))
				_ = os.Remove(media.FilePath)
				return
			}
		}()
	}
}
