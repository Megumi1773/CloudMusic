  # CloudMusic功能与接口文档

## 项目概述

CloudMusic是一个音乐播放器应用程序的后端服务，提供了用户管理、音乐播放、歌单管理等功能。该项目使用Golang语言和Gin框架构建，采用GORM作为ORM框架与MySQL数据库交互，并使用MinIO进行对象存储。

## 系统架构

1. **语言框架**: Golang + Gin
2. **数据库**: MySQL + GORM
3. **对象存储**: MinIO
4. **认证方式**: JWT Token

## 核心功能模块

### 1. 用户管理模块

- 用户注册
- 用户登录
- 用户登出
- 获取用户信息
- 更新用户信息

### 2. 音乐播放模块

- 获取歌曲详细信息
- 获取歌曲播放地址
- 获取歌词

### 3. 歌手管理模块

- 获取歌手详情
- 获取歌手的歌曲列表
- 获取歌手的专辑列表

### 4. 专辑管理模块

- 获取专辑详情
- 获取专辑歌曲列表

### 5. 歌单管理模块

- 创建歌单
- 获取歌单详情
- 获取用户全部歌单
- 更新歌单信息
- 删除歌单
- 向歌单添加歌曲
- 从歌单删除歌曲
- 获取歌单的歌曲列表

## 数据库模型

### 用户模型 (User)

```go
type User struct {
    gorm.Model
    Username string
    Password string
    Nickname string
    Email    string
    Phone    string
    Avatar   string
    Playlists []Playlist
    Comments []Comment
}
```

### 歌曲模型 (Song)

```go
type Song struct {
    Id       int
    Name     string
    ArtistId uint64
    AlbumId  uint64
    Duration uint32
    Url      string
    Lyric    string
    CreateAt time.Time
    UpdateAt time.Time
    
    Artist Artist
    Album Album
    Playlists []Playlist
}
```

### 歌单模型 (Playlist)

```go
type Playlist struct {
    gorm.Model
    Name        string
    UserId      uint
    Cover       string
    Description string
    IsPublic    uint

    User User
    Songs []Song
}
```

### 歌手模型 (Artist)

```go
type Artist struct {
    Id          uint64
    Name        string
    Description string
    Avatar      string
}
```

### 专辑模型 (Album)

```go
type Album struct {
    Id          uint64
    Name        string
    ArtistId    uint64
    Cover       string
    Description string
    PublishTime time.Time
    
    Artist Artist
    Songs []Song
}
```

## API接口文档

### 用户相关接口

#### 用户注册

- **URL**: `/api/user/register`
- **方法**: `POST`
- **请求体**:

  ```json
  {
    "username": "用户名",
    "password": "密码"
  }
  ```

- **响应**:

  ```json
  {
    "code": 200,
    "message": "注册成功",
    "data": {
      "username": "用户名",
      "nickname": "随机生成的昵称"
    }
  }
  ```

#### 用户登录

- **URL**: `/api/user/login`
- **方法**: `POST`
- **请求体**:

  ```json
  {
    "username": "用户名",
    "password": "密码"
  }
  ```

- **响应**:

  ```json
  {
    "code": 200,
    "message": "登录成功",
    "data": {
      "user": {
        "username": "用户名",
        "nickname": "昵称"
      },
      "token": "JWT令牌"
    }
  }
  ```

#### 用户登出

- **URL**: `/api/user/logout`
- **方法**: `POST`
- **请求头**: `Authorization: {token}`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "退出成功",
    "data": null
  }
  ```

#### 获取用户信息

- **URL**: `/api/user/info`
- **方法**: `GET`
- **请求头**: `Authorization: {token}`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "username": "用户名",
      "nickname": "昵称",
      "email": "邮箱",
      "phone": "手机号"
    }
  }
  ```

#### 更新用户信息

- **URL**: `/api/user/info`
- **方法**: `PUT`
- **请求头**: `Authorization: {token}`
- **请求体**:

  ```json
  {
    "nickName": "新昵称",
    "email": "新邮箱",
    "phone": "新手机号"
  }
  ```

- **响应**:

  ```json
  {
    "code": 200,
    "message": "更新成功",
    "data": {
      "nickName": "新昵称",
      "email": "新邮箱",
      "phone": "新手机号"
    }
  }
  ```

### 音乐相关接口

#### 获取歌曲详情

- **URL**: `/api/v1/songs/{id}`
- **方法**: `GET`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "id": 1,
      "name": "歌曲名",
      "duration": 180,
      "artist_name": "歌手名",
      "album_name": "专辑名",
      "album_cover": "专辑封面URL"
    }
  }
  ```

#### 获取歌曲播放地址

- **URL**: `/api/v1/songs/url/{id}`
- **方法**: `GET`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "url": "歌曲播放地址"
    }
  }
  ```

#### 获取歌词

- **URL**: `/api/v1/songs/lyric/{id}`
- **方法**: `GET`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "lyric": "歌词文本"
    }
  }
  ```

### 歌手相关接口

#### 获取歌手详情

- **URL**: `/api/artists/{id}`
- **方法**: `GET`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "id": 1,
      "name": "歌手名",
      "description": "歌手描述",
      "avatar": "歌手头像URL"
    }
  }
  ```

#### 获取歌手的歌曲

- **URL**: `/api/artists/{id}/songs`
- **方法**: `GET`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "成功",
    "data": [
      {
        "id": 1,
        "name": "歌曲名",
        "duration": 180,
        "album_name": "专辑名",
        "album_cover": "专辑封面URL"
      }
    ]
  }
  ```

#### 获取歌手的专辑

- **URL**: `/api/artists/{id}/albums`
- **方法**: `GET`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "成功",
    "data": [
      {
        "id": 1,
        "name": "专辑名",
        "cover": "专辑封面URL",
        "description": "专辑描述",
        "publish_time": "发行时间"
      }
    ]
  }
  ```

### 专辑相关接口

#### 获取专辑详情

- **URL**: `/api/albums/{id}`
- **方法**: `GET`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "id": 1,
      "name": "专辑名",
      "artist_id": 1,
      "artist_name": "歌手名",
      "cover": "专辑封面URL",
      "description": "专辑描述",
      "publish_time": "发行时间"
    }
  }
  ```

#### 获取专辑的歌曲

- **URL**: `/api/albums/{id}/songs`
- **方法**: `GET`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "成功",
    "data": [
      {
        "id": 1,
        "name": "歌曲名",
        "duration": 180
      }
    ]
  }
  ```

### 歌单相关接口

#### 获取歌单详情

- **URL**: `/api/playlists/{id}`
- **方法**: `GET`
- **请求头**: `Authorization: {token}`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "id": 1,
      "name": "歌单名",
      "nickname": "创建者昵称",
      "user_avatar": "创建者头像",
      "description": "歌单描述",
      "cover": "歌单封面",
      "is_public": 1,
      "created_at": "创建时间"
    }
  }
  ```

#### 获取用户全部歌单

- **URL**: `/api/playlists`
- **方法**: `GET`
- **请求头**: `Authorization: {token}`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "成功",
    "data": [
      {
        "id": 1,
        "name": "歌单名",
        "cover": "歌单封面",
        "is_public": 1,
        "created_at": "创建时间"
      }
    ]
  }
  ```

#### 获取歌单的歌曲

- **URL**: `/api/playlists/{id}/songs`
- **方法**: `GET`
- **请求头**: `Authorization: {token}`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "成功",
    "data": [
      {
        "id": 1,
        "name": "歌曲名",
        "duration": 180,
        "artist_name": "歌手名",
        "album_name": "专辑名",
        "album_cover": "专辑封面URL"
      }
    ]
  }
  ```

#### 创建歌单

- **URL**: `/api/playlists`
- **方法**: `POST`
- **请求头**: `Authorization: {token}`
- **请求体**:

  ```json
  {
    "name": "歌单名",
    "description": "歌单描述",
    "cover": "歌单封面URL",
    "is_public": 1
  }
  ```

- **响应**:

  ```json
  {
    "code": 200,
    "message": "创建成功",
    "data": {
      "id": 1,
      "name": "歌单名",
      "description": "歌单描述",
      "cover": "歌单封面URL",
      "is_public": 1,
      "created_at": "创建时间"
    }
  }
  ```

#### 更新歌单

- **URL**: `/api/playlists/{id}`
- **方法**: `PUT`
- **请求头**: `Authorization: {token}`
- **请求体**:

  ```json
  {
    "name": "新歌单名",
    "description": "新歌单描述",
    "cover": "新歌单封面URL",
    "is_public": 0
  }
  ```

- **响应**:

  ```json
  {
    "code": 200,
    "message": "更新成功",
    "data": {
      "id": 1,
      "name": "新歌单名",
      "description": "新歌单描述",
      "cover": "新歌单封面URL",
      "is_public": 0
    }
  }
  ```

#### 删除歌单

- **URL**: `/api/playlists/{id}`
- **方法**: `DELETE`
- **请求头**: `Authorization: {token}`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "删除成功",
    "data": null
  }
  ```

#### 添加歌曲到歌单

- **URL**: `/api/playlists/{id}/songs`
- **方法**: `POST`
- **请求头**: `Authorization: {token}`
- **请求体**:

  ```json
  {
    "song_id": 1
  }
  ```

- **响应**:

  ```json
  {
    "code": 200,
    "message": "添加成功",
    "data": null
  }
  ```

#### 从歌单中删除歌曲

- **URL**: `/api/playlists/{id}/songs/{songId}`
- **方法**: `DELETE`
- **请求头**: `Authorization: {token}`
- **响应**:

  ```json
  {
    "code": 200,
    "message": "删除成功",
    "data": null
  }
  ```

## 前端开发指南

### 推荐技术栈

1. **框架**: Vue.js 或 React
2. **UI库**: Element UI, Ant Design 或 Material UI
3. **状态管理**: Vuex (Vue) 或 Redux (React)
4. **路由**: Vue Router 或 React Router
5. **HTTP请求**: Axios

### 身份验证流程

1. 用户登录后，后端返回JWT令牌
2. 前端将令牌保存在localStorage或sessionStorage中
3. 对于需要身份验证的API请求，在请求头中添加Authorization字段
4. 用户登出时，从存储中删除令牌并将令牌添加到黑名单中

### 主要页面建议

1. **登录/注册页面**
2. **首页/推荐页面**
3. **歌曲播放页面**
4. **歌单页面**
5. **歌手页面**
6. **专辑页面**
7. **用户个人中心页面**
8. **搜索结果页面**

### 数据获取流程

1. 应用启动时获取用户信息（如果已登录）
2. 获取用户的歌单列表
3. 按需加载歌曲、歌手和专辑详情

### 注意事项

1. 请求API时应处理各种错误情况
2. 播放音乐时注意音频资源的加载和缓存
3. 对于可能较大的数据（如歌单中的歌曲列表）应实现分页加载
4. 实现响应式设计，适应不同屏幕尺寸
