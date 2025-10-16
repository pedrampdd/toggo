package loader

import (
	"strings"
	"testing"

	"github.com/pedram/toggo"
)

func TestJSONLoader_LoadFromFile(t *testing.T) {
	loader := NewJSONFile("../testdata/flags.json")

	flags, err := loader.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(flags) == 0 {
		t.Error("expected flags to be loaded")
	}

	// Find new_checkout flag
	var checkoutFlag *toggo.Flag
	for _, flag := range flags {
		if flag.Name == "new_checkout" {
			checkoutFlag = flag
			break
		}
	}

	if checkoutFlag == nil {
		t.Fatal("expected to find new_checkout flag")
	}

	if !checkoutFlag.Enabled {
		t.Error("expected new_checkout to be enabled")
	}

	if checkoutFlag.Rollout != 50 {
		t.Errorf("expected rollout 50, got %d", checkoutFlag.Rollout)
	}
}

func TestJSONLoader_LoadFromReader(t *testing.T) {
	jsonData := `{
		"flags": [
			{
				"name": "test_flag",
				"enabled": true,
				"rollout": 100
			}
		]
	}`

	reader := strings.NewReader(jsonData)
	loader := NewJSONReader(reader)

	flags, err := loader.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(flags) != 1 {
		t.Errorf("expected 1 flag, got %d", len(flags))
	}

	if flags[0].Name != "test_flag" {
		t.Errorf("expected test_flag, got %s", flags[0].Name)
	}
}

func TestYAMLLoader_LoadFromFile(t *testing.T) {
	loader := NewYAMLFile("../testdata/flags.yaml")

	flags, err := loader.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(flags) == 0 {
		t.Error("expected flags to be loaded")
	}
}

func TestYAMLLoader_LoadFromReader(t *testing.T) {
	yamlData := `
flags:
  - name: test_flag
    enabled: true
    rollout: 100
`

	reader := strings.NewReader(yamlData)
	loader := NewYAMLReader(reader)

	flags, err := loader.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(flags) != 1 {
		t.Errorf("expected 1 flag, got %d", len(flags))
	}

	if flags[0].Name != "test_flag" {
		t.Errorf("expected test_flag, got %s", flags[0].Name)
	}
}

func TestLoader_LoadIntoStore(t *testing.T) {
	jsonData := `{
		"flags": [
			{
				"name": "flag1",
				"enabled": true,
				"rollout": 100
			},
			{
				"name": "flag2",
				"enabled": false,
				"rollout": 50
			}
		]
	}`

	reader := strings.NewReader(jsonData)
	loader := NewJSONReader(reader)

	store := toggo.NewStore()
	err := loader.LoadIntoStore(store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if store.Size() != 2 {
		t.Errorf("expected 2 flags in store, got %d", store.Size())
	}

	flag1, _ := store.GetFlag("flag1")
	if !flag1.Enabled {
		t.Error("expected flag1 to be enabled")
	}

	flag2, _ := store.GetFlag("flag2")
	if flag2.Enabled {
		t.Error("expected flag2 to be disabled")
	}
}

func TestLoader_InvalidJSON(t *testing.T) {
	invalidJSON := `{ invalid json `

	reader := strings.NewReader(invalidJSON)
	loader := NewJSONReader(reader)

	_, err := loader.Load()
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoader_InvalidFlag(t *testing.T) {
	// Rollout > 100 should fail validation
	jsonData := `{
		"flags": [
			{
				"name": "bad_flag",
				"enabled": true,
				"rollout": 150
			}
		]
	}`

	reader := strings.NewReader(jsonData)
	loader := NewJSONReader(reader)

	_, err := loader.Load()
	if err == nil {
		t.Error("expected error for invalid flag")
	}
}
