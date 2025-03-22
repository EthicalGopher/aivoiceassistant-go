package main

import (
	"fmt"
	"os"
	

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	// "github.com/google/uuid"
)
var Path string
func TTS(input, dir,path string) (string, error) {
	// Generate unique filename
	uuid := path
	filename := uuid
	fullPath:=dir+`\`+filename+`.mp3`
	
	os.Truncate(Path,0)
	
	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("directory creation failed: %w", err)
	}
	file,_:=os.Create(path)
	
	// Generate speech file
	speech := htgotts.Speech{
		Folder:   dir,
		Language: voices.English,
		Handler:  &handlers.Native{},
	}

	if _, err := speech.CreateSpeechFile(input, filename); err != nil {
		return "", fmt.Errorf("speech generation failed: %w", err)
	}
	Path=fullPath
	file.Close()
	return fullPath, nil
}



