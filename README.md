# 使用 Gin + Gorm + mysql 实现的个人博客web应用

## 1.界面截图

### （1）登录界面

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/login.png)

###   (2) 注册用户界面

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/register.png)

### （3）修改密码页面

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/update_passw.png)

### （4） 用户主页

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/home.png)

### （5）我的盆摘页面

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/my_plant01.png)
!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/my_plant02.png)



### (6）朋友圈页面

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/my_friends.png)



### （7）管理员后台管理界面权限表

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/mgr_permission.png)



### 	（8）管理员后台管理用户表

​			!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/mgr_user.png)

###       (9) 管理员后台管理用户组表

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/mgr_usergroup.png)

### （10）管理员后台管理设备表

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/mgr_device.png)

### （11）管理员后台管理朋友圈动态表

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/mgr_friendinfo.png)

### （12）管理员后台管理朋友圈动态评论表

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/mgr_friendcomments.png)

### （13）管理员后台管理朋友圈点赞信息表

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/mgr_friendlikes.png)



## 2.环境配置



### （1）安装Golang ，及gin，gorm



### （2） 安装虚拟环境并安装依赖包

​	!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/windows_env.png)

​			

![_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/windows_env_1.png)

​	

Linux下：

​			创建虚拟环境：virtualenv -p python3.6 garden_ENV  或 virtualenv garden_ENV

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/linux_env.png)

​			激活虚拟环境: source garden_ENV/bin/activate

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/linux_env_1.png)

​			除了上述步骤与windows不一样以外其他步骤均相同。

### （3）创建表结构

使用的数据库为:Mysql  

首先需要启用mysql服务: windows :net start mysql   linux: service mysql start

在项目代码中的app.py中将下面的代码运行即可自动创建表结构.

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/create_table.png)

## 3. 运行web应用程序

在权限表中导入权限数据：

权限数据在：database文件夹中(permission.sql,permissions.xlsx)

执行sql脚本或者在网页的后台管理权限表中批量导入。

进入虚拟环境后：python app.py 即可运行该应用程序。

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/run.png)

## 4. 使用nginx + gunicorn 部署

### （1）nginx的配置

####         安装nginx ：yum install nginx

​		（1）安装完成后查看nginx版本：nginx -v!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/nginx_version.png)



​         （2）在/etc/nginx/nginx.conf中修改配置：vim /etc/nginx/nginx.conf    注：server中的端口为nginx的监听端口， location中的地址为gunicorn 运行flask app服务的地址,如:127.0.0.1:5001 ,下面两个location为配置的静态资源地址。如果出现加载静态资源报403的情况，需要把nginx中的配置：user nginx; 改为：user root; 修改完成后按esc 冒号 输入：wq 保存退出vim。

![/flask/练习项目/智能浇水平台new/GardenPlatform/img/nginx_config.png)

​             (3) 重启nginx : nginx -s reload  查看nginx的服务 ：lsof -i:5000 查看到对应的进程即为配置成功 。运行nginx命令：nginx

!(D:/python_web/flask/练习项目/智能浇水平台new/GardenPlatform/img/nginx_status.png)



## 5.感谢

[![](./img/jetbrains.png)](https://www.jetbrains.com/)







