package service

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"

	"tzh.com/web/model"
	"tzh.com/web/util"
)

// 业务处理函数, 获取用户列表
func ListUser(username string, offset, limit int) ([]*model.UserInfo, uint, error) {
	infos := make([]*model.UserInfo, 0)
	//调用 本地缓存提高性能
	var localCacheRet []*model.UserInfo
	localCacheKey := fmt.Sprintf("%s%s%d%d", "user_cache_", username)
	cacheRet, err := model.LocalCache.Self.Get(localCacheKey)
	if cacheRet != nil && err == nil {
		err = json.Unmarshal(cacheRet, &localCacheRet)
		if err != nil {
			logrus.Errorf("json Unmarshal err:[%v]", err)
		}

		//localCacheRet.List = localCacheRet.List[(page-1)*size : size]
		return localCacheRet, uint(len(localCacheRet)), nil
	}

	users, count, err := model.ListUser(username, offset, limit)
	if err != nil {
		return nil, count, err
	}

	ids := []uint{}
	for _, user := range users {
		ids = append(ids, user.ID)
	}

	wg := sync.WaitGroup{}
	userList := model.UserList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint]*model.UserInfo, len(users)),
	}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	//开启临时对象池
	pool := sync.Pool{
		New: func() interface{} {
			return model.UserModel{}
		},
	}

	// 并行转换
	for _, u := range users {
		wg.Add(1)
		go func(u *model.UserModel) {
			defer wg.Done()

			shortID, err := util.GenShortID()
			if err != nil {
				errChan <- err
				return
			}

			// 更新数据时加锁, 保持一致性
			userList.Lock.Lock()
			defer userList.Lock.Unlock()

			//先从临时对象池获取，减少内存申请频率
			info := pool.Get().(model.UserInfo)
			info = model.UserInfo{
				ID:        u.ID,
				Username:  u.Username,
				SayHello:  fmt.Sprintf("Hello %s", shortID),
				Password:  u.Password,
				CreatedAt: util.TimeToStr(&u.CreatedAt),
				UpdatedAt: util.TimeToStr(&u.UpdatedAt),
				DeletedAt: util.TimeToStr(u.DeletedAt),
			}
			userList.IdMap[u.ID] = &info
		}(u)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	// 等待完成
	select {
	case <-finished:
	case err := <-errChan:
		return nil, count, err
	}

	for _, id := range ids {
		infos = append(infos, userList.IdMap[id])
	}

	return infos, count, nil
}
