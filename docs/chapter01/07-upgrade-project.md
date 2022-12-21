# 升级项目

之前的项目中只有 prod 服务，具有版本的区分。 现在对项目进行一些升级， 模拟一个多服务的项目。

1. 两个服务， `review / prod`
2. 服务之前还有调用关系。 `prod -> review`

## review

这次新加入了 `review` 评论服务。 


```json5
{
  "1": {
    "id": "1",
    "name": "zhangsan",
    "commment": "istio 功能很强大， 就是配置太麻烦"
  },
  "2": {
    "id": "1",
    "name": "wangwu",
    "commment": "《istio in action》 真是一本了不起的书"
  }
}
```

## prod

升级 prod 服务， 除了之前返回本身的数据信息之外，还需要返回关联的评论信息。

```go
type Product struct {
	Name    string
	Price   int
	Reviews interface{}  // 评论信息
}
```

这部分评论信息的来源就是上面新添加的评论服务。

```go
func getReivews() (map[string]model.Review, error) {

	reviews := make(map[string]model.Review)

	resp, err := http.Get("http://svc-review/review/all")
	if err != nil {
		return nil, fmt.Errorf("reqeust svc-review failed: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %v", err)
	}

	err = json.Unmarshal(data, &reviews)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal data failed: %v", err)
	}

	return reviews, nil
}
```

完整结果如下

```json5
{
  "data": {
    "Name": "istio in action",
    "Price": 300,
    "Reviews": {
      "1": {
        "id": "1",
        "name": "zhangsan",
        "commment": "istio 功能很强大， 就是配置太麻烦"
      },
      "2": {
        "id": "1",
        "name": "wangwu",
        "commment": "《istio in action》 真是一本了不起的书"
      }
    }
  },
  "version": "v1.1.0"
}
```

