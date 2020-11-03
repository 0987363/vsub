package models

import (
	"encoding/base64"
	"encoding/json"

	"github.com/0987363/mgo"
	"github.com/0987363/mgo/bson"
)

type NodeV2ray struct {
	Port string `json:"port" bson:"port"`
	Ps   string `json:"ps" bson:"ps"`
	Tls  string `json:"tls" bson:"tls"`
	ID   string `json:"id" bson:"id"`
	Aid  string `json:"aid" bson:"aid"`
	V    string `json:"v" bson:"v"`
	Host string `json:"host" bson:"host"`
	Type string `json:"type" bson:"type"`
	Path string `json:"path" bson:"path"`
	Net  string `json:"net" bson:"net"`
	Add  string `json:"add" bson:"add"`
}

type Node struct {
	ID     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserID bson.ObjectId `json:"user_id" bson:"user_id"`
	Class  string        `json:"class" bson:"class"` // v2ray, ss, ssr

	V2ray *NodeV2ray `json:"v2ray,omitempty" bson:"v2ray,omitempty"`
}

func DecodeV2ray(data string) (*NodeV2ray, error) {
	j, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	v := NodeV2ray{}
	if err := json.Unmarshal(j, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (node *Node) Create(db *mgo.Session) error {
	c := db.DB(VSub).C(NodeCollection)

	if err := c.Insert(node); err != nil {
		return err
	}

	return nil
}
func (node *Node) Update(db *mgo.Session) error {
	c := db.DB(VSub).C(NodeCollection)

	if err := c.Update(bson.M{"_id": node.ID, "user_id": node.UserID}, bson.M{"$set": node}); err != nil {
		return err
	}

	return nil
}

func FindNodeByID(db *mgo.Session, id bson.ObjectId) (*Node, error) {
	c := db.DB(VSub).C(NodeCollection)
	node := Node{}

	if err := c.FindId(id).One(&node); err != nil {
		return nil, err
	}

	return &node, nil
}

func (node *Node) Delete(db *mgo.Session) error {
	c := db.DB(VSub).C(NodeCollection)

	if err := c.RemoveId(node.ID); err != nil {
		return err
	}

	return nil
}

func DeleteNodeByUserID(db *mgo.Session, id bson.ObjectId) error {
	c := db.DB(VSub).C(NodeCollection)

	if _, err := c.RemoveAll(bson.M{"user_id": id}); err != nil {
		return err
	}

	return nil
}

func ListNodeByUserID(db *mgo.Session, id bson.ObjectId) ([]*Node, error) {
	c := db.DB(VSub).C(NodeCollection)
	nodes := []*Node{}
	if err := c.Find(bson.M{"user_id": id}).All(&nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func ListNodeByFilter(db *mgo.Session, m bson.M) ([]*Node, error) {
	c := db.DB(VSub).C(NodeCollection)
	nodes := []*Node{}
	if err := c.Find(m).All(&nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}
