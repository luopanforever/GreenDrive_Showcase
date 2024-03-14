package model

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
	ResourceName string             `json:"resourceName"`
	FileID       primitive.ObjectID `json:"fileId"`
}
