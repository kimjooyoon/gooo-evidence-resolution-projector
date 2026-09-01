package gooo

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kimjooyoon/gooo-evidence-resolution-projector/internal/model"
)

func Load(path string) (model.SourceDocument, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return model.SourceDocument{}, nil, fmt.Errorf("read Gooo source: %w", err)
	}
	var document model.SourceDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return model.SourceDocument{}, nil, fmt.Errorf("parse Gooo source: %w", err)
	}
	return document, data, nil
}
