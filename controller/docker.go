package controller

import (
	"tx-interview/structs"
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/satori/go.uuid"
	log "tx-interview/common/formatlog"
    "sync"
    "strconv"
    "fmt"
)

var DockerController *DockerCtrl

type DockerCtrl struct {
	cli *client.Client
    Lock *sync.RWMutex
    PortUsed map[string]string
}

func init() {
	dockercli, err := client.NewClient("tcp://192.168.159.133:1234", "v1.26", nil, nil)
	if err != nil {
		log.Errorf("[docker]获取docker cli错误, %s", err.Error())
	}
    portUsed := make(map[string]string)
    for i := 10; i< 30; i++ {
        tmpPort := "330" + strconv.Itoa(i)
        portUsed[tmpPort] = "free"
    }
    for j := 10; j< 30; j++ {
        tmpPort := "62" + strconv.Itoa(j)
        portUsed[tmpPort] = "free"
    }    
	DockerController = &DockerCtrl{
		cli: dockercli,
        Lock: new(sync.RWMutex),
        PortUsed: portUsed,
	}
}

// 拉取image清单
func (d DockerCtrl) ListImage() ([]string, error) {
	result := []string{}
	images, err := d.cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Errorf("[docker]获取docker image 错误, %s", err.Error())
		return []string{}, err
	} else {
		for _, image := range images {
			for _, repoTag := range image.RepoTags {
				result = append(result, repoTag)
			}
		}
	}
	return result, nil
}

// 创建容器
func (d DockerCtrl) CreateContainer(imageName string) (*structs.ContainerInfo, error) {
	// containter config
	randomId := uuid.NewV4().String()[:9]
	containerName := imageName + "-" + randomId

	env := []string{}
	var containerPort string
	var hostPort string



	if imageName == "mysql" {
		envPasswd := "MYSQL_ROOT_PASSWORD=" + randomId
		env = append(env, envPasswd)
		containerPort = "3306"
        for k ,v := range d.PortUsed {
            if v == "free" {
                hostPort = k
                break
            }
        }
        if len(hostPort) == 0 {
            return nil, fmt.Errorf("端口已被用完")
        }

	}
	exports := make(nat.PortSet, 10)
	port, _ := nat.NewPort("tcp", containerPort)
	exports[port] = struct{}{}

	container_config := &container.Config{
		Image:        imageName,
		ExposedPorts: exports,
		Env:          env,
	}

	// host config
	ports := make(nat.PortMap)
	pb := make([]nat.PortBinding, 0)
	pb = append(pb, nat.PortBinding{
		HostPort: hostPort,
	})
	ports[port] = pb
	host_config := &container.HostConfig{
		PortBindings: ports,
	}

	resp, err := d.cli.ContainerCreate(context.Background(), container_config, host_config, nil, nil, containerName)
	if err != nil {
		log.Errorf("创建container 失败，%s", err.Error())
		return nil, err
	}

	containerID := resp.ID

	err = d.cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	result := &structs.ContainerInfo{}
	result.ContainerID = containerID
	result.UserName = "root"
	result.PassWord = randomId
	result.Port = hostPort
    d.Lock.Lock()
    d.PortUsed[hostPort] = containerID
    d.Lock.Unlock()
	return result, nil

}
