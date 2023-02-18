package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemWithHash interface {
	GetHash() string
}

func GetItemById[T ItemWithHash](items []T, itemId primitive.ObjectID) (*T, bool) {
	for _, item := range items {
		if item.GetHash() == itemId.Hex() {
			return &item, true
		}
	}
	return nil, false
}

func RemoveItemById[T ItemWithHash](items []T, itemId string) {
	for i, item := range items {
		if item.GetHash() == itemId {
			items = append(items[:i], items[i+1:]...)
		}
	}
}

func UpdateItem[T ItemWithHash](items []T, newItem T) {
	for i, item := range items {
		if item.GetHash() == newItem.GetHash() {
			items[i] = newItem
		}
	}
}
