package embeddingservice

import (
	"fmt"
	"log"

	pretrained "github.com/sugarme/tokenizer/pretrained"
	ort "github.com/yalue/onnxruntime_go"
)

// const modelPath = `../pyenv/paraphrase-multilingual-mpnet-base-v2.onnx`
const modelPath = `..\models\all-MiniLM-L6-v2\all-MiniLM-L6-v2.onnx`
const tokenizerpath = `..\models\all-MiniLM-L6-v2\tokenizer.json`

func CheckExecution() (string, error) {
	//return "OPENVINO", nil
	return "", nil
}

func Tokenize(text string) ([]int, error) {
	//TODO: real implementation
	tk, err := pretrained.FromFile(tokenizerpath)
	if err != nil {
		return nil, fmt.Errorf("%v: for path :%v", err, tokenizerpath)
	}

	// // Set padding manually
	// padding := tokenizers.PaddingParams{}
	// padding.Strategy = *tokenizers.NewPaddingStrategy(tokenizers.WithFixed(128)) // or tok.PaddingStrategyBatchLongest
	// padding.PadId = 0
	// padding.PadToken = "[PAD]"
	// tk.WithPadding(&padding)

	//defer tk.Close()
	ids, _ := tk.EncodeSingle(text, false)

	return ids.GetIds(), err
	//return []float32{0.75, 0.55, -0.83, -0.76}, nil
}

func Onnxservice(text string) ([]float32, error) {
	err := ort.InitializeEnvironment()
	if err != nil {
		log.Fatalf("Failed to initialize ORT: %v ", err)
	}

	tokens, err := Tokenize(text)
	if err != nil {
		err = fmt.Errorf("error tokenizing :%v ", err)
		return nil, err
	}
	ids := make([]int64, 128)
	for i, t := range tokens {
		ids[i] = int64(t)
	}

	for len(ids) < 128 {
		ids = append(ids, 0)
	}

	executor, err := CheckExecution()
	if err != nil {
		err = fmt.Errorf("error checking executor :%v ", err)
		return nil, err
	}

	inputShape := ort.NewShape(1, 128)
	//inputTensor, err := ort.NewEmptyTensor[float32](inputShape)
	inputTensor, err := ort.NewTensor[int64](inputShape, ids)
	if err != nil {
		err = fmt.Errorf("Error creating input tensor: %w ", err)
		return nil, err
	}
	mask := make([]int64, len(ids))
	for i := range ids {
		if ids[i] != 0 { // 0 is pad token id
			mask[i] = 1
		}
	}
	maskShape := ort.NewShape(1, int64(len(ids)))
	maskTensor, err := ort.NewTensor[int64](maskShape, mask)
	if err != nil {
		log.Fatalf("Error creating attention_mask tensor: %v", err)
	}
	tokenTypeIds := make([]int64, 128) // all zeros (single sentence)
	tokenTypeIdsTensor, err := ort.NewTensor[int64](maskShape, tokenTypeIds)
	if err != nil {
		log.Fatalf("Error creating tokenTypeIds tensor: %v", err)
	}

	// todo: add chunking/truncation

	outputShape := ort.NewShape(1, 128, 384)
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

	session, err := ort.NewAdvancedSession(modelPath,
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

	// outputTensor.GetData() returns []float32 of shape [1,128,384]
	data := outputTensor.GetData()

	meanEmbedding := make([]float32, 384)
	validCount := 0

	for i := range 128 {
		if mask[i] == 1 {
			for j := 0; j < 384; j++ {
				meanEmbedding[j] += data[i*384+j]
			}
			validCount++
		}
	}

	for j := range meanEmbedding {
		meanEmbedding[j] /= float32(validCount)
	}

	return meanEmbedding, nil
	// file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Fatalf("failed to open file: %v", err)
	// }
	// // The defer statement ensures the file is closed at the end of the main function.
	// defer file.Close()

	// Iterate over the slice and write each float to a new line in the file.
	// for _, number := range data {
	// 	_, err := fmt.Fprintf(file, "%.5f\n", number)
	// 	if err != nil {
	// 		log.Fatalf("failed to write to file: %v", err)
	// 	}
	// }
}
