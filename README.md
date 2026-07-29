## Categories

- GET       /category               -> возвращает все категории
- POST      /category               -> создает новую категорию
- GET       /category/id/{name}     -> возвращает id по названию
- DELETE    /category/{name}        -> удаляет категорию

## Materials

- POST      /material               -> сохраняет материал и возвращает id
- GET       /materials/{id}         -> возвращает материалы категории id
- DELETE    /material/{id}          -> удаляет материал