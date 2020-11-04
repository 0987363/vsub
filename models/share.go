package models

import (
	"github.com/0987363/mgo"
	"github.com/0987363/mgo/bson"
)

type Share struct {
	ID     bson.ObjectId   `json:"id" bson:"_id,omitempty"`
	UserID bson.ObjectId   `json:"user_id" bson:"user_id"`
	Key    string          `json:"key" bson:"key"`
	Name   string          `json:"name" bson:"name"`
	NodeID []bson.ObjectId `json:"node_id,omitempty" bson:"node_id,omitempty"`
}

func (share *Share) Create(db *mgo.Session) error {
	c := db.DB(VSub).C(ShareCollection)

	if err := c.Insert(share); err != nil {
		return err
	}

	return nil
}

func (share *Share) Update(db *mgo.Session) error {
	c := db.DB(VSub).C(ShareCollection)

	if err := c.Update(bson.M{
		"_id":     share.ID,
		"user_id": share.UserID,
	}, bson.M{"$set": bson.M{
		"name":    share.Name,
		"node_id": share.NodeID,
	}}); err != nil {
		return err
	}

	return nil
}

func (share *Share) Delete(db *mgo.Session) error {
	c := db.DB(VSub).C(ShareCollection)

	if err := c.RemoveId(share.ID); err != nil {
		return err
	}

	return nil
}

func DeleteShareByUserID(db *mgo.Session, id bson.ObjectId) error {
	c := db.DB(VSub).C(ShareCollection)

	if _, err := c.RemoveAll(bson.M{"user_id": id}); err != nil {
		return err
	}

	return nil
}

func RemoveNodeFromShare(db *mgo.Session, userID, nodeID bson.ObjectId) error {
	c := db.DB(VSub).C(ShareCollection)

	if _, err := c.UpdateAll(bson.M{"user_id": userID}, bson.M{"$pull": bson.M{"node_id": nodeID}}); err != nil {
		return err
	}

	return nil
}

func FindShareByKey(db *mgo.Session, key string) (*Share, error) {
	c := db.DB(VSub).C(ShareCollection)
	share := Share{}
	if err := c.Find(bson.M{"key": key}).One(&share); err != nil {
		return nil, err
	}
	return &share, nil
}

func FindShareByID(db *mgo.Session, id bson.ObjectId) (*Share, error) {
	c := db.DB(VSub).C(ShareCollection)
	share := Share{}
	if err := c.FindId(id).One(&share); err != nil {
		return nil, err
	}
	return &share, nil
}

func ListShareByUserID(db *mgo.Session, id bson.ObjectId) ([]*Share, error) {
	c := db.DB(VSub).C(ShareCollection)
	shares := []*Share{}
	if err := c.Find(bson.M{"user_id": id}).All(&shares); err != nil {
		return nil, err
	}
	return shares, nil
}
