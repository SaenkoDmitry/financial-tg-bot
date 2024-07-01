# Financial Telegram Bot

## First Run
```
docker-compose up -d
```

Also you need store local config:
```
rates_cache_default_expiration: 720h
calc_cache_default_expiration: 24h
rates_cache_cleanup_interval: 1h
token: 
abstract_api_key: 
postgres_user:
postgres_password:
postgres_db:
postgres_port:
postgres_host: db
cache_host: memcached:11211
```

## Функционал

### Главное меню приложения / Добавление новой траты / Отчет о тратах по категориям :
<p align="left">
  <img width="200" height="500" src="/screenshots/main_menu.png">  
  <img width="200" height="500" src="/screenshots/add_operation.png">  
  <img width="200" height="500" src="/screenshots/show_report.png">
</p>

### Установка лимита трат и сообщение о его превышении:
<p align="left">
  <img width="400" height="200" src="/screenshots/set_limit.jpg">
  <img width="400" height="350" src="/screenshots/limit_notification.jpg">
</p>

