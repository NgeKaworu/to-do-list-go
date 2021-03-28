package engine

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/NgeKaworu/to-do-list-go/src/models"
	"github.com/NgeKaworu/to-do-list-go/src/parsup"
	"github.com/NgeKaworu/to-do-list-go/src/resultor"
	"github.com/NgeKaworu/to-do-list-go/src/utils"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddTask 添加记录
func (d *DbEngine) AddTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		resultor.RetFail(w, err)
		return
	}
	if len(body) == 0 {
		resultor.RetFail(w, errors.New("not has body"))
		return
	}

	p, err := parsup.ParSup().ConvJSON(body)
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	err = utils.Required(p, map[string]string{
		"title": "请填写任务名",
		"level": "请选一个优先级",
	})

	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	t := d.GetColl(models.TTask)

	p["uid"] = uid
	p["createAt"] = time.Now().Local()

	res, err := t.InsertOne(context.Background(), p)
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	resultor.RetOk(w, res.InsertedID.(primitive.ObjectID).Hex())
}

// SetTask 更新记录
func (d *DbEngine) SetTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		resultor.RetFail(w, err)
		return
	}
	if len(body) == 0 {
		resultor.RetFail(w, errors.New("not has body"))
		return
	}

	p, err := parsup.ParSup().ConvJSON(body)
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	err = utils.Required(p, map[string]string{
		"event": "请填写发生了什么",
		"tid":   "请至少选一个标签",
		"id":    "ID不能为空",
	})

	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	t := d.GetColl(models.TTask)

	p["uid"] = uid
	p["updateAt"] = time.Now().Local()

	id := p["id"]
	delete(p, "id")

	res := t.FindOneAndUpdate(context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": p},
	)
	if res.Err() != nil {
		resultor.RetFail(w, res.Err())
		return
	}

	resultor.RetOk(w, "修改成功")
}

// RemoveTask 删除记录
func (d *DbEngine) RemoveTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		resultor.RetFail(w, err)
		return
	}
	id, err := primitive.ObjectIDFromHex(ps.ByName("id"))
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	t := d.GetColl(models.TTask)

	res := t.FindOneAndDelete(context.Background(), bson.M{"_id": id, "uid": uid})

	if res.Err() != nil {
		resultor.RetFail(w, res.Err())
		return
	}

	resultor.RetOk(w, "删除成功")
}

// ListTask record列表
func (d *DbEngine) ListTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q := r.URL.Query()
	l := q.Get("limit")
	s := q.Get("skip")

	uid, err := primitive.ObjectIDFromHex(r.Header.Get("uid"))
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	limit, _ := strconv.ParseInt(l, 10, 64)
	skip, _ := strconv.ParseInt(s, 10, 64)

	t := d.GetColl(models.TTask)

	cur, err := t.Find(context.Background(), bson.M{
		"uid": uid,
	}, options.Find().SetSort(bson.M{"level": -1}).SetSkip(skip).SetLimit(limit))

	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	list := make([]models.MainTask, 0)

	err = cur.All(context.Background(), &list)
	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	total, err := t.CountDocuments(context.Background(), bson.M{})

	if err != nil {
		resultor.RetFail(w, err)
		return
	}

	resultor.RetOkWithTotal(w, list, total)
}
