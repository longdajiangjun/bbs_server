package model

import (
	"bbs_server/config"
	"bbs_server/database"
	"log"

	"gopkg.in/mgo.v2/bson"
)

// Post 贴子结构
type Post struct {
	TopStorey  `json:"topStorey"`
	ReList1    []Reply1      `json:"reList1"`
	ReList2    []Reply2      `json:"reList2"`
	UpdateTime string        `json:"time"`
	TID        bson.ObjectId `json:"tid"`
}

// TopStorey .
type TopStorey struct {
	Title    string   `json:"title"`
	ImgList  []string `json:"imgList"`
	ReplyNum uint32   `json:"replyNum"`
	Support  uint32   `json:"support"`
	ReadNum  uint32   `json:"readNum"`
	ShareMsg `bson:",inline"`
}

// Reply1 .
type Reply1 struct {
	ID       bson.ObjectId `json:"id"`
	Show     bool          `json:"show"`
	RName    string        `json:"rName"`
	ShareMsg `bson:",inline"`
}

// Reply2 .
type Reply2 struct {
	ID       bson.ObjectId `json:"id"`
	RID      bson.ObjectId `json:"rid"`
	RName    string        `json:"rName"`
	Show     bool          `json:"show"`
	ShareMsg `bson:",inline"`
}

// ShareMsg .
type ShareMsg struct {
	HeadImg    string        `json:"headImg"`
	UName      string        `json:"uName"`
	CreateTime string        `json:"createTime"`
	Content    string        `json:"content"`
	TID        bson.ObjectId `json:"tid"`
	Type       int           `json:"type"`
	Tag        int           `json:"tag"`
	Topic      string        `json:"topic"`
}

// HeadPost 置顶帖
type HeadPost struct {
	TID   bson.ObjectId `json:"tid"`
	Title string        `json:"title"`
}

// Save 保存贴子信息.
func (p *Post) Save() bool {
	session := database.Session.Clone()
	defer session.Close()
	c := session.DB(config.DbName).C("bbs_user")
	c.Update(bson.M{"uname": p.TopStorey.UName}, bson.M{"$inc": bson.M{"exp": 15}})
	c.Update(bson.M{"uname": p.TopStorey.UName}, bson.M{"$inc": bson.M{"integral": 15}})
	c = session.DB(config.DbName).C("bbs_topics")
	c.Update(bson.M{"name": p.ShareMsg.Topic}, bson.M{"$inc": bson.M{"num": 1}})
	c = session.DB(config.DbName).C("bbs_posts")
	err := c.Insert(p)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// UpdatePosts 获取所有贴子.
func UpdatePosts(postsPool *[]Post, topic string, Type int) error {
	session := database.Session.Clone()
	defer session.Close()
	c := session.DB(config.DbName).C("bbs_posts")
	return c.Find(bson.M{"topstorey.topic": topic, "topstorey.type": Type}).Sort("-_id").All(postsPool)
}

// Get 获取单个贴子详情.
func (p *Post) Get(tid bson.ObjectId) bool {
	session := database.Session.Clone()
	defer session.Close()
	c := session.DB(config.DbName).C("bbs_posts")
	err := c.Find(bson.M{"tid": tid}).One(p)
	if err != nil {
		log.Println(err)
		return false
	}
	c = session.DB(config.DbName).C("bbs_user")
	c.Update(bson.M{"uname": p.TopStorey.UName}, bson.M{"$inc": bson.M{"readnum": 1}})
	return true
}

// Save 保存一级回复信息.
func (reply1 *Reply1) Save(tid bson.ObjectId) bool {
	session := database.Session.Clone()
	defer session.Close()
	c := session.DB(config.DbName).C("bbs_user")
	c.Update(bson.M{"uname": reply1.UName}, bson.M{"$inc": bson.M{"exp": 5}})
	c.Update(bson.M{"uname": reply1.UName}, bson.M{"$inc": bson.M{"integral": 5}})
	c.Update(bson.M{"uname": reply1.RName}, bson.M{"$inc": bson.M{"replynum": 1}})
	c = session.DB(config.DbName).C("bbs_posts")
	err := c.Update(bson.M{"tid": tid}, bson.M{"$push": bson.M{"relist1": reply1}})
	err = c.Update(bson.M{"tid": tid}, bson.M{"$inc": bson.M{"topstorey.replynum": 1}})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Save 保存二级回复信息
func (reply2 *Reply2) Save(id bson.ObjectId) bool {
	session := database.Session.Clone()
	defer session.Close()
	c := session.DB(config.DbName).C("bbs_user")
	c.Update(bson.M{"uname": reply2.UName}, bson.M{"$inc": bson.M{"exp": 5}})
	c.Update(bson.M{"uname": reply2.UName}, bson.M{"$inc": bson.M{"integral": 5}})
	c.Update(bson.M{"uname": reply2.RName}, bson.M{"$inc": bson.M{"replynum": 1}})
	c = session.DB(config.DbName).C("bbs_posts")
	err := c.Update(bson.M{"tid": id}, bson.M{"$push": bson.M{"relist2": reply2}})
	if err != nil {
		log.Println(err)
		return false
	}     
	return true
}

// Del 删除贴子
func (p *Post) Del(tid bson.ObjectId, name string) bool {
	session := database.Session.Clone()
	defer session.Close()
	c := session.DB(config.DbName).C("bbs_posts")
	err := c.Find(bson.M{"tid": tid}).One(p)
	if err != nil {
		log.Println(err)
		return false
	}
	if p.TopStorey.UName == name || name == "admin" {
		c.Remove(bson.M{"tid": tid})
		c = session.DB(config.DbName).C("bbs_feedback")
		c.Remove(bson.M{"tid": tid})
		c = session.DB(config.DbName).C("bbs_zhiding")
		c.Remove(bson.M{"tid": tid})
		c = session.DB(config.DbName).C("bbs_topics")
		c.Update(bson.M{"name": p.ShareMsg.Topic}, bson.M{"$inc": bson.M{"num": -1}})
		return true
	}
	return false
}

// AddSupport 增加点赞数
func (p *Post) AddSupport() bool {
	session := database.Session.Clone()
	defer session.Close()
	c := session.DB(config.DbName).C("bbs_posts")
	err := c.Update(bson.M{"tid": p.TID}, bson.M{"$inc": bson.M{"topstorey.support": 1}})
	if err != nil {
		return false
	}
	return true
}

// ReduceSupport 减少点赞数
func (p *Post) ReduceSupport() bool {
	session := database.Session.Clone()
	defer session.Close()
	c := session.DB(config.DbName).C("bbs_posts")
	err := c.Update(bson.M{"tid": p.TID}, bson.M{"$inc": bson.M{"topstorey.support": -1}})
	if err != nil {
		return false
	}
	return true
}

// AgreeZhiDIng 同意贴子置顶
func (p *Post) AgreeZhiDIng(tid string) bool {
	session := database.Session.Clone()
	defer session.Close()
	c := session.DB(config.DbName).C("bbs_feedback")
	id := bson.ObjectIdHex(tid)
	_, err := c.UpdateAll(bson.M{"tid": id}, bson.M{"$set": bson.M{"status": 1}})
	if err != nil {
		return false
	}
	p.Get(id)
	c = session.DB(config.DbName).C("bbs_zhiding")
	headPost := &HeadPost{}
	err2 := c.Find(bson.M{"tid": id}).One(headPost)
	if err2 != nil {
		headPost.TID = p.TID
		headPost.Title = p.Title
		c.Insert(headPost)
	}
	return true
}

// GetHeadPost 获取置顶帖子
func (p *HeadPost) GetHeadPost() *[]HeadPost {
	session := database.Session.Clone()
	defer session.Close()
	c := session.DB(config.DbName).C("bbs_zhiding")
	result := &[]HeadPost{}
	err := c.Find(nil).Sort("-_id").Limit(3).All(result)
	if err != nil {
		return nil
	}
	return result
}
