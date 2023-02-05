const express = require("express");
const mongoose = require("mongoose");
const router = require("./routes/index");

const { PORT = 8080 } = process.env;

const app = express();

const allowedCors = ["https://localhost:3001", "http://localhost:3001"];

app.use(express.urlencoded({
  extended: true
}))
app.use(express.json());

// CORS
app.use(function (req, res, next) {
  const { origin } = req.headers; // Сохраняем источник запроса в переменную origin
  // проверяем, что источник запроса есть среди разрешённых
  if (allowedCors.includes(origin)) {
    res.header("Access-Control-Allow-Origin", origin);
  }
  const { method } = req; // Сохраняем тип запроса (HTTP-метод) в соответствующую переменную

  // Значение для заголовка Access-Control-Allow-Methods по умолчанию (разрешены все типы запросов)
  const DEFAULT_ALLOWED_METHODS = "GET,HEAD,PUT,PATCH,POST,DELETE";

  // Если это предварительный запрос, добавляем нужные заголовки
  if (method === "OPTIONS") {
    // разрешаем кросс-доменные запросы любых типов (по умолчанию)
    res.header("Access-Control-Allow-Methods", DEFAULT_ALLOWED_METHODS);
  }

  const requestHeaders = req.headers["access-control-request-headers"];
  if (method === "OPTIONS") {
    // разрешаем кросс-доменные запросы с этими заголовками
    res.header("Access-Control-Allow-Headers", requestHeaders);
    // завершаем обработку запроса и возвращаем результат клиенту
    return res.end();
  }
  next();
});

// app.use(auth);

app.use(router);

mongoose.connect("mongodb://127.0.0.1/diploma");

app.listen(PORT, () => {
  console.log(`App listening on port ${PORT}`);
});
