package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"io"
	"log"
	"time"
)

type Client struct {
	*client.Client
}

func NewClient() *Client {
	cli, _ := client.NewClientWithOpts(client.WithVersion("1.41"))
	return &Client{
		cli,
	}
}

func (cli *Client) RunProgram(resp container.CreateResponse, input string, binaryName string) (string, error) {
	execResp, err := cli.ContainerExecCreate(context.Background(), resp.ID, types.ExecConfig{
		Cmd:          []string{"bash", "-c", fmt.Sprintf("cd /app && ./%s", binaryName)},
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		log.Println("Error while run program")
		return "", err
	}

	respp, err := cli.ContainerExecAttach(context.Background(), execResp.ID, types.ExecStartCheck{})
	if err != nil {
		log.Println("Error while Attach container")
		return "", err
	}
	defer respp.Close()
	_, err = respp.Conn.Write([]byte(input))
	if err != nil {
		log.Println("Error while write data in programm")
		return "", err
	}
	out, err := io.ReadAll(respp.Reader)
	if err != nil {
		log.Println("Error while read output")
		return "", err
	}
	return string(out), nil
}

func (cli *Client) CompileProgram(resp container.CreateResponse, compileCommand string) (types.IDResponse, error) {
	execResp, err := cli.ContainerExecCreate(context.Background(), resp.ID, types.ExecConfig{
		Cmd: []string{"bash", "-c", fmt.Sprintf("cd /app && %s", compileCommand)},
	})
	if err != nil {
		if err != nil {
			log.Println("Error while Exec to container")
			return execResp, err
		}
	}

	err = cli.ContainerExecStart(context.Background(), execResp.ID, types.ExecStartCheck{})
	if err != nil {
		if err != nil {
			log.Println("Error while Start Compile")
			return execResp, err
		}
	}

	time.Sleep(2 * time.Second)
	return execResp, nil
}

func (cli *Client) CreateContainer(image string, AbsPath string) (container.CreateResponse, error) {
	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: image,
		Tty:   true,
	}, &container.HostConfig{
		Resources: container.Resources{
			Memory:   1024 * 1024 * 1024, // 1 GB
			NanoCPUs: 1000000000,         // 1 CPU
		},
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
		log.Println("Error while create a container")
		return resp, err
	}

	err = cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Println("Error while Start a container")
		return resp, err
	}

	return resp, nil
}

func (cli *Client) RemoveContainer(resp container.CreateResponse) {
	err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		log.Println("Error while remove container")
	}

}
