# mis
mis-background  
mis项目后端仓库

### 主要目录文件
- api: 接口层。业务逻辑代码，各个模块的接口方法
    - v1: v1版本接口
    - api.go: 接口层通用方法
- conf: 存放项目配置文件
    - app.ini: 项目配置文件
    - app-example.ini: 脱敏，示例配置文件
- middleware: 项目中间件
    - jwt.go: jwt中间件
    - power.go: 接口权限控制中间件
- models: 数据库层。数据库模型定义、数据操作，建议每个模块都新建不同的model
    - models.go: 数据库初始化
- pkg: 工具库，包括一些基本的方法log日志之类的
    - e: 定义错误码以及对应的错误信息
        - code.go: 错误码
        - msg.go: 错误信息
    - logging: 日志记录相关方法
        - file.go: 日志文件相关方法
        - log.go: 日志记录api
    - settings: 项目配置文件读取
        - settings.go
    - util: 其他一些小工具
        - jwt.go: jwt生成和解析token相关方法
        - pagination.go: 分页相关方法
        - utils.go: 其他工具方法
        - validate.go: 接口参数校验
- routers: 路由层。
- runtime: 运行时日志文件存放目录
- .gitignore: 在git push时忽略一些不必要上传的文件
- go.mod: 项目依赖
- main.go: 项目入口
- README.md: 项目说明