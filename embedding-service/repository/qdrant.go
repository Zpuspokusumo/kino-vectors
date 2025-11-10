package repository

import (
	"context"
	"strconv"
	"strings"

	moviepb "github.com/Zpuspokusumo/kino-vectors/contract/golang/movie-services"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

// type MovieInfo struct {
// 	Id       string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
// 	Title    string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
// 	Director string   `protobuf:"bytes,3,opt,name=director,proto3" json:"director,omitempty"`
// 	Year     int64    `protobuf:"bytes,4,opt,name=year,proto3" json:"year,omitempty"`
// 	Genre    []string `protobuf:"bytes,5,rep,name=genre,proto3" json:"genre,omitempty"`
// 	Actors   []string `protobuf:"bytes,6,rep,name=actors,proto3" json:"actors,omitempty"`
// 	Summary  string   `protobuf:"bytes,7,opt,name=summary,proto3" json:"summary,omitempty"`
// }

var x moviepb.RecommendMoviesRequest

func MovieToPayload(m *moviepb.MovieInfo) map[string]*qdrant.Value {
	return map[string]*qdrant.Value{
		"id":       qdrant.NewValueString(m.Id),
		"title":    qdrant.NewValueString(m.Title),
		"director": qdrant.NewValueString(m.Director),
		"year":     qdrant.NewValueInt(int64(m.Year)),
		"genre":    qdrant.NewValueList(stringSliceToValues(m.Genre)),
		"actors":   qdrant.NewValueList(stringSliceToValues(m.Actors)),
		"summary":  qdrant.NewValueString(m.Summary),
	}
}

func stringSliceToValues(items []string) *qdrant.ListValue {
	list := &qdrant.ListValue{}
	for _, v := range items {
		list.Values = append(list.Values, qdrant.NewValueString(v))
	}
	return list
}

type QdrantRepository struct {
	Client *qdrant.Client
}

// localhost:6333:<your-api-key>
func NewClient(key string) (*qdrant.Client, error) {
	credentials := strings.Split(key, ":")
	Host := credentials[0]
	Port, _ := strconv.Atoi(credentials[1])
	var APIKey string
	if len(credentials) > 2 {
		APIKey = credentials[2]
	}

	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   Host,
		Port:   Port,
		APIKey: APIKey,
		UseTLS: false, // local doesnt need tls
		// TLSConfig: &tls.Config{...},
		// GrpcOptions: []grpc.DialOption{},
	})

	return client, err
}

func (repo *QdrantRepository) NewCollection(ctx context.Context, CollectionReq *qdrant.CreateCollection) error {
	return repo.Client.CreateCollection(ctx, CollectionReq)
}

func (repo *QdrantRepository) ListCollections(ctx context.Context) ([]string, error) {
	return repo.Client.ListCollections(ctx)
}

func (repo *QdrantRepository) DeleteCollection(ctx context.Context, collectionName string) error {
	return repo.Client.DeleteCollection(ctx, collectionName)
}

func (repo *QdrantRepository) GetCollection(ctx context.Context, collectionName string) (*qdrant.CollectionInfo, error) {
	return repo.Client.GetCollectionInfo(ctx, collectionName)
}
func (repo *QdrantRepository) UpsertPoints(ctx context.Context, collectionName string, embeddings []float32, payload map[string]*qdrant.Value) (*qdrant.UpdateResult, error) {
	return repo.Client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points: []*qdrant.PointStruct{
			{
				Id:      &qdrant.PointId{PointIdOptions: &qdrant.PointId_Uuid{Uuid: uuid.New().String()}},
				Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: embeddings}}},
				Payload: payload,
			},
		},
	})
}
func (repo *QdrantRepository) SearchGeneral(v []float32) {
	repo.Client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "{collection_name}",
		Query:          qdrant.NewQuery(0.2, 0.1, 0.9, 0.7),
		Filter: &qdrant.Filter{
			Must: []*qdrant.Condition{
				qdrant.NewMatch("city", "London"),
			},
		},
		Params: &qdrant.SearchParams{
			Exact:  qdrant.PtrOf(false),
			HnswEf: qdrant.PtrOf(uint64(128)),
		},
	})
}
func (repo *QdrantRepository) SearchMovie(v []float32, genres []string) ([]*qdrant.ScoredPoint, error) {
	var Filter *qdrant.Filter
	if len(genres) > 0 {
		Filter = &qdrant.Filter{
			Must: []*qdrant.Condition{
				qdrant.NewMatchKeywords("genres", genres...),
			},
		}
	} else {
		Filter = nil
	}
	results, err := repo.Client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: "testmovie01",
		Query:          qdrant.NewQuery(v...),
		Filter:         Filter,
		Params: &qdrant.SearchParams{
			Exact:  qdrant.PtrOf(false),
			HnswEf: qdrant.PtrOf(uint64(128)),
		},
		//WithVectors: &qdrant.WithVectorsSelector{SelectorOptions: &qdrant.WithVectorsSelector_Enable{Enable: true}},
		WithPayload: &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
	})

	return results, err
}
func (repo *QdrantRepository) ScrollMovie(v []float32, genres []string, qty, offset uint32) ([]*qdrant.RetrievedPoint, error) {
	var Filter *qdrant.Filter
	if len(genres) > 0 {
		Filter = &qdrant.Filter{
			Must: []*qdrant.Condition{
				qdrant.NewMatchKeywords("genres", genres...),
			},
		}
	} else {
		Filter = nil
	}
	results, err := repo.Client.Scroll(context.Background(), &qdrant.ScrollPoints{
		CollectionName:   "",
		Filter:           Filter,
		Offset:           &qdrant.PointId{},
		Limit:            new(uint32),
		WithPayload:      &qdrant.WithPayloadSelector{},
		WithVectors:      &qdrant.WithVectorsSelector{},
		ReadConsistency:  &qdrant.ReadConsistency{},
		ShardKeySelector: &qdrant.ShardKeySelector{},
		OrderBy:          &qdrant.OrderBy{},
		Timeout:          new(uint64),
	})

	return results, err
}
