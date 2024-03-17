package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CarMetadata represents the metadata of a 3D car model stored in GridFS.
type CarMetadata struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Length     int64              `bson:"length"`
	ChunkSize  int32              `bson:"chunkSize"`
	UploadDate time.Time          `bson:"uploadDate"`
	Filename   string             `bson:"filename"`
}

type ResourceInfo struct {
	Name   string             `json:"name"`
	FileId primitive.ObjectID `json:"fileId"`
}

type ModelData struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ModelName   string             `bson:"modelName"`
	ModelFileId primitive.ObjectID `bson:"modelFileId"`
	Resources   []ResourceInfo     `bson:"resources"`
}

type CarNames struct {
	UsedNames []string `bson:"usedNames"`
}
