package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luopanforever/backgreendrive/entity"
	"github.com/luopanforever/backgreendrive/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NameRepository struct {
	DB *mongo.Database
}

func newNameRepository() *NameRepository {
	return (*NameRepository)(NewRepository())
}

func GetNameRepository() *NameRepository {
	return newNameRepository()
}

func (r *NameRepository) FindAvailableName() (string, error) {
	var result struct {
		UsedNames []string `bson:"usedNames"`
	}
	if err := r.DB.Collection("carNames").FindOne(context.Background(), bson.D{}).Decode(&result); err != nil {
		return "", err
	}

	// [car1, car3] return car4
	maxNum := 0
	for _, name := range result.UsedNames {
		if len(name) > 3 {
			// 从名称中提取编号部分，并转换为整数
			if num, err := strconv.Atoi(name[3:]); err == nil {
				// 更新最大编号
				if num > maxNum {
					maxNum = num
				}
			}
		}
	}

	// 生成下一个可用的名称，即最大编号加1
	return fmt.Sprintf("car%d", maxNum+1), nil

	// [car1, car3] return car2
	// nameMap := make(map[int]bool)
	// for _, name := range result.UsedNames {
	// 	if len(name) > 3 {
	// 		if num, err := strconv.Atoi(name[3:]); err == nil {
	// 			nameMap[num] = true
	// 		}
	// 	}
	// }
	// for i := 1; ; i++ {
	// 	if !nameMap[i] {
	// 		return fmt.Sprintf("car%02d", i), nil
	// 	}
	// }
}

// name管理
// 获取汽车名字列表
func (r *NameRepository) GetNameList() ([]string, error) {
	var result struct {
		UsedNames []string `bson:"usedNames"`
	}
	if err := r.DB.Collection("carNames").FindOne(context.Background(), bson.D{}).Decode(&result); err != nil {
		return nil, err
	}
	return result.UsedNames, nil
}

// AddCarName adds a new car name to the usedNames array in carNames collection.
func (r *NameRepository) AddCarNameTest(c *gin.Context) {
	name := c.Param("carName")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$push": bson.M{"usedNames": name}}
	_, err := r.DB.Collection("carNames").UpdateOne(ctx, bson.M{}, update)

	if err != nil {
		response.Fail(c, "Failed to add car name", gin.H{"error": err.Error()})
		return
	}
	response.Success(c, gin.H{"carname": name}, "add success")

}
func (r *NameRepository) AddCarName(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$push": bson.M{"usedNames": name}}
	_, err := r.DB.Collection("carNames").UpdateOne(ctx, bson.M{}, update)
	return err
}

// RemoveCarName removes a car name from the usedNames array in carNames collection.
func (r *NameRepository) RemoveCarNameTest(c *gin.Context) {
	name := c.Param("carName")
	println(name)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$pull": bson.M{"usedNames": name}}
	_, err := r.DB.Collection("carNames").UpdateOne(ctx, bson.M{}, update)
	if err != nil {
		response.Fail(c, "Failed to delete car name", gin.H{"error": err.Error()})
		return
	}
	response.Success(c, gin.H{"carname": name}, "delete success")
}
func (r *NameRepository) RemoveCarName(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$pull": bson.M{"usedNames": name}}
	_, err := r.DB.Collection("carNames").UpdateOne(ctx, bson.M{}, update)
	return err
}

// CarNameExists checks if the given car name already exists in the carNames array.
func (r *NameRepository) CarNameExists(carName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var carNames entity.CarNames
	err := r.DB.Collection("carNames").FindOne(ctx, bson.D{}).Decode(&carNames)
	if err != nil {
		return false, err
	}

	for _, name := range carNames.UsedNames {
		if name == carName {
			return true, nil
		}
	}
	return false, nil
}
