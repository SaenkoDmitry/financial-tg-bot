# Financial Telegram Bot
@finance_for_you_bot

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

### Главное меню приложения / Добавление расходов / Отчет о расходах по категориям :
<p align="left">
  <img width="200" height="450" src="/screenshots/main_menu.png">  
  <img width="200" height="450" src="/screenshots/add_operation.png">  
  <img width="200" height="450" src="/screenshots/show_report.png">
</p>

### Установка лимита трат и сообщение о его превышении:
<p align="left">
  <img width="350" height="200" src="/screenshots/set_limit.jpg">
  <img width="250" height="200" src="/screenshots/limit_notification.jpg">
</p>
