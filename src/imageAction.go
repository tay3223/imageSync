package src

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"strings"
)

// Pull 拉取docker镜像
func Pull(imageName string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()
	fmt.Println("\n\n\n拉取镜像...")
	FormatOut(out, "pull")
}

// Push 推送docker镜像
func Push(imageName string) {
	//实例化一个Docker认证客户端
	authConfig := types.AuthConfig{
		Username:      UserConfig.Username,
		Password:      UserConfig.Password,
		ServerAddress: UserConfig.ServerAddress,
	}
	encodedJSON, _ := json.Marshal(authConfig)                //结构体转换为json
	authStr := base64.URLEncoding.EncodeToString(encodedJSON) //进行Base64编码

	//实例化一个docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	//docker image tag
	imageNameTag := imageTagRename(imageName)
	if err := cli.ImageTag(context.Background(), imageName, imageNameTag); err != nil {
		return
	}

	//push image
	out, err := cli.ImagePush(context.TODO(), imageNameTag, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		return
	}
	fmt.Println("\n\n\n推送镜像...")
	FormatOut(out, "push")

	//判断imageTag是否是以阿里云公有仓库地址为开头
	if strings.HasPrefix(imageNameTag, "registry.cn-shanghai.aliyuncs.com") {
		urlList := strings.Split(imageNameTag, ".cn-shanghai.aliyuncs.com")
		urlList[0] = urlList[0] + "-vpc.cn-shanghai.aliyuncs.com"
		urlList2 := strings.Join(urlList, "")
		fmt.Printf("\n\n新的镜像地址为（公网）：" + imageNameTag + "\n")
		fmt.Printf("新的镜像地址为（VPC）：" + urlList2 + "\n\n\n")
	} else {
		fmt.Printf("\n\n新的镜像地址为（公网）：" + imageNameTag + "\n\n\n")
	}

}

func imageTagRename(imageName string) string {
	//切割字符串
	repoTagList := strings.Split(imageName, "/")

	//取切割后数组的最后一个
	repoTag := repoTagList[len(repoTagList)-1]

	if !strings.Contains(repoTag, ":") {
		repoTag += ":latest"
	}

	//repoTagList2 := strings.Split(repoTag, ":")
	//repoTagList3 := repoTagList2[1:]
	//var repoTag1 string
	//for _, v := range repoTagList3 {
	//	repoTag1 += string(v)
	//}
	//把:号替换为-号（n小于0时表示不限制替换的次数）
	//repoTag2 := strings.Replace(repoTag1, ":", "-", -1)

	//获取配置文件中的image_tag内容
	imageTag := UserConfig.ImageTag
	//把:号替换为空
	imageTag2 := strings.Replace(imageTag, ":", "", -1)

	//拼接新的imageTag内容
	//newImageTag := imageTag2 + ":" + repoTag2
	newImageTag := imageTag2 + "/" + repoTag

	return newImageTag
}
