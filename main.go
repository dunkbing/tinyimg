package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mholt/archiver/v3"
)

type RequestBody struct {
    Keys []string `json:"keys"`
}

var bucketName = "optipic"
var accountId = "<account_id>"
var accessKeyId = "<access_key_id>"
var accessKeySecret = "<access_key_secret>"

func main() {
    http.HandleFunc("/zip-and-upload", func(w http.ResponseWriter, r *http.Request) {
        var body RequestBody
        err := json.NewDecoder(r.Body).Decode(&body)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

		r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
			}, nil
		})
        cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithEndpointResolverWithOptions(r2Resolver),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, ""),
		),
	)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        client := s3.NewFromConfig(cfg)

        tempFiles := make([]string, len(body.Keys))

        for i, key := range body.Keys {
            file, err := os.Create(key)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            defer file.Close()

            input := &s3.GetObjectInput{
                Bucket: aws.String(bucketName),
                Key:    aws.String(key),
            }

            downloader := manager.NewDownloader(client)
            _, err = downloader.Download(context.TODO(), file, input)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }

            tempFiles[i] = file.Name()
        }

        zipName := "archive.zip"
        err = archiver.Archive(tempFiles, zipName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        file, err := os.Open(zipName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer file.Close()

        uploader := manager.NewUploader(client)
        _, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
            Bucket: aws.String("myBucket"),
            Key:    aws.String(zipName),
            Body:   file,
        })
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, "Successfully zipped and uploaded: %s\n", zipName)
    })

    http.ListenAndServe(":8080", nil)
}
