// Copyright (c) 2021 rookie-ninja
//
// Use of this source code is governed by an Apache-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rookie-ninja/rk-boot/v2"
	"github.com/rookie-ninja/rk-db/mongodb"
	"github.com/rookie-ninja/rk-gin/v2/boot"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

var (
	hiplotTaskCollection         *mongo.Collection
	hiplotDownloadFileCollection *mongo.Collection
	hiplotUploadFileCollection   *mongo.Collection
)

func createCollection(db *mongo.Database, name string) {
	opts := options.CreateCollection()
	err := db.CreateCollection(context.TODO(), name, opts)
	if err != nil {
		fmt.Println("collection exists may be, continue")
	}
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample rk-demo server.
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.basic BasicAuth

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	// Create a new boot instance.
	boot := rkboot.NewBoot()

	boot.Bootstrap(context.TODO())

	// Auto migrate database and init global userDb variable
	db := rkmongo.GetMongoDB("hiplot", "hiplot")
	createCollection(db, "hiplot")

	hiplotTaskCollection = db.Collection("hiplot_task")
	hiplotDownloadFileCollection = db.Collection("hiplot_download_file")
	hiplotUploadFileCollection = db.Collection("hiplot_upload_file")

	// Add shutdown hook function
	boot.AddShutdownHookFunc("shutdown-hook", func() {
		fmt.Println("shutting down")
	})

	// Register APIs
	ginEntry := rkgin.GetGinEntry("hiplot")
	ginEntry.Router.GET("/v1/greeter", Greeter)
	ginEntry.Router.GET("/api/v1/hiplot/:id", GetHiplotTask)

	boot.WaitForShutdownSig(context.TODO())
}

// Greeter handler
// @Summary Greeter
// @Id 1
// @Tags Hello
// @version 1.0
// @Param name query string true "name"
// @produce application/json
// @Success 200 {object} GreeterResponse
// @Router /v1/greeter [get]
func Greeter(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, &GreeterResponse{
		Message: fmt.Sprintf("Hello %s!", ctx.Query("name")),
	})
}

type GreeterResponse struct {
	Message string
}

// *************************************
// *************** Model ***************
// *************************************

type HiplotTask struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

func RetrieveHiplotTask(ctx *gin.Context) {

}

func GetHiplotTask(ctx *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(ctx.Param("id"))
	res := hiplotTaskCollection.FindOne(context.Background(), bson.M{"_id": id})

	if res.Err() != nil {
		ctx.AbortWithError(http.StatusInternalServerError, res.Err())
		return
	}

	task := &HiplotTask{}
	err := res.Decode(task)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, task)
}
