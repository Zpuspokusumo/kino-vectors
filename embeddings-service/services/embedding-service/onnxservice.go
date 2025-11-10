package embeddingservice

import (
	"context"
	"fmt"
	"kino-vectors/repository"
	"log"
	"math"
	"strings"

	moviepb "github.com/Zpuspokusumo/kino-vectors/contract/golang/movie-services"
	"github.com/qdrant/go-client/qdrant"
	"github.com/sugarme/tokenizer"
	pretrained "github.com/sugarme/tokenizer/pretrained"
	ort "github.com/yalue/onnxruntime_go"
)

// const modelPath = `../pyenv/paraphrase-multilingual-mpnet-base-v2.onnx`
const modelPath = `..\models\all-MiniLM-L6-v2\all-MiniLM-L6-v2.onnx`
const tokenizerpath = `..\models\all-MiniLM-L6-v2\tokenizer.json`

type EmbeddingServiceONNX struct {
	modelPath     string
	tokenizerPath string
	repo          repository.QdrantRepository
	seqSize       int
	weightcount   int
	tk            *tokenizer.Tokenizer
}

func (s *EmbeddingServiceONNX) GetSeqsize() int {
	return s.seqSize
}

func MakeServiceONNX(repo repository.QdrantRepository) (*EmbeddingServiceONNX, error) {
	tk, err := pretrained.FromFile(tokenizerpath)
	if err != nil {
		return nil, fmt.Errorf("%v: for path :%v", err, tokenizerpath)
	}

	if tk == nil {
		log.Fatal("tokenizer is nil")
	}
	return &EmbeddingServiceONNX{
		modelPath:     modelPath,
		tokenizerPath: tokenizerpath,
		repo:          repo,
		seqSize:       128,
		weightcount:   3,
		tk:            tk,
	}, nil
}

func CheckExecution() (string, error) {
	//return "OPENVINO", nil
	return "", nil
}

func (service *EmbeddingServiceONNX) Tokenize(text string) ([]int, error) {

	//defer tk.Close()
	ids, err := service.tk.EncodeSingle(text, false)
	if err != nil {
		log.Fatal("tokenizer error:", err)
	}
	if ids == nil || ids.Len() == 0 {
		log.Fatalf("tokenizer returned empty ids for text: %q", text)
	}
	fmt.Println("text ", text)
	fmt.Println(ids.GetIds())

	return ids.GetIds(), err
}

func (service *EmbeddingServiceONNX) Embed(chunklen int, ids, mask []int64) ([]float32, error) {
	executor, err := CheckExecution()
	if err != nil {
		err = fmt.Errorf("error checking executor :%v ", err)
		return nil, err
	}

	inputShape := ort.NewShape(int64(chunklen), 128)
	//inputTensor, err := ort.NewEmptyTensor[float32](inputShape)
	inputTensor, err := ort.NewTensor[int64](inputShape, ids)
	if err != nil {
		err = fmt.Errorf("Error creating input tensor: %w ", err)
		return nil, err
	}
	for i := range ids {
		if ids[i] != 0 { // 0 is pad token id
			mask[i] = 1
		}
	}
	maskShape := ort.NewShape(int64(chunklen), 128)
	maskTensor, err := ort.NewTensor[int64](maskShape, mask)
	if err != nil {
		log.Fatalf("Error creating attention_mask tensor: %v", err)
	}
	//tokentypeids is for cls/sep tokens, unused for embeddings, just return all zero
	tokenTypeIds := make([]int64, chunklen*128) // all zeros (single sentence) un
	tokenTypeIdsTensor, err := ort.NewTensor[int64](maskShape, tokenTypeIds)
	if err != nil {
		log.Fatalf("Error creating tokenTypeIds tensor: %v", err)
	}

	// todo: add chunking/truncation

	outputShape := ort.NewShape(int64(chunklen), 128, 384)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		inputTensor.Destroy()
		err = fmt.Errorf("Error creating output tensor: %w", err)
		return nil, err
	}
	options, err := ort.NewSessionOptions()
	if err != nil {
		inputTensor.Destroy()
		outputTensor.Destroy()
		err = fmt.Errorf("Error creating ORT session options: %w", err)
		return nil, err
	}
	defer options.Destroy()

	if executor == "OPENVINO" {
		options.AppendExecutionProviderOpenVINO(map[string]string{"device_type": "GPU"})
	}

	session, err := ort.NewAdvancedSession(service.modelPath,
		[]string{"input_ids", "attention_mask", "token_type_ids"}, []string{"last_hidden_state"},
		[]ort.ArbitraryTensor{inputTensor, maskTensor, tokenTypeIdsTensor},
		[]ort.ArbitraryTensor{outputTensor},
		options)
	if err != nil {
		inputTensor.Destroy()
		outputTensor.Destroy()
		err = fmt.Errorf("Error creating ORT session: %w", err)
		return nil, err
	}
	defer session.Destroy()

	err = session.Run()
	if err != nil {
		err = fmt.Errorf("Error running ORT session: %w", err)
		return nil, err
	}

	// outputTensor.GetData() returns []float32 of shape [chunklen,128,384]
	data := outputTensor.GetData()

	return data, nil
}

func (service *EmbeddingServiceONNX) GenerateEmbeddings(text string) ([]float32, error) {
	err := ort.InitializeEnvironment()
	if err != nil {
		log.Fatalf("Failed to initialize ORT: %v ", err)
	}
	tokens, err := service.Tokenize(text)
	if err != nil {
		err = fmt.Errorf("error tokenizing :%v ", err)
		return nil, err
	}
	chunklen := int(math.Ceil(float64(len(tokens)) / 128.0))
	ids := make([]int64, chunklen*128)
	for i := range tokens {
		ids[i] = int64(tokens[i])
	}
	mask := make([]int64, len(ids))

	data, err := service.Embed(chunklen, ids, mask)
	if err != nil {
		return nil, err
	}

	meanEmbedding := make([]float32, 384)
	validCount := 0 //number of valid tokens

	// change from iterating over mask for shape safety
	for tokenIdx := 0; tokenIdx < chunklen*128; tokenIdx++ {
		if mask[tokenIdx] == 1 {
			base := tokenIdx * 384
			for j := 0; j < 384; j++ {
				meanEmbedding[j] += data[base+j]
			}
			validCount++
		}
	}

	for j := range meanEmbedding {
		meanEmbedding[j] /= float32(validCount)
	}

	return meanEmbedding, nil
}

func (service *EmbeddingServiceONNX) GetMovieEmbeddings(texts []string) ([]float32, error) {
	if len(texts) != service.weightcount {
		return nil, fmt.Errorf("texts is of improper batch count")
	}
	err := ort.InitializeEnvironment()
	if err != nil {
		log.Fatalf("Failed to initialize ORT: %v ", err)
	}

	//tokenize for each batch
	// need chunk offsets because model will ingest the text all as one argument to save resources
	tokensbatch := [][]int{}
	chunkoffsets := []int{} //for later
	totalchunklen := int(0)
	for i := range texts {
		tokens, err := service.Tokenize(texts[i])
		if err != nil {
			err = fmt.Errorf("error tokenizing :%v ", err)
			return nil, err
		}
		tokensbatch = append(tokensbatch, tokens)

		chunklen := int(math.Ceil(float64(len(tokensbatch[i])) / 128.0))
		chunkoffsets = append(chunkoffsets, chunklen)
		totalchunklen += chunklen
	}
	totallength := totalchunklen * 128
	// determine which parts to weight?
	// for now texts at index 0 Actors, 1 Genres, 2 Summary (incls title, year, dir, etc)

	//like calloc, the rest is already zero
	//cant move this to the previous loop without making ids variable-length. might save on memory
	ids := make([]int64, totallength) //literally all padding
	start := 0
	for i, offset := range chunkoffsets {
		//tokensbatch doesnt have pad here
		for j := range tokensbatch[i] { // j counts index of tokens in inner list
			ids[start*128+j] = int64(tokensbatch[i][j])
		}
		start += offset //move up the belt
	}
	/* scrutinize this code
	idx := 0
	for _, toks := range tokensbatch {
	    for _, t := range toks {
	        ids[idx] = int64(t)
	        idx++
	    }
	    // pad this chunk to 128 tokens
	    pad := 128 - (len(toks) % 128)
	    if pad < 128 {
	        idx += pad
	    }
	}

	*/

	mask := make([]int64, len(ids))

	// embed after fixing tokensbatch
	data, err := service.Embed(totalchunklen, ids, mask)
	if err != nil {
		return nil, err
	}

	//mean pooling here
	meanEmbedding := make([]float32, 384)
	validCount := 0 //number of valid tokens

	// change from iterating over mask for shape safety
	for tokenIdx := 0; tokenIdx < totalchunklen*128; tokenIdx++ {
		if mask[tokenIdx] == 1 {
			base := tokenIdx * 384
			for j := 0; j < 384; j++ {
				meanEmbedding[j] += data[base+j]
			}
			validCount++
		}
	}

	for j := range meanEmbedding {
		meanEmbedding[j] /= float32(validCount)
	}

	return meanEmbedding, nil
}

func MovieDatatoString(m *moviepb.MovieInfo) string {
	if m == nil {
		return ""
	}
	return fmt.Sprintf("%v %d dir. %v \n\nGenres: %v \n\nActors: %v \n\nPlot: %v", m.Title, m.Year, m.Director,
		strings.Join(m.Genre, ","), strings.Join(m.Actors, ","), m.Summary)
}

func MovieDataByWeights3(m *moviepb.MovieInfo) []string {
	if m == nil {
		return []string{"", "", ""}
	}
	actors := fmt.Sprintf("Actors: %v\n", strings.Join(m.Actors, ","))
	genres := fmt.Sprintf("Genres: %v\n", strings.Join(m.Genre, ","))
	summary := fmt.Sprintf("%v %d dir. %v \n\nPlot: %v", m.Title, m.Year, m.Director, m.Summary)
	return []string{actors, genres, summary}
}

func (service *EmbeddingServiceONNX) ProcessMovieData(ctx context.Context, movie *moviepb.MovieInfo, collection string) ([]*qdrant.ScoredPoint, error) {
	//movietext := MovieDatatoString(movie)
	v, err := service.GetMovieEmbeddings(MovieDataByWeights3(movie))
	if err != nil {
		return nil, err
	}

	res, err := service.repo.UpsertPoints(ctx, collection, v, repository.MovieToPayload(movie))
	if err != nil {
		return nil, fmt.Errorf("upsert failed: %v", err.Error())
	}
	fmt.Println(res.Status.Descriptor(), "op ID", res.OperationId)
	if res.Status == qdrant.UpdateStatus_ClockRejected {
		return nil, fmt.Errorf("upsert rejected")
	}

	points, err := service.repo.SearchMovie(v, []string{})
	if err != nil {
		return nil, err
	}
	return points, nil
}

func (service *EmbeddingServiceONNX) SearchMovieFromText(ctx context.Context, movie *moviepb.MovieInfo, collection string) ([]*qdrant.ScoredPoint, error) {
	//movietext := MovieDatatoString(movie)
	v, err := service.GetMovieEmbeddings(MovieDataByWeights3(movie))
	if err != nil {
		return nil, err
	}

	res, err := service.repo.UpsertPoints(ctx, collection, v, repository.MovieToPayload(movie))
	if err != nil {
		return nil, fmt.Errorf("upsert failed: %v", err.Error())
	}
	fmt.Println(res.Status.Descriptor(), "op ID", res.OperationId)
	if res.Status == qdrant.UpdateStatus_ClockRejected {
		return nil, fmt.Errorf("upsert rejected")
	}

	points, err := service.repo.SearchMovie(v, []string{})
	if err != nil {
		return nil, err
	}
	return points, nil
}
