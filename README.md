# bookstore-go

前端启动方式
- 用户端
```bash

# 停止并删除旧容器
docker rm -f bookstore-frontend

# 构建镜像
docker build -t bookstore-frontend:latest ./frontend/bookstore-fronted-master

# 运行容器
docker run -d --name bookstore-frontend -p 3000:3000 bookstore-frontend:latest

# 启动容器 bookstore-frontend
docker start bookstore-frontend

# 查看日志
docker logs -f bookstore-frontend
```

- 管理端
```bash
docker build -t bookstore-admin-frontend:latest ./frontend/bookstore-admin-fronted-master
docker run -d --name bookstore-admin-frontend -p 3001:3000 bookstore-admin-frontend:latest
```

- 维护
```bash

# 查看日志
docker logs -f bookstore-frontend
docker logs -f bookstore-admin-frontend

# 停止与删除容器
docker rm -f bookstore-frontend bookstore-admin-frontend
```