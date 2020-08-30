refs:
- https://github.com/blankjee/file-storage-system
- https://github.com/samtake/filestore-server/blob/master/filestore-server-v0.2/static/view/signin.html


# Arch

- Front end : React + Ant D
- Back end: Golang + Gin

## Coding Progress
###  Simple file upload:
- [ ] [Http form action url `#`](https://developer.mozilla.org/zh-CN/docs/Learn/HTML/Forms/Sending_and_retrieving_form_data)
- [ ] [Http code 302](https://www.cnblogs.com/woshimrf/p/http-code-302.html)
- [ ] [HandleFunc `func http.HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))`]()
- [ ] `os.Create()` when the dir not exist will error.


### Keep meta data of the file
1. Use the hash of the file as it's id.
2. Generate file meta info before save files.
3. FileMeta info query interface.
   1. Single file query.
   2. Range file query.
4. Download file using fileHash.

   A better way is: when the file is uploaded, and handle these files as static source, NGINX, for example, deploy it as _Reverse Proxy_ to download the file, also we can do simple stream limiting and permission check, which will reduce the pressure of the back end.
5. Update file name by file hash.
6. Del file both in the fileMetaMap and disk.

Current work flow:
1. Using get to get the file upload page.
2. Choose the local file and, use form to upload file.
3. Back end get file upload stream, and save to local disk.
4. Update the the `metaFileMap`.

### Using MySQL to save the metaInf
#### MySQL table design

<details> <summary>SQL</summary>

```sql
-- 创建文件表
CREATE TABLE `tbl_file` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `file_sha1` char(40) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `file_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
  `create_at` datetime default NOW() COMMENT '创建日期',
  `update_at` datetime default NOW() on update current_timestamp() COMMENT '更新日期',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)',
  `ext1` int(11) DEFAULT '0' COMMENT '备用字段1',
  `ext2` text COMMENT '备用字段2',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_hash` (`file_sha1`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 创建用户表
CREATE TABLE `tbl_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT '用户encoded密码',
  `email` varchar(64) DEFAULT '' COMMENT '邮箱',
  `phone` varchar(128) DEFAULT '' COMMENT '手机号',
  `email_validated` tinyint(1) DEFAULT 0 COMMENT '邮箱是否已验证',
  `phone_validated` tinyint(1) DEFAULT 0 COMMENT '手机号是否已验证',
  `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
  `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间戳',
  `profile` text COMMENT '用户属性',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`user_name`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;

-- 创建用户token表
CREATE TABLE `tbl_user_token` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
  `user_token` char(40) NOT NULL DEFAULT '' COMMENT '用户登录token',
    PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 创建用户文件表
CREATE TABLE `tbl_user_file` (
  `id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
  `user_name` varchar(64) NOT NULL,
  `file_sha1` varchar(64) NOT NULL DEFAULT '' COMMENT '文件hash',
  `file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
  `file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
  `upload_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `last_update` datetime DEFAULT CURRENT_TIMESTAMP
          ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态(0正常1已删除2禁用)',
  UNIQUE KEY `idx_user_file` (`user_name`, `file_sha1`),
  KEY `idx_status` (`status`),
  KEY `idx_user_id` (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```
 </details>

#### Sub-database and sub-table
- Horizontal sub-table
  - If we need to divide a large table into 256 tables
  - We could use the last two digits of `Sha1` of a record to put it into certain table
  - Use the rule of `tbl_${file_sha1}[:-2]`, e.g. `tbl_01` ... `tal_ff`

- *Vertical split*
  - Divide by all fields of a large table into small table
  - Use a field to connect the small table

### DB ops in go
- [ ] mysql_prepared_statements
- [ ] https://blog.biezhi.me/2018/10/values-or-pointers-in-golang.html

# What I need ?

1. RateLimit from stream video
2. Watch src code of golang
3. QR-Code info