# ISCTF 开发文档



***

背景：此 CTF 平台主要用于多校联合比赛

***



## 一、参赛学校模块

### 1. 数据库设计

***

**数据库：**`dalictf`

**表名：**`dalictf_school`

***

#### 字段设计

|    字段名    |             类型              |             约束条件              |                   说明                   |
| :----------: | :---------------------------: | :-------------------------------: | :--------------------------------------: |
|      id      |         `BIGINT(20)`          |  `PRIMARY KEY AUTO_INCREMENT`   |                   主键                   |
| school_name  |        `VARCHAR(255)`         |        `NOT NULL UNIQUE`        |             参赛学校名字             |
| school_admin |         `BIGINT(20)`          |            `DEFAULT NULL`             | 院校负责人（外键关联用户表 user_id） |
|  user_count  |            `INT(11)`            |        `NOT NULL DEFAULT 0`         |      学校参赛人数（实时统计）      |
|    status    | `ENUM('active', 'suspended')` |    `NOT NULL DEFAULT 'active'`    |     学校状态（active:正常/suspended:封禁）     |
|  created_at  |          `DATETIME`           |        `NOT NULL DEFAULT CURRENT_TIMESTAMP`         |                 创建时间                 |
|  updated_at  |          `DATETIME`           | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` |                 更新时间                 |
|  deleted_at  |          `DATETIME`           |            `DEFAULT NULL`             |      软删除时间（NULL表示未删除）      |

#### **建表 `sql`** 

```sql
CREATE TABLE `dalictf_school` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '学校主键ID',
  `school_name` VARCHAR(255) NOT NULL COMMENT '参赛学校名称',
  `school_admin` BIGINT(20) DEFAULT NULL COMMENT '院校负责人ID（外键关联用户表）',
  `user_count` INT(11) NOT NULL DEFAULT 0 COMMENT '学校参赛人数',
  `status` ENUM('active', 'suspended') NOT NULL DEFAULT 'active' COMMENT '学校状态：active-正常, suspended-封禁',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_school_name` (`school_name`),
  KEY `idx_school_admin` (`school_admin`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='参赛学校表';

-- 添加外键约束（需要先创建用户表）
-- ALTER TABLE `dalictf_school` 
-- ADD CONSTRAINT `fk_school_admin` 
-- FOREIGN KEY (`school_admin`) REFERENCES `dalictf_user`(`id`) 
-- ON DELETE SET NULL ON UPDATE CASCADE;
```

***



### 2. API 设计

> **统一前缀：**`/api/v1/schools`
>
> **认证机制：**
> - 需要权限的接口必须在请求头中携带 JWT Token
> - Token 通过用户登录接口获取，有效期为 24 小时
> - Token 失效后需要重新登录获取新 Token
>
> **统一请求头：**
>
> ```
> Content-Type: application/json
> Authorization: Bearer <JWT_TOKEN>  // 需要权限的接口必须携带
> ```
>
> **返回格式：**
>
> ```json
> {
>     "code": "<数字错误码，200 表示成功>",
>     "msg": "<提示信息>",
>     "data": "<返回数据对象或列表>"
> }
> ```

#### 新增学校

- **URL：** `POST /api/v1/schools`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求体：**

  ```json
  {
    "school_name": "大理大学",
    "school_admin": 10001  // 可选，院校负责人用户ID
  }
  ```

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "School created successfully",
    "data": {
      "id": 1,
      "school_name": "大理大学",
      "user_count": 0,
      "status": "active"
    }
  }
  ```

#### 查询学校列表

- **URL**：`GET /api/v1/schools`

- **权限**：公开

- **请求头：**

  ```
  Content-Type: application/json
  // 公开接口，无需 Authorization
  ```

- **请求参数（Query）：**
  - `page` 页码（默认 1）
  - `limit` 每页数量（默认 20）
  - `search` 模糊搜索字段（学校名称）
  - `sort_by` 排序字段（如 `user_count`、`created_at`）
  - `order` 排序方式（asc/desc，默认 desc）
  - `status` 学校状态（`active` / `suspended`，管理员可选）

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total": 100,
      "page": 1,
      "limit": 20,
      "list": [
        {
          "id": 1,
          "school_name": "大理大学",
          "school_admin": 10001,
          "admin_name": "张老师",
          "user_count": 45,
          "status": "active",
          "created_at": "2025-01-01 10:00:00",
          "updated_at": "2025-01-10 15:30:00"
        },
        {
          "id": 2,
          "school_name": "云南大学",
          "school_admin": null,
          "admin_name": null,
          "user_count": 38,
          "status": "active",
          "created_at": "2025-01-02 11:00:00",
          "updated_at": "2025-01-12 16:20:00"
        }
      ]
    }
  }
  ```

#### 查询单个学校详情

- **URL**：`GET /api/v1/schools/:id` 或 `GET /api/v1/schools/name/:school_name`

- **权限**：公开

- **请求头：**

  ```
  Content-Type: application/json
  // 公开接口，无需 Authorization
  ```

- **路径参数：**
  - `id` 学校ID
  - `school_name` 学校名称（URL编码）

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "id": 1,
      "school_name": "大理大学",
      "school_admin": 10001,
      "admin_name": "张老师",
      "admin_email": "zhang@example.com",
      "user_count": 45,
      "status": "active",
      "created_at": "2025-01-01 10:00:00",
      "updated_at": "2025-01-10 15:30:00"
    }
  }
  ```

#### 修改学校信息

- **URL**：`PUT /api/v1/schools/:id`

- **权限**：管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 学校ID

- **请求体：**

  ```json
  {
    "school_name": "大理大学（修改后）",  // 可选
    "school_admin": 10002  // 可选，院校负责人用户ID，传null则清空负责人
  }
  ```

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "School updated successfully",
    "data": {
      "id": 1,
      "school_name": "大理大学（修改后）",
      "school_admin": 10002,
      "user_count": 45,
      "status": "active",
      "updated_at": "2025-01-15 14:20:00"
    }
  }
  ```

#### 删除学校（软删除）

- **URL**：`DELETE /api/v1/schools/:id`

- **权限**：管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 学校ID

- **请求体：**无

- **说明**：软删除，设置 `deleted_at` 字段为当前时间，不会真正删除数据

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "School deleted successfully",
    "data": null
  }
  ```

#### 修改学校状态

- **URL**：`PATCH /api/v1/schools/:id/status`

- **权限**：管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 学校ID

- **请求体：**

  ```json
  {
    "status": "suspended"  // active 或 suspended
  }
  ```

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "School status updated successfully",
    "data": {
      "id": 1,
      "school_name": "大理大学",
      "status": "suspended",
      "updated_at": "2025-01-16 09:15:00"
    }
  }
  ```

#### 批量分配学校负责人

- **URL**：`POST /api/v1/schools/batch-assign-admin`

- **权限**：管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求体：**

  ```json
  {
    "assignments": [
      {
        "school_id": 1,
        "admin_id": 10001
      },
      {
        "school_id": 2,
        "admin_id": 10002
      }
    ]
  }
  ```

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Batch assignment completed",
    "data": {
      "success_count": 2,
      "failed_count": 0,
      "failed_items": []
    }
  }
  ```

#### 获取学校统计信息

- **URL**：`GET /api/v1/schools/statistics`

- **权限**：管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求体：**无

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total_schools": 50,
      "active_schools": 48,
      "suspended_schools": 2,
      "total_participants": 1250,
      "schools_with_admin": 45,
      "schools_without_admin": 5
    }
  }
  ```

***

### 3. 错误码定义

| 错误码 |               说明               |
| :----: | :------------------------------: |
|  200   |               成功               |
|  400   |           请求参数错误           |
|  401   | 未授权（Token缺失或格式错误）  |
|  402   |       Token已过期，请重新登录       |
|  403   |     无权限操作（权限不足）     |
|  404   |           学校不存在           |
|  409   |         学校名称已存在         |
|  500   |           服务器错误           |

#### 认证错误返回示例

```json
{
  "code": 401,
  "msg": "Unauthorized: Token is missing or invalid",
  "data": null
}
```

```json
{
  "code": 402,
  "msg": "Token expired, please login again",
  "data": null
}
```

```json
{
  "code": 403,
  "msg": "Forbidden: Insufficient permissions",
  "data": null
}
```

***

### 4. 业务逻辑说明

#### 4.1 参赛人数统计

- `user_count` 字段应根据用户表实时统计更新
- 建议实现方式：
  - 用户注册/加入学校时：`user_count + 1`
  - 用户退出学校时：`user_count - 1`
  - 定时任务：每小时校准一次，从用户表统计实际人数

#### 4.2 软删除机制

- 删除学校时设置 `deleted_at` 为当前时间
- 所有查询接口默认过滤 `deleted_at IS NULL` 的记录
- 管理员可查看已删除学校（添加 `include_deleted=true` 参数）
- 已删除学校的参赛用户数据保留，但不可再参赛

#### 4.3 学校状态管理

- **active（正常）**：学校可正常参赛，用户可加入
- **suspended（封禁）**：
  - 学校暂停参赛资格
  - 已加入的用户无法提交答案
  - 新用户无法加入该学校
  - 可由管理员解除封禁

#### 4.4 负责人权限

- 院校负责人应具有以下权限：
  - 查看本校所有参赛人员
  - 管理本校参赛队伍
  - 查看本校比赛成绩统计
  - **审核本校学生提交的附加信息**（如学号、班级、身份证明等）
  - 不可修改学校基本信息（由平台管理员管理）

#### 4.5 Token 认证机制

- **Token 生成**：
  - 用户登录成功后，系统生成 JWT Token
  - Token 中包含用户信息（user_id、role、school_id 等）
  - Token 有效期为 24 小时

- **JWT Token 结构**：
  
  ```json
  {
    "header": {
      "alg": "HS256",
      "typ": "JWT"
    },
    "payload": {
      "user_id": 10001,
      "username": "zhangsan",
      "role": "school_admin",  // 角色：admin(平台管理员) / school_admin(院校负责人) / user(普通用户)
      "school_id": 1,          // 所属学校ID
      "exp": 1700000000,       // 过期时间戳
      "iat": 1699913600        // 签发时间戳
    },
    "signature": "..."
  }
  ```

- **Token 验证**：
  - 需要权限的接口必须验证 Token
  - 验证内容：Token 格式、签名、有效期、用户权限
  - 验证失败返回对应错误码（401/402/403）

- **Token 刷新**：
  - Token 过期后用户需重新登录
  - 可选：实现刷新 Token 机制（refresh_token）

- **权限级别**：
  - **公开**：无需 Token，任何人可访问
  - **用户**：需要 Token，登录用户可访问
  - **院校负责人**：需要 Token 且角色为院校负责人
  - **管理员**：需要 Token 且角色为平台管理员

#### 4.6 学生附加信息审核流程

> 这是院校负责人的核心功能之一

**审核场景**：
- 学生注册时填写学校、学号、班级等附加信息
- 学生提交后，信息状态为"待审核"
- 院校负责人登录后可查看待审核列表
- 负责人审核通过后，学生才能正常参赛
- 负责人可驳回并要求学生重新提交

**审核流程**：
1. 学生提交附加信息 → 状态：`pending`
2. 院校负责人审核 → 通过：`approved` / 驳回：`rejected`
3. 如被驳回，学生可修改后重新提交
4. 只有审核通过的学生才能参加比赛

**相关API**（后续需要在用户模块中实现）：
- `GET /api/v1/schools/:id/pending-users` - 查看待审核学生列表
- `POST /api/v1/schools/users/:user_id/review` - 审核学生信息
- `GET /api/v1/schools/:id/users` - 查看本校所有学生

***



## 二、用户模块

### 1. 数据库设计

***

**数据库：**`dalictf`

**表名：**`dalictf_user`

***

#### 字段设计

|      字段名      |                 类型                  |                约束条件                 |                            说明                            |
| :--------------: | :-----------------------------------: | :-------------------------------------: | :--------------------------------------------------------: |
|        id        |            `BIGINT(20)`             |      `PRIMARY KEY AUTO_INCREMENT`       |                          用户主键                          |
|     username     |           `VARCHAR(50)`           |            `NOT NULL UNIQUE`            |                   用户昵称（登录账号）                   |
|     password     |          `VARCHAR(255)`           |               `NOT NULL`                |                密码（加密存储，建议bcrypt）                |
|      email       |          `VARCHAR(100)`           |            `NOT NULL UNIQUE`            |                          注册邮箱                          |
|       role       | `ENUM('user', 'school_admin', 'admin', 'super_admin')` |      `NOT NULL DEFAULT 'user'`      | 用户角色：user-普通用户/school_admin-院校负责人/admin-管理员/super_admin-究极管理员 |
|      track       | `ENUM('social', 'school')` |          `NOT NULL DEFAULT 'social'`          |      参赛赛道：social-社会赛道/school-联合院校赛道      |
|    school_id     |            `BIGINT(20)`             |            `DEFAULT NULL`             |      所属学校ID（外键关联学校表，社会赛道用户为NULL）      |
|   school_name    |          `VARCHAR(255)`           |            `DEFAULT NULL`             |              所属院校名称（冗余字段，便于查询）              |
|    user_name     |           `VARCHAR(50)`           |            `DEFAULT NULL`             |         用户真实姓名（联合院校赛道必填，社会赛道可选）         |
| student_number |          `VARCHAR(50)`          |            `DEFAULT NULL`             |              学号（联合院校赛道必填）              |
|  school_grade  |          `VARCHAR(10)`          |            `DEFAULT NULL`             |         年级（如"2022"或"22"，联合院校赛道必填）         |
| student_nature | `ENUM('undergraduate', 'graduate')` |            `DEFAULT NULL`             | 学生性质：undergraduate-本科生/graduate-研究生（联合院校赛道必填） |
| email_verified |            `TINYINT(1)`             |       `NOT NULL DEFAULT 0`        |          邮箱验证状态：0-未验证/1-已验证          |
| email_verify_code |          `VARCHAR(10)`          |            `DEFAULT NULL`             |          邮箱验证码（6位数字）          |
| verify_code_expires_at |            `DATETIME`             |            `DEFAULT NULL`             |          验证码过期时间          |
| register_fail_count |            `INT(11)`            |        `NOT NULL DEFAULT 0`         | 注册失败次数（最多3次） |
|  verify_status   |   `ENUM('pending', 'approved', 'rejected')`   | `NOT NULL DEFAULT 'pending'` | 学生信息审核状态：pending-待审核/approved-已通过/rejected-已驳回（仅联合院校赛道） |
|  verify_reason   |          `VARCHAR(500)`           |            `DEFAULT NULL`             |               审核备注（驳回原因或其他说明）               |
|   verified_by    |            `BIGINT(20)`             |            `DEFAULT NULL`             |           审核人ID（关联用户表，记录审核的负责人）           |
|   verified_at    |            `DATETIME`             |            `DEFAULT NULL`             |                          审核时间                          |
|      status      |   `ENUM('active', 'suspended')`   |    `NOT NULL DEFAULT 'active'`    |   用户状态：active-正常/suspended-封禁   |
| last_login_time |            `DATETIME`             |            `DEFAULT NULL`             |                        最后登录时间                        |
|   last_login_ip   |           `VARCHAR(50)`           |            `DEFAULT NULL`             |                        最后登录IP                        |
|    created_at    |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                          创建时间                          |
|    updated_at    |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` |                          更新时间                          |
|    deleted_at    |            `DATETIME`             |            `DEFAULT NULL`             |               软删除时间（NULL表示未删除）               |

#### **建表 `sql`**

```sql
CREATE TABLE `dalictf_user` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '用户主键ID',
  `username` VARCHAR(50) NOT NULL COMMENT '用户昵称（登录账号）',
  `password` VARCHAR(255) NOT NULL COMMENT '密码（bcrypt加密）',
  `email` VARCHAR(100) NOT NULL COMMENT '注册邮箱',
  `role` ENUM('user', 'school_admin', 'admin', 'super_admin') NOT NULL DEFAULT 'user' COMMENT '用户角色',
  `track` ENUM('social', 'school') NOT NULL DEFAULT 'social' COMMENT '参赛赛道：social-社会/school-联合院校',
  `school_id` BIGINT(20) DEFAULT NULL COMMENT '所属学校ID（外键）',
  `school_name` VARCHAR(255) DEFAULT NULL COMMENT '所属院校名称',
  `user_name` VARCHAR(50) DEFAULT NULL COMMENT '用户真实姓名',
  `student_number` VARCHAR(50) DEFAULT NULL COMMENT '学号',
  `school_grade` VARCHAR(10) DEFAULT NULL COMMENT '年级（如2022或22）',
  `student_nature` ENUM('undergraduate', 'graduate') DEFAULT NULL COMMENT '学生性质：undergraduate-本科/graduate-研究生',
  `email_verified` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '邮箱验证状态：0-未验证/1-已验证',
  `email_verify_code` VARCHAR(10) DEFAULT NULL COMMENT '邮箱验证码',
  `verify_code_expires_at` DATETIME DEFAULT NULL COMMENT '验证码过期时间',
  `register_fail_count` INT(11) NOT NULL DEFAULT 0 COMMENT '注册失败次数',
  `verify_status` ENUM('pending', 'approved', 'rejected') NOT NULL DEFAULT 'pending' COMMENT '学生信息审核状态',
  `verify_reason` VARCHAR(500) DEFAULT NULL COMMENT '审核备注',
  `verified_by` BIGINT(20) DEFAULT NULL COMMENT '审核人ID',
  `verified_at` DATETIME DEFAULT NULL COMMENT '审核时间',
  `status` ENUM('active', 'suspended') NOT NULL DEFAULT 'active' COMMENT '用户状态',
  `last_login_time` DATETIME DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` VARCHAR(50) DEFAULT NULL COMMENT '最后登录IP',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_school_id` (`school_id`),
  KEY `idx_track` (`track`),
  KEY `idx_verify_status` (`verify_status`),
  KEY `idx_status` (`status`),
  KEY `idx_role` (`role`),
  KEY `idx_email_verified` (`email_verified`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 添加外键约束
ALTER TABLE `dalictf_user` 
ADD CONSTRAINT `fk_user_school` 
FOREIGN KEY (`school_id`) REFERENCES `dalictf_school`(`id`) 
ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `dalictf_user` 
ADD CONSTRAINT `fk_verified_by` 
FOREIGN KEY (`verified_by`) REFERENCES `dalictf_user`(`id`) 
ON DELETE SET NULL ON UPDATE CASCADE;
```

#### 字段说明补充

**关于赛道（track）字段：**
- **社会赛道（social）**：
  - 只需填写：`username`、`password`、`email`
  - 其他学生信息字段均为 NULL
  - `verify_status` 默认为 `pending` 但无需审核，注册即可参赛
  
- **联合院校赛道（school）**：
  - 必填字段：`username`、`password`、`email`、`school_id`、`school_name`、`user_name`、`student_number`、`school_grade`、`student_nature`
  - 注册流程：填写信息 → 获取邮件验证码 → 验证邮箱 → 提交注册申请 → 等待院校负责人审核
  - 审核通过后 `verify_status` 变为 `approved`，才能正常参赛

**关于邮件验证机制：**
- **email_verified（邮箱验证状态）**：
  - 0 - 未验证（默认值）
  - 1 - 已验证（验证码验证通过）
  - 联合院校赛道注册时必须先验证邮箱
  
- **email_verify_code（邮箱验证码）**：
  - 6位数字验证码
  - 发送到用户填写的邮箱
  - 有效期5分钟（`verify_code_expires_at`字段记录过期时间）
  
- **register_fail_count（注册失败次数）**：
  - 同一邮箱最多允许失败3次
  - 每次审核被驳回时 +1
  - 达到3次后，该邮箱将被禁止再次注册

**关于审核状态（verify_status）：**
- **pending（待审核）**：用户提交注册申请，等待院校负责人审核
- **approved（已通过）**：院校负责人审核通过，用户可正常参赛
- **rejected（已驳回）**：信息有误被驳回，记录失败次数，发送邮件通知原因

**关于角色（role）字段：**
- **user（普通用户）**：参赛用户，默认角色
- **school_admin（院校负责人）**：管理本校学生审核
- **admin（管理员）**：平台管理员，可管理学校和用户
- **super_admin（究极管理员）**：最高权限，可赋予/取消管理员身份

***



### 2. API 设计

> **统一前缀：**`/api/v1/users` 或 `/api/v1/auth`
>
> **认证机制：**JWT Token（详见参赛学校模块说明）
>
> **统一请求头：**
>
> ```
> Content-Type: application/json
> Authorization: Bearer <JWT_TOKEN>  // 需要权限的接口必须携带
> ```
>
> **返回格式：**统一使用标准响应格式（同学校模块）

#### 用户注册（社会赛道）

- **URL：** `POST /api/v1/auth/register/social`

- **权限：**公开

- **请求头：**

  ```
  Content-Type: application/json
  // 公开接口，无需 Authorization
  ```

- **请求体：**

  ```json
  {
    "username": "zhangsan",
    "password": "Password123!",
    "email": "zhangsan@example.com"
  }
  ```

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Registration successful",
    "data": {
      "id": 10001,
      "username": "zhangsan",
      "email": "zhangsan@example.com",
      "role": "user",
      "track": "social",
      "status": "active",
      "created_at": "2025-11-20 20:00:00"
    }
  }
  ```

#### 发送邮件验证码（联合院校注册第一步）

- **URL：** `POST /api/v1/auth/send-verify-code`

- **权限：**公开

- **请求头：**

  ```
  Content-Type: application/json
  // 公开接口，无需 Authorization
  ```

- **请求体：**

  ```json
  {
    "email": "lisi@dali.edu.cn"
  }
  ```

- **说明**：
  - 系统生成6位数字验证码
  - 发送到指定邮箱
  - 验证码有效期5分钟
  - 同一邮箱60秒内只能发送一次

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Verification code sent successfully",
    "data": {
      "email": "lisi@dali.edu.cn",
      "expires_at": "2025-11-20 20:10:00"
    }
  }
  ```

#### 用户注册（联合院校赛道）

- **URL：** `POST /api/v1/auth/register/school`

- **权限：**公开

- **请求头：**

  ```
  Content-Type: application/json
  // 公开接口，无需 Authorization
  ```

- **请求体：**

  ```json
  {
    "username": "lisi",
    "password": "Password123!",
    "email": "lisi@dali.edu.cn",
    "email_verify_code": "123456",
    "school_id": 1,
    "school_name": "大理大学",
    "user_name": "李四",
    "student_number": "2022110101",
    "school_grade": "2022",
    "student_nature": "undergraduate"
  }
  ```

- **说明**：
  - 必须先通过 `send-verify-code` 接口获取验证码
  - 验证码必须在5分钟内使用
  - 学校只能从后台已导入的学校列表中选择
  - 注册成功后，系统提示"您的注册申请已提交，我们将在一个工作日内处理，请耐心等待"
  - 同时发送邮件通知用户注册申请已提交

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Registration request submitted successfully. We will process it within one business day",
    "data": {
      "id": 10002,
      "username": "lisi",
      "email": "lisi@dali.edu.cn",
      "role": "user",
      "track": "school",
      "school_name": "大理大学",
      "user_name": "李四",
      "verify_status": "pending",
      "status": "active",
      "created_at": "2025-11-20 20:05:00"
    }
  }
  ```

#### 用户登录

- **URL：** `POST /api/v1/auth/login`

- **权限：**公开

- **请求头：**

  ```
  Content-Type: application/json
  // 公开接口，无需 Authorization
  ```

- **请求体：**

  ```json
  {
    "username": "zhangsan",  // 用户名或邮箱
    "password": "Password123!"
  }
  ```

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Login successful",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 86400,  // Token有效期（秒）
      "user": {
        "id": 10001,
        "username": "zhangsan",
        "email": "zhangsan@example.com",
        "role": "user",
        "track": "social",
        "status": "active"
      }
    }
  }
  ```

#### 获取当前用户信息

- **URL：** `GET /api/v1/users/me`

- **权限：**需要登录

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求体：**无

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "id": 10002,
      "username": "lisi",
      "email": "lisi@dali.edu.cn",
      "role": "user",
      "track": "school",
      "school_id": 1,
      "school_name": "大理大学",
      "user_name": "李四",
      "student_number": "2022110101",
      "school_grade": "2022",
      "student_nature": "undergraduate",
      "verify_status": "approved",
      "status": "active",
      "last_login_time": "2025-11-20 21:00:00",
      "created_at": "2025-11-20 20:05:00"
    }
  }
  ```

#### 修改用户信息

- **URL：** `PUT /api/v1/users/me`

- **权限：**需要登录

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求体：**

  ```json
  {
    "email": "newemail@example.com",  // 可选
    "user_name": "李四改",            // 可选，联合院校用户修改后需重新审核
    "student_number": "2022110102"    // 可选，联合院校用户修改后需重新审核
  }
  ```

- **说明**：
  - 社会赛道用户可修改：email
  - 联合院校用户修改学生信息（user_name、student_number、school_grade、student_nature）后，`verify_status` 将重置为 `pending`，需重新审核

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "User information updated successfully",
    "data": {
      "id": 10002,
      "username": "lisi",
      "email": "newemail@example.com",
      "verify_status": "pending",
      "updated_at": "2025-11-20 21:30:00"
    }
  }
  ```

#### 修改密码

- **URL：** `PUT /api/v1/users/me/password`

- **权限：**需要登录

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求体：**

  ```json
  {
    "old_password": "Password123!",
    "new_password": "NewPassword456!"
  }
  ```

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Password changed successfully",
    "data": null
  }
  ```

#### 查询用户列表（管理员）

- **URL：** `GET /api/v1/users`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求参数（Query）：**
  - `page` 页码（默认 1）
  - `limit` 每页数量（默认 20）
  - `search` 模糊搜索（用户名、邮箱、真实姓名）
  - `track` 赛道筛选（social/school）
  - `role` 角色筛选（user/school_admin/admin）
  - `status` 用户状态（active/suspended）
  - `verify_status` 审核状态（pending/approved/rejected）
  - `school_id` 学校ID筛选
  - `sort_by` 排序字段（created_at/last_login_time）
  - `order` 排序方式（asc/desc，默认 desc）

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total": 500,
      "page": 1,
      "limit": 20,
      "list": [
        {
          "id": 10001,
          "username": "zhangsan",
          "email": "zhangsan@example.com",
          "role": "user",
          "track": "social",
          "status": "active",
          "last_login_time": "2025-11-20 21:00:00",
          "created_at": "2025-11-20 20:00:00"
        },
        {
          "id": 10002,
          "username": "lisi",
          "email": "lisi@dali.edu.cn",
          "role": "user",
          "track": "school",
          "school_name": "大理大学",
          "user_name": "李四",
          "verify_status": "approved",
          "status": "active",
          "last_login_time": "2025-11-20 20:30:00",
          "created_at": "2025-11-20 20:05:00"
        }
      ]
    }
  }
  ```

#### 查询单个用户详情（管理员）

- **URL：** `GET /api/v1/users/:id`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 用户ID

- **请求体：**无

- **成功返回：**（返回完整用户信息）

#### 修改用户状态（管理员）

- **URL：** `PATCH /api/v1/users/:id/status`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 用户ID

- **请求体：**

  ```json
  {
    "status": "suspended"  // active 或 suspended
  }
  ```

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "User status updated successfully",
    "data": {
      "id": 10001,
      "username": "zhangsan",
      "status": "suspended",
      "updated_at": "2025-11-20 22:00:00"
    }
  }
  ```

#### 删除用户（管理员，软删除）

- **URL：** `DELETE /api/v1/users/:id`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 用户ID

- **请求体：**无

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "User deleted successfully",
    "data": null
  }
  ```

#### 查看待审核学生列表（院校负责人）

- **URL：** `GET /api/v1/schools/:school_id/pending-users`

- **权限：**院校负责人权限（且只能查看自己所属学校的）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `school_id` 学校ID

- **请求参数（Query）：**
  - `page` 页码（默认 1）
  - `limit` 每页数量（默认 20）

- **请求体：**无

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total": 15,
      "page": 1,
      "limit": 20,
      "list": [
        {
          "id": 10002,
          "username": "lisi",
          "email": "lisi@dali.edu.cn",
          "user_name": "李四",
          "student_number": "2022110101",
          "school_grade": "2022",
          "student_nature": "undergraduate",
          "verify_status": "pending",
          "created_at": "2025-11-20 20:05:00"
        }
      ]
    }
  }
  ```

#### 审核学生信息（院校负责人）

- **URL：** `POST /api/v1/schools/users/:user_id/review`

- **权限：**院校负责人权限（且只能审核自己所属学校的学生）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `user_id` 用户ID

- **请求体：**

  ```json
  {
    "verify_status": "approved",  // approved 或 rejected
    "verify_reason": "信息核实无误"  // 驳回时必填，通过时可选
  }
  ```

- **说明**：
  - **审核通过**：
    - 发送邮件通知用户注册成功
    - 邮件内容包含：用户名、学校信息、登录链接
  - **审核驳回**：
    - `register_fail_count` 自动 +1
    - 发送邮件通知用户驳回原因
    - 如果失败次数达到3次，禁止该邮箱再次注册
    - 邮件内容包含：驳回原因、剩余可尝试次数

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Review completed successfully",
    "data": {
      "id": 10002,
      "username": "lisi",
      "user_name": "李四",
      "verify_status": "approved",
      "verify_reason": "信息核实无误",
      "verified_by": 10005,
      "verified_at": "2025-11-20 22:30:00",
      "register_fail_count": 0
    }
  }
  ```

#### 查看本校所有学生（院校负责人）

- **URL：** `GET /api/v1/schools/:school_id/users`

- **权限：**院校负责人权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `school_id` 学校ID

- **请求参数（Query）：**
  - `page` 页码（默认 1）
  - `limit` 每页数量（默认 20）
  - `verify_status` 审核状态筛选（pending/approved/rejected）
  - `search` 模糊搜索（姓名、学号）

- **请求体：**无

- **成功返回：**（同待审核列表，但包含所有状态的学生）

#### 赋予管理员权限（究极管理员）

- **URL：** `POST /api/v1/users/:user_id/grant-admin`

- **权限：**究极管理员权限（super_admin）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `user_id` 用户ID

- **请求体：**

  ```json
  {
    "role": "admin"  // 可选：admin 或 school_admin
  }
  ```

- **说明**：
  - 只有究极管理员可以调用此接口
  - 可以将普通用户提升为管理员或院校负责人
  - 不能将用户提升为究极管理员
  - 操作记录将被记录到日志

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Admin permission granted successfully",
    "data": {
      "id": 10001,
      "username": "zhangsan",
      "role": "admin",
      "updated_at": "2025-11-20 23:00:00"
    }
  }
  ```

#### 取消管理员权限（究极管理员）

- **URL：** `POST /api/v1/users/:user_id/revoke-admin`

- **权限：**究极管理员权限（super_admin）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `user_id` 用户ID

- **请求体：**无

- **说明**：
  - 只有究极管理员可以调用此接口
  - 将管理员或院校负责人降级为普通用户
  - 不能取消究极管理员的权限
  - 操作记录将被记录到日志

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Admin permission revoked successfully",
    "data": {
      "id": 10001,
      "username": "zhangsan",
      "role": "user",
      "updated_at": "2025-11-20 23:05:00"
    }
  }
  ```

***

### 3. 错误码定义

| 错误码 |                  说明                  |
| :----: | :------------------------------------: |
|  200   |                  成功                  |
|  400   |      请求参数错误（如必填字段缺失）      |
|  401   |      未授权（Token缺失或格式错误）      |
|  402   |         Token已过期，请重新登录         |
|  403   |        无权限操作（权限不足）        |
|  404   |             用户不存在             |
|  409   |   用户名或邮箱已存在   |
|  410   |      用户账号已被封禁      |
|  411   | 学生信息未通过审核，无法参赛 |
|  412   |       密码错误       |
|  413   |       旧密码错误       |
|  414   |    邮件验证码错误或已过期    |
|  415   | 注册失败次数已达上限（3次），该邮箱已被禁止注册 |
|  416   |    验证码发送过于频繁，请60秒后再试    |
|  417   |  邮箱未验证，请先获取验证码  |
|  418   |  只有究极管理员可以执行此操作  |
|  419   |  不能操作究极管理员账户  |
|  500   |              服务器错误              |

#### 用户相关错误返回示例

```json
{
  "code": 409,
  "msg": "Username already exists",
  "data": null
}
```

```json
{
  "code": 410,
  "msg": "User account has been suspended",
  "data": null
}
```

```json
{
  "code": 411,
  "msg": "Student information is pending verification, cannot participate in competition",
  "data": {
    "verify_status": "pending"
  }
}
```

***

### 4. 业务逻辑说明

#### 4.1 用户注册流程

**社会赛道注册：**
1. 用户填写基本信息（用户名、密码、邮箱）
2. 系统验证信息格式和唯一性
3. 密码使用 bcrypt 加密存储
4. 注册成功，`track` 设为 `social`，`verify_status` 为 `pending`（但无需审核）
5. 用户可立即登录参赛

**联合院校赛道注册（完整流程）：**
1. **第一步：填写注册信息**
   - 用户填写完整学生信息（用户名、密码、邮箱、学校、真实姓名、学号、年级、学生性质）
   - 学校只能从后台已导入的学校列表中选择
   - 系统验证学校是否存在且状态为 `active`
   - 检查该邮箱是否已达到失败次数上限（3次）
   
2. **第二步：邮箱验证**
   - 用户点击"获取邮件验证码"
   - 系统检查60秒内是否已发送过验证码（防止频繁发送）
   - 生成6位数字验证码，设置5分钟过期时间
   - 发送验证码到用户邮箱
   - 用户输入验证码并验证
   
3. **第三步：提交注册申请**
   - 验证码验证通过后，`email_verified` 设为 1
   - 用户提交注册申请
   - 系统创建用户记录，`verify_status` 为 `pending`
   - 显示提示："您的注册申请已提交，我们将在一个工作日内处理，请耐心等待"
   - 发送邮件通知用户注册申请已提交
   - 对应学校的 `user_count` 自动 +1
   
4. **第四步：院校负责人审核**
   - 院校负责人在管理面板查看待审核学生列表
   - 核实学生信息真实性
   - **审核通过**：
     - `verify_status` 变为 `approved`
     - 发送邮件通知用户注册成功，可以登录参赛
   - **审核驳回**：
     - `verify_status` 变为 `rejected`
     - `register_fail_count` 自动 +1
     - 发送邮件通知用户驳回原因和剩余可尝试次数
     - 如果失败次数达到3次，该邮箱将被永久禁止注册

#### 4.2 密码安全策略

- **密码强度要求**：
  - 最小长度 8 位
  - 必须包含大小写字母和数字
  - 建议包含特殊字符
  
- **密码存储**：
  - 使用 bcrypt 算法加密
  - 成本因子建议设置为 10-12
  - 永远不存储明文密码

- **密码修改**：
  - 必须验证旧密码
  - 新密码不能与旧密码相同
  - 修改密码后建议重新登录

#### 4.3 审核流程详解

**谁可以审核：**
- 院校负责人只能审核本校学生
- 平台管理员可以审核所有学校的学生

**审核规则：**
1. 学生注册后，`verify_status` 为 `pending`
2. 院校负责人查看待审核列表
3. 核实学生信息（姓名、学号、年级等）
4. 审核操作：
   - **通过**：`verify_status` 变为 `approved`，学生可参赛
   - **驳回**：`verify_status` 变为 `rejected`，需填写驳回原因
5. 记录审核人和审核时间
6. 学生被驳回后，可修改信息重新提交，状态重置为 `pending`

**审核权限验证：**
- 院校负责人必须满足：
  - `role` 为 `school_admin`
  - 审核的学生 `school_id` 与负责人的 `school_id` 一致
- 否则返回 403 权限不足

#### 4.4 学生信息修改机制

**社会赛道用户：**
- 可随时修改 `email`
- 修改后无需审核

**联合院校赛道用户：**
- 修改 `email`：无需重新审核
- 修改 `user_name`、`student_number`、`school_grade`、`student_nature`（学生关键信息）：
  - `verify_status` 自动重置为 `pending`
  - 需要院校负责人重新审核
  - 未通过审核期间无法参赛
  - 不影响 `register_fail_count`（失败次数只在注册驳回时增加）

#### 4.5 用户状态管理

**active（正常）：**
- 可以正常登录
- 可以参加比赛（联合院校需审核通过）
- 可以提交答案

**suspended（封禁）：**
- 可以登录但无法参赛
- 无法提交答案
- 无法修改信息（除密码外）
- 由平台管理员操作，需明确封禁原因

#### 4.6 角色权限说明

**user（普通用户）：**
- 查看和修改自己的信息
- 参加比赛（需满足审核条件）
- 提交答案
- 查看比赛题目和排行榜

**school_admin（院校负责人）：**
- 拥有普通用户的所有权限
- 查看本校所有学生信息
- 审核本校学生的附加信息
- 查看本校比赛统计数据
- 不能修改学校基本信息

**admin（平台管理员）：**
- 管理所有学校和用户
- 修改任何用户的状态（除了究极管理员）
- 查看所有数据和统计信息
- 管理比赛题目和配置
- **不能**赋予或取消其他用户的管理员权限
- **不能**修改究极管理员的任何信息

**super_admin（究极管理员）：**
- 拥有最高权限，凌驾于普通管理员之上
- 可以赋予用户管理员权限（admin 或 school_admin）
- 可以取消用户的管理员权限（降级为普通用户）
- **不能**将用户提升为究极管理员
- **不能**被普通管理员修改或删除
- 操作日志会被完整记录

**权限层级（从高到低）：**
```
super_admin（究极管理员）
    ↓ 可以赋予/取消
admin（平台管理员）
    ↓ 可以管理（部分权限）
school_admin（院校负责人）
    ↓ 可以审核
user（普通用户）
```

**权限约束规则：**
1. 究极管理员只能通过数据库直接创建，不能通过API创建
2. 管理员不能修改比自己权限高的用户
3. 管理员不能赋予或取消管理员权限
4. 院校负责人只能管理本校学生
5. 所有权限变更操作都需要记录日志

#### 4.7 登录逻辑

1. 用户提交用户名（或邮箱）和密码
2. 系统查询用户是否存在
3. 验证密码是否正确（bcrypt.compare）
4. 检查用户状态：
   - `status` 为 `suspended`：返回 410 账号已封禁
   - `deleted_at` 不为 NULL：返回 404 用户不存在
5. 生成 JWT Token，包含用户基本信息
6. 更新 `last_login_time` 和 `last_login_ip`
7. 返回 Token 和用户信息

#### 4.8 数据一致性保证

**学校人数统计：**
- 用户注册联合院校赛道时：`school.user_count + 1`
- 用户删除（软删除）时：`school.user_count - 1`
- 用户切换学校时：原学校 `-1`，新学校 `+1`
- 建议每天凌晨运行定时任务校准数据

**学校名称冗余字段：**
- 用户注册时自动填充 `school_name`
- 学校名称修改时，需同步更新用户表的 `school_name`
- 提高查询效率，避免每次都关联学校表

**软删除处理：**
- 用户删除时设置 `deleted_at`
- 所有查询默认过滤 `deleted_at IS NULL`
- 用户删除后，用户名和邮箱应允许重新注册（可考虑在唯一索引中排除已删除记录）

#### 4.9 邮件通知机制

**邮件验证码：**
- **主题**：【ISCTF】邮箱验证码
- **内容**：
  ```
  您好，
  
  您正在注册 ISCTF 竞赛平台账号，您的验证码是：123456
  
  验证码有效期为 5 分钟，请尽快完成验证。
  如果这不是您的操作，请忽略此邮件。
  
  ISCTF 竞赛平台
  ```

**注册申请已提交：**
- **主题**：【ISCTF】注册申请已提交
- **内容**：
  ```
  您好，{user_name}
  
  您的注册申请已成功提交！
  
  学校：{school_name}
  姓名：{user_name}
  学号：{student_number}
  
  我们将在一个工作日内完成审核，请耐心等待。
  审核结果将通过邮件通知您。
  
  ISCTF 竞赛平台
  ```

**审核通过通知：**
- **主题**：【ISCTF】注册审核通过
- **内容**：
  ```
  您好，{user_name}
  
  恭喜您！您的注册申请已通过审核。
  
  用户名：{username}
  学校：{school_name}
  
  您现在可以登录参加比赛了！
  登录地址：https://isctf.example.com/login
  
  祝您在比赛中取得好成绩！
  
  ISCTF 竞赛平台
  ```

**审核驳回通知：**
- **主题**：【ISCTF】注册审核未通过
- **内容**：
  ```
  您好，{user_name}
  
  很抱歉，您的注册申请未通过审核。
  
  驳回原因：{verify_reason}
  
  剩余可尝试次数：{3 - register_fail_count} 次
  
  {如果失败次数达到3次，追加：该邮箱已达到注册失败次数上限，无法再次注册。}
  
  如有疑问，请联系您所在院校的负责人。
  
  ISCTF 竞赛平台
  ```

**邮件发送注意事项：**
- 使用异步队列发送邮件，避免阻塞接口响应
- 记录邮件发送日志，包括发送时间、收件人、邮件类型、发送状态
- 邮件发送失败时应重试（最多3次）
- 建议使用专业的邮件服务（如阿里云邮件推送、SendGrid等）
- 邮件内容应支持HTML格式，提升用户体验

***



## 三、团队模块

### 1. 数据库设计

***

**数据库：**`dalictf`

**表名：**`dalictf_team`

***

#### 字段设计

|      字段名      |                 类型                  |                约束条件                 |                            说明                            |
| :--------------: | :-----------------------------------: | :-------------------------------------: | :--------------------------------------------------------: |
|        id        |            `BIGINT(20)`             |      `PRIMARY KEY AUTO_INCREMENT`       |                          团队主键                          |
|     team_name    |           `VARCHAR(100)`           |            `NOT NULL UNIQUE`            |                   团队名称（唯一）                   |
|  team_password   |          `VARCHAR(255)`           |               `NOT NULL`                |          团队密码（加密存储，用于成员加入团队）          |
|    captain_id    |            `BIGINT(20)`             |            `NOT NULL`             |      队长用户ID（外键关联用户表）      |
|   captain_name   |           `VARCHAR(50)`           |            `NOT NULL`             |              队长昵称（冗余字段，便于查询）              |
|    member1_id    |            `BIGINT(20)`             |            `DEFAULT NULL`             |      成员1用户ID（外键关联用户表）      |
|   member1_name   |           `VARCHAR(50)`           |            `DEFAULT NULL`             |              成员1昵称（冗余字段）              |
|    member2_id    |            `BIGINT(20)`             |            `DEFAULT NULL`             |      成员2用户ID（外键关联用户表）      |
|   member2_name   |           `VARCHAR(50)`           |            `DEFAULT NULL`             |              成员2昵称（冗余字段）              |
|    team_track    | `ENUM('social', 'freshman', 'advanced')` |      `NOT NULL`      | 团队赛道：social-社会赛道/freshman-新生赛道/advanced-进阶赛道 |
|    school_id     |            `BIGINT(20)`             |            `DEFAULT NULL`             |      所属学校ID（联合院校赛道必填，外键关联学校表）      |
|   school_name    |          `VARCHAR(255)`           |            `DEFAULT NULL`             |              所属院校名称（冗余字段，便于查询）              |
|    team_score    |            `INT(11)`            |        `NOT NULL DEFAULT 0`         |         团队总得分         |
|  member_count  |            `TINYINT(2)`             |        `NOT NULL DEFAULT 1`         | 团队成员数量（1-3人） |
|      status      |   `ENUM('active', 'disbanded', 'banned')`   |    `NOT NULL DEFAULT 'active'`    |   团队状态：active-正常/disbanded-已解散/banned-作弊封禁   |
|    created_at    |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                          创建时间                          |
|    updated_at    |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` |                          更新时间                          |
|    deleted_at    |            `DATETIME`             |            `DEFAULT NULL`             |               软删除时间（NULL表示未删除）               |

#### **建表 `sql`**

```sql
CREATE TABLE `dalictf_team` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '团队主键ID',
  `team_name` VARCHAR(100) NOT NULL COMMENT '团队名称',
  `team_password` VARCHAR(255) NOT NULL COMMENT '团队密码（加密存储）',
  `captain_id` BIGINT(20) NOT NULL COMMENT '队长用户ID',
  `captain_name` VARCHAR(50) NOT NULL COMMENT '队长昵称',
  `member1_id` BIGINT(20) DEFAULT NULL COMMENT '成员1用户ID',
  `member1_name` VARCHAR(50) DEFAULT NULL COMMENT '成员1昵称',
  `member2_id` BIGINT(20) DEFAULT NULL COMMENT '成员2用户ID',
  `member2_name` VARCHAR(50) DEFAULT NULL COMMENT '成员2昵称',
  `team_track` ENUM('social', 'freshman', 'advanced') NOT NULL COMMENT '团队赛道',
  `school_id` BIGINT(20) DEFAULT NULL COMMENT '所属学校ID',
  `school_name` VARCHAR(255) DEFAULT NULL COMMENT '所属院校名称',
  `team_score` INT(11) NOT NULL DEFAULT 0 COMMENT '团队总得分',
  `member_count` TINYINT(2) NOT NULL DEFAULT 1 COMMENT '团队成员数量',
  `status` ENUM('active', 'disbanded', 'banned') NOT NULL DEFAULT 'active' COMMENT '团队状态',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_team_name` (`team_name`),
  KEY `idx_captain_id` (`captain_id`),
  KEY `idx_school_id` (`school_id`),
  KEY `idx_team_track` (`team_track`),
  KEY `idx_status` (`status`),
  KEY `idx_team_score` (`team_score`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='团队表';

-- 添加外键约束
ALTER TABLE `dalictf_team` 
ADD CONSTRAINT `fk_team_captain` 
FOREIGN KEY (`captain_id`) REFERENCES `dalictf_user`(`id`) 
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `dalictf_team` 
ADD CONSTRAINT `fk_team_member1` 
FOREIGN KEY (`member1_id`) REFERENCES `dalictf_user`(`id`) 
ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `dalictf_team` 
ADD CONSTRAINT `fk_team_member2` 
FOREIGN KEY (`member2_id`) REFERENCES `dalictf_user`(`id`) 
ON DELETE SET NULL ON UPDATE CASCADE;

ALTER TABLE `dalictf_team` 
ADD CONSTRAINT `fk_team_school` 
FOREIGN KEY (`school_id`) REFERENCES `dalictf_school`(`id`) 
ON DELETE SET NULL ON UPDATE CASCADE;
```

#### 字段说明补充

**关于团队赛道（team_track）自动判断规则：**

**社会赛道（social）：**
- 队长为社会赛道用户（`track = 'social'`）
- 所有成员必须都是社会赛道用户
- `school_id` 和 `school_name` 为 NULL

**新生赛道（freshman）：**
- 队长为联合院校赛道且年级为大一或研一
- 判断规则：
  - `track = 'school'`
  - `school_grade` 为当前年份或当前年份-1（如2024或2023表示大一）
  - `student_nature = 'undergraduate'` 且 `school_grade` 为大一年级
  - `student_nature = 'graduate'` 且 `school_grade` 为研一年级
- **重要限制**：
  - 大一学生不能与研一学生组队
  - 必须来自同一学校
  - 不能与进阶赛道成员组队

**进阶赛道（advanced）：**
- 队长为联合院校赛道且年级非大一、研一
- 包括大二、大三、大四、研二、研三等
- 必须来自同一学校
- 不能与新生赛道成员组队

**组队限制规则总结：**
1. 一个用户同一时间只能加入一个团队
2. 团队最多3人（1队长+2成员）
3. 社会赛道 ≠ 联合院校赛道（不可混合组队）
4. 新生赛道 ≠ 进阶赛道（不可混合组队）
5. 大一 ≠ 研一（新生赛道内也不可混合）
6. 联合院校成员必须同一学校（不可跨校组队）
7. 团队赛道由队长决定，成员必须符合队长赛道要求

**关于团队状态（status）：**
- **active（正常）**：团队正常运行，可以参赛、查看题目、提交答案
- **disbanded（已解散）**：队长解散团队，成员自动退出，团队数据保留但不可再使用
- **banned（作弊封禁）**：团队因作弊被管理员封禁，无法查看题目、无法提交答案、无法参与排行榜

**关于成员数量（member_count）：**
- 创建团队时默认为 1（只有队长）
- 成员加入时自动 +1
- 成员退出时自动 -1
- 最小值 1，最大值 3

***

### 2. API 设计

> **统一前缀：**`/api/v1/teams`
>
> **认证机制：**JWT Token（详见参赛学校模块说明）
>
> **通用请求头：**
> ```
> Content-Type: application/json
> Authorization: Bearer <JWT_TOKEN>
> ```

***

#### 创建团队

- **URL：** `POST /api/v1/teams`

- **权限：**普通用户（user及以上）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求体：**

  ```json
  {
    "team_name": "银河战队",
    "team_password": "Galaxy@2024"
  }
  ```

- **说明**：
  - 用户创建团队时自动成为队长
  - 系统根据队长的 `track`、`school_grade`、`student_nature` 自动判断团队赛道
  - 联合院校用户创建的团队自动关联学校信息
  - 一个用户同一时间只能加入一个团队（包括作为队长）
  - 团队密码使用 bcrypt 加密存储

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Team created successfully",
    "data": {
      "id": 1001,
      "team_name": "银河战队",
      "captain_id": 10001,
      "captain_name": "zhangsan",
      "team_track": "social",
      "school_id": null,
      "school_name": null,
      "team_score": 0,
      "member_count": 1,
      "status": "active",
      "created_at": "2025-11-20 21:00:00"
    }
  }
  ```

#### 加入团队

- **URL：** `POST /api/v1/teams/:team_id/join`

- **权限：**普通用户（user及以上）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `team_id` 团队ID

- **请求体：**

  ```json
  {
    "team_password": "Galaxy@2024"
  }
  ```

- **说明**：
  - 需要提供正确的团队密码
  - 系统自动验证用户是否符合该团队的赛道要求
  - 验证规则：
    - 检查团队是否已满（3人）
    - 检查用户是否已在其他团队
    - 检查赛道是否匹配（社会/新生/进阶）
    - 检查学校是否一致（联合院校）
    - 检查学生性质是否一致（新生赛道：大一≠研一）
  - 加入成功后，`member_count` 自动 +1

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Joined team successfully",
    "data": {
      "id": 1001,
      "team_name": "银河战队",
      "captain_name": "zhangsan",
      "member1_name": "lisi",
      "member2_name": null,
      "team_track": "social",
      "member_count": 2,
      "your_position": "member1"
    }
  }
  ```

#### 退出团队

- **URL：** `POST /api/v1/teams/:team_id/quit`

- **权限：**普通用户（team成员）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `team_id` 团队ID

- **请求体：**无

- **说明**：
  - 只有队员可以退出，队长不能退出（只能解散团队）
  - 退出后，该位置设为 NULL，昵称也清空
  - `member_count` 自动 -1
  - 如果 member2 退出，直接清空 member2 的位置
  - 如果 member1 退出且 member2 存在，则 member2 自动移动到 member1 位置

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Quit team successfully",
    "data": null
  }
  ```

#### 解散团队（队长）

- **URL：** `POST /api/v1/teams/:team_id/disband`

- **权限：**队长

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `team_id` 团队ID

- **请求体：**无

- **说明**：
  - 只有队长可以解散团队
  - 解散后，`status` 变为 `disbanded`
  - 所有成员自动退出
  - 团队数据保留但不可再使用
  - 解散后的团队不会在列表中显示

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Team disbanded successfully",
    "data": null
  }
  ```

#### 踢出成员（队长）

- **URL：** `POST /api/v1/teams/:team_id/kick/:user_id`

- **权限：**队长

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `team_id` 团队ID
  - `user_id` 要踢出的用户ID

- **请求体：**无

- **说明**：
  - 只有队长可以踢出成员
  - 不能踢出队长自己
  - 踢出后，该位置设为 NULL
  - `member_count` 自动 -1
  - 位置调整规则同退出团队

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Member kicked out successfully",
    "data": {
      "team_id": 1001,
      "team_name": "银河战队",
      "member_count": 1
    }
  }
  ```

#### 修改团队密码（队长）

- **URL：** `PUT /api/v1/teams/:team_id/password`

- **权限：**队长

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `team_id` 团队ID

- **请求体：**

  ```json
  {
    "old_password": "Galaxy@2024",
    "new_password": "NewGalaxy@2024"
  }
  ```

- **说明**：
  - 只有队长可以修改团队密码
  - 必须提供正确的旧密码
  - 新密码使用 bcrypt 加密存储

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Team password updated successfully",
    "data": null
  }
  ```

#### 查询我的团队

- **URL：** `GET /api/v1/teams/my`

- **权限：**普通用户（user及以上）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求参数：**无

- **说明**：
  - 返回当前用户所在的团队信息
  - 如果用户未加入任何团队，返回 null

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "id": 1001,
      "team_name": "银河战队",
      "captain_id": 10001,
      "captain_name": "zhangsan",
      "member1_id": 10002,
      "member1_name": "lisi",
      "member2_id": null,
      "member2_name": null,
      "team_track": "social",
      "school_id": null,
      "school_name": null,
      "team_score": 1500,
      "member_count": 2,
      "status": "active",
      "created_at": "2025-11-20 21:00:00",
      "your_role": "captain"  // captain, member1, member2
    }
  }
  ```

#### 查询团队详情

- **URL：** `GET /api/v1/teams/:team_id`

- **权限：**公开（任何人都可以查看）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>  // 可选
  ```

- **路径参数：**
  - `team_id` 团队ID

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "id": 1001,
      "team_name": "银河战队",
      "captain_name": "zhangsan",
      "member1_name": "lisi",
      "member2_name": "wangwu",
      "team_track": "freshman",
      "school_name": "大理大学",
      "team_score": 2500,
      "member_count": 3,
      "created_at": "2025-11-20 21:00:00"
    }
  }
  ```

#### 查询团队列表

- **URL：** `GET /api/v1/teams`

- **权限：**公开

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>  // 可选
  ```

- **请求参数（Query）：**
  - `page` 页码（默认 1）
  - `limit` 每页数量（默认 20）
  - `team_track` 赛道筛选（social/freshman/advanced）
  - `school_id` 学校ID筛选
  - `search` 模糊搜索（团队名称）
  - `order_by` 排序字段（score/created_at，默认 score）

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total": 150,
      "page": 1,
      "limit": 20,
      "list": [
        {
          "id": 1001,
          "team_name": "银河战队",
          "captain_name": "zhangsan",
          "member_count": 3,
          "team_track": "social",
          "school_name": null,
          "team_score": 3500,
          "created_at": "2025-11-20 21:00:00"
        },
        {
          "id": 1002,
          "team_name": "代码骑士",
          "captain_name": "lisi",
          "member_count": 2,
          "team_track": "freshman",
          "school_name": "大理大学",
          "team_score": 3200,
          "created_at": "2025-11-20 21:10:00"
        }
      ]
    }
  }
  ```

#### 团队排行榜

- **URL：** `GET /api/v1/teams/rankings`

- **权限：**公开

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>  // 可选
  ```

- **请求参数（Query）：**
  - `team_track` 赛道筛选（social/freshman/advanced，必选）
  - `limit` 返回数量（默认 50，最大 100）

- **说明**：
  - 按照 `team_score` 降序排列
  - 只显示 `status = 'active'` 且 `deleted_at IS NULL` 的团队
  - 分赛道进行排名

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "track": "social",
      "total": 50,
      "list": [
        {
          "rank": 1,
          "team_id": 1001,
          "team_name": "银河战队",
          "captain_name": "zhangsan",
          "member_count": 3,
          "team_score": 5000
        },
        {
          "rank": 2,
          "team_id": 1003,
          "team_name": "网络精英",
          "captain_name": "wangwu",
          "member_count": 3,
          "team_score": 4800
        }
      ]
    }
  }
  ```

#### 封禁团队（管理员）

- **URL：** `POST /api/v1/teams/:team_id/ban`

- **权限：**管理员（admin及以上）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `team_id` 团队ID

- **请求体：**

  ```json
  {
    "reason": "使用非法工具作弊"
  }
  ```

- **说明**：
  - 只有管理员和究极管理员可以封禁团队
  - 封禁后，团队 `status` 变为 `banned`
  - 被封禁的团队无法查看题目、提交答案
  - 被封禁的团队不参与排行榜
  - 必须提供封禁原因
  - 封禁操作记录到日志

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Team banned successfully",
    "data": {
      "team_id": 1001,
      "team_name": "银河战队",
      "status": "banned",
      "ban_reason": "使用非法工具作弊",
      "banned_at": "2025-11-20 23:00:00"
    }
  }
  ```

#### 解禁团队（管理员）

- **URL：** `POST /api/v1/teams/:team_id/unban`

- **权限：**管理员（admin及以上）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `team_id` 团队ID

- **请求体：**无

- **说明**：
  - 只有管理员和究极管理员可以解禁团队
  - 解禁后，团队 `status` 恢复为 `active`
  - 团队可以正常参赛
  - 解禁操作记录到日志

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Team unbanned successfully",
    "data": {
      "team_id": 1001,
      "team_name": "银河战队",
      "status": "active",
      "unbanned_at": "2025-11-20 23:30:00"
    }
  }
  ```

***

### 3. 错误码定义

| 错误码 |                  说明                  |
| :----: | :------------------------------------: |
|  200   |                  成功                  |
|  400   |      请求参数错误（如必填字段缺失）      |
|  401   |      未授权（Token缺失或格式错误）      |
|  402   |         Token已过期，请重新登录         |
|  403   |        无权限操作（权限不足）        |
|  404   |             团队不存在             |
|  409   |   团队名称已存在   |
|  420   |      您已经加入了其他团队      |
|  421   |       团队密码错误       |
|  422   | 团队已满员（最多3人） |
|  423   |  赛道不匹配，无法加入该团队  |
|  424   |  学校不一致，无法跨校组队  |
|  425   | 学生性质不一致（大一不能与研一组队） |
|  426   |   只有队长可以执行此操作   |
|  427   |   队长不能退出团队，只能解散   |
|  428   |   团队已解散，无法操作   |
|  429   |   您不在该团队中   |
|  430   |   团队已被封禁，无法查看题目或提交答案   |
|  500   |              服务器错误              |

#### 团队相关错误返回示例

```json
{
  "code": 420,
  "msg": "You have already joined another team",
  "data": {
    "current_team_id": 1002,
    "current_team_name": "代码骑士"
  }
}
```

```json
{
  "code": 423,
  "msg": "Track mismatch, cannot join this team",
  "data": {
    "your_track": "social",
    "team_track": "freshman",
    "reason": "Social track users cannot join school track teams"
  }
}
```

```json
{
  "code": 425,
  "msg": "Student nature mismatch, freshmen undergraduates cannot team with freshmen graduates",
  "data": {
    "your_nature": "undergraduate",
    "your_grade": "2024",
    "team_requirement": "Freshman track requires same student nature"
  }
}
```

```json
{
  "code": 430,
  "msg": "Team has been banned for cheating, cannot view or submit challenges",
  "data": {
    "team_status": "banned",
    "ban_reason": "使用非法工具作弊",
    "banned_at": "2025-11-20 23:00:00"
  }
}
```

***

### 4. 业务逻辑说明

#### 4.1 团队赛道自动判断逻辑

创建团队时，系统根据队长信息自动判断团队赛道：

**判断流程：**
```
1. 检查队长的 track 字段
   ├─ track = 'social' → 团队赛道 = social
   └─ track = 'school' → 继续判断
       ├─ 检查 student_nature 和 school_grade
       ├─ student_nature = 'undergraduate' 且 school_grade 为大一 → freshman（本科新生）
       ├─ student_nature = 'graduate' 且 school_grade 为研一 → freshman（研究生新生）
       └─ 其他情况 → advanced（进阶赛道）
```

**年级判断规则（以2025年为例）：**
- 大一：`school_grade` 为 "2024" 或 "24"
- 研一：`school_grade` 为 "2024" 或 "24" 且 `student_nature = 'graduate'`
- 其他年级：如 "2023"、"2022"、"2021" 等

#### 4.2 加入团队验证规则

用户申请加入团队时，系统按以下顺序验证：

**1. 基础验证：**
- 团队是否存在且未解散
- 团队密码是否正确
- 团队是否已满（3人）
- 用户是否已在其他团队

**2. 赛道验证：**
- **社会赛道团队**：
  - 用户必须是社会赛道（`track = 'social'`）
  
- **新生赛道团队**：
  - 用户必须是联合院校学生（`track = 'school'`）
  - 必须与队长来自同一学校（`school_id` 一致）
  - 必须是大一或研一
  - **重要**：大一学生（undergraduate）不能与研一学生（graduate）组队
  
- **进阶赛道团队**：
  - 用户必须是联合院校学生（`track = 'school'`）
  - 必须与队长来自同一学校（`school_id` 一致）
  - 不能是大一或研一（必须是高年级）

**3. 验证通过后：**
- 将用户ID和昵称填入第一个空位（member1 或 member2）
- `member_count` 自动 +1
- 返回加入成功信息

#### 4.3 成员退出与踢出机制

**退出团队（成员主动）：**
- 只有 member1 和 member2 可以退出
- 队长不能退出，只能解散团队
- 退出后位置自动调整：
  - member1 退出 且 member2 存在 → member2 移至 member1 位置
  - member2 退出 → 直接清空 member2 位置

**踢出成员（队长操作）：**
- 只有队长可以踢出成员
- 不能踢出自己
- 踢出后位置调整规则同退出

**位置调整示例：**
```
调整前：captain + member1 + member2
member1 退出 → captain + member2 (移至member1位置) + null

调整前：captain + member1 + member2
member2 退出 → captain + member1 + null
```

#### 4.4 团队解散机制

**解散条件：**
- 只有队长可以解散团队
- 解散后 `status` 变为 `disbanded`
- 所有成员位置保留但团队不可再使用

**解散后的限制：**
- 已解散的团队不出现在团队列表中
- 已解散的团队不参与排行榜
- 成员无法再加入已解散的团队
- 队长可以创建新团队
- 成员可以加入其他团队

#### 4.5 团队封禁机制

**封禁条件：**
- 只有管理员（admin）和究极管理员（super_admin）可以封禁团队
- 封禁原因必须明确记录
- 封禁后 `status` 变为 `banned`
- 封禁操作记录到操作日志

**封禁后的限制：**
- **无法查看题目**：被封禁团队的所有成员无法查看比赛题目列表和题目详情
- **无法提交答案**：被封禁团队无法提交任何答案
- **不参与排行榜**：被封禁团队不会出现在排行榜中
- **团队信息可见**：团队信息仍然可以被查询，但会显示 `banned` 状态
- **成员不能退出**：被封禁期间，成员无法主动退出团队
- **队长不能解散**：被封禁期间，队长无法解散团队

**解禁机制：**
- 只有管理员和究极管理员可以解禁团队
- 解禁后 `status` 恢复为 `active`
- 团队恢复正常权限，可以查看题目、提交答案
- 解禁操作记录到操作日志

**与解散的区别：**
- **解散（disbanded）**：队长主动操作，成员自动退出，数据保留
- **封禁（banned）**：管理员惩罚操作，成员无法退出，可解禁恢复

#### 4.6 团队得分机制

**得分来源：**
- 解题得分：团队成员解出题目后，分数计入团队总分
- 得分实时更新到 `team_score` 字段
- 只有活跃团队（`status = 'active'`）的得分才会计入排行榜

**排行榜规则：**
- 按赛道分别排名（social/freshman/advanced）
- 按 `team_score` 降序排列
- 分数相同时，按解题时间先后排序（需结合题目模块）

#### 4.7 数据一致性保证

**成员数量统计：**
- 创建团队：`member_count = 1`
- 成员加入：`member_count + 1`
- 成员退出/被踢：`member_count - 1`
- 最小值 1（只有队长），最大值 3

**成员昵称冗余字段：**
- 加入团队时自动填充成员昵称
- 用户修改昵称时，需同步更新团队表的成员昵称字段
- 提高查询效率，避免频繁关联用户表

**学校信息同步：**
- 联合院校团队创建时，自动填充 `school_id` 和 `school_name`
- 学校名称修改时，需同步更新团队表的 `school_name`

**外键约束处理：**
- 队长不能被删除（`ON DELETE RESTRICT`）
- 成员被删除时，对应位置自动设为 NULL（`ON DELETE SET NULL`）
- 学校被删除时，团队的 `school_id` 设为 NULL

***

## 四、题目类型模块

### 1. 数据库设计

***

**数据库：**`dalictf`

**表名：**`dalictf_challenge_category`

***

#### 字段设计

|      字段名      |                 类型                  |                约束条件                 |                            说明                            |
| :--------------: | :-----------------------------------: | :-------------------------------------: | :--------------------------------------------------------: |
|        id        |            `BIGINT(20)`             |      `PRIMARY KEY AUTO_INCREMENT`       |                          分类主键                          |
|     direction    |           `VARCHAR(50)`           |            `NOT NULL UNIQUE`            |          类型方向（Web、Misc、Crypto、Reverse、Pwn等）          |
|     name_zh     |           `VARCHAR(50)`           |            `NOT NULL`            |          类型中文名称          |
|     name_en     |           `VARCHAR(50)`           |            `NOT NULL`            |          类型英文名称          |
|   description    |          `VARCHAR(500)`           |            `DEFAULT NULL`             |              类型描述（介绍该方向的内容）              |
|      icon        |          `VARCHAR(100)`           |            `DEFAULT NULL`             |              图标标识（用于前端展示，如 'icon-web'）              |
|      color       |           `VARCHAR(20)`           |            `DEFAULT NULL`             |          主题颜色（HEX格式，如 '#FF6B6B'，用于UI区分）          |
|   sort_order   |            `INT(11)`            |        `NOT NULL DEFAULT 0`         | 排序顺序（数字越小越靠前） |
|      status      |   `ENUM('active', 'inactive')`   |    `NOT NULL DEFAULT 'active'`    |   分类状态：active-启用/inactive-停用   |
|    created_at    |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                          创建时间                          |
|    updated_at    |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` |                          更新时间                          |
|    deleted_at    |            `DATETIME`             |            `DEFAULT NULL`             |               软删除时间（NULL表示未删除）               |

#### **建表 `sql`**

```sql
CREATE TABLE `dalictf_challenge_category` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '分类主键ID',
  `direction` VARCHAR(50) NOT NULL COMMENT '类型方向',
  `name_zh` VARCHAR(50) NOT NULL COMMENT '类型中文名称',
  `name_en` VARCHAR(50) NOT NULL COMMENT '类型英文名称',
  `description` VARCHAR(500) DEFAULT NULL COMMENT '类型描述',
  `icon` VARCHAR(100) DEFAULT NULL COMMENT '图标标识',
  `color` VARCHAR(20) DEFAULT NULL COMMENT '主题颜色',
  `sort_order` INT(11) NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `status` ENUM('active', 'inactive') NOT NULL DEFAULT 'active' COMMENT '分类状态',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_direction` (`direction`),
  KEY `idx_status` (`status`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='题目类型分类表';
```

#### 字段说明补充

**关于类型方向（direction）：**
- 使用英文标识，全局唯一
- 常见方向包括：
  - `Web` - Web安全
  - `Misc` - 杂项
  - `Crypto` - 密码学
  - `Reverse` - 逆向工程
  - `Pwn` - 二进制漏洞利用
  - `Forensics` - 取证分析
  - `Mobile` - 移动安全
  - `Blockchain` - 区块链安全
  - `IoT` - 物联网安全
  - `AI` - AI安全

**关于中英文名称：**
- `name_zh`：中文显示名称，如"Web安全"、"密码学"
- `name_en`：英文显示名称，如"Web Security"、"Cryptography"
- 支持国际化，前端可根据语言设置显示

**关于图标和颜色：**
- `icon`：图标class或图标名称，用于前端UI展示
- `color`：HEX颜色代码，用于区分不同类型的视觉效果
- 示例配色方案：
  - Web: `#FF6B6B` (红色)
  - Misc: `#4ECDC4` (青色)
  - Crypto: `#FFE66D` (黄色)
  - Reverse: `#A8E6CF` (绿色)
  - Pwn: `#FF8B94` (粉红)

**关于排序（sort_order）：**
- 数字越小，排序越靠前
- 用于前端显示题目分类列表的顺序
- 可以通过管理后台调整

**关于状态（status）：**
- `active`：启用状态，该类型的题目可以被创建和显示
- `inactive`：停用状态，该类型暂时不可用（但不删除数据）

#### 初始化数据示例

```sql
INSERT INTO `dalictf_challenge_category` 
(`direction`, `name_zh`, `name_en`, `description`, `icon`, `color`, `sort_order`, `status`) 
VALUES
('Web', 'Web安全', 'Web Security', 'Web应用安全，包括SQL注入、XSS、CSRF等常见Web漏洞', 'icon-web', '#FF6B6B', 1, 'active'),
('Misc', '杂项', 'Miscellaneous', '杂项题目，包括编码、隐写、社工等多种技巧', 'icon-misc', '#4ECDC4', 2, 'active'),
('Crypto', '密码学', 'Cryptography', '密码学相关题目，包括古典密码、现代加密算法等', 'icon-crypto', '#FFE66D', 3, 'active'),
('Reverse', '逆向工程', 'Reverse Engineering', '二进制程序逆向分析，包括软件破解、协议分析等', 'icon-reverse', '#A8E6CF', 4, 'active'),
('Pwn', '二进制漏洞', 'Binary Exploitation', '二进制漏洞利用，包括栈溢出、堆漏洞等', 'icon-pwn', '#FF8B94', 5, 'active');
```

***

### 2. API 设计

> **统一前缀：**`/api/v1/categories`
>
> **认证机制：**JWT Token（详见参赛学校模块说明）

***

#### 创建题目类型（管理员）

- **URL：** `POST /api/v1/categories`

- **权限：**管理员（admin及以上）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求体：**

  ```json
  {
    "direction": "Web",
    "name_zh": "Web安全",
    "name_en": "Web Security",
    "description": "Web应用安全，包括SQL注入、XSS、CSRF等常见Web漏洞",
    "icon": "icon-web",
    "color": "#FF6B6B",
    "sort_order": 1
  }
  ```

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Category created successfully",
    "data": {
      "id": 1,
      "direction": "Web",
      "name_zh": "Web安全",
      "name_en": "Web Security",
      "description": "Web应用安全，包括SQL注入、XSS、CSRF等常见Web漏洞",
      "icon": "icon-web",
      "color": "#FF6B6B",
      "sort_order": 1,
      "status": "active",
      "created_at": "2025-11-20 21:00:00"
    }
  }
  ```

#### 查询题目类型列表

- **URL：** `GET /api/v1/categories`

- **权限：**公开（任何人都可以查看）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>  // 可选
  ```

- **请求参数（Query）：**
  - `status` 状态筛选（active/inactive，默认只显示active）
  - `include_inactive` 是否包含停用的分类（true/false，默认false）

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total": 5,
      "list": [
        {
          "id": 1,
          "direction": "Web",
          "name_zh": "Web安全",
          "name_en": "Web Security",
          "description": "Web应用安全，包括SQL注入、XSS、CSRF等常见Web漏洞",
          "icon": "icon-web",
          "color": "#FF6B6B",
          "sort_order": 1,
          "status": "active"
        },
        {
          "id": 2,
          "direction": "Misc",
          "name_zh": "杂项",
          "name_en": "Miscellaneous",
          "description": "杂项题目，包括编码、隐写、社工等多种技巧",
          "icon": "icon-misc",
          "color": "#4ECDC4",
          "sort_order": 2,
          "status": "active"
        }
      ]
    }
  }
  ```

#### 查询单个题目类型

- **URL：** `GET /api/v1/categories/:id`

- **权限：**公开

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>  // 可选
  ```

- **路径参数：**
  - `id` 分类ID

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "id": 1,
      "direction": "Web",
      "name_zh": "Web安全",
      "name_en": "Web Security",
      "description": "Web应用安全，包括SQL注入、XSS、CSRF等常见Web漏洞",
      "icon": "icon-web",
      "color": "#FF6B6B",
      "sort_order": 1,
      "status": "active",
      "created_at": "2025-11-20 21:00:00",
      "updated_at": "2025-11-20 21:00:00"
    }
  }
  ```

#### 修改题目类型（管理员）

- **URL：** `PUT /api/v1/categories/:id`

- **权限：**管理员（admin及以上）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 分类ID

- **请求体：**

  ```json
  {
    "name_zh": "Web应用安全",
    "name_en": "Web Application Security",
    "description": "Web应用安全，包括SQL注入、XSS、CSRF、SSRF等常见Web漏洞",
    "icon": "icon-web-new",
    "color": "#FF5555",
    "sort_order": 1
  }
  ```

- **说明**：
  - `direction` 不可修改（作为唯一标识）
  - 其他字段均可选修改

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Category updated successfully",
    "data": {
      "id": 1,
      "direction": "Web",
      "name_zh": "Web应用安全",
      "name_en": "Web Application Security",
      "description": "Web应用安全，包括SQL注入、XSS、CSRF、SSRF等常见Web漏洞",
      "icon": "icon-web-new",
      "color": "#FF5555",
      "sort_order": 1,
      "status": "active",
      "updated_at": "2025-11-20 22:00:00"
    }
  }
  ```

#### 修改题目类型状态（管理员）

- **URL：** `PUT /api/v1/categories/:id/status`

- **权限：**管理员（admin及以上）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 分类ID

- **请求体：**

  ```json
  {
    "status": "inactive"  // active 或 inactive
  }
  ```

- **说明**：
  - 停用分类后，该分类的题目仍然存在，但前端可能不显示
  - 停用分类不影响已有题目，只是新建题目时不能选择该分类

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Category status updated successfully",
    "data": {
      "id": 1,
      "direction": "Web",
      "status": "inactive",
      "updated_at": "2025-11-20 22:30:00"
    }
  }
  ```

#### 删除题目类型（管理员）

- **URL：** `DELETE /api/v1/categories/:id`

- **权限：**管理员（admin及以上）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 分类ID

- **说明**：
  - 软删除，设置 `deleted_at` 字段
  - 如果该分类下有题目，不允许删除（需先转移或删除题目）

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Category deleted successfully",
    "data": null
  }
  ```

#### 批量调整排序（管理员）

- **URL：** `PUT /api/v1/categories/sort`

- **权限：**管理员（admin及以上）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求体：**

  ```json
  {
    "sort_list": [
      { "id": 1, "sort_order": 1 },
      { "id": 2, "sort_order": 2 },
      { "id": 3, "sort_order": 3 }
    ]
  }
  ```

- **说明**：
  - 批量更新多个分类的排序顺序
  - 用于拖拽调整分类顺序后的保存

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Sort order updated successfully",
    "data": null
  }
  ```

***

### 3. 错误码定义

| 错误码 |                  说明                  |
| :----: | :------------------------------------: |
|  200   |                  成功                  |
|  400   |      请求参数错误（如必填字段缺失）      |
|  401   |      未授权（Token缺失或格式错误）      |
|  403   |        无权限操作（权限不足）        |
|  404   |             分类不存在             |
|  409   |   分类direction已存在   |
|  431   |      该分类下存在题目，无法删除      |
|  500   |              服务器错误              |

***



## 五、题目模块

### 1. 数据库设计

#### 1.1 题目基础信息表

***

**数据库：**`dalictf`

**表名：**`dalictf_challenge`

***

#### 字段设计

|       字段名        |                 类型                  |                约束条件                 |                            说明                            |
| :-----------------: | :-----------------------------------: | :-------------------------------------: | :--------------------------------------------------------: |
|         id          |            `BIGINT(20)`             |      `PRIMARY KEY AUTO_INCREMENT`       |                          题目主键                          |
|   challenge_name    |          `VARCHAR(255)`           |            `NOT NULL`             |                   题目名称                   |
|      direction      |           `VARCHAR(50)`           |            `NOT NULL`             |      题目类型（外键关联 challenge_category 表的 direction）      |
|       author        |           `VARCHAR(100)`           |            `NOT NULL`             |      出题人网名（不是 user 表中的名字，是出题人的网名）      |
|     description     |              `TEXT`               |            `NOT NULL`             |              题目描述（支持Markdown格式）              |
|        hint         |              `TEXT`               |            `DEFAULT NULL`             |              题目提示（可选，支持Markdown格式）              |
|        state        |   `ENUM('visible', 'hidden')`   |    `NOT NULL DEFAULT 'visible'`    |   题目状态：visible-显示/hidden-隐藏   |
|        mode         |   `ENUM('static', 'dynamic')`   |    `NOT NULL DEFAULT 'static'`    |   题目模式：static-静态题/dynamic-动态题   |
|     static_flag     |          `VARCHAR(500)`           |            `DEFAULT NULL`             |      静态题 flag（mode=static 时必填）      |
|    docker_image     |          `VARCHAR(255)`           |            `DEFAULT NULL`             |      动态题 Docker 镜像（mode=dynamic 时必填）      |
|    docker_ports     |              `JSON`               |            `DEFAULT NULL`             |      容器端口映射（JSON格式，如 {"80":"tcp", "3306":"tcp"}）      |
|     difficulty      | `ENUM('easy', 'medium', 'hard', 'expert')` |      `NOT NULL DEFAULT 'medium'`      | 题目难度：easy-简单/medium-中等/hard-困难/expert-专家级 |
|   initial_score    |            `INT(11)`            |        `NOT NULL DEFAULT 100`         | 初始分值 |
|     min_score      |            `INT(11)`            |        `NOT NULL DEFAULT 50`         | 最低分值（动态衰减的最低分） |
|   current_score    |            `INT(11)`            |        `NOT NULL DEFAULT 100`         | 当前分值（动态衰减时更新） |
|    decay_ratio     |        `DECIMAL(5, 2)`        |        `NOT NULL DEFAULT 0.90`         | 分数衰减比率（如0.90表示每次解出后得分×0.90） |
|   solved_count    |            `INT(11)`            |        `NOT NULL DEFAULT 0`         | 解出次数（统计） |
|     created_at     |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                          创建时间                          |
|     updated_at     |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` |                          更新时间                          |
|     deleted_at     |            `DATETIME`             |            `DEFAULT NULL`             |               软删除时间（NULL表示未删除）               |

#### **建表 `sql`**

```sql
CREATE TABLE `dalictf_challenge` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '题目主键ID',
  `challenge_name` VARCHAR(255) NOT NULL COMMENT '题目名称',
  `direction` VARCHAR(50) NOT NULL COMMENT '题目类型',
  `author` VARCHAR(100) NOT NULL COMMENT '出题人网名',
  `description` TEXT NOT NULL COMMENT '题目描述',
  `hint` TEXT DEFAULT NULL COMMENT '题目提示',
  `state` ENUM('visible', 'hidden') NOT NULL DEFAULT 'visible' COMMENT '题目状态',
  `mode` ENUM('static', 'dynamic') NOT NULL DEFAULT 'static' COMMENT '题目模式',
  `static_flag` VARCHAR(500) DEFAULT NULL COMMENT '静态题flag',
  `docker_image` VARCHAR(255) DEFAULT NULL COMMENT '动态题Docker镜像',
  `docker_ports` JSON DEFAULT NULL COMMENT '容器端口映射',
  `difficulty` ENUM('easy', 'medium', 'hard', 'expert') NOT NULL DEFAULT 'medium' COMMENT '题目难度',
  `initial_score` INT(11) NOT NULL DEFAULT 100 COMMENT '初始分值',
  `min_score` INT(11) NOT NULL DEFAULT 50 COMMENT '最低分值',
  `current_score` INT(11) NOT NULL DEFAULT 100 COMMENT '当前分值',
  `decay_ratio` DECIMAL(5, 2) NOT NULL DEFAULT 0.90 COMMENT '分数衰减比率',
  `solved_count` INT(11) NOT NULL DEFAULT 0 COMMENT '解出次数',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_direction` (`direction`),
  KEY `idx_state` (`state`),
  KEY `idx_mode` (`mode`),
  KEY `idx_difficulty` (`difficulty`),
  KEY `idx_current_score` (`current_score`),
  KEY `idx_solved_count` (`solved_count`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='题目基础信息表';

-- 添加外键约束（确保题目类型存在）
ALTER TABLE `dalictf_challenge` 
ADD CONSTRAINT `fk_challenge_direction` 
FOREIGN KEY (`direction`) REFERENCES `dalictf_challenge_category`(`direction`) 
ON DELETE RESTRICT ON UPDATE CASCADE;
```

#### 字段说明补充

**关于题目模式（mode）：**

- **静态题（static）**：
  - Flag 固定不变，所有用户的 flag 都相同
  - 必须填写 `static_flag` 字段
  - `docker_image` 和 `docker_ports` 为 NULL
  - 适合密码学、编码、杂项等题型
  - 无需容器部署，节省资源

- **动态题（dynamic）**：
  - 每个团队/用户有独立的容器实例
  - Flag 动态生成，每个团队的 flag 不同
  - 必须填写 `docker_image` 字段
  - 可选填写 `docker_ports`（端口映射配置）
  - `static_flag` 为 NULL
  - 适合 Web、Pwn、Reverse 等需要环境的题型
  - 需要容器编排平台（如 Docker、Kubernetes）

**关于 Docker 端口映射（docker_ports）：**
- JSON 格式存储端口映射关系
- 格式示例：`{"80": "tcp", "3306": "tcp", "8080": "http"}`
- Key：容器内部端口
- Value：协议类型（tcp、udp、http、https 等）
- 系统会自动分配宿主机端口并映射到容器端口
- 前端展示时会显示：`http://xxx.xxx.xxx.xxx:随机端口`

**关于题目状态（state）：**
- **visible（显示）**：题目对参赛用户可见，可以查看和提交答案
- **hidden（隐藏）**：题目暂时不可见，用于比赛开始前准备或临时下线

**关于难度（difficulty）：**
- **easy（简单）**：适合新手，基础知识即可解决
- **medium（中等）**：需要一定技巧和经验
- **hard（困难）**：需要深入理解和多种技术结合
- **expert（专家级）**：极具挑战性，需要创新思维和高级技术

**关于分数机制：**
- **initial_score（初始分值）**：题目刚发布时的分数
- **current_score（当前分值）**：根据解出人数动态衰减后的分数
- **min_score（最低分值）**：分数衰减的下限，不会低于此值
- **decay_ratio（衰减比率）**：每有一个团队解出后，分数按此比率衰减
- **solved_count（解出次数）**：统计有多少团队/用户解出该题

**分数衰减公式：**
```
当前分值 = max(初始分值 × (衰减比率 ^ 解出次数), 最低分值)

示例：
初始分值 = 500
衰减比率 = 0.90
最低分值 = 100

第1个团队解出：500 分
第2个团队解出：500 × 0.90 = 450 分
第3个团队解出：500 × 0.90^2 = 405 分
...
直到降至最低分值 100 分
```

**关于出题人（author）：**
- 存储出题人的网名或昵称（不是平台用户名）
- 可以是出题人在CTF圈的常用ID
- 用于致谢和版权标识
- 可以在题目详情页显示

***

#### 1.2 题目附件表

***

**数据库：**`dalictf`

**表名：**`dalictf_challenge_attachment`

***

#### 字段设计

|       字段名        |                 类型                  |                约束条件                 |                            说明                            |
| :-----------------: | :-----------------------------------: | :-------------------------------------: | :--------------------------------------------------------: |
|         id          |            `BIGINT(20)`             |      `PRIMARY KEY AUTO_INCREMENT`       |                          附件主键                          |
|    challenge_id     |            `BIGINT(20)`             |            `NOT NULL`             |      所属题目ID（外键关联 challenge 表）      |
|       storage       |   `ENUM('url', 'object')`   |    `NOT NULL DEFAULT 'object'`    |   存储形态：url-外链/object-对象存储   |
|         url         |          `VARCHAR(1000)`           |            `DEFAULT NULL`             |      外链地址（storage=url 时必填）      |
|   object_bucket    |          `VARCHAR(100)`           |            `DEFAULT NULL`             |      对象存储桶名称（storage=object 时必填）      |
|    object_key     |          `VARCHAR(500)`           |            `DEFAULT NULL`             |      对象存储 Key（storage=object 时必填）      |
|     file_name      |          `VARCHAR(255)`           |            `NOT NULL`             |              文件名（展示用）              |
|   content_type    |          `VARCHAR(100)`           |            `DEFAULT NULL`             |              MIME 类型（如 application/zip）              |
|     file_size      |            `BIGINT(20)`             |            `DEFAULT NULL`             |              文件大小（字节）              |
|       sha256       |           `CHAR(64)`           |            `DEFAULT NULL`             |              SHA256 校验值（完整性校验）              |
|       status       | `ENUM('pending', 'active', 'infected', 'error')` |      `NOT NULL DEFAULT 'pending'`      | 安全状态：pending-待检查/active-正常/infected-病毒/error-错误 |
|    visibility     | `ENUM('public', 'private', 'team')` |      `NOT NULL DEFAULT 'private'`      | 可见性：public-公开/private-私有/team-团队可见 |
|      version       |          `VARCHAR(50)`           |        `DEFAULT '1.0'`         | 版本号 |
|    sort_order     |            `INT(11)`            |        `NOT NULL DEFAULT 0`         | 排序顺序（同一题目多个附件时的显示顺序） |
|    created_by     |            `BIGINT(20)`             |            `DEFAULT NULL`             |      上传者用户ID（外键关联 user 表）      |
|     created_at     |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                          创建时间                          |
|     updated_at     |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` |                          更新时间                          |
|     deleted_at     |            `DATETIME`             |            `DEFAULT NULL`             |               软删除时间（NULL表示未删除）               |

#### **建表 `sql`**

```sql
CREATE TABLE `dalictf_challenge_attachment` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '附件主键ID',
  `challenge_id` BIGINT(20) NOT NULL COMMENT '所属题目ID',
  `storage` ENUM('url', 'object') NOT NULL DEFAULT 'object' COMMENT '存储形态',
  `url` VARCHAR(1000) DEFAULT NULL COMMENT '外链地址',
  `object_bucket` VARCHAR(100) DEFAULT NULL COMMENT '对象存储桶名称',
  `object_key` VARCHAR(500) DEFAULT NULL COMMENT '对象存储Key',
  `file_name` VARCHAR(255) NOT NULL COMMENT '文件名',
  `content_type` VARCHAR(100) DEFAULT NULL COMMENT 'MIME类型',
  `file_size` BIGINT(20) DEFAULT NULL COMMENT '文件大小（字节）',
  `sha256` CHAR(64) DEFAULT NULL COMMENT 'SHA256校验值',
  `status` ENUM('pending', 'active', 'infected', 'error') NOT NULL DEFAULT 'pending' COMMENT '安全状态',
  `visibility` ENUM('public', 'private', 'team') NOT NULL DEFAULT 'private' COMMENT '可见性',
  `version` VARCHAR(50) DEFAULT '1.0' COMMENT '版本号',
  `sort_order` INT(11) NOT NULL DEFAULT 0 COMMENT '排序顺序',
  `created_by` BIGINT(20) DEFAULT NULL COMMENT '上传者用户ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_challenge_id` (`challenge_id`),
  KEY `idx_storage` (`storage`),
  KEY `idx_status` (`status`),
  KEY `idx_visibility` (`visibility`),
  KEY `idx_created_by` (`created_by`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='题目附件表';

-- 添加外键约束
ALTER TABLE `dalictf_challenge_attachment` 
ADD CONSTRAINT `fk_attachment_challenge` 
FOREIGN KEY (`challenge_id`) REFERENCES `dalictf_challenge`(`id`) 
ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE `dalictf_challenge_attachment` 
ADD CONSTRAINT `fk_attachment_creator` 
FOREIGN KEY (`created_by`) REFERENCES `dalictf_user`(`id`) 
ON DELETE SET NULL ON UPDATE CASCADE;
```

#### 字段说明补充

**关于存储形态（storage）：**

- **外链（url）**：
  - 附件存储在外部平台（如网盘、CDN）
  - 必须填写 `url` 字段
  - `object_bucket` 和 `object_key` 为 NULL
  - 适合超大文件或已有外部资源
  - 示例：`https://pan.baidu.com/s/xxxxx`

- **对象存储（object）**：
  - 附件存储在对象存储服务（如阿里云OSS、AWS S3、MinIO）
  - 必须填写 `object_bucket` 和 `object_key`
  - `url` 为 NULL（下载时动态生成临时URL）
  - 更安全可控，支持权限管理
  - 示例：bucket=`ctf-attachments`, key=`challenges/2025/web/challenge1.zip`

**关于安全状态（status）：**
- **pending（待检查）**：附件刚上传，等待安全扫描
- **active（正常）**：已通过安全检查，可以正常下载
- **infected（病毒）**：检测到病毒或恶意代码，禁止下载
- **error（错误）**：上传或处理过程出错，需要重新上传

**关于可见性（visibility）：**
- **public（公开）**：任何人都可以下载（包括未登录用户）
- **private（私有）**：需要登录且参赛的用户才能下载
- **team（团队可见）**：只有已组队的用户才能下载

**关于文件校验（sha256）：**
- 上传时自动计算文件的 SHA256 值
- 用户下载后可以验证文件完整性
- 防止文件在传输过程中被篡改
- 可以在前端展示校验值供用户验证

**关于版本号（version）：**
- 同一题目的附件可能会更新
- 版本号用于标识附件的迭代版本
- 格式建议：1.0、1.1、2.0 等
- 方便追踪附件变更历史

**附件上传流程：**
1. 管理员/出题人上传附件
2. 系统计算 SHA256 校验值
3. 状态设为 `pending`，触发安全扫描
4. 安全扫描通过后，状态更新为 `active`
5. 对象存储模式：生成临时下载URL（有效期如1小时）
6. 外链模式：直接返回外链地址

**附件下载流程：**
1. 用户请求下载附件
2. 验证用户权限（根据 `visibility` 判断）
3. 检查附件状态（必须是 `active`）
4. 对象存储：生成临时签名URL（防盗链）
5. 外链：直接返回URL或重定向
6. 记录下载日志（可选）

***

### 2. API 设计

> **统一前缀：**`/api/v1/challenges`
>
> **认证机制：**JWT Token（详见参赛学校模块说明）
>
> **统一请求头：**
>
> ```
> Content-Type: application/json
> Authorization: Bearer <JWT_TOKEN>  // 需要权限的接口必须携带
> ```
>
> **返回格式：**统一使用标准响应格式（同学校模块）

***

#### 创建题目（管理员）

- **URL：** `POST /api/v1/challenges`

- **权限：**管理员权限（admin及以上）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求体：**

  ```json
  {
    "challenge_name": "简单的SQL注入",
    "direction": "Web",
    "author": "Hacker123",
    "description": "这是一个简单的SQL注入题目，请找到flag\n\n**目标：**\n- 获取管理员密码\n- 提交flag",
    "hint": "尝试使用union注入",
    "state": "visible",
    "mode": "dynamic",
    "static_flag": null,
    "docker_image": "ctf/web-sqli:latest",
    "docker_ports": {"80": "tcp"},
    "difficulty": "easy",
    "initial_score": 500,
    "min_score": 100,
    "decay_ratio": 0.90
  }
  ```

- **说明**：
  - `mode=static` 时必须填写 `static_flag`
  - `mode=dynamic` 时必须填写 `docker_image`
  - `description` 和 `hint` 支持 Markdown 格式
  - `direction` 必须是已存在的题目类型

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Challenge created successfully",
    "data": {
      "id": 1001,
      "challenge_name": "简单的SQL注入",
      "direction": "Web",
      "author": "Hacker123",
      "mode": "dynamic",
      "difficulty": "easy",
      "initial_score": 500,
      "current_score": 500,
      "solved_count": 0,
      "state": "visible",
      "created_at": "2025-11-20 22:00:00"
    }
  }
  ```

#### 查询题目列表

- **URL：** `GET /api/v1/challenges`

- **权限：**
  - 公开：可查看 `state=visible` 的题目列表
  - 管理员：可查看所有题目（包括 hidden）

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>  // 可选，登录后获取更多信息
  ```

- **请求参数（Query）：**
  - `page` 页码（默认 1）
  - `limit` 每页数量（默认 20）
  - `direction` 题目类型筛选（Web、Misc、Crypto 等）
  - `difficulty` 难度筛选（easy/medium/hard/expert）
  - `mode` 模式筛选（static/dynamic）
  - `state` 状态筛选（visible/hidden，仅管理员可用）
  - `search` 模糊搜索（题目名称、作者）
  - `sort_by` 排序字段（score/solved_count/created_at，默认 created_at）
  - `order` 排序方式（asc/desc，默认 desc）

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total": 50,
      "page": 1,
      "limit": 20,
      "list": [
        {
          "id": 1001,
          "challenge_name": "简单的SQL注入",
          "direction": "Web",
          "author": "Hacker123",
          "difficulty": "easy",
          "current_score": 450,
          "initial_score": 500,
          "solved_count": 3,
          "mode": "dynamic",
          "state": "visible",
          "is_solved": false,
          "created_at": "2025-11-20 22:00:00"
        },
        {
          "id": 1002,
          "challenge_name": "Base64编码",
          "direction": "Misc",
          "author": "Alice",
          "difficulty": "easy",
          "current_score": 100,
          "initial_score": 100,
          "solved_count": 25,
          "mode": "static",
          "state": "visible",
          "is_solved": true,
          "created_at": "2025-11-19 10:00:00"
        }
      ]
    }
  }
  ```

- **说明**：
  - `is_solved`：当前用户/团队是否已解出该题（需登录）
  - 未登录用户看不到 `is_solved` 字段

#### 查询单个题目详情

- **URL：** `GET /api/v1/challenges/:id`

- **权限：**
  - 公开：可查看 `state=visible` 的题目
  - 管理员：可查看所有题目

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>  // 可选
  ```

- **路径参数：**
  - `id` 题目ID

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "id": 1001,
      "challenge_name": "简单的SQL注入",
      "direction": "Web",
      "author": "Hacker123",
      "description": "这是一个简单的SQL注入题目，请找到flag...",
      "hint": "尝试使用union注入",
      "difficulty": "easy",
      "current_score": 450,
      "initial_score": 500,
      "min_score": 100,
      "decay_ratio": 0.90,
      "solved_count": 3,
      "mode": "dynamic",
      "state": "visible",
      "is_solved": false,
      "solve_time": null,
      "attachments": [
        {
          "id": 2001,
          "file_name": "source_code.zip",
          "file_size": 2048576,
          "sha256": "abc123...",
          "download_url": "/api/v1/challenges/1001/attachments/2001/download"
        }
      ],
      "container_info": {
        "status": "running",
        "url": "http://192.168.1.100:32768",
        "expires_at": "2025-11-20 23:00:00"
      },
      "created_at": "2025-11-20 22:00:00",
      "updated_at": "2025-11-20 22:30:00"
    }
  }
  ```

- **说明**：
  - `is_solved`：当前用户/团队是否已解出
  - `solve_time`：当前用户/团队解出时间（未解出为null）
  - `container_info`：动态题目的容器信息（静态题为null）
  - `attachments`：题目附件列表

#### 修改题目信息（管理员）

- **URL：** `PUT /api/v1/challenges/:id`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 题目ID

- **请求体：**

  ```json
  {
    "challenge_name": "中级SQL注入",
    "description": "更新后的题目描述...",
    "hint": "更新后的提示",
    "difficulty": "medium",
    "initial_score": 600,
    "min_score": 150,
    "decay_ratio": 0.85
  }
  ```

- **说明**：
  - 所有字段均为可选
  - 不能修改 `mode`（静态/动态模式创建后不可更改）
  - 修改 `initial_score` 会重置 `current_score`
  - 不能修改 `solved_count`（由系统自动统计）

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Challenge updated successfully",
    "data": {
      "id": 1001,
      "challenge_name": "中级SQL注入",
      "difficulty": "medium",
      "initial_score": 600,
      "current_score": 600,
      "updated_at": "2025-11-20 23:00:00"
    }
  }
  ```

#### 修改题目状态（管理员）

- **URL：** `PATCH /api/v1/challenges/:id/state`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 题目ID

- **请求体：**

  ```json
  {
    "state": "hidden"
  }
  ```

- **说明**：
  - 将题目设为 `hidden` 后，普通用户无法查看
  - 已开启的动态容器不会自动关闭
  - 管理员仍然可以查看和管理

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Challenge state updated successfully",
    "data": {
      "id": 1001,
      "challenge_name": "简单的SQL注入",
      "state": "hidden",
      "updated_at": "2025-11-20 23:10:00"
    }
  }
  ```

#### 删除题目（管理员）

- **URL：** `DELETE /api/v1/challenges/:id`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 题目ID

- **请求体：**无

- **说明**：
  - 软删除，设置 `deleted_at` 字段
  - 删除题目会级联删除所有附件（软删除）
  - 已有的解题记录保留
  - 已开启的动态容器会自动关闭

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Challenge deleted successfully",
    "data": null
  }
  ```

***

#### 提交答案（Flag）

- **URL：** `POST /api/v1/challenges/:id/submit`

- **权限：**需要登录且已组队

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 题目ID

- **请求体：**

  ```json
  {
    "flag": "flag{th1s_is_a_test_fl4g}"
  }
  ```

- **说明**：
  - 提交前会验证：
    - 用户是否已登录
    - 用户是否已加入团队
    - 团队状态是否为 `active`（不能是 `banned` 或 `disbanded`）
    - 题目状态是否为 `visible`
    - 是否已经解出该题（避免重复提交）
  - Flag 比对不区分大小写
  - 提交成功后：
    - 团队分数增加当前题目分值
    - 题目 `solved_count + 1`
    - 题目 `current_score` 按衰减比率更新
    - 记录解题时间和团队信息

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Correct flag! Congratulations!",
    "data": {
      "challenge_id": 1001,
      "challenge_name": "简单的SQL注入",
      "earned_score": 450,
      "team_total_score": 2450,
      "solve_time": "2025-11-20 22:45:00",
      "rank": 5
    }
  }
  ```

- **错误返回（Flag错误）：**

  ```json
  {
    "code": 432,
    "msg": "Incorrect flag, please try again",
    "data": {
      "attempts_left": null
    }
  }
  ```

#### 上传题目附件（管理员）

- **URL：** `POST /api/v1/challenges/:id/attachments`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: multipart/form-data
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 题目ID

- **请求体（FormData）：**
  - `file` 附件文件（必填）
  - `storage` 存储类型：object/url（可选，默认 object）
  - `url` 外链地址（storage=url 时必填）
  - `visibility` 可见性：public/private/team（可选，默认 private）
  - `version` 版本号（可选，默认 1.0）

- **说明**：
  - 对象存储模式：上传文件到服务器，系统自动存储到对象存储
  - 外链模式：只提供URL，不上传文件
  - 上传后自动触发安全扫描
  - 支持的文件类型：zip、rar、7z、tar、gz、txt、pdf、docx 等
  - 单个文件大小限制：100MB

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Attachment uploaded successfully",
    "data": {
      "id": 2001,
      "challenge_id": 1001,
      "file_name": "source_code.zip",
      "file_size": 2048576,
      "content_type": "application/zip",
      "sha256": "abc123def456...",
      "storage": "object",
      "status": "pending",
      "visibility": "private",
      "version": "1.0",
      "created_at": "2025-11-20 22:50:00"
    }
  }
  ```

#### 下载题目附件

- **URL：** `GET /api/v1/challenges/:challenge_id/attachments/:attachment_id/download`

- **权限：**根据附件的 `visibility` 判断

- **请求头：**

  ```
  Authorization: Bearer <JWT_TOKEN>  // 根据 visibility 决定是否必须
  ```

- **路径参数：**
  - `challenge_id` 题目ID
  - `attachment_id` 附件ID

- **说明**：
  - `visibility=public`：任何人可下载
  - `visibility=private`：需要登录
  - `visibility=team`：需要登录且已组队
  - 对象存储模式：返回临时签名URL（有效期1小时）
  - 外链模式：重定向到外链地址
  - 检查附件状态必须为 `active`

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "download_url": "https://oss.example.com/ctf-attachments/xxx?signature=xxx&expires=xxx",
      "file_name": "source_code.zip",
      "file_size": 2048576,
      "sha256": "abc123def456...",
      "expires_at": "2025-11-20 23:50:00"
    }
  }
  ```

#### 删除题目附件（管理员）

- **URL：** `DELETE /api/v1/challenges/:challenge_id/attachments/:attachment_id`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `challenge_id` 题目ID
  - `attachment_id` 附件ID

- **请求体：**无

- **说明**：
  - 软删除，设置 `deleted_at` 字段
  - 对象存储模式：不会立即删除文件，保留30天后清理

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Attachment deleted successfully",
    "data": null
  }
  ```

#### 启动动态题目容器

- **URL：** `POST /api/v1/challenges/:id/container/start`

- **权限：**需要登录且已组队

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 题目ID

- **请求体：**无

- **说明**：
  - 只有 `mode=dynamic` 的题目才能启动容器
  - 每个团队同一题目同时只能有一个运行中的容器
  - 容器默认运行时间：2小时，到期自动关闭
  - 容器关闭后可以重新启动（会生成新的flag）
  - 系统根据 `docker_image` 和 `docker_ports` 启动容器
  - 自动生成该团队的唯一 flag 并注入容器

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Container started successfully",
    "data": {
      "container_id": "abc123xyz",
      "challenge_id": 1001,
      "team_id": 1001,
      "status": "running",
      "url": "http://192.168.1.100:32768",
      "ports": {
        "80/tcp": "32768"
      },
      "started_at": "2025-11-20 22:55:00",
      "expires_at": "2025-11-21 00:55:00"
    }
  }
  ```

#### 停止动态题目容器

- **URL：** `POST /api/v1/challenges/:id/container/stop`

- **权限：**需要登录且已组队

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 题目ID

- **请求体：**无

- **说明**：
  - 手动停止当前团队的容器
  - 停止后可以重新启动
  - 管理员可以停止任何团队的容器

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Container stopped successfully",
    "data": {
      "container_id": "abc123xyz",
      "stopped_at": "2025-11-20 23:00:00"
    }
  }
  ```

#### 查询容器状态

- **URL：** `GET /api/v1/challenges/:id/container/status`

- **权限：**需要登录且已组队

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 题目ID

- **说明**：
  - 查询当前团队在该题目的容器状态
  - 如果没有运行中的容器，返回 `status: "not_running"`

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "container_id": "abc123xyz",
      "challenge_id": 1001,
      "team_id": 1001,
      "status": "running",
      "url": "http://192.168.1.100:32768",
      "ports": {
        "80/tcp": "32768"
      },
      "started_at": "2025-11-20 22:55:00",
      "expires_at": "2025-11-21 00:55:00",
      "time_remaining": "1h 55m"
    }
  }
  ```

#### 延长容器运行时间

- **URL：** `POST /api/v1/challenges/:id/container/renew`

- **权限：**需要登录且已组队

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 题目ID

- **请求体：**

  ```json
  {
    "duration": 60
  }
  ```

- **说明**：
  - `duration`：延长时间（分钟），可选值：30、60、120
  - 每个容器最多延长3次
  - 最长运行时间不超过6小时

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Container time extended successfully",
    "data": {
      "container_id": "abc123xyz",
      "old_expires_at": "2025-11-21 00:55:00",
      "new_expires_at": "2025-11-21 01:55:00",
      "renewals_left": 2
    }
  }
  ```

#### 查询题目解题记录（管理员）

- **URL：** `GET /api/v1/challenges/:id/solves`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **路径参数：**
  - `id` 题目ID

- **请求参数（Query）：**
  - `page` 页码（默认 1）
  - `limit` 每页数量（默认 20）

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total": 15,
      "page": 1,
      "limit": 20,
      "list": [
        {
          "team_id": 1001,
          "team_name": "银河战队",
          "solve_time": "2025-11-20 22:10:00",
          "earned_score": 500,
          "rank": 1
        },
        {
          "team_id": 1003,
          "team_name": "代码骑士",
          "solve_time": "2025-11-20 22:45:00",
          "earned_score": 450,
          "rank": 2
        }
      ]
    }
  }
  ```

#### 获取题目统计信息（管理员）

- **URL：** `GET /api/v1/challenges/statistics`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求体：**无

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total_challenges": 50,
      "visible_challenges": 45,
      "hidden_challenges": 5,
      "static_challenges": 30,
      "dynamic_challenges": 20,
      "by_direction": {
        "Web": 15,
        "Misc": 12,
        "Crypto": 10,
        "Reverse": 8,
        "Pwn": 5
      },
      "by_difficulty": {
        "easy": 20,
        "medium": 18,
        "hard": 10,
        "expert": 2
      },
      "total_solves": 520,
      "running_containers": 38
    }
  }
  ```

***

### 3. 错误码定义

| 错误码 |                  说明                  |
| :----: | :------------------------------------: |
|  200   |                  成功                  |
|  400   |      请求参数错误（如必填字段缺失）      |
|  401   |      未授权（Token缺失或格式错误）      |
|  402   |         Token已过期，请重新登录         |
|  403   |        无权限操作（权限不足）        |
|  404   |             题目不存在             |
|  409   |   题目名称已存在   |
|  430   |   团队已被封禁，无法提交答案   |
|  432   |       Flag错误       |
|  433   |      该题目已解出，无法重复提交      |
|  434   |   只有动态题目可以启动容器   |
|  435   | 容器已在运行中 |
|  436   |  容器启动失败，请稍后重试  |
|  437   |    容器不存在或已停止    |
|  438   |  附件状态异常，无法下载  |
|  439   |   附件被检测到病毒，禁止下载   |
|  440   |  文件大小超过限制（100MB）  |
|  441   |    不支持的文件类型    |
|  442   |  容器延长次数已达上限  |
|  443   |   静态题目必须填写static_flag   |
|  444   |   动态题目必须填写docker_image   |
|  445   |    题目类型不存在    |
|  500   |              服务器错误              |

#### 题目相关错误返回示例

```json
{
  "code": 432,
  "msg": "Incorrect flag, please try again",
  "data": {
    "attempts_left": null
  }
}
```

```json
{
  "code": 433,
  "msg": "Challenge already solved by your team",
  "data": {
    "solve_time": "2025-11-20 22:10:00",
    "earned_score": 500
  }
}
```

```json
{
  "code": 435,
  "msg": "Container is already running for this challenge",
  "data": {
    "container_id": "abc123xyz",
    "url": "http://192.168.1.100:32768",
    "expires_at": "2025-11-21 00:55:00"
  }
}
```

```json
{
  "code": 439,
  "msg": "Attachment contains virus or malicious code, download forbidden",
  "data": {
    "attachment_id": 2001,
    "file_name": "infected.zip",
    "detected_at": "2025-11-20 22:50:00"
  }
}
```

***

### 4. 业务逻辑说明

#### 4.1 题目创建流程

**静态题目创建：**
1. 管理员填写题目信息
2. `mode` 设为 `static`
3. 必须填写 `static_flag`
4. `docker_image` 和 `docker_ports` 为 NULL
5. 创建成功后，题目立即可用（如果 `state=visible`）

**动态题目创建：**
1. 管理员填写题目信息
2. `mode` 设为 `dynamic`
3. 必须填写 `docker_image`（Docker镜像名称）
4. 可选填写 `docker_ports`（端口映射配置）
5. `static_flag` 为 NULL
6. 系统验证 Docker 镜像是否存在
7. 创建成功后，用户启动容器时才会生成flag

#### 4.2 分数衰减机制

**衰减触发时机：**
- 每当有团队首次解出该题目时触发
- 使用公式：`current_score = max(initial_score × (decay_ratio ^ solved_count), min_score)`

**衰减示例：**
```
题目设置：
- initial_score: 500
- min_score: 100
- decay_ratio: 0.90

解题记录：
第1个团队：获得 500 分（初始分值）
第2个团队：获得 450 分（500 × 0.90）
第3个团队：获得 405 分（500 × 0.90²）
第4个团队：获得 364 分（500 × 0.90³）
...
第N个团队：获得 100 分（已降至最低分值）
```

**衰减规则：**
- 衰减比率必须在 0.50 ~ 0.99 之间
- 最低分值不能低于初始分值的 10%
- 分值衰减后向下取整
- 一旦达到最低分值，后续团队都获得最低分值

#### 4.3 动态容器管理

**容器生命周期：**
```
启动 → 运行中 → 到期/手动停止 → 已停止
  ↑                              ↓
  └────────────── 可重新启动 ──────┘
```

**容器启动流程：**
1. 团队请求启动容器
2. 系统检查：
   - 题目模式是否为 `dynamic`
   - 团队状态是否为 `active`
   - 是否已有运行中的容器（同一题目）
3. 生成唯一flag：`flag{team_${team_id}_${challenge_id}_${random_string}}`
4. 启动Docker容器，注入flag（通过环境变量或文件）
5. 分配随机端口映射
6. 返回容器访问URL
7. 设置过期时间（默认2小时）
8. 创建定时任务，到期自动清理容器

**容器管理规则：**
- 每个团队同一题目同时只能有1个运行中的容器
- 容器默认运行时间：2小时
- 可延长3次，每次30/60/120分钟
- 最长运行时间：6小时
- 容器停止后可重新启动（会生成新flag）
- 管理员可以查看和管理所有容器

**容器清理策略：**
- 到期自动停止并删除
- 用户手动停止
- 团队解散时停止所有容器
- 团队被封禁时保留容器但禁止访问
- 题目被删除时停止并删除所有相关容器

**Flag生成规则：**
- 格式：`flag{team_${team_id}_${challenge_id}_${timestamp}_${random}}`
- `team_id`：团队ID
- `challenge_id`：题目ID
- `timestamp`：容器启动时间戳（秒）
- `random`：8位随机字符串
- 示例：`flag{team_1001_1001_1732118400_a7b9c3d5}`

#### 4.4 附件管理

**附件上传流程：**
1. 管理员选择题目并上传附件
2. 系统接收文件并计算SHA256
3. 保存文件信息到数据库，状态为 `pending`
4. 对象存储模式：
   - 上传文件到对象存储服务（OSS/S3/MinIO）
   - 使用路径：`challenges/{challenge_id}/{timestamp}_{filename}`
5. 外链模式：
   - 只保存URL，不上传文件
6. 触发异步安全扫描任务（病毒扫描）
7. 扫描通过：状态更新为 `active`
8. 扫描不通过：状态更新为 `infected` 或 `error`

**附件下载流程：**
1. 用户点击下载附件
2. 系统验证权限：
   - `visibility=public`：无需登录
   - `visibility=private`：需要登录
   - `visibility=team`：需要登录且已组队
3. 检查附件状态必须为 `active`
4. 对象存储模式：
   - 生成临时签名URL（有效期1小时）
   - 返回URL给前端
   - 前端直接下载（或重定向）
5. 外链模式：
   - 直接返回外链URL
   - 或重定向到外链地址
6. 记录下载日志（可选）

**附件版本管理：**
- 同一题目可上传多个附件
- 同一附件可更新版本（如1.0 → 1.1 → 2.0）
- 旧版本保留但隐藏，只展示最新版本
- `sort_order` 字段控制显示顺序

**附件安全策略：**
- 自动病毒扫描（使用ClamAV或第三方API）
- 限制文件类型（白名单机制）
- 限制文件大小（默认100MB）
- 对象存储使用私有权限，防止盗链
- 临时URL有效期限制（1小时）
- SHA256校验防篡改

#### 4.5 答案提交验证

**提交前置检查：**
1. 用户是否已登录
2. 用户是否已加入团队
3. 团队状态是否为 `active`
   - `banned`：提示团队已被封禁
   - `disbanded`：提示团队已解散
4. 题目是否存在且状态为 `visible`
5. 是否已经解出该题（查询解题记录表）

**Flag验证规则：**
- **静态题目**：
  - 与 `static_flag` 字段比对
  - 不区分大小写
  - 去除首尾空格
  
- **动态题目**：
  - 查询该团队在该题目的容器记录
  - 比对容器生成的flag
  - 不区分大小写
  - 去除首尾空格
  - 容器停止后flag仍然有效（允许记录下来后提交）

**提交成功后的操作：**
1. 计算当前题目分值：`current_score`
2. 团队总分增加：`team_score += current_score`
3. 题目解出次数 +1：`solved_count += 1`
4. 更新题目分值：`current_score = max(initial_score × (decay_ratio ^ solved_count), min_score)`
5. 记录解题信息到解题记录表：
   - `team_id`：团队ID
   - `challenge_id`：题目ID
   - `solve_time`：解题时间
   - `earned_score`：获得分数
6. 触发排行榜更新
7. 可选：发送解题通知（Webhook、邮件等）

#### 4.6 题目可见性控制

**题目状态（state）：**
- **visible（显示）**：
  - 用户可以查看题目列表
  - 可以查看题目详情
  - 可以下载附件
  - 可以启动容器（动态题）
  - 可以提交答案
  
- **hidden（隐藏）**：
  - 用户无法在列表中看到
  - 无法查看题目详情
  - 无法下载附件
  - 无法启动容器
  - 无法提交答案
  - **管理员仍然可以查看和管理**

**使用场景：**
- 比赛开始前：所有题目设为 `hidden`
- 比赛开始时：批量设为 `visible`
- 题目临时下线：设为 `hidden`
- 题目出现问题需修复：设为 `hidden`

#### 4.7 数据一致性保证

**题目解出次数统计：**
- 每次提交成功后 `solved_count + 1`
- 定时任务每小时校准一次（从解题记录表统计）
- 避免并发问题：使用数据库事务和行锁

**团队分数更新：**
- 使用数据库事务确保原子性
- 题目分值计算和团队分数更新在同一事务中
- 失败自动回滚

**容器状态同步：**
- 定时任务每分钟检查容器状态
- 已停止的容器从运行列表中移除
- 到期的容器自动清理
- 异常容器自动重启或删除

**附件存储一致性：**
- 软删除不立即删除物理文件
- 定时任务30天后清理已删除的附件
- 对象存储使用版本控制（如果支持）

#### 4.8 性能优化建议

**数据库优化：**
- 题目列表查询：使用索引（direction、difficulty、state、current_score）
- 解题记录查询：使用复合索引（team_id、challenge_id）
- 分数计算：使用 Redis 缓存当前分值

**容器管理优化：**
- 使用容器编排平台（Kubernetes、Docker Swarm）
- 限制每个节点最大容器数
- 使用资源配额（CPU、内存限制）
- 预热常用镜像，减少启动时间

**附件下载优化：**
- 使用CDN加速附件下载
- 对象存储启用跨域配置
- 临时URL缓存到Redis（避免重复生成）
- 大文件支持断点续传

**缓存策略：**
- 题目列表：Redis缓存5分钟
- 题目详情：Redis缓存10分钟
- 排行榜：Redis缓存1分钟
- 容器状态：Redis实时更新

***



## 六、解题与日志模块

### 1. 数据库设计

***

**数据库：**`dalictf`

**表名：**`dalictf_solve`

***

#### 字段设计

|      字段名      |      类型      |           约束条件           |                  说明                   |
| :--------------: | :------------: | :--------------------------: | :-------------------------------------: |
|        id        |  `BIGINT(20)`  | `PRIMARY KEY AUTO_INCREMENT` |                解题主键                 |
|   challenge_id   |  `BIGINT(20)`  |          `NOT NULL`          |      题目ID（外键关联题目表）       |
|     team_id      |  `BIGINT(20)`  |          `NOT NULL`          |      队伍ID（外键关联团队表）       |
|     user_id      |  `BIGINT(20)`  |          `NOT NULL`          |      用户ID（外键关联用户表）       |
|   earned_score   |   `INT(11)`    |          `NOT NULL`          | 解题时获得的实际分数（用于快照记录） |
|       rank       |   `INT(11)`    |          `NOT NULL`          |   解题排名（第几个解出该题的队伍）    |
|   is_first_blood |  `TINYINT(1)`  |     `NOT NULL DEFAULT 0`     |      是否为一血（First Blood）      |
|   is_second_blood|  `TINYINT(1)`  |     `NOT NULL DEFAULT 0`     |      是否为二血（Second Blood）     |
|   is_third_blood |  `TINYINT(1)`  |     `NOT NULL DEFAULT 0`     |      是否为三血（Third Blood）      |
|   solving_time   |   `DATETIME`   | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                解题时间                 |
|    created_at    |   `DATETIME`   | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                创建时间                 |
|    updated_at    |   `DATETIME`   | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` |                更新时间                 |
|    deleted_at    |   `DATETIME`   |        `DEFAULT NULL`        |      软删除时间（NULL表示未删除）       |

#### **建表 `sql`**

```sql
CREATE TABLE `dalictf_solve` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '解题主键ID',
  `challenge_id` BIGINT(20) NOT NULL COMMENT '题目ID',
  `team_id` BIGINT(20) NOT NULL COMMENT '队伍ID',
  `user_id` BIGINT(20) NOT NULL COMMENT '用户ID',
  `earned_score` INT(11) NOT NULL COMMENT '获得分数',
  `rank` INT(11) NOT NULL COMMENT '解题排名',
  `is_first_blood` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否一血',
  `is_second_blood` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否二血',
  `is_third_blood` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否三血',
  `solving_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '解题时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_team_challenge` (`team_id`, `challenge_id`), -- 每个队伍每道题只能解出一次
  KEY `idx_challenge_id` (`challenge_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_solving_time` (`solving_time`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='解题记录表';

-- 添加外键约束
ALTER TABLE `dalictf_solve` 
ADD CONSTRAINT `fk_solve_challenge` 
FOREIGN KEY (`challenge_id`) REFERENCES `dalictf_challenge`(`id`) 
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `dalictf_solve` 
ADD CONSTRAINT `fk_solve_team` 
FOREIGN KEY (`team_id`) REFERENCES `dalictf_team`(`id`) 
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `dalictf_solve` 
ADD CONSTRAINT `fk_solve_user` 
FOREIGN KEY (`user_id`) REFERENCES `dalictf_user`(`id`) 
ON DELETE RESTRICT ON UPDATE CASCADE;
```

***

**表名：**`dalictf_submission_log`

***

#### 字段设计

|      字段名      |                 类型                  |                约束条件                 |                            说明                            |
| :--------------: | :-----------------------------------: | :-------------------------------------: | :--------------------------------------------------------: |
|        id        |            `BIGINT(20)`             |      `PRIMARY KEY AUTO_INCREMENT`       |                          日志主键                          |
|   challenge_id   |            `BIGINT(20)`             |            `NOT NULL`             |                  题目ID（外键关联题目表）                  |
|     team_id      |            `BIGINT(20)`             |            `NOT NULL`             |                  队伍ID（外键关联团队表）                  |
|     user_id      |            `BIGINT(20)`             |            `NOT NULL`             |                  用户ID（外键关联用户表）                  |
|  submitted_flag  |          `VARCHAR(500)`           |            `NOT NULL`             |                      用户提交的Flag                      |
|   flag_result    | `ENUM('correct', 'wrong', 'duplicate')` |            `NOT NULL`             | 提交结果：correct-正确/wrong-错误/duplicate-重复提交 |
|  challenge_type  |     `ENUM('static', 'dynamic')`     |            `NOT NULL`             |             题目类型：static-静态/dynamic-动态             |
|    ip_address    |           `VARCHAR(50)`           |            `NOT NULL`             |                        提交IP地址                        |
|    user_agent    |          `VARCHAR(500)`           |            `DEFAULT NULL`             |                用户浏览器UA（用于反作弊分析）                |
| submission_time  |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                          提交时间                          |
|    created_at    |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                          创建时间                          |
|    deleted_at    |            `DATETIME`             |            `DEFAULT NULL`             |               软删除时间（NULL表示未删除）               |

#### **建表 `sql`**

```sql
CREATE TABLE `dalictf_submission_log` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '日志主键ID',
  `challenge_id` BIGINT(20) NOT NULL COMMENT '题目ID',
  `team_id` BIGINT(20) NOT NULL COMMENT '队伍ID',
  `user_id` BIGINT(20) NOT NULL COMMENT '用户ID',
  `submitted_flag` VARCHAR(500) NOT NULL COMMENT '提交的Flag',
  `flag_result` ENUM('correct', 'wrong', 'duplicate') NOT NULL COMMENT '提交结果',
  `challenge_type` ENUM('static', 'dynamic') NOT NULL COMMENT '题目类型',
  `ip_address` VARCHAR(50) NOT NULL COMMENT '提交IP',
  `user_agent` VARCHAR(500) DEFAULT NULL COMMENT '用户UserAgent',
  `submission_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_challenge_id` (`challenge_id`),
  KEY `idx_team_id` (`team_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_ip_address` (`ip_address`),
  KEY `idx_submission_time` (`submission_time`),
  KEY `idx_flag_result` (`flag_result`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Flag提交日志表';
```

***

### 2. API 设计

> **统一前缀：**`/api/v1/logs`
>
> **认证机制：**JWT Token

#### 查询提交日志（管理员）

- **URL：** `GET /api/v1/logs/submissions`

- **权限：**管理员权限

- **请求头：**

  ```
  Content-Type: application/json
  Authorization: Bearer <JWT_TOKEN>
  ```

- **请求参数（Query）：**
  - `page` 页码
  - `limit` 每页数量
  - `challenge_id` 题目筛选
  - `team_id` 队伍筛选
  - `user_id` 用户筛选
  - `flag_result` 结果筛选（correct/wrong）
  - `search` 模糊搜索（Flag内容、IP）

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total": 1000,
      "page": 1,
      "limit": 20,
      "list": [
        {
          "id": 5001,
          "challenge_name": "Web签到",
          "team_name": "银河战队",
          "username": "zhangsan",
          "submitted_flag": "flag{wrong_flag}",
          "flag_result": "wrong",
          "ip_address": "192.168.1.101",
          "submission_time": "2025-11-20 22:05:00"
        }
      ]
    }
  }
  ```

#### 导出提交日志（管理员）

- **URL：** `GET /api/v1/logs/export`

- **权限：**管理员权限

- **请求参数（Query）：**
  - 同查询接口，用于筛选导出范围

- **说明**：
  - 导出格式为 CSV 或 Excel
  - 包含字段：ID、题目、队伍、用户、提交Flag、结果、IP、UA、提交时间

- **成功返回：**
  - 直接返回文件流（Content-Type: application/vnd.ms-excel 或 text/csv）

#### 查询可疑作弊记录（管理员）

- **URL：** `GET /api/v1/logs/suspicious`

- **权限：**管理员权限

- **请求参数（Query）：**
  - `type` 作弊类型（ip_sharing/flag_brute_force/flag_sharing）
  - `page` 页码
  - `limit` 每页数量

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total": 5,
      "list": [
        {
          "team_name": "作弊小队",
          "user_name": "hacker",
          "type": "flag_sharing",
          "description": "Submitted flag belonging to Team B",
          "detected_at": "2025-11-20 23:00:00"
        }
      ]
    }
  }
  ```

#### 查询解题动态（公开/实时）

- **URL：** `GET /api/v1/logs/solves`

- **权限：**公开

- **请求参数（Query）：**
  - `limit` 返回数量（默认10，用于轮播展示）
  - `challenge_id` 筛选特定题目

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": [
      {
        "team_name": "银河战队",
        "challenge_name": "Web签到",
        "earned_score": 100,
        "rank": 1,
        "is_first_blood": true,
        "solving_time": "2025-11-20 22:00:05"
      }
    ]
  }
  ```

#### 查询特定队伍解题记录（公开）

- **URL：** `GET /api/v1/teams/:team_id/solves`

- **权限：**公开

- **路径参数：**
  - `team_id` 队伍ID

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": [
      {
        "challenge_id": 1001,
        "challenge_name": "Web签到",
        "earned_score": 100,
        "solving_time": "2025-11-20 22:00:05",
        "is_first_blood": true
      }
    ]
  }
  ```

#### 查询特定用户解题记录（公开）

- **URL：** `GET /api/v1/users/:user_id/solves`

- **权限：**公开

- **路径参数：**
  - `user_id` 用户ID

- **成功返回：**（格式同队伍解题记录）

#### 获取提交统计信息（管理员）

- **URL：** `GET /api/v1/logs/statistics`

- **权限：**管理员权限

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total_submissions": 5000,
      "correct_submissions": 500,
      "wrong_submissions": 4500,
      "correct_rate": "10.0%",
      "submissions_by_hour": [ ... ] // 用于绘制图表
    }
  }
  ```

***

### 3. 业务逻辑说明

#### 3.1 一血、二血、三血机制

- **定义**：每道题目最先解出的前三名队伍分别获得一血（First Blood）、二血（Second Blood）、三血（Third Blood）荣誉。
- **实现逻辑**：
  1. 用户提交正确 Flag。
  2. 系统查询 `dalictf_solve` 表中该题目的解题记录数量 `count`。
  3. 如果 `count == 0`，标记 `is_first_blood = 1`。
  4. 如果 `count == 1`，标记 `is_second_blood = 1`。
  5. 如果 `count == 2`，标记 `is_third_blood = 1`。
  6. 记录 `rank = count + 1`。
- **展示**：前端在题目列表和动态中高亮显示一二三血成就，给予特殊UI标识。

#### 3.2 防作弊日志分析

`dalictf_submission_log` 表是反作弊系统的核心数据源，主要分析维度包括：

1.  **IP 异常分析**：
    -   同一队伍不同成员 IP 跨度过大（异地登录）。
    -   同一 IP 登录了多个不同队伍的账号（疑似小号刷分）。
    -   短时间内 IP 频繁变动。

2.  **Flag 爆破检测**：
    -   同一用户/队伍短时间内提交大量错误 Flag。
    -   监控 `flag_result = 'wrong'` 的频率，超过阈值（如 1分钟 10次）触发警告或临时封禁。

3.  **Flag 撞库/分享检测**：
    -   检测是否有队伍提交了**其他队伍的动态 Flag**（动态题 Flag 是唯一的）。
    -   如果 Team A 提交了 Team B 的 Flag，系统应立即告警，判定为作弊（Flag 分享）。
    -   记录 `submitted_flag` 即使是错误的，也用于后续比对分析。

4.  **User-Agent 分析**：
    -   检测使用自动化脚本工具的 UA（如 Python-requests, sqlmap 等），除非题目允许。

#### 3.3 分数快照

- `dalictf_solve` 表中的 `earned_score` 字段记录的是**解题时刻**该题目的分值。
- 即使后续该题目因解题人数增加而发生了分数衰减，**已解题队伍的得分不应改变**（或者根据赛制决定是否动态更新）。
- **通常赛制（动态积分）**：所有解出该题的队伍得分都会随题目当前分值变化而变化。
  - 这种情况下，`earned_score` 仅作为参考记录，实际计算总分时应关联 `dalictf_challenge.current_score`。
  - **公式**：`Team Total Score = Sum(Challenge.current_score) where Team solved Challenge`。
- **本系统采用动态积分制**，因此排行榜计算时应实时聚合。

***



## 七、容器管理模块

### 1. 数据库设计

***

**数据库：**`dalictf`

**表名：**`dalictf_container`

***

#### 字段设计

|      字段名      |                 类型                  |                约束条件                 |                            说明                            |
| :--------------: | :-----------------------------------: | :-------------------------------------: | :--------------------------------------------------------: |
|        id        |            `BIGINT(20)`             |      `PRIMARY KEY AUTO_INCREMENT`       |                          容器主键                          |
|   challenge_id   |            `BIGINT(20)`             |            `NOT NULL`             |                  题目ID（外键关联题目表）                  |
|     team_id      |            `BIGINT(20)`             |            `NOT NULL`             |                  队伍ID（外键关联团队表）                  |
|     user_id      |            `BIGINT(20)`             |            `NOT NULL`             |              启动用户ID（外键关联用户表，记录操作人）              |
|  container_name  |          `VARCHAR(255)`           |            `NOT NULL`             |           容器实际名称（如 `ctf_challenge_201_5`）           |
|   docker_image   |          `VARCHAR(255)`           |            `NOT NULL`             |                          启动镜像                          |
|   docker_ports   |              `JSON`               |            `DEFAULT NULL`             |          暴露端口（JSON 格式，如 `{"1337":"tcp"}`）          |
|   host_mapping   |              `JSON`               |            `DEFAULT NULL`             |          主机端口映射（JSON 格式，如 `{"1337":32768}`）          |
|  container_flag  |          `VARCHAR(500)`           |            `NOT NULL`             |             动态 flag（系统生成，写入容器环境变量）             |
|      state       | `ENUM('running', 'stopped', 'destroyed')` |    `NOT NULL DEFAULT 'running'`    |        容器状态：running-运行中/stopped-已停止/destroyed-已销毁        |
|    start_time    |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                          启动时间                          |
|     end_time     |            `DATETIME`             |            `NOT NULL`             |             到期销毁时间（默认启动后 1 小时）             |
|  extended_count  |            `TINYINT(4)`             |        `NOT NULL DEFAULT 0`         |                 续期次数（限制最多续期3次）                 |
|    created_at    |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                          创建时间                          |
|    updated_at    |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` |                          更新时间                          |
|    deleted_at    |            `DATETIME`             |            `DEFAULT NULL`             |               软删除时间（NULL表示未删除）               |

#### **建表 `sql`**

```sql
CREATE TABLE `dalictf_container` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '容器主键ID',
  `challenge_id` BIGINT(20) NOT NULL COMMENT '题目ID',
  `team_id` BIGINT(20) NOT NULL COMMENT '队伍ID',
  `user_id` BIGINT(20) NOT NULL COMMENT '启动用户ID',
  `container_name` VARCHAR(255) NOT NULL COMMENT '容器实际名称',
  `docker_image` VARCHAR(255) NOT NULL COMMENT '启动镜像',
  `docker_ports` JSON DEFAULT NULL COMMENT '暴露端口配置',
  `host_mapping` JSON DEFAULT NULL COMMENT '主机端口映射',
  `container_flag` VARCHAR(500) NOT NULL COMMENT '动态Flag',
  `state` ENUM('running', 'stopped', 'destroyed') NOT NULL DEFAULT 'running' COMMENT '容器状态',
  `start_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '启动时间',
  `end_time` DATETIME NOT NULL COMMENT '到期时间',
  `extended_count` TINYINT(4) NOT NULL DEFAULT 0 COMMENT '续期次数',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_challenge_id` (`challenge_id`),
  KEY `idx_team_id` (`team_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_state` (`state`),
  KEY `idx_end_time` (`end_time`),
  KEY `idx_container_name` (`container_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='容器实例表';

-- 添加外键约束
ALTER TABLE `dalictf_container` 
ADD CONSTRAINT `fk_container_challenge` 
FOREIGN KEY (`challenge_id`) REFERENCES `dalictf_challenge`(`id`) 
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `dalictf_container` 
ADD CONSTRAINT `fk_container_team` 
FOREIGN KEY (`team_id`) REFERENCES `dalictf_team`(`id`) 
ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE `dalictf_container` 
ADD CONSTRAINT `fk_container_user` 
FOREIGN KEY (`user_id`) REFERENCES `dalictf_user`(`id`) 
ON DELETE RESTRICT ON UPDATE CASCADE;
```

***

### 2. API 设计

> **统一前缀：**`/api/v1`
>
> **认证机制：**JWT Token

#### 2.1 用户侧 API

> 用户操作自己的容器实例。

#### 启动容器

- **URL：** `POST /api/v1/containers/start`

- **权限：**需登录，且已加入团队（如需）

- **请求体：**

  ```json
  {
    "challenge_id": 1001
  }
  ```

- **说明**：
  - 检查用户/团队是否已有运行中容器（通常限制每队每题或全局只能开一个）。
  - 分配端口，启动 Docker，写入 Flag。
  - 返回连接信息。

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "container_id": 123,
      "host": "1.2.3.4",
      "port": 32768,
      "remaining_time": 3600
    }
  }
  ```

#### 销毁容器

- **URL：** `POST /api/v1/containers/destroy`

- **权限：**需登录，且是容器所有者

- **请求体：** `{"challenge_id": 1001}` 或 `{"container_id": 123}`

- **说明**：
  - 用户主动销毁容器，释放资源。
  - 销毁后无法恢复，需重新启动。

#### 续期容器

- **URL：** `POST /api/v1/containers/renew`

- **权限：**需登录，且是容器所有者

- **请求体：** `{"container_id": 123}`

- **说明**：
  - 延长容器生命周期（如增加 1 小时）。
  - 检查剩余续期次数。

#### 查询容器状态

- **URL：** `GET /api/v1/containers/status`

- **权限：**需登录

- **请求参数：** `challenge_id`

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "status": "running",
      "host": "1.2.3.4",
      "port": 32768,
      "remaining_time": 2500,
      "extended_count": 1
    }
  }
  ```

#### 2.2 管理员 API

> **前缀：**`/api/v1/admin/containers`
>
> **权限：**管理员权限

#### 查询容器列表（管理员）

- **URL：** `GET /api/v1/admin/containers`

- **请求参数（Query）：**
  - `page` 页码
  - `limit` 每页数量
  - `state` 状态筛选（running/stopped/destroyed）
  - `challenge_id` 题目筛选
  - `team_id` 队伍筛选
  - `search` 模糊搜索（容器名、Flag、镜像名）

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total": 50,
      "page": 1,
      "limit": 20,
      "list": [
        {
          "id": 1001,
          "container_name": "ctf_web_1001_team_5",
          "challenge_name": "Web签到",
          "team_name": "银河战队",
          "docker_image": "ctf/web-sign:latest",
          "state": "running",
          "host_mapping": {"80": 32768},
          "start_time": "2025-11-20 22:00:00",
          "end_time": "2025-11-20 23:00:00",
          "time_remaining": "45m"
        }
      ]
    }
  }
  ```

#### 强制停止容器（管理员）

- **URL：** `POST /api/v1/admin/containers/:id/stop`

- **说明**：
  - 管理员强制停止某个异常或违规的容器
  - 停止后状态变为 `stopped`
  - 并不立即删除数据，但容器实例会被 kill

#### 强制销毁容器（管理员）

- **URL：** `DELETE /api/v1/admin/containers/:id`

- **说明**：
  - 强制删除容器实例和相关数据
  - 状态变为 `destroyed`
  - 释放占用的端口和资源

#### 清理过期容器（系统任务）

- **说明**：这是一个内部定时任务，不通过 HTTP API 触发（或仅限 localhost 调用）。
- **逻辑**：
  1. 扫描 `dalictf_container` 表中 `state = 'running'` 且 `end_time < NOW()` 的记录。
  2. 调用 Docker API 停止并删除容器。
  3. 更新数据库状态为 `destroyed`。
  4. 释放端口资源。

***

### 3. 业务逻辑说明

#### 3.1 容器生命周期管理

1.  **启动（Start）**：
    -   用户请求启动 -> 检查权限和资源 -> 分配端口 -> 生成 Flag -> 启动 Docker 容器 -> 记录数据库。
    -   **端口分配**：维护一个可用端口池（如 30000-40000），启动时随机选取可用端口，避免冲突。

2.  **运行（Running）**：
    -   容器处于活跃状态，用户可通过映射端口访问。
    -   系统定期（如每分钟）同步 Docker 实际状态与数据库状态，防止状态不一致。

3.  **续期（Renew）**：
    -   用户请求续期 -> 检查剩余次数 (`extended_count < 3`) -> 增加 `end_time` -> 更新数据库。
    -   每次续期增加时长通常为 30 或 60 分钟。

4.  **停止（Stop）**：
    -   用户主动停止或管理员强制停止。
    -   容器被 kill，但保留数据库记录以便审计。
    -   用户可重新启动（Re-start），此时通常会创建新容器实例（新 Flag）。

5.  **销毁（Destroy）**：
    -   过期自动销毁或管理员删除。
    -   彻底移除 Docker 容器实例。
    -   数据库记录标记为 `destroyed`。

#### 3.2 动态 Flag 生成策略

-   **格式**：`flag{uuid}` 或 `flag{team_hash_challenge_hash_random}`。
-   **注入方式**：
    -   **环境变量**：启动容器时通过 `-e FLAG=xxx` 传入。题目镜像需编写启动脚本读取该环境变量并写入 flag 文件（如 `/flag`）。
    -   **文件挂载**：系统生成 flag 文件，通过 `-v /host/path/flag:/flag` 挂载到容器内。

#### 3.3 资源限制与安全

-   **资源配额**：
    -   CPU：限制如 0.5 core。
    -   内存：限制如 512MB。
    -   防止恶意容器耗尽服务器资源。
-   **网络隔离**：
    -   容器应运行在独立的 Docker Network 中。
    -   禁止容器访问宿主机内网（通过 iptables 规则）。
    -   禁止容器间相互访问（除非特定题目需要）。
-   **只读文件系统**：
    -   对于不需要写入的题目，可挂载为只读模式，防止被篡改。

***



## 八、公告模块

### 1. 数据库设计

***

**数据库：**`dalictf`

**表名：**`dalictf_notice`

***

#### 字段设计

|      字段名      |                 类型                  |                约束条件                 |                            说明                            |
| :--------------: | :-----------------------------------: | :-------------------------------------: | :--------------------------------------------------------: |
|        id        |            `BIGINT(20)`             |      `PRIMARY KEY AUTO_INCREMENT`       |                          公告主键                          |
|      title       |          `VARCHAR(255)`           |            `NOT NULL`             |                          公告标题                          |
|     content      |              `TEXT`               |            `NOT NULL`             |                  公告内容（支持 Markdown）                  |
|      is_top      |            `TINYINT(1)`             |        `NOT NULL DEFAULT 0`         |                 是否置顶（1-置顶，0-普通）                 |
|      status      | `ENUM('published', 'draft', 'archived')` |   `NOT NULL DEFAULT 'published'`    |        状态：published-发布/draft-草稿/archived-归档        |
|    created_by    |            `BIGINT(20)`             |            `NOT NULL`             |                 创建人ID（外键关联用户表）                 |
|    created_at    |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                          创建时间                          |
|    updated_at    |            `DATETIME`             | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` |                          更新时间                          |
|    deleted_at    |            `DATETIME`             |            `DEFAULT NULL`             |               软删除时间（NULL表示未删除）               |

#### **建表 `sql`**

```sql
CREATE TABLE `dalictf_notice` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '公告主键ID',
  `title` VARCHAR(255) NOT NULL COMMENT '公告标题',
  `content` TEXT NOT NULL COMMENT '公告内容',
  `is_top` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否置顶',
  `status` ENUM('published', 'draft', 'archived') NOT NULL DEFAULT 'published' COMMENT '公告状态',
  `created_by` BIGINT(20) NOT NULL COMMENT '创建人ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '软删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_top` (`is_top`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统公告表';

-- 添加外键约束
ALTER TABLE `dalictf_notice` 
ADD CONSTRAINT `fk_notice_user` 
FOREIGN KEY (`created_by`) REFERENCES `dalictf_user`(`id`) 
ON DELETE RESTRICT ON UPDATE CASCADE;
```

***

### 2. API 设计

> **统一前缀：**`/api/v1/notices`
>
> **认证机制：**JWT Token（部分接口公开）

#### 获取公告列表（公开）

- **URL：** `GET /api/v1/notices`

- **权限：**公开

- **请求参数（Query）：**
  - `page` 页码
  - `limit` 每页数量
  - `search` 关键词搜索

- **排序规则**：
  - 优先按 `is_top` 降序（置顶在前）。
  - 其次按 `created_at` 降序（最新在前）。

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "total": 5,
      "list": [
        {
          "id": 1,
          "title": "比赛延期通知",
          "is_top": 1,
          "created_at": "2025-11-20 10:00:00"
        }
      ]
    }
  }
  ```

#### 获取公告详情（公开）

- **URL：** `GET /api/v1/notices/:id`

- **权限：**公开

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "id": 1,
      "title": "比赛延期通知",
      "content": "由于不可抗力...",
      "is_top": 1,
      "created_by_name": "Admin",
      "created_at": "2025-11-20 10:00:00"
    }
  }
  ```

#### 发布/修改公告（管理员）

- **URL：** `POST /api/v1/admin/notices` （新增）
- **URL：** `PUT /api/v1/admin/notices/:id` （修改）

- **权限：**管理员权限

- **请求体：**

  ```json
  {
    "title": "新公告",
    "content": "公告内容...",
    "is_top": 0,
    "status": "published"
  }
  ```

#### 删除公告（管理员）

- **URL：** `DELETE /api/v1/admin/notices/:id`

- **权限：**管理员权限

- **说明**：软删除。

#### 置顶/取消置顶（管理员）

- **URL：** `PATCH /api/v1/admin/notices/:id/top`

- **权限：**管理员权限

- **请求体：** `{"is_top": 1}`

***

## 九、比赛配置模块

### 1. 数据库设计

***

**数据库：**`dalictf`

**表名：**`dalictf_config`

***

#### 字段设计

|      字段名      |      类型      |           约束条件           |                  说明                   |
| :--------------: | :------------: | :--------------------------: | :-------------------------------------: |
|        id        |  `BIGINT(20)`  | `PRIMARY KEY AUTO_INCREMENT` |                配置主键                 |
|    config_key    | `VARCHAR(100)` |          `NOT NULL`          |      配置键名（如 `start_time`）      |
|   config_value   |     `TEXT`     |        `DEFAULT NULL`        |                配置键值                 |
|   description    | `VARCHAR(255)` |        `DEFAULT NULL`        |                配置说明                 |
|    created_at    |   `DATETIME`   | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |                创建时间                 |
|    updated_at    |   `DATETIME`   | `NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP` |                更新时间                 |

#### **建表 `sql`**

```sql
CREATE TABLE `dalictf_config` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '配置主键ID',
  `config_key` VARCHAR(100) NOT NULL COMMENT '配置键名',
  `config_value` TEXT DEFAULT NULL COMMENT '配置键值',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '配置说明',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 初始化基础配置数据
INSERT INTO `dalictf_config` (`config_key`, `config_value`, `description`) VALUES
('competition_name', 'ISCTF 2025', '比赛名称'),
('competition_start_time', '2025-11-20 09:00:00', '比赛开始时间'),
('competition_end_time', '2025-11-22 09:00:00', '比赛结束时间'),
('registration_start_time', '2025-11-01 00:00:00', '注册开始时间'),
('registration_end_time', '2025-11-22 09:00:00', '注册结束时间'),
('is_paused', 'false', '比赛是否暂停'),
('announcement', '欢迎参加 ISCTF 2025！', '首页滚动公告');
```

***

### 2. API 设计

> **统一前缀：**`/api/v1/config`
>
> **认证机制：**JWT Token（部分接口公开）

#### 获取公开配置信息（公开）

- **URL：** `GET /api/v1/config`

- **权限：**公开

- **说明**：返回前端展示所需的基础信息（比赛时间、名称等），敏感配置不返回。

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "competition_name": "ISCTF 2025",
      "competition_start_time": "2025-11-20 09:00:00",
      "competition_end_time": "2025-11-22 09:00:00",
      "registration_start_time": "2025-11-01 00:00:00",
      "registration_end_time": "2025-11-22 09:00:00",
      "is_paused": "false",
      "server_time": "2025-11-15 12:00:00" // 返回服务器当前时间，用于前端倒计时校准
    }
  }
  ```

#### 获取全部配置（管理员）

- **URL：** `GET /api/v1/admin/config`

- **权限：**管理员权限

- **成功返回：** 返回所有键值对。

#### 修改配置（管理员）

- **URL：** `PUT /api/v1/admin/config`

- **权限：**管理员权限

- **请求体：**

  ```json
  {
    "competition_start_time": "2025-12-01 09:00:00",
    "is_paused": "true"
  }
  ```

- **说明**：支持批量修改。

***

### 3. 业务逻辑说明

#### 3.1 时间控制拦截器

系统应在中间件层（Middleware）实现基于配置的时间控制逻辑：

1.  **注册拦截**：
    -   在调用注册接口时，检查 `NOW()` 是否在 `[registration_start_time, registration_end_time]` 区间内。
    -   如果不在区间，返回 `403 Forbidden`，提示“注册未开放”或“注册已截止”。

2.  **比赛拦截**：
    -   在调用获取题目详情、启动容器、提交 Flag 等比赛相关接口时，检查 `NOW()` 是否在 `[competition_start_time, competition_end_time]` 区间内。
    -   **未开始**：只能查看题目列表（名称/分值），不能查看详情/附件/启动容器。
    -   **已结束**：可以查看题目详情，但**不能提交 Flag**（或提交但不计分，视赛制而定），**不能启动新容器**。
    -   **暂停中** (`is_paused = true`)：全站进入维护模式或暂停计分。

#### 3.2 前端倒计时同步

-   前端应使用 `/api/v1/config` 返回的 `server_time` 计算与本地时间的差值（Offset），确保倒计时以服务器时间为准，防止客户端时间篡改导致显示异常。
#### 3.2 前端倒计时同步

-   前端应使用 `/api/v1/config` 返回的 `server_time` 计算与本地时间的差值（Offset），确保倒计时以服务器时间为准，防止客户端时间篡改导致显示异常。

***



## 十、比赛大屏模块

### 1. 概述

比赛大屏（Dashboard）主要用于比赛现场或线上直播展示实时战况。该模块的核心需求是**高并发读取**和**实时性**。

**设计原则：**
1.  **读写分离**：大屏接口只读，且数据应高度聚合。
2.  **缓存优先**：所有大屏接口必须经过 Redis 缓存，缓存时间根据业务实时性要求设置（1s - 10s）。
3.  **数据脱敏**：大屏展示的数据通常是公开的，需注意隐藏敏感信息（如用户手机号、邮箱）。

### 2. API 设计

> **统一前缀：**`/api/v1/dashboard`
>
> **认证机制：**公开（或仅限特定 Token 访问，视部署需求而定，默认公开）

#### 获取大屏聚合数据（核心接口）

- **URL：** `GET /api/v1/dashboard/overview`

- **权限：**公开

- **缓存策略**：Redis 缓存 5 秒

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "competition_status": "running", // running, paused, ended
      "countdown": 3600, // 距离结束剩余秒数
      "stats": {
        "total_teams": 120,
        "active_teams": 85, // 有提交记录的队伍
        "total_solves": 450,
        "total_flags_submitted": 1200
      },
      "top_teams": [ // 前 10 名队伍简要信息
        {
          "rank": 1,
          "team_name": "银河战队",
          "avatar": "url...",
          "score": 2500,
          "school_name": "大理大学"
        },
        ...
      ],
      "recent_solves": [ // 最近 5 条解题记录（滚动播报）
        {
          "team_name": "银河战队",
          "challenge_name": "Web签到",
          "time": "10:05:00",
          "is_first_blood": true
        }
      ],
      "category_progress": [ // 各分类解题进度
        {
          "name": "Web",
          "solved_count": 50,
          "total_count": 10
        },
        {
          "name": "Pwn",
          "solved_count": 12,
          "total_count": 5
        }
      ]
    }
  }
  ```

#### 获取得分趋势图数据

- **URL：** `GET /api/v1/dashboard/trend`

- **权限：**公开

- **缓存策略**：Redis 缓存 30 秒

- **请求参数：**
  - `top`: 返回前几名队伍的趋势（默认 10）

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": {
      "times": ["09:00", "10:00", "11:00", "12:00"], // X轴时间点
      "series": [
        {
          "team_name": "银河战队",
          "scores": [0, 500, 1200, 2500] // Y轴分数点
        },
        {
          "team_name": "星际穿越",
          "scores": [0, 300, 800, 1500]
        }
      ]
    }
  }
  ```

#### 获取地图可视化数据（可选）

- **URL：** `GET /api/v1/dashboard/map`

- **权限：**公开

- **说明**：返回各参赛学校的地理位置分布和活跃度，用于绘制热力图。

- **成功返回：**

  ```json
  {
    "code": 200,
    "msg": "Success",
    "data": [
      {
        "school_name": "大理大学",
        "latitude": 25.6,
        "longitude": 100.2,
        "team_count": 5,
        "total_score": 3000
      }
    ]
  }
  ```

***

### 3. 业务逻辑说明

#### 3.1 数据聚合与缓存

-   **聚合计算**：大屏数据涉及多表关联（Team, User, Solve, Challenge, School），直接查询数据库会造成巨大压力。
-   **缓存更新机制**：
    1.  **被动更新**：API 请求时检查 Redis 缓存，过期则重新计算并写入。
    2.  **主动更新（推荐）**：
        -   当有 `Submit Flag` 成功事件发生时，触发异步任务更新 `Overview` 缓存。
        -   或设置定时任务（每 5 秒）在后台计算好 JSON 存入 Redis，API 直接读取，实现 0 延迟响应。

#### 3.2 趋势图采样

-   由于解题记录可能非常密集，前端绘制折线图不需要所有点。
-   **采样算法**：
    -   按固定时间间隔（如每 10 分钟）取一个快照点。
    -   或者只记录分数发生变化的时刻。
    -   建议后端处理好采样数据，直接返回给前端渲染。

#### 3.3 实时推送（WebSocket 扩展）

-   如果对实时性要求极高（毫秒级），可引入 WebSocket 服务。
-   **Topic 设计**：
    -   `dashboard:overview`：推送总分、排名变化。
    -   `dashboard:log`：推送最新解题日志（弹幕效果）。
-   **实现**：当用户提交 Flag 成功后，后端向 MQ 发送消息，WebSocket 服务订阅 MQ 并广播给所有连接的大屏客户端。

## 十一、反作弊检测模块

### 1. 数据库设计

反作弊需要多个子表来记录行为轨迹。

***

#### (1) 容器抓包记录表

**表名**：`dalictf_container_packet_log`

| 字段名            | 类型         | 约束条件                   | 说明                              |
| ----------------- | ------------ | -------------------------- | --------------------------------- |
| `id`              | `INT UNSIGNED` | `PK, AUTO_INCREMENT`         | 唯一 ID                           |
| `container_id`    | `BIGINT(20)` | `NOT NULL` | 容器 ID（外键关联容器表）                           |
| `team_id`         | `BIGINT(20)`      | `NOT NULL`      | 队伍 ID（外键关联团队表）                           |
| `pcap_path`       | `VARCHAR(255)` | `NOT NULL`                   | 抓包文件存储路径                  |
| `analysis_result` | `TEXT`         | `DEFAULT NULL`                       | 自动流量分析摘要（如外联可疑 IP） |
| `created_at`      | `DATETIME`    | `NOT NULL DEFAULT CURRENT_TIMESTAMP`  | 生成时间                          |

```sql
CREATE TABLE `dalictf_container_packet_log` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '唯一ID',
  `container_id` BIGINT(20) NOT NULL COMMENT '容器ID',
  `team_id` BIGINT(20) NOT NULL COMMENT '队伍ID',
  `pcap_path` VARCHAR(255) NOT NULL COMMENT '抓包文件路径',
  `analysis_result` TEXT DEFAULT NULL COMMENT '分析摘要',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '生成时间',
  PRIMARY KEY (`id`),
  KEY `idx_container_id` (`container_id`),
  KEY `idx_team_id` (`team_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='容器抓包记录表';
```

***

#### (2) 附件下载记录表

**表名**：`dalictf_attachment_download`

| 字段名          | 类型            | 约束条件                   | 说明     |
| --------------- | --------------- | -------------------------- | -------- |
| `id`            | `BIGINT UNSIGNED` | `PK, AUTO_INCREMENT`         | 唯一 ID  |
| `challenge_id`  | `BIGINT(20)`    | `NOT NULL` | 题目 ID（外键关联题目表）  |
| `team_id`       | `BIGINT(20)`    | `NOT NULL`      | 队伍 ID（外键关联团队表）  |
| `user_id`       | `BIGINT(20)`    | `NOT NULL`      | 用户 ID（外键关联用户表）  |
| `download_time` | `DATETIME`       | `NOT NULL DEFAULT CURRENT_TIMESTAMP`  | 下载时间 |
| `ip_address`    | `VARCHAR(45)`     | `NOT NULL`                   | 下载 IP  |

```sql
CREATE TABLE `dalictf_attachment_download` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '唯一ID',
  `challenge_id` BIGINT(20) NOT NULL COMMENT '题目ID',
  `team_id` BIGINT(20) NOT NULL COMMENT '队伍ID',
  `user_id` BIGINT(20) NOT NULL COMMENT '用户ID',
  `download_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '下载时间',
  `ip_address` VARCHAR(45) NOT NULL COMMENT '下载IP',
  PRIMARY KEY (`id`),
  KEY `idx_challenge_id` (`challenge_id`),
  KEY `idx_team_id` (`team_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_ip_address` (`ip_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='附件下载记录表';
```

***

#### (3) 用户 IP 变动记录表

**表名**：`dalictf_user_ip_log`

| 字段名       | 类型                   | 约束条件                  | 说明             |
| ------------ | ---------------------- | ------------------------- | ---------------- |
| `id`         | `BIGINT UNSIGNED`        | `PK, AUTO_INCREMENT`        | 唯一 ID          |
| `user_id`    | `BIGINT(20)`           | `NOT NULL`     | 用户 ID（外键关联用户表）          |
| `ip_address` | `VARCHAR(45)`            | `NOT NULL`                  | 登录/提交时的 IP |
| `action`     | `ENUM('login','submit')` | `NOT NULL`                  | 行为类型         |
| `time`       | `DATETIME`              | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 发生时间         |

```sql
CREATE TABLE `dalictf_user_ip_log` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '唯一ID',
  `user_id` BIGINT(20) NOT NULL COMMENT '用户ID',
  `ip_address` VARCHAR(45) NOT NULL COMMENT 'IP地址',
  `action` ENUM('login','submit') NOT NULL COMMENT '行为类型',
  `time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发生时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_ip_address` (`ip_address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户IP变动记录表';
```

***

#### (4) 可疑行为告警表

**表名**：`dalictf_suspicious_activity`

| 字段名             | 类型                                                         | 约束条件                  | 说明                |
| ------------------ | ------------------------------------------------------------ | ------------------------- | ------------------- |
| `id`               | `BIGINT UNSIGNED`                                              | `PK, AUTO_INCREMENT`        | 唯一 ID             |
| `activity_type`    | `ENUM('flag_sharing','multi_ip_download','same_device','traffic_anomaly')` | `NOT NULL`                  | 可疑类型            |
| `description`      | `TEXT`                                                         | `NOT NULL`                  | 详细描述            |
| `related_team_ids` | `VARCHAR(255)`                                                 | `DEFAULT NULL`                      | 相关队伍 ID（JSON） |
| `created_at`       | `DATETIME`                                                    | `NOT NULL DEFAULT CURRENT_TIMESTAMP` | 记录时间            |
| `status`           | `ENUM('pending','reviewed','dismissed','confirmed_cheating')`  | `NOT NULL DEFAULT 'pending'`         | 审核状态            |

```sql
CREATE TABLE `dalictf_suspicious_activity` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '唯一ID',
  `activity_type` ENUM('flag_sharing','multi_ip_download','same_device','traffic_anomaly') NOT NULL COMMENT '可疑类型',
  `description` TEXT NOT NULL COMMENT '详细描述',
  `related_team_ids` VARCHAR(255) DEFAULT NULL COMMENT '相关队伍ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录时间',
  `status` ENUM('pending','reviewed','dismissed','confirmed_cheating') NOT NULL DEFAULT 'pending' COMMENT '审核状态',
  PRIMARY KEY (`id`),
  KEY `idx_activity_type` (`activity_type`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='可疑行为告警表';
```

***

### 2. API 设计

> **接口前缀：**`/api/v1/admin/anti-cheat`
>
> **权限：**仅管理员

#### (1) 查询可疑 Flag 提交（相同动态 flag 不同队伍提交）

- **URL**：`GET /api/v1/admin/anti-cheat/suspicious-flags`
- **请求参数**：`flag` (Flag值)
- **返回示例**：

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "flag_value": "flag{dyn_123xyz}",
    "submissions": [
      {
        "team_id": 201,
        "team_name": "Hackers",
        "user_id": 101,
        "username": "alice",
        "submission_time": "2025-08-26 10:00:00",
        "ip_address": "192.168.1.10"
      },
      {
        "team_id": 202,
        "team_name": "CryptoWarriors",
        "user_id": 102,
        "username": "bob",
        "submission_time": "2025-08-26 10:01:30",
        "ip_address": "192.168.1.20"
      }
    ]
  }
}
```

#### (2) 查询附件下载行为

- **URL**：`GET /api/v1/admin/anti-cheat/attachment-downloads`
- **请求参数**：`challenge_id` (题目ID)
- **返回示例**：

```json
{
  "code": 200,
  "msg": "success",
  "data": [
    {
      "team_id": 201,
      "team_name": "Hackers",
      "user_id": 101,
      "username": "alice",
      "download_time": "2025-08-26 09:50:00",
      "ip_address": "10.0.0.5"
    },
    {
      "team_id": 202,
      "team_name": "CryptoWarriors",
      "user_id": 102,
      "username": "bob",
      "download_time": "2025-08-26 09:51:00",
      "ip_address": "10.0.0.5"
    }
  ]
}
```

> ⚠️ 系统可以自动检测“同一 IP 短时间内下载多个队伍的附件”，并生成一条 `dalictf_suspicious_activity` 记录。

#### (3) 查询用户 IP 变动

- **URL**：`GET /api/v1/admin/anti-cheat/user-ips`
- **请求参数**：`user_id` (用户ID)
- **返回示例**：

```json
{
  "code": 200,
  "msg": "success",
  "data": [
    {
      "action": "login",
      "ip_address": "192.168.1.10",
      "time": "2025-08-26 08:30:00"
    },
    {
      "action": "submit",
      "ip_address": "10.0.0.5",
      "time": "2025-08-26 09:55:00"
    }
  ]
}
```

#### (4) 查询容器抓包分析结果

- **URL**：`GET /api/v1/admin/anti-cheat/containers/:id/pcap`
- **返回示例**：

```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "pcap_path": "/var/log/ctf-pcap/2025/301.pcap",
    "analysis_result": "Detected suspicious connection to 8.8.8.8"
  }
}
```

#### (5) 查询可疑行为告警列表

- **URL**：`GET /api/v1/admin/anti-cheat/alerts`
- **返回示例**：

```json
{
  "code": 200,
  "msg": "success",
  "data": [
    {
      "id": 1,
      "activity_type": "flag_sharing",
      "description": "flag{dyn_123xyz} submitted by teams 201 and 202 within 2 minutes",
      "related_team_ids": "[201,202]",
      "status": "pending",
      "created_at": "2025-08-26 10:05:00"
    }
  ]
}
```

#### (6) 管理员标记告警状态

- **URL**：`PUT /api/v1/admin/anti-cheat/alerts/:id`
- **请求体**：

```json
{
  "status": "confirmed_cheating"
}
```

- **返回**：

```json
{
  "code": 200,
  "msg": "Alert updated successfully"
}
```

***

### 3. 业务逻辑说明

1. **容器流量分析**
   - 抓包在宿主机进行，容器销毁后文件仍存档。
   - 定时分析（tshark/zeek/suricata）→ 写入 `analysis_result`。
   - 管理员可人工下载 PCAP 文件进一步调查。
2. **附件下载检测**
   - 短时间（如 1 分钟）内，同一 IP 下载了不同队伍的附件 → 标记为可疑。
   - 与提交日志对比时，若提交 IP 和下载 IP 重合，进一步提高作弊嫌疑。
3. **IP 变动跟踪**
   - 记录用户的所有登录 / 提交 IP。
   - 用于检测账号共享、VPN/代理频繁切换等情况。
4. **可疑行为告警**
   - 系统自动检测以下行为并写入 `dalictf_suspicious_activity`：
     - **flag_sharing**：不同队伍提交了相同的动态 flag
     - **multi_ip_download**：同一 IP 短时间内下载多个队伍的附件
     - **same_device**（可选扩展）：同一设备指纹登录不同账号
     - **traffic_anomaly**：容器流量异常（如可疑外联 IP）
   - 管理员人工审核，标记结果。
5. **权限严格**
   - 所有反作弊接口 **仅管理员可用**。
   - 前端用户完全无法访问这些数据。
## 十二、运维与监控模块

### 1. 数据库设计

***

#### (1) 系统操作日志表

**表名**：`dalictf_system_log`

| 字段名        | 类型           | 约束条件                             | 说明                               |
| :-----------: | :------------: | :----------------------------------: | :--------------------------------: |
|      id       |  `BIGINT(20)`  |     `PRIMARY KEY AUTO_INCREMENT`     |              日志主键              |
|    user_id    |  `BIGINT(20)`  |            `DEFAULT NULL`            | 操作人ID（NULL表示系统自动操作） |
|   username    | `VARCHAR(50)`  |            `DEFAULT NULL`            |          操作人用户名          |
|    action     | `VARCHAR(100)` |              `NOT NULL`              |      操作类型（如 Login）      |
|    method     |  `VARCHAR(10)` |            `DEFAULT NULL`            |       请求方法（GET/POST）       |
|      url      | `VARCHAR(255)` |            `DEFAULT NULL`            |              请求路径              |
|    params     |     `TEXT`     |            `DEFAULT NULL`            |              请求参数              |
|      ip       | `VARCHAR(50)`  |            `DEFAULT NULL`            |              操作IP              |
|  status_code  |    `INT(11)`   |              `NOT NULL`              |              响应状态码              |
| error_message |     `TEXT`     |            `DEFAULT NULL`            |              错误信息              |
|  created_at   |   `DATETIME`   | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |              操作时间              |

```sql
CREATE TABLE `dalictf_system_log` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '日志主键ID',
  `user_id` BIGINT(20) DEFAULT NULL COMMENT '操作人ID',
  `username` VARCHAR(50) DEFAULT NULL COMMENT '操作人用户名',
  `action` VARCHAR(100) NOT NULL COMMENT '操作类型',
  `method` VARCHAR(10) DEFAULT NULL COMMENT '请求方法',
  `url` VARCHAR(255) DEFAULT NULL COMMENT '请求路径',
  `params` TEXT DEFAULT NULL COMMENT '请求参数',
  `ip` VARCHAR(50) DEFAULT NULL COMMENT '操作IP',
  `status_code` INT(11) NOT NULL COMMENT '响应状态码',
  `error_message` TEXT DEFAULT NULL COMMENT '错误信息',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_action` (`action`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统操作日志表';
```

***

#### (2) 数据库备份记录表

**表名**：`dalictf_backup_log`

| 字段名     | 类型           | 约束条件                             | 说明                               |
| :--------: | :------------: | :----------------------------------: | :--------------------------------: |
|     id     |  `BIGINT(20)`  |     `PRIMARY KEY AUTO_INCREMENT`     |              备份主键              |
|  filename  | `VARCHAR(255)` |              `NOT NULL`              |              备份文件名              |
|    size    |  `BIGINT(20)`  |              `NOT NULL`              |           文件大小（字节）           |
|    path    | `VARCHAR(255)` |              `NOT NULL`              |              存储路径              |
|   status   | `ENUM(...)`    |     `NOT NULL DEFAULT 'success'`     | 状态：success-成功/failed-失败 |
| created_by |  `BIGINT(20)`  |            `DEFAULT NULL`            | 操作人ID（NULL表示自动备份） |
| created_at |   `DATETIME`   | `NOT NULL DEFAULT CURRENT_TIMESTAMP` |              备份时间              |

```sql
CREATE TABLE `dalictf_backup_log` (
  `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '备份主键ID',
  `filename` VARCHAR(255) NOT NULL COMMENT '备份文件名',
  `size` BIGINT(20) NOT NULL COMMENT '文件大小',
  `path` VARCHAR(255) NOT NULL COMMENT '存储路径',
  `status` ENUM('success', 'failed') NOT NULL DEFAULT 'success' COMMENT '备份状态',
  `created_by` BIGINT(20) DEFAULT NULL COMMENT '操作人ID',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '备份时间',
  PRIMARY KEY (`id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据库备份记录表';
```

***

### 2. API 设计

> **接口前缀：**`/api/v1/admin/system`
>
> **权限：**仅超级管理员（Super Admin）

#### 获取系统状态

- **URL**：`GET /api/v1/admin/system/status`
- **说明**：返回服务器当前的资源使用情况。
- **成功返回**：

```json
{
  "code": 200,
  "msg": "Success",
  "data": {
    "cpu_usage": "15.5%",
    "memory_usage": "4.2GB / 16GB",
    "disk_usage": "45%",
    "uptime": "10 days 5 hours",
    "database_status": "connected",
    "redis_status": "connected"
  }
}
```

#### 查询操作日志

- **URL**：`GET /api/v1/admin/system/logs`
- **请求参数**：
  - `page`, `limit`
  - `user_id`
  - `action`
  - `start_time`, `end_time`
- **成功返回**：日志列表。

#### 触发数据库备份

- **URL**：`POST /api/v1/admin/system/backup`
- **说明**：手动触发一次全量数据库备份。
- **成功返回**：

```json
{
  "code": 200,
  "msg": "Backup started successfully",
  "data": {
    "job_id": "backup_20251120_100000"
  }
}
```

#### 获取备份列表

- **URL**：`GET /api/v1/admin/system/backups`
- **成功返回**：备份文件列表。

***

### 3. 业务逻辑说明

#### 3.1 日志审计

-   **AOP 拦截**：建议使用 AOP（面向切面编程）拦截所有 `/api/v1/admin/*` 的写操作请求，自动记录到 `dalictf_system_log`。
-   **异常捕获**：全局异常处理器捕获的 500 错误也应记录到日志中，方便排查。

#### 3.2 备份策略

-   **自动备份**：配置 Linux Crontab 或应用内定时任务，每天凌晨 2:00 执行 `mysqldump`。
-   **保留策略**：保留最近 7 天的备份，更早的备份自动删除或归档到对象存储（OSS）。

#### 3.3 健康检查

-   系统应提供 `/health` 接口（不鉴权或简单鉴权），供负载均衡器或 K8s 探针调用。
-   检查内容：MySQL 连接、Redis 连接、磁盘剩余空间。
