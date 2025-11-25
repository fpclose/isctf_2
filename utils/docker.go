package utils

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var dockerClient *client.Client

// InitDocker 初始化 Docker 客户端
func InitDocker() error {
	var err error
	// 使用环境变量初始化 (DOCKER_HOST, DOCKER_API_VERSION 等)
	dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return err
}

// EnsureDockerClient 确保 Docker 客户端已初始化
func EnsureDockerClient() error {
	if dockerClient == nil {
		return InitDocker()
	}
	return nil
}

// StartContainer 启动容器
// image: 镜像名
// ports: 容器内部端口 -> 协议 (例如 {"80": "tcp"})
// env: 环境变量列表 (例如 ["FLAG=ctf{...}"])
// 返回: containerID, hostMapping(容器端口->宿主机端口), error
func StartContainer(image string, ports map[string]string, env []string) (string, map[string]string, error) {
	if err := EnsureDockerClient(); err != nil {
		return "", nil, fmt.Errorf("docker client init failed: %v", err)
	}

	ctx := context.Background()

	// 1. 尝试拉取镜像（如果本地不存在）
	// 注意：生产环境建议预先拉取或配置私有仓库认证
	reader, err := dockerClient.ImagePull(ctx, image, types.ImagePullOptions{})
	if err == nil {
		io.Copy(io.Discard, reader) // 读取输出以完成拉取
		reader.Close()
	} else {
		// 如果拉取失败，尝试直接使用本地镜像
		fmt.Printf("Image pull failed (using local if available): %v\n", err)
	}

	// 2. 配置端口映射
	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}

	for port, proto := range ports {
		p := nat.Port(fmt.Sprintf("%s/%s", port, proto))
		exposedPorts[p] = struct{}{}
		// 绑定到所有接口的随机端口
		portBindings[p] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: "", // 让 Docker 随机分配宿主机端口
			},
		}
	}

	// 3. 创建容器
	resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
		Image:        image,
		Env:          env,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		PortBindings: portBindings,
		// 可以在这里添加资源限制，例如 Memory: 512 * 1024 * 1024
		AutoRemove: false, // 不自动删除，以便排查问题
	}, nil, nil, "")
	if err != nil {
		return "", nil, err
	}

	// 4. 启动容器
	if err := dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		// 启动失败尝试清理
		_ = dockerClient.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		return "", nil, err
	}

	// 5. 获取分配的宿主机端口
	inspect, err := dockerClient.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return resp.ID, nil, err
	}

	hostMapping := make(map[string]string)
	for p, bindings := range inspect.NetworkSettings.Ports {
		if len(bindings) > 0 {
			// 格式: 80/tcp -> 32768
			hostMapping[string(p)] = bindings[0].HostPort
		}
	}

	return resp.ID, hostMapping, nil
}

// StopContainer 停止容器
func StopContainer(containerID string) error {
	if err := EnsureDockerClient(); err != nil {
		return err
	}
	ctx := context.Background()
	timeout := 5 // 5秒超时
	return dockerClient.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
}

// RemoveContainer 强制删除容器
func RemoveContainer(containerID string) error {
	if err := EnsureDockerClient(); err != nil {
		return err
	}
	ctx := context.Background()
	return dockerClient.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
}

// GetContainerStatus 获取容器运行状态
func GetDockerContainerStatus(containerID string) (string, error) {
	if err := EnsureDockerClient(); err != nil {
		return "unknown", err
	}
	ctx := context.Background()
	inspect, err := dockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", err
	}
	return inspect.State.Status, nil // running, exited, dead, etc.
}

// GenerateDynamicFlag 生成动态 Flag
func GenerateDynamicFlag(teamID, challengeID int64) string {
	randPart := fmt.Sprintf("%08x", rand.Int63())
	return fmt.Sprintf("ISCTF{team_%d_%d_%s}", teamID, challengeID, randPart)
}
