package handler

import (
	"Builder/lib"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Program struct {
	Code  string `json:"code"`
	Lang  string `json:"lang"`
	Input string `json:"input"`
}

func (h *Handler) Compile(c *gin.Context) {

	var program Program
	err := c.BindJSON(&program)
	if err != nil {
		log.Println("Error convert from JSON")
		c.String(http.StatusInternalServerError, "error 500 Server")
		return
	}

	var image, compileCommand, AbsPath string
	codePath, err := lib.CreateTempCFile(program.Code)
	if err != nil {
		log.Println("Error while create a TempFile")
		c.String(http.StatusInternalServerError, "error 500 Server")
		return
	}
	defer os.Remove(codePath.Name())

	binaryName := path.Base(codePath.Name())[:len(path.Base(codePath.Name()))-4]
	defer os.Remove(path.Join(path.Dir(codePath.Name()), binaryName))

	switch program.Lang {
	case "c++":
		image = "gcc:latest"
		compileCommand = fmt.Sprintf("g++ -o %s %s", binaryName, path.Base(codePath.Name()))
	}

	AbsPath, _ = filepath.Abs(filepath.Dir(codePath.Name()))

	cliContainer, err := h.services.CreateContainer(image, AbsPath)
	if err != nil {
		log.Println("Error while create a container")
		c.String(http.StatusInternalServerError, "error: 500 Server")
		return
	}
	_, err = h.services.CompileProgram(cliContainer, compileCommand)
	if err != nil {
		log.Println("Error while compile program")
		c.String(http.StatusInternalServerError, "error: 500 Server")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	done := make(chan bool)
	go func() {
		out, err := h.services.RunProgram(cliContainer, program.Input, binaryName)
		if err != nil {
			log.Println("Error while Run program")
			c.String(http.StatusInternalServerError, "error: while run program")
			return
		}
		h.services.RemoveContainer(cliContainer)
		c.String(http.StatusOK, out)
		done <- true
	}()
	select {
	case <-ctx.Done():
		h.services.RemoveContainer(cliContainer)
		log.Println("Execution time limit exceeded")
		c.String(http.StatusInternalServerError, "error: Execution time limit exceeded")
	case <-done:
	}

}
