package handler

import (
	"Builder/lib"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type Program struct {
	Code  string `json:"code"`
	Lang  string `json:"lang"`
	Input string `json:"input"`
}

const binaryName = "program"

func (h *Handler) Compile(c *gin.Context) {

	var program Program
	err := c.BindJSON(&program)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error convert from JSON")
		return
	}

	var image, compileCommand, AbsPath string
	codePath, err := lib.CreateTempCFile(program.Code)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while create a TempFile")
		return
	}
	defer os.Remove(codePath.Name())

	switch program.Lang {
	case "c++":
		image = "gcc:latest"
		compileCommand = fmt.Sprintf("g++ -o %s %s", binaryName, path.Base(codePath.Name()))
	}

	AbsPath, _ = filepath.Abs(filepath.Dir(codePath.Name()))

	cliContainer, err := h.services.CreateContainer(image, AbsPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while create a container")
	}
	_, err = h.services.CompileProgram(cliContainer, compileCommand)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while compile program")
	}
	out, err := h.services.RunProgram(cliContainer, program.Input)
	h.services.RemoveContainer(cliContainer)

	c.String(http.StatusOK, string(out))
}
