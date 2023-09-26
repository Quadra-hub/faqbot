package gpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"github.com/Quadra-hub/go-chatgpt/config"
	"github.com/gofrs/uuid"
	"github.com/pgvector/pgvector-go"
)

type apiRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

func FetchEmbeddings(input []string) (pgvector.Vector, error) {
	apiKey := config.GptApiKey()
	if apiKey == "" {
		return pgvector.Vector{}, fmt.Errorf("GPT_API_KEY is not set")
	}

	url := "https://api.openai.com/v1/embeddings"
	data := &apiRequest{
		Input: input,
		Model: "text-embedding-ada-002",
	}

	b, err := json.Marshal(data)
	if err != nil {
		return pgvector.Vector{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return pgvector.Vector{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return pgvector.Vector{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return pgvector.Vector{}, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return pgvector.Vector{}, err
	}

	var embeddings []float32
	for _, v := range result["embeddings"].([]interface{}) {
		embeddings = append(embeddings, float32(v.(float64)))
	}
	return pgvector.NewVector(embeddings), nil

}

type Embeddding struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v7()"`
	Context   []string
	Embedding pgvector.Vector `pg:"type:vector(1536)"`
}

type EmbedddingRepo interface {
	InsertEmbedding(*Embeddding) error
	RenewEmbedding(*Embeddding) error
}

type EmbedddingService struct {
	Db *gorm.DB
}

func NewEmbedddingService(db *gorm.DB) *EmbedddingService {
	return &EmbedddingService{Db: db}
}

func (s *EmbedddingService) InsertEmbedding(e *Embeddding) error {
	return s.Db.Save(e).Error
}

func (s *EmbedddingService) RenewEmbedding(e *Embeddding) error {
	return s.Db.Model(e).Update("embedding", e.Embedding).Error
}
