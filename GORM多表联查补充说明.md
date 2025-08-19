# GORM多表联查补充说明

## 一、关联类型声明的区别

虽然`has many`、`has one`和`belongs to`的声明语法看起来相似，但它们在概念和使用上有明显区别：

### 1.1 声明方式的本质区别

#### belongs to（从属关系）

- **关系方向**：从子表指向父表
- **外键位置**：外键在当前模型中
- **典型场景**：歌曲属于歌手，外键(artist_id)在歌曲表中

```go
// 在Song模型中声明belongs to关系
type Song struct {
    gorm.Model
    Name     string
    ArtistID uint      // 外键在当前模型
    Artist   Artist    `gorm:"foreignKey:ArtistID"` // 指向父表
}
```

#### has one（拥有一个）

- **关系方向**：从父表指向子表
- **外键位置**：外键在关联模型中
- **典型场景**：用户有一个个人资料，外键(user_id)在个人资料表中

```go
// 在User模型中声明has one关系
type User struct {
    gorm.Model
    Name    string
    // 外键在关联模型(Profile)中，不在当前模型
    Profile Profile `gorm:"foreignKey:UserID"` // 指向子表
}

type Profile struct {
    gorm.Model
    UserID   uint   // 外键在这里
    Address  string
}
```

#### has many（拥有多个）

- **关系方向**：从父表指向多个子表记录
- **外键位置**：外键在关联模型中
- **典型场景**：歌手有多首歌曲，外键(artist_id)在歌曲表中

```go
// 在Artist模型中声明has many关系
type Artist struct {
    gorm.Model
    Name  string
    // 外键在关联模型(Song)中，不在当前模型
    Songs []Song `gorm:"foreignKey:ArtistID"` // 指向多个子表记录
}

type Song struct {
    gorm.Model
    Name     string
    ArtistID uint   // 外键在这里
}
```

### 1.2 关联关系的反向引用

GORM允许在两个模型中同时声明关联关系的两个方向：

```go
// 双向关联示例
type Artist struct {
    gorm.Model
    Name  string
    Songs []Song `gorm:"foreignKey:ArtistID"` // has many关系
}

type Song struct {
    gorm.Model
    Name     string
    ArtistID uint    // 外键
    Artist   Artist  `gorm:"foreignKey:ArtistID"` // belongs to关系
}
```

### 1.3 参考字段设置

如果外键引用的不是主键，可以使用`references`标签指定：

```go
type Artist struct {
    gorm.Model
    UniqueCode string `gorm:"uniqueIndex"` // 不是主键，但有唯一索引
    Name       string
    Songs      []Song `gorm:"foreignKey:ArtistCode;references:UniqueCode"` 
}

type Song struct {
    gorm.Model
    Name       string
    ArtistCode string // 外键引用Artist的UniqueCode字段，而不是ID
}
```

## 二、手动Joins与模型关系的关系

### 2.1 手动Joins与模型关系的独立性

**手动Joins不需要配置表关系模型**，它们是完全独立的功能：

- **模型关系声明**：用于定义对象之间的关联，主要用于Preload和Association操作
- **手动Joins**：直接在SQL级别执行连接查询，不依赖于模型关系声明

### 2.2 手动Joins的详细用法

#### 基本JOIN查询

```go
type Result struct {
    ID        uint
    SongName  string
    ArtistName string
}

var results []Result
db.Table("songs").
    Select("songs.id, songs.name as song_name, artists.name as artist_name").
    Joins("LEFT JOIN artists ON songs.artist_id = artists.id").
    Scan(&results)
```

#### 多表JOIN

```go
db.Table("songs").
    Select("songs.id, songs.name, artists.name as artist_name, albums.name as album_name").
    Joins("LEFT JOIN artists ON songs.artist_id = artists.id").
    Joins("LEFT JOIN albums ON songs.album_id = albums.id").
    Where("artists.name LIKE ?", "%周杰伦%").
    Scan(&results)
```

#### JOIN时使用条件

```go
// INNER JOIN with conditions
db.Table("songs").
    Joins("JOIN artists ON songs.artist_id = artists.id AND artists.is_popular = ?", true).
    Scan(&results)
```

#### 使用原生SQL

```go
db.Raw(`
    SELECT s.id, s.name, a.name as artist_name 
    FROM songs s 
    LEFT JOIN artists a ON s.artist_id = a.id 
    WHERE s.name LIKE ?
`, "%爱%").Scan(&results)
```

### 2.3 何时使用手动Joins而非模型关系

- **复杂查询条件**：需要在JOIN条件中添加额外筛选
- **性能优化**：只查询必要字段，减少数据传输
- **多表聚合**：需要聚合多个表的数据
- **特定SQL需求**：需要使用特定数据库的SQL特性

## 三、子查询详解

子查询是在一个查询中嵌套另一个查询，GORM提供了多种方式实现子查询。

### 3.1 基本子查询

#### IN子查询

```go
// 查找所有流行歌手的歌曲
db.Where("artist_id IN (?)", 
    db.Table("artists").Select("id").Where("is_popular = ?", true),
).Find(&songs)
```

生成的SQL：

```sql
SELECT * FROM songs WHERE artist_id IN (
    SELECT id FROM artists WHERE is_popular = true
)
```

#### EXISTS子查询

```go
// 查找至少有一首歌的歌手
db.Where("EXISTS (?)", 
    db.Table("songs").Select("1").Where("songs.artist_id = artists.id"),
).Find(&artists)
```

生成的SQL：

```sql
SELECT * FROM artists WHERE EXISTS (
    SELECT 1 FROM songs WHERE songs.artist_id = artists.id
)
```

### 3.2 子查询作为表

```go
// 使用子查询作为表
db.Table("(?) as u", 
    db.Table("songs").Select("artist_id, COUNT(*) as song_count").Group("artist_id"),
).Joins("JOIN artists ON artists.id = u.artist_id").
  Select("artists.name, u.song_count").
  Scan(&results)
```

生成的SQL：

```sql
SELECT artists.name, u.song_count FROM (
    SELECT artist_id, COUNT(*) as song_count FROM songs GROUP BY artist_id
) as u JOIN artists ON artists.id = u.artist_id
```

### 3.3 子查询作为字段

```go
// 查询每个歌手及其歌曲数量
type ArtistWithSongCount struct {
    ID        uint
    Name      string
    SongCount int
}

var results []ArtistWithSongCount
db.Model(&Artist{}).
    Select("artists.id, artists.name, (?)", 
        db.Table("songs").Select("COUNT(*)").Where("artist_id = artists.id"),
    ).
    Scan(&results)
```

生成的SQL：

```sql
SELECT artists.id, artists.name, (
    SELECT COUNT(*) FROM songs WHERE artist_id = artists.id
) FROM artists
```

### 3.4 子查询与关联查询的选择

- **子查询优势**：可以实现复杂的筛选和聚合
- **关联查询优势**：代码更简洁，更符合对象关系模型
- **选择依据**：根据查询复杂度和性能需求选择

## 四、Association关联详解

Association是GORM提供的关联操作API，用于处理模型之间的关联关系。

### 4.1 获取关联

#### 基本用法

```go
// 获取歌单中的所有歌曲
var songs []Song
db.Model(&playlist).Association("Songs").Find(&songs)

// 带条件的关联查询
db.Model(&playlist).Association("Songs").Find(&songs, "duration > ?", 180)
```

#### 关联计数

```go
// 获取歌单中的歌曲数量
count := db.Model(&playlist).Association("Songs").Count()
```

### 4.2 添加关联

#### 添加单个关联

```go
// 添加一首歌到歌单
db.Model(&playlist).Association("Songs").Append(&song)
```

#### 添加多个关联

```go
// 添加多首歌到歌单
songs := []Song{song1, song2, song3}
db.Model(&playlist).Association("Songs").Append(&songs)
```

#### 通过主键添加

```go
// 通过ID添加
db.Model(&playlist).Association("Songs").Append(&Song{ID: 1})

// 添加多个ID
db.Model(&playlist).Association("Songs").Append(&[]Song{{ID: 1}, {ID: 2}})
```

### 4.3 替换关联

```go
// 替换歌单中的所有歌曲
newSongs := []Song{song1, song2}
db.Model(&playlist).Association("Songs").Replace(&newSongs)
```

### 4.4 删除关联

#### 删除单个关联

```go
// 从歌单中删除一首歌
db.Model(&playlist).Association("Songs").Delete(&song)
```

#### 删除多个关联

```go
// 从歌单中删除多首歌
songsToDelete := []Song{song1, song2}
db.Model(&playlist).Association("Songs").Delete(&songsToDelete)
```

#### 通过条件删除

```go
// 删除所有时长超过5分钟的歌曲
db.Model(&playlist).Association("Songs").Delete(&Song{}, "duration > ?", 300)
```

### 4.5 清空关联

```go
// 清空歌单中的所有歌曲
db.Model(&playlist).Association("Songs").Clear()
```

### 4.6 批量处理关联

```go
// 批量处理多个歌单的关联
playlists := []Playlist{playlist1, playlist2}
db.Model(&playlists).Association("Songs").Append(&song)
db.Model(&playlists).Association("Songs").Replace(&[]Song{song1, song2})
```

### 4.7 Association与事务

```go
// 在事务中处理关联
err := db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Model(&playlist).Association("Songs").Append(&song); err != nil {
        return err
    }
    
    // 更新歌单信息
    if err := tx.Model(&playlist).Update("updated_at", time.Now()).Error; err != nil {
        return err
    }
    
    return nil
})
```

## 五、实际应用场景示例

### 5.1 复杂查询示例：获取用户收藏的所有歌手的热门歌曲

```go
// 使用子查询和关联查询结合
type Result struct {
    ArtistID   uint
    ArtistName string
    SongName   string
    PlayCount  int
}

var results []Result

// 1. 获取用户收藏的歌手ID
var favoriteArtistIDs []uint
db.Table("user_favorites").
    Where("user_id = ? AND type = ?", userId, "artist").
    Pluck("target_id", &favoriteArtistIDs)

// 2. 获取这些歌手的热门歌曲
db.Table("songs").
    Select("songs.artist_id, artists.name as artist_name, songs.name as song_name, songs.play_count").
    Joins("JOIN artists ON artists.id = songs.artist_id").
    Where("songs.artist_id IN ? AND songs.play_count > ?", favoriteArtistIDs, 1000).
    Order("songs.play_count DESC").
    Limit(5).
    Scan(&results)
```

### 5.2 多层嵌套关联示例：获取歌单详情（包含歌曲、歌手和专辑信息）

```go
type PlaylistDetail struct {
    ID          uint
    Name        string
    Description string
    Songs       []SongWithDetails
}

type SongWithDetails struct {
    ID         uint
    Name       string
    Duration   int
    ArtistName string
    AlbumName  string
    AlbumCover string
}

func GetPlaylistWithDetails(playlistID uint) (*PlaylistDetail, error) {
    var playlist Playlist
    
    // 使用嵌套预加载
    if err := db.Preload("Songs").
              Preload("Songs.Artist").
              Preload("Songs.Album").
              First(&playlist, playlistID).Error; err != nil {
        return nil, err
    }
    
    // 构建响应
    result := &PlaylistDetail{
        ID:          playlist.ID,
        Name:        playlist.Name,
        Description: playlist.Description,
        Songs:       make([]SongWithDetails, 0, len(playlist.Songs)),
    }
    
    for _, song := range playlist.Songs {
        result.Songs = append(result.Songs, SongWithDetails{
            ID:         song.ID,
            Name:       song.Name,
            Duration:   song.Duration,
            ArtistName: song.Artist.Name,
            AlbumName:  song.Album.Name,
            AlbumCover: song.Album.Cover,
        })
    }
    
    return result, nil
}
```

### 5.3 使用Association处理多对多关系：管理用户歌单

```go
// 创建歌单并添加歌曲
func CreatePlaylistWithSongs(userID uint, name string, songIDs []uint) error {
    // 开始事务
    return db.Transaction(func(tx *gorm.DB) error {
        // 创建歌单
        playlist := Playlist{
            Name:   name,
            UserID: userID,
        }
        
        if err := tx.Create(&playlist).Error; err != nil {
            return err
        }
        
        // 准备歌曲引用
        var songs []Song
        for _, id := range songIDs {
            songs = append(songs, Song{ID: id})
        }
        
        // 添加歌曲到歌单
        return tx.Model(&playlist).Association("Songs").Append(&songs)
    })
}

// 更新歌单中的歌曲
func UpdatePlaylistSongs(playlistID uint, songIDs []uint) error {
    var playlist Playlist
    if err := db.First(&playlist, playlistID).Error; err != nil {
        return err
    }
    
    // 准备歌曲引用
    var songs []Song
    for _, id := range songIDs {
        songs = append(songs, Song{ID: id})
    }
    
    // 替换歌单中的所有歌曲
    return db.Model(&playlist).Association("Songs").Replace(&songs)
}
