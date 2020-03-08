package gom

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/davidddw/gopj/gocms/app/domain"
)

// ArticleDao dao
type ArticleDao struct {
}

// NewArticleDao new dao
func NewArticleDao() *ArticleDao {
	return &ArticleDao{}
}

// getCol get collection
func (a *ArticleDao) getCol(client *mongo.Client) *mongo.Collection {
	return getCollect(client, "article")
}

// Insert insert function
func (a *ArticleDao) Insert(article string, chapter *domain.Article) (err error) {
	client, context := GetSession()
	coll := a.getCol(client)

	// insert article
	if article == "" {
		_, err = coll.InsertOne(context, chapter)
		return
	}

	// insert chapter
	coll.UpdateOne(context, bson.M{"name": article}, bson.M{
		"$push": bson.M{"chapters": chapter},
	})

	return
}

// UpdateTypeName update Type
func (a *ArticleDao) UpdateTypeName(oldName, newName string) (err error) {
	client, context := GetSession()
	coll := a.getCol(client)
	_, err = coll.UpdateMany(context,
		bson.M{"type": oldName},
		bson.M{
			"$set": bson.M{
				"type": newName,
			},
		},
	)
	return
}

// UpdateArticle update article
func (a *ArticleDao) UpdateArticle(oldType, oldArticle string, article *domain.Article) error {
	client, context := GetSession()
	coll := a.getCol(client)
	_, err := coll.UpdateOne(context,
		bson.M{"type": oldType, "name": oldArticle},
		bson.M{
			"$set": bson.M{
				"name":        article.Name,
				"title":       article.Title,
				"type":        article.Type,
				"description": article.Description,
				"content":     article.Content,
				"sort":        article.Sort,
				"prev":        article.Prev,
				"next":        article.Next,
				"good":        article.Good,
				"top":         article.Top,
				"tags":        article.Tags,
				"hits":        article.Hits,
				"author":      article.Author,
				"createdat":   article.CreatedAt,
			},
		},
	)
	return err
}

// UpdateChapter update chapter
func (a *ArticleDao) UpdateChapter(oldType, oldArticle, oldChapter string, article *domain.Article) error {
	client, context := GetSession()
	coll := a.getCol(client)
	_, err := coll.UpdateOne(context,
		bson.M{"type": oldType, "name": oldArticle, "chapters.name": oldChapter},
		bson.M{
			"$set": bson.M{
				"chapters.0.name":        article.Name,
				"chapters.0.title":       article.Title,
				"chapters.0.type":        article.Type,
				"chapters.0.description": article.Description,
				"chapters.0.content":     article.Content,
				"chapters.0.sort":        article.Sort,
				"chapters.0.prev":        article.Prev,
				"chapters.0.next":        article.Next,
				"chapters.0.good":        article.Good,
				"chapters.0.top":         article.Top,
				"chapters.0.tags":        article.Tags,
				"chapters.0.hits":        article.Hits,
				"chapters.0.author":      article.Author,
				"chapters.0.createdat":   article.CreatedAt,
			},
		})
	return err
}

// Delete article
func (a *ArticleDao) Delete(typeName, articleName string) error {
	client, context := GetSession()
	coll := a.getCol(client)
	_, err := coll.DeleteOne(context, bson.M{"type": typeName, "name": articleName})
	return err
}

// DeleteChapter chapter
func (a *ArticleDao) DeleteChapter(typeName, articleName, chapterName string) error {
	client, context := GetSession()
	coll := a.getCol(client)
	return coll.Update(context, bson.M{"type": typeName, "name": articleName}, bson.M{
		"$pull": bson.M{"chapters": bson.M{"name": chapterName}},
	})
}

func (this *ArticleDao) Get(ty, article, chapter string) (a *domain.Article, err error) {
	client, context := GetSession()
	coll := this.getCol(client)
	if chapter == "" {
		err = coll.FindOne(context, bson.M{"type": ty, "name": article}).Decode(&a)
	} else {
		err = coll.Find(context, bson.M{"type": ty, "name": article, "chapters.name": chapter}).Select(bson.M{"name": 1, "type": 1, "title": 1, "chapters.$": 1}).One(&a)
	}
	if err != nil {
		return
	}
	return
}

func (this *ArticleDao) SelectChapter(article string, page, pagesize int, condition bson.M, sorts []string) (t domain.Article, err error) {
	client, context := GetSession()
	coll := this.getCol(client)
	err = coll.Find(condition).Select(
		bson.M{"name": 1, "title": 1, "type": 1, "description": 1, "sort": 1, "prev": 1, "next": 1, "good": 1, "top": 1, "tags": 1, "hits": 1, "author": 1, "createdat": 1,
			"chapters": 1,
		},
	).One(&t)
	if err != nil {
		return
	}
	return
}

// Select select
func (this *ArticleDao) Select(page, pagesize int, condition bson.M, sorts []string) (t []domain.Article, err error) {
	client, context := GetSession()
	coll := this.getCol(client)
	project := bson.M{"name": 1, "title": 1, "type": 1, "description": 1, "sort": 1, "prev": 1, "next": 1, "good": 1, "top": 1, "tags": 1, "hits": 1, "author": 1, "createdat": 1}
	err = coll.Find(context, bson.D{{condition, project}}, &options.FindOptions{
		Limit: pagesize,
		Skip:  (page - 1) * pagesize,
	})
	return
}

// DeleteByTypeName delete by type
func (a *ArticleDao) DeleteByTypeName(typeName string) (err error) {
	client, context := GetSession()
	coll := a.getCol(client)
	_, err = coll.DeleteMany(context, bson.M{"type": typeName})
	return
}
