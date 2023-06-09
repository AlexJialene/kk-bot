# kk-bot

是一个自用的微信机器人集成了chatGPT

## 背景

在chatGPT刚发布就入坑开荒了，chatGPT本身就不支持国内访问，所以一开始就尝试用GPT写代码（~~还挺爽~~）。直到后面突然的爆火、以及国内众小中大型公司纷纷跟进/接入后，开始各种封IP导致后续体验直线下滑。

其次由于工作关系，在访问openAI的时候体验也并不是很友好。故有此项目。



## 效果

<table rules="none" align="center">
	<tr>
		<td>
			<center>
				<img src="https://github.com/AlexJialene/kk-bot/blob/feature/message_handler_by_lamkeizyi/doc/chat.jpg?raw=true" width="60%" />
			</center>
		</td>
	</tr>
</table>




## 准备

### pandora

项目与GPT交互采用的是zhile大佬的**pandora** 作为代理（docker 安装）

#### 拉取

```bash
docker pull pengzhile/pandora
```

#### 启动

```bash
docker run -it --rm -p 19090:19090 -d \
 -e PANDORA_SERVER=0.0.0.0:19090 \
 -e PANDORA_VERBOSE \
 -e PANDORA_ACCESS_TOKEN=your gpt token \
 pengzhile/pandora 
```

> 注意：这里PANDORA_ACCESS_TOKEN是你的chatGPT accessToken。 可访问[ChatGPT Auth (fakeopen.com)](https://ai.fakeopen.com/auth)获取

## 开始

下载release压缩包

直接启动



`linux：`

```bash
nohup ./kkbot > ./kkbot.log & 
```



## 鸣谢

- [pengzhile/pandora](https://github.com/pengzhile/pandora)
- [openwechat](https://github.com/eatmoreapple/openwechat)
