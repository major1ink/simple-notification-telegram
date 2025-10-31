# simple-notification-telegram

## Содержание

- [О проекте](#о-проекте)
- [Стуктура конфигурационного файла](#структура-конфигурационного-файла-необходимо-соблюдать-вложенность)

## О проекте
Проект для отправки уведомлений в telegram из topic событий в kafka.

По умолчанию сервис ищет файл конфигурации рядом с собой. При необходимости можно указать, где находится файл с помощью ключа `--configPath=Path`. 
Где `Path`- путь до файла.

```bash
--configPath=
```

## Структура конфигурационного файла (Необходимо соблюдать вложенность)

```YAML
# Параметры логгера
logger:
  # Уровень логирования (INFO|ERROR|DEBUG)
  logLevel: DEBUG
  # Директория сохранения лог файла (рекомендуется не использовать / в конце)
  logDir: ./
  # Режим логирования (stdout(вывод в консоль), file(запись в файл), оставить поле пустым(stdout+file))
  logMode: 
  # Перезапись лог файла при старте работы
  rewriteLog: true
# Конфигурация брокера kafka
kafkaConfig:
  brokers:
    - 127.0.0.1:32000
    - 127.0.0.1:32001
# Конфигурация topic
consumerConfig:
  topic: notification-assembled
  group_id: notification-assembled-1
# Конфигурация telegram
telegramConfig:
  telegram_bot_token:
  telegram_chat_id: 
```

