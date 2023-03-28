package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

const code = "#include <iostream>\n#include <cstdlib>\n#include <cmath>\nint main()\n{\n  using namespace std;\n  double a = 0, b = 0;\n  cin >> a;\n cin >> b;\n  cout.precision(16); \n  cout << \"a to b power  = \" << pow(a, b) << endl;\n  return 0;\n}"
const lang = "c++"
const stdin = "2\n3\n"
const binaryName = "program"

type Container struct {
	*client.Client
}

func NewContainer() *Container {
	cli, _ := client.NewClientWithOpts(client.WithVersion("1.41"))
	return &Container{
		Client: cli,
	}
}

func main() {

	cli, err := client.NewClientWithOpts(client.WithVersion("1.41"))
	if err != nil {
		// обрабатываем ошибку создания клиента
		fmt.Println(11)
	}

	codePath, err := CreateTempCFile(code)
	if err != nil {
		fmt.Println(12)
	}
	defer os.Remove(codePath.Name())

	var image, compileCommand string
	switch lang {
	case "c++":
		image = "gcc:latest"
		compileCommand = fmt.Sprintf("g++ -o %s %s", binaryName, path.Base(codePath.Name()))
	}

	AbsPath, _ := filepath.Abs(filepath.Dir(codePath.Name()))

	// создаем контейнер Docker
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: image,
		Tty:   true,
	}, &container.HostConfig{
		Privileged: true,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: AbsPath,
				Target: "/app",
			},
		},
	}, nil, nil, "")
	if err != nil {
		// обрабатываем ошибку создания контейнера
		fmt.Println(10)
		fmt.Println(err)
	}

	err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		// обрабатываем ошибку запуска контейнера
		fmt.Println(14)
	}

	// компилируем программный код
	execResp, err := cli.ContainerExecCreate(context.Background(), resp.ID, types.ExecConfig{
		Cmd: []string{"bash", "-c", fmt.Sprintf("cd /app && %s", compileCommand)},
	})
	if err != nil {
		// обрабатываем ошибку создания команды компиляции
		fmt.Println(err)
		fmt.Println(9)
	}

	err = cli.ContainerExecStart(context.Background(), execResp.ID, types.ExecStartCheck{})
	if err != nil {
		// обрабатываем ошибку выполнения команды компиляции
		fmt.Println(8)
	}

	// ждем завершения выполнения команды компиляции
	inspect, err := cli.ContainerExecInspect(context.Background(), execResp.ID)
	if err != nil {
		fmt.Println(err)
		fmt.Println(15)
	}
	for !inspect.Running {
		inspect, err = cli.ContainerExecInspect(context.Background(), execResp.ID)
		if err != nil {
			fmt.Println(err)
			fmt.Println(15)
		}
	}
	time.Sleep(5 * time.Second)
	// запускаем контейнер Docker и передаем входные данные
	execResp, err = cli.ContainerExecCreate(context.Background(), resp.ID, types.ExecConfig{
		Cmd:          []string{"bash", "-c", fmt.Sprintf("cd /app && %s", "./program")},
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		// обрабатываем ошибку создания команды выполнения
		fmt.Println(7)
	}

	respp, err := cli.ContainerExecAttach(context.Background(), execResp.ID, types.ExecStartCheck{})
	if err != nil {
		// обрабатываем ошибку выполнения команды выполнения
		fmt.Println(2)
	}
	defer respp.Close()

	_, err = respp.Conn.Write([]byte(stdin))
	if err != nil {
		// обрабатываем ошибку записи входных данных в контейнер
		fmt.Println(4)
	}

	out, err := io.ReadAll(respp.Reader)
	if err != nil {
		// обрабатываем ошибку чтения выходных данных из контейнера
		fmt.Println(6)
	}

	fmt.Println(string(out))
}

func CreateContainer() cli {
	// создаем контейнер Docker
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: image,
		Tty:   true,
	}, &container.HostConfig{
		Privileged: true,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: AbsPath,
				Target: "/app",
			},
		},
	}, nil, nil, "")
	if err != nil {
		// обрабатываем ошибку создания контейнера
		fmt.Println(10)
		fmt.Println(err)
	}

	err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		// обрабатываем ошибку запуска контейнера
		fmt.Println(14)
	}
}

func CreateTempCFile(code string) (*os.File, error) {
	file, err := os.CreateTemp("./Temp", "code.*.cpp")
	if err != nil {
		return nil, err
	}

	_, err = file.WriteString(code)
	if err != nil {
		return nil, err
	}

	err = file.Chmod(0755)
	if err != nil {
		panic(err)
	}

	err = file.Sync()
	if err != nil {
		return nil, err
	}

	return file, nil
}
