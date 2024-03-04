package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Car struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	ModelPath string             `bson:"modelPath"` // 存储3D模型文件路径
	Metadata  map[string]string  `bson:"metadata"`  // 存储额外元数据，如作者、创建日期等
}
