# Онлайн Библиотека Песен 🎶

## Описание проекта

Задание включало реализацию онлайн библиотеки песен с возможностью фильтрации, получения текста песни, работы с API и интеграции внешних сервисов. В ходе работы возникло несколько сложных моментов, которые я хотел бы выделить.

---

## Проблемы, с которыми столкнулся

### 1. Интеграция внешних API

Здесь можно целую легенду расписывать… Изначально в планах было использовать обычный API какого-нибудь музыкального сайта, однако тестируя запросы, я пришел к выводу, что обрабатывая ответы, я потрачу очень много времени, которого у меня и так не было. Поэтому я решил интегрировать ИИ — идея была безумно хорошей, горжусь собой, что вообще догадался. Однако подводных камней было в десятки раз больше...

Для начала я искал бесплатный API с ИИ, который работает в СНГ без проблем. Нашел, начал использовать — бесплатный оказался слишком глупым, и запрос обрабатывал крайне плохо. Купил новую версию, умнее — та же проблема, они безумно глупые.

Пошел по тяжёлому пути, купил европейский номер телефона для валидации ключа для OpenAI. Купил мощный ВПН. Но этого мало, нужна ещё и карта для пополнения счёта... С этим тоже огромные проблемы, по итогу я её так и не нашёл. И просто попросил у знакомого его API Key...

Продолжая сие великолепие, я уже думал, что всё, самое сложное кончилось... Ага, щас... OpenAI не блокирует авторские права :)

И вот я в том месте, где был и сначала... Пришлось всё-таки брать GeniusAPI, обрабатывать огромный ответ и парсить его в текст песни. Но есть нюанс. Текст песни получается без каких-либо пагинаций, куплетов и т. д. Тут-то OpenAI и пригодился. Сделал пересылку этого текста из GeniusAPI в ChatGPT3.5_16k.

Также для нахождения URL ссылки на клип интегрировал YouTube API. Тут, слава богу, справился быстро, без особых проблем.

Итого, на один только внешний API я потратил почти все степендию и часов так 15 работы :)

Хорошо, с этим закончил.

---

### 2. Проблемы с .env файлом

Здесь же и объясню, почему его нет. Дело в том, что GoLand (мой IDE) не поддерживает .env конфигурацию. Да, безусловно, можно установить плагин, но нет, нельзя :) Почему? Да потому что JetBrains контора... Которая блокирует все, что можно для СНГ. Решением могло бы быть переход на VSCode... Но я почему-то и не подумал об этом тогда...

В общем, по итогу там стоит .yaml конфигурация, надеюсь, это сильно не повлияет на ваше решение...

---

### 3. Идеи, что бы я сделал, но не успел

Хотел бы выделить свои идеи по этому проекту, что бы я ещё сделал, но не сделал, потому что уже не успею (IRL жизнь):

1. **Redis** — в первую очередь, это мой любимый способ оптимизации. Я до сих пор не понимаю, почему все им не пользуются.Обернуть им PostgreSQL, и запросы будут обрабатываться в разы быстрее, а то и в десятки раз. Безумно мощный инструмент оптимизации, наверное моя любимая БД в целом.Обычно я пихаю её везде, куда только можно.

2. **Логирование** — я сделал его, мягко говоря, хреново. Мне самому очень не нравится... Пощадите, пожалуйста. Обычно делаю это сильно структурирование...

3. **Индексы и транзакции в PostgreSQL** — без комментариев, оптимизация. 
#### В целом это основное, что приходило в голову.

---

### 4. Обоснование выбора библиотек

1. **Tracerr** — врапер, обработчик ошибок. Обожаю его, пихаю куда только можно, да и не можно тоже. Безумно удобный инструмент.

2. **Logrus** — логирование. Здесь надо выделить, почему не дефолтный `slog`. Почему? Да просто личная симпатия, не больше. Привычка и всё.

3. **OpenAI** — в отличие от GeniusAPI (где запросы писать мягко говоря больно), в OpenAI сделали библиотеку для работы с их API. Спасибо им большое за это.

---

## Вывод

Проект очень интересный, даже если вас не устроит решение, я вам очень благодарен. Объективно — мне понравилось.

> **Автор:** [Данила Пиварович]
