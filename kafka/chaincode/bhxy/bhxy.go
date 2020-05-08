package main

import (
	"encoding/json"
	"fmt"
	//shim包作用 客户端需要与fabric框架通信
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//客户端与fabric请求的响应
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type BhxyCC struct{}

/*
	链码即对键值对的操作
	存入对象时需将对象序列化为[]byte 并设定key 存储数据库是一个kv数据库,所以存的时候带对象类型如 user_key是一个好选择
	取出对象时 同样的key取出来 并反序列化
	用于操作的json 第一个参数指定了要调用的方法名 后面的参数是一个[]string 切片
	接收到参数后先将参数判空 再带入查询数据存在性 最后再操作 以保证操作数据的安全
*/
//角色
type User struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Level    int      `json:"level"`
	LevelExp int      `json:"levelexp"`
	VipLevel int      `json:"viplevel"`
	VipExp   int      `json:"vipexp"`
	Coin     int      `json:"coin"`
	Crystal  int      `json:"crystal"`
	Assets   []string `json:"assets"`
}

//资产
type Asset struct {
	Id        string `json:"id"`
	UserID    int    `json:"user_id"`
	Level     int    `json:"level"`
	LevelExp  int    `json:"level_exp"`
	Attribute string `json:"attribute"`
	Sign      int    `json:"sign"`
	ArticleID int    `json:"article_id"`
}

//物品
type Article struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Properties string `json:"properties"`
}

type UserInfo struct {
	User   User        `json:"user"`
	Assets []AssetInfo `json:"assets"`
}
type AssetInfo struct {
	Asset   Asset   `json:"asset"`
	Article Article `json:"article"`
}

//将id转化一下, stub.GetState stub.PutState 两个方法都使用转化后的id
func constructUserKey(userId int) string {
	return fmt.Sprintf("user_%d", userId)
}

func constructAssetKey(assetId string) string {
	return fmt.Sprintf("asset_%d", assetId)
}

func constructArticleKey(articleID int) string {
	return fmt.Sprintf("article_%d", articleID)
}

// 用户开户
func userRegister(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 套路1：检查参数的个数
	if len(args) != 2 {
		return shim.Error("not enough args,need 2 args ")
	}

	// 套路2：验证参数的正确性
	name := args[1]
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	if name == "" || id == 0 {
		return shim.Error("invalid args")
	}

	// 套路3：验证数据是否存在 应该存在 or 不应该存在
	if userBytes, err := stub.GetState(constructUserKey(id)); err == nil && len(userBytes) != 0 {
		return shim.Error("userId already exist")
	}

	// 套路4：写入状态
	user := &User{
		Name:   name,
		Id:     id,
		Assets: make([]string, 0),
	}

	// 序列化对象 因为stub.PutState 接收的是[]byte
	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error %s", err))
	}

	if err := stub.PutState(constructUserKey(id), userBytes); err != nil {
		return shim.Error(fmt.Sprintf("put user error %s", err))
	}

	// 成功返回
	return shim.Success(nil)
}

// 用户销户
func userDestroy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 套路1：检查参数的个数
	if len(args) != 1 {
		return shim.Error("not enough args")
	}

	// 套路2：验证参数的正确性
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	if id == 0 {
		return shim.Error("invalid args")
	}

	// 套路3：验证数据是否存在 应该存在 or 不应该存在
	userBytes, err := stub.GetState(constructUserKey(id))
	if err != nil || len(userBytes) == 0 {
		return shim.Error("user not found")
	}

	// 套路4：写入状态
	if err := stub.DelState(constructUserKey(id)); err != nil {
		return shim.Error(fmt.Sprintf("delete user error: %s", err))
	}

	// 删除用户名下的资产
	user := new(User)
	if err := json.Unmarshal(userBytes, user); err != nil {
		return shim.Error(fmt.Sprintf("unmarshal user error: %s", err))
	}
	for _, assetid := range user.Assets {
		if err := stub.DelState(constructAssetKey(assetid)); err != nil {
			return shim.Error(fmt.Sprintf("delete asset error: %s", err))
		}
	}

	return shim.Success(nil)
}

// 用户查询
func queryUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 套路1：检查参数的个数
	if len(args) != 1 {
		return shim.Error("not enough args")
	}

	// 套路2：验证参数的正确性
	ownerId, err := strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	if ownerId == 0 {
		return shim.Error("invalid args")
	}

	// 套路3：验证数据是否存在 应该存在 or 不应该存在
	userBytes, err := stub.GetState(constructUserKey(ownerId))
	if err != nil || len(userBytes) == 0 {
		return shim.Error("user not found")
	}
	userInfo := new(UserInfo)
	user := new(User)
	if err := json.Unmarshal(userBytes, &user); err != nil {
		shim.Error("user Unmarshal fail")
	}
	//user放进去
	userInfo.User = *user
	//进入循环
	for _, v := range user.Assets {
		if assetBytes, err := stub.GetState(constructAssetKey(v)); err != nil || len(assetBytes) == 0 {
			shim.Error("one asset not find")
		} else {
			assetInfo := new(AssetInfo)
			asset := new(Asset)
			if err := json.Unmarshal(assetBytes, &asset); err != nil {
				shim.Error(" asset Unmarshal fail")
			}
			//asset放进assetInfo
			assetInfo.Asset = *asset
			article := new(Article)
			if articleBytes, err := stub.GetState(constructArticleKey(asset.ArticleID)); err != nil || len(articleBytes) == 0 {
				shim.Error("one asset not find")
			} else {
				if err := json.Unmarshal(articleBytes, &article); err != nil {
					shim.Error("article Unmarshal fail")
				} else {
					//article放进去
					assetInfo.Article = *article
				}
			}
			//assetInfo放进userInfo
			userInfo.Assets = append(userInfo.Assets, *assetInfo)
		}
	}
	userInfoBytes, err := json.Marshal(userInfo)
	if err != nil || len(userInfoBytes) == 0 {
		shim.Error("userInfo marshal fail")
	}
	return shim.Success(userInfoBytes)
}

// 资产登记
func assetEnroll(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 套路1：检查参数的个数
	if len(args) != 7 {
		return shim.Error("not enough args")
	}

	// 套路2：验证参数的正确性
	id := args[0]
	userID, _ := strconv.Atoi(args[1])
	level, _ := strconv.Atoi(args[2])
	levelExp, _ := strconv.Atoi(args[3])
	attribute := args[4]
	sign, _ := strconv.Atoi(args[5])
	articleID, _ := strconv.Atoi(args[6])
	if id == "" || userID == 0 || level == 0 || articleID == 0 {
		return shim.Error("invalid args")
	}

	// 套路3：验证数据是否存在 应该存在 or 不应该存在
	userBytes, err := stub.GetState(constructUserKey(userID))
	if err != nil || len(userBytes) == 0 {
		return shim.Error("user not found")
	}

	if assetBytes, err := stub.GetState(constructAssetKey(id)); err == nil && len(assetBytes) != 0 {
		return shim.Error("asset already exist")
	}

	if assetBytes, err := stub.GetState(constructArticleKey(articleID)); err != nil || len(assetBytes) == 0 {
		return shim.Error("unknown Article in args[6]")
	}
	// 套路4：写入状态
	// 1. 写入资产对象 2. 更新用户对象 3. 写入资产变更记录
	asset := &Asset{
		Id:        id,
		UserID:    userID,
		Level:     level,
		LevelExp:  levelExp,
		Attribute: attribute,
		Sign:      sign,
		ArticleID: articleID,
	}
	//序列化资产,存盘
	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal asset error: %s", err))
	}
	if err := stub.PutState(constructAssetKey(id), assetBytes); err != nil {
		return shim.Error(fmt.Sprintf("save asset error: %s", err))
	}

	user := new(User)
	// 反序列化user
	if err := json.Unmarshal(userBytes, user); err != nil {
		return shim.Error(fmt.Sprintf("unmarshal user error: %s", err))
	}
	user.Assets = append(user.Assets, id)
	// 序列化user
	userBytes, err = json.Marshal(user)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error: %s", err))
	}
	if err := stub.PutState(constructUserKey(user.Id), userBytes); err != nil {
		return shim.Error(fmt.Sprintf("update user error: %s", err))
	}
	return shim.Success(nil)
}

//查询资产
func queryAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//判断参数个数 -> 判空 -> 查询 ->返回结果
	if len(args) != 1 {
		return shim.Error("args number error")
	}

	id := args[0]
	if id == "" {
		return shim.Error("invalid args")
	}
	assetBytes, err := stub.GetState(constructAssetKey(id))
	if err != nil || len(assetBytes) == 0 {
		return shim.Error("asset not exist")
	}

	//TODO 查询物品详情

	return shim.Success(assetBytes)
}

//删除资产
func assetDestroy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//判断参数是否正确
	if len(args) != 1 {
		shim.Error("not enough args ")
	}
	assetID := args[0]
	if assetID == "" {
		shim.Error("invalid args !")
	}
	//判断资产是否存在
	if assetBytes, err := stub.GetState(assetID); err != nil || len(assetBytes) == 0 {
		shim.Error("asset not exist")
	}
	//删除操作
	if err := stub.DelState(constructAssetKey(assetID)); err != nil {
		shim.Error(fmt.Sprintf("del asset error: %s", err))
	}
	return shim.Success(nil)
}

//登记物品
func addArticle(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 套路1：检查参数的个数
	if len(args) != 3 {
		return shim.Error("not enough args")
	}

	id, _ := strconv.Atoi(args[0])
	name := args[1]
	properties := args[2]
	if id == 0 || name == "" || properties == "" {
		return shim.Error("invalid args")
	}

	//判断存在
	if articleByte, err := stub.GetState(constructArticleKey(id)); err == nil && len(articleByte) != 0 {
		return shim.Error("article is exist")
	}

	article := &Article{
		Id:         id,
		Name:       name,
		Properties: properties,
	}
	//序列化存盘
	articleByte, err := json.Marshal(article)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal article error: %s", err))
	}
	if err := stub.PutState(constructArticleKey(id), articleByte); err != nil {
		return shim.Error(fmt.Sprintf("save article err: %s", err))
	}

	return shim.Success(nil)
}

//查询物品
func queryArticle(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//判断参数个数 -> 判空 -> 查询 ->返回结果
	if len(args) != 1 {
		return shim.Error("args number error")
	}

	id, _ := strconv.Atoi(args[0])

	if id == 0 {
		return shim.Error("invalid args")
	}
	articleBytes, err := stub.GetState(constructArticleKey(id))
	if err != nil || len(articleBytes) == 0 {
		return shim.Error("asset not exist")
	}

	return shim.Success(articleBytes)
}

//删除物品
func articleDestroy(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		shim.Error("not enough args")
	}
	articleID, err := strconv.Atoi(args[0])
	if err != nil || articleID == 0 {
		shim.Error("invalid args")
	}
	if articleBytes, err := stub.GetState(constructArticleKey(articleID)); err != nil || len(articleBytes) == 0 {
		shim.Error("article not exist")
	}
	if err := stub.DelState(constructArticleKey(articleID)); err != nil {
		shim.Error(fmt.Sprintf("del article err: %d", err))
	}
	return shim.Success(nil)
}

/*初始化方法 */
func (c *BhxyCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	user := &User{
		Id:       1,
		Name:     "崩坏娘",
		Level:    480,
		LevelExp: 0,
		VipLevel: 18,
		VipExp:   0,
		Coin:     100000000,
		Crystal:  100000000,
		Assets:   nil,
	}

	userByte, _ := json.Marshal(user)
	if err := stub.PutState(constructUserKey(user.Id), userByte); err != nil {
		shim.Error(fmt.Sprintf("init Error %s", err))
	}

	return shim.Success(nil)
}

/*调用方法入口*/
func (c *BhxyCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "userRegister":
		return userRegister(stub, args)
	case "userDestroy":
		return userDestroy(stub, args)
	case "queryUser":
		return queryUser(stub, args)
	case "assetEnroll":
		return assetEnroll(stub, args)
	case "queryAsset":
		return queryAsset(stub, args)
	case "assetDestroy":
		return assetDestroy(stub, args)
	case "addArticle":
		return addArticle(stub, args)
	case "queryArticle":
		return queryArticle(stub, args)
	case "articleDestroy":
		return articleDestroy(stub, args)
	case "xyh":
		s, _ := json.Marshal("辣鸡咸鱼皇!")
		return shim.Success(s)
	default:
		return shim.Error(fmt.Sprintf("unknow function: %s", function))
	}

	//return shim.Error(fmt.Sprintf("unsupported function: %s", c))
}

func main() {
	err := shim.Start(new(BhxyCC))
	if err != nil {
		fmt.Printf("Error starting bhxy chaincode: %s", err)
	}
}
