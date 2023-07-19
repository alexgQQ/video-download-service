package main

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/lithammer/shortuuid/v4"
	"github.com/redis/go-redis/v9"
)

// Need to read up on reflection
type Download struct {
	OriginalUrl string `datastore:"original_url" redis:"original_url"`
	DownloadUrl string `datastore:"download_url" redis:"download_url"`
	Complete    bool   `datastore:"complete" redis:"complete"`
}

var GCP_PROJECT string
var GCP_TOPIC_ID string
var DOWNLOAD_HOST string
var REDIS_HOST string
var REDIS_PORT string
var GCP_BUCKET string

func loadEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return value, fmt.Errorf("Environment value `%s` is not set", key)
	}
	return value, nil
}

func loadConfig() error {
	var err error
	GCP_PROJECT, err = loadEnvVar("GCP_PROJECT")
	GCP_TOPIC_ID, err = loadEnvVar("GCP_TOPIC_ID")
	GCP_BUCKET, err = loadEnvVar("GCP_BUCKET")
	DOWNLOAD_HOST, err = loadEnvVar("DOWNLOAD_HOST")
	REDIS_HOST, err = loadEnvVar("REDIS_HOST")
	REDIS_PORT, err = loadEnvVar("REDIS_PORT")
	if err != nil {
		return err
	}
	return nil
}

func http500(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func http405(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func getFromRedis(uuid string) (*Download, error) {
	redisAddr := fmt.Sprintf("%s:%s", REDIS_HOST, REDIS_PORT)
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	data, err := client.HGetAll(ctx, uuid).Result()

	// No error is raised if there is no record
	// so instead we check for an empty map
	if len(data) == 0 {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("Error reading from redis: %w", err)
	}
	return &Download{data["original_url"], data["download_url"], true}, nil
}

func publishMessage(url, jobId string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, GCP_PROJECT)
	if err != nil {
		return fmt.Errorf("Error creating pubsub client: %w", err)
	}
	defer client.Close()

	t := client.Topic(GCP_TOPIC_ID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte("Video Download Request"),
		Attributes: map[string]string{
			"url":  url,
			"uuid": jobId,
		},
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("Error publishing message: %w", err)
	}
	return nil
}

func requestDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http405(w)
		return
	}

	targetUrl := r.FormValue("download_link")
	if len(targetUrl) == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	jobId := shortuuid.New()
	err := publishMessage(targetUrl, jobId)
	if err != nil {
		log.Fatal("main: publishMessage: %w", err)
		http500(w)
		return
	}

	cookie := http.Cookie{
		Name:     "downloadId",
		Value:    jobId,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/download", http.StatusFound)
}

func checkDownloadStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http405(w)
		return
	}
	cookie, err := r.Cookie("downloadId")
	// Not too sure how to handle these
	if errors.Is(err, http.ErrNoCookie) {
		//
	} else if err != nil {
		//
	}
	jobId := cookie.Value

	data, err := getFromRedis(jobId)
	if err != nil {
		http500(w)
		return
	} else if data == nil {
		// jobId is not found and the download may still be processing
		// return a NoContent (204) to signal htmx to ignore the response
		// https://htmx.org/docs/#requests
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Return this specific code to signal htmx to stop polling
	// https://htmx.org/docs/#polling
	w.WriteHeader(286)

	if data.Complete == false {
		// Render error page or component
	}

	link := fmt.Sprintf("%s/%s.mp4", DOWNLOAD_HOST, jobId)
	tmpl, err := template.ParseFiles(filepath.Join("templates", "html", "available.html"))
	err = tmpl.Execute(w, struct{ DownloadLink string }{link})
}

func renderIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http405(w)
		return
	}

	lp := filepath.Join("templates", "html", "layout.html")
	fp := filepath.Join("templates", "html", "index.html")

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Fatal("template: ParseFiles: %w", err)
		http500(w)
		return
	}
	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Fatal("template: ExecuteTemplate: %w", err)
		http500(w)
		return
	}
}

func renderDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http405(w)
		return
	}

	lp := filepath.Join("templates", "html", "layout.html")
	fp := filepath.Join("templates", "html", "download.html")

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		log.Fatal("template: ParseFiles: %w", err)
		http500(w)
		return
	}
	err = tmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Fatal("template: ExecuteTemplate: %w", err)
		http500(w)
		return
	}
}

func signedDownloadUrl(uuid string) (string, error) {
	object := fmt.Sprintf("vdl-test/%s", uuid)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(15 * time.Minute),
	}

	u, err := client.Bucket(GCP_BUCKET).SignedURL(object, opts)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %w", GCP_BUCKET, err)
	}

	return u, nil
}

func fileRedirect(w http.ResponseWriter, r *http.Request) {
	filename := path.Base(r.URL.Path)
	target, err := signedDownloadUrl(filename)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func main() {
	err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public"))))
	http.HandleFunc("/index.html", renderIndex)
	http.HandleFunc("/request-download", requestDownload)
	http.HandleFunc("/download-status", checkDownloadStatus)
	http.HandleFunc("/download", renderDownload)
	http.HandleFunc("/files/", fileRedirect)

	fmt.Printf("Listening on port 8080\n")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
