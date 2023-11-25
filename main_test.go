package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestIntegration(t *testing.T) {
	// 创建临时目录作为测试输入和输出目录
	tempDirRoot := os.TempDir()
	tempDir, err := os.MkdirTemp(tempDirRoot, "hashcard")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建测试输入文件
	inputFile := filepath.Join(tempDir, "test.md")
	inputContent := `Hello
Front of Card 1
#card <!--2023/11/25/card1-->
Back of Card 1
---
Front of Card 2
#card <!--2023/11/25/card2-->
Back of Card 2`
	err = os.WriteFile(inputFile, []byte(inputContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// 运行程序
	dir := tempDir
	strategy := "md-json"
	outDir := filepath.Join(tempDir, "out")
	os.Args = []string{"cmd", "-dir", dir, "-strategy", strategy, "-out-dir", outDir}
	main()

	// 检查生成的输出文件是否符合预期
	expectedOutput1 := Card{
		ID:    "2023/11/25/card1",
		Front: "Front of Card 1",
		Back:  "Back of Card 1\n",
	}
	checkOutputFile(t, outDir, "2023/11/25/card1.md", expectedOutput1)

	expectedOutput2 := Card{
		ID:    "2023/11/25/card2",
		Front: "Front of Card 2",
		Back:  "Back of Card 2\n",
	}
	checkOutputFile(t, outDir, "2023/11/25/card2.md", expectedOutput2)
}

func checkOutputFile(t *testing.T, outDir, filename string, expected Card) {
	outputFile := filepath.Join(outDir, filename)
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %s does not exist", outputFile)
		return
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Errorf("Failed to read output file: %v", err)
		return
	}

	var card Card
	err = json.Unmarshal(content, &card)
	if err != nil {
		t.Errorf("Failed to unmarshal output file: %v", err)
		return
	}

	if card.ID != expected.ID {
		t.Errorf("Output file %s has incorrect ID. Expected: %s, Got: %s", outputFile, expected.ID, card.ID)
	}
	if card.Front != expected.Front {
		t.Errorf("Output file %s has incorrect front content. Expected: %s, Got: %s", outputFile, expected.Front, card.Front)
	}
	if card.Back != expected.Back {
		t.Errorf("Output file %s has incorrect back content. Expected: %s, Got: %s", outputFile, expected.Back, card.Back)
	}
}
