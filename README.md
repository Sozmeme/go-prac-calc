# go-prac-calc
Простой калькулятор на Go для проработки основных особенностей языка, конкурентности. Подробности задания в `task.md`.
# Установка
1. Клонируйте репозиторий
   ```bash
   git clone https://github.com/Sozmeme/go-prac-calc.git
   cd ./go-prac-calc
   ```
2. Поднимите контейнер
   ```bash
   docker-compose up -d
   ```
# Использование
1. curl Windows:
   ```bash
   curl -X POST http://localhost:8080/calculate -H "Content-Type: application/json" -d "[{\"type\":\"calc\",\"op\":\"+\",\"var\":\"x\",\"left\":10,\"right\":2},{\"type\":\"calc\",\"op\":\"*\",\"var\":\"y\",\"left\":\"x\",\"right\":5},{\"type\":\"calc\",\"op\":\"-\",\"var\":\"q\",\"left\":\"y\",\"right\":20},{\"type\":\"calc\",\"op\":\"+\",\"var\":\"unusedA\",\"left\":\"y\",\"right\":100},{\"type\":\"calc\",\"op\":\"*\",\"var\":\"unusedB\",\"left\":\"unusedA\",\"right\":2},{\"type\":\"print\",\"var\":\"q\"},{\"type\":\"calc\",\"op\":\"-\",\"var\":\"z\",\"left\":\"x\",\"right\":15},{\"type\":\"print\",\"var\":\"z\"},{\"type\":\"calc\",\"op\":\"+\",\"var\":\"ignoreC\",\"left\":\"z\",\"right\":\"y\"},{\"type\":\"print\",\"var\":\"x\"}]"
   ```
2. curl Linux
   ```bash
   curl -X POST http://localhost:8080/calculate -H "Content-Type: application/json" -d '[{"type":"calc","op":"+","var":"x","left":10,"right":2},{"type":"calc","op":"*","var":"y","left":"x","right":5},{"type":"calc","op":"-","var":"q","left":"y","right":20},{"type":"calc","op":"+","var":"unusedA","left":"y","right":100},{"type":"calc","op":"*","var":"unusedB","left":"unusedA","right":2},{"type":"print","var":"q"},{"type":"calc","op":"-","var":"z","left":"x","right":15},{"type":"print","var":"z"},{"type":"print","var":"x"}]'
   ```
- Ответ:
  `{"items":[{"var":"q","value":40},{"var":"z","value":-3},{"var":"x","value":12}]}`
3. Возможно использование swagger
    ```bash
    http://localhost:8080/swagger/index.html
    ```
  
