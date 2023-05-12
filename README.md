# Metrics Consumer

Este proyecto es un Cron Job que corre cada 1 minuto.
Es el Consumer de las metricas guardadas en Redis por el proyecto `metrics-queue`.

## Arquitectura de Métricas

Los distintos servicios del backend producen varios tipos de métricas. Entre ellas, tenemos:
- Métricas de Entrenamiento (Resultados de un entrenamiento completado)
- Métricas Personales (Actualización de la cantidad de pasos, distancia, o calorías quemadas para un usuario)
- Métricas de Sistema (creación de nuevo usuario, creación de plan de entrenamiento, login de usuario, etc)

A medida que los Servicios producen estas métricas, las envían a `metrics-queue` vía petición HTTP, que es
un servidor que se encarga de acumular estas métricas para que luego sean procesadas asincrónicamente.


Las métricas, son guardadas en una de 3 Queues en una base de datos Redis según los tupos mencionados arriba. Estas queues son:
- `training-metrics` para métricas de entrenamiento
- `personal-metrics` para métricas personales
- `system-metrics` para métricas de sistemas

Luego, el `metrics-consumer` se encarga de procesarlas, es un Cron que se ejecuta cada 1 minuto.

Para cada una de estas colas, aplica el procesamiento necesario y las guarda en una base de datos MongoDB, que luego será consultada por
la interfaz gráfica de admins y de la aplicación.

- Para las `training-metrics`, se guardan en mongo tal como fueron encoladas en redis
- Las `personal-metrics` deben ser agregadas con una función de `count`, ya que podría haber más de un documento para un mismo usuario en un minuto.
- Las `system-metrics` son también contadas por nombre, y almacenadas con el count en Mongo

## Formato de los Documentos JSON recibidos

### Métricas de Entrenamiento

Las métricas de entrenamiento incluyen cual fue el usuario y el entrenamiento completado, con los valores del resultado del entrenamiento

```
{
    "type": "training"
    "name": "training_completed"
    "user_id": 123,
    "training_id": 456,
    "data": {
        "calories": 555,
        "time_in_minutes": 45,
        "distance_in_meters": 3000  
    },
    "timestamp": "16:08:55-07/05/2023"
}
```

### Métricas de Sistema

Marcan solamente el nombre de la métrica. Dos ejemplos para usuarios creados y entrenamiento creado:

```
{
    "type": "system",
    "name": "user_created"
}
```

```
{
    "type": "system",
    "name": "training_plan.created.running"   
}
```

### Métricas Personales

Marcan para un usuario y una métrica personal (puede ser `calories`, `distance` o `steps`), la cantidad a aumentar

```
{
    "type": "personal",
    "user_id": 123,
    "metric": "steps",
    "amount": 10000
}
```

## Formato de los documentos JSON a guardar en MongoDB

### Métricas Personales

Agregan cada tipo de métrica para cada usuario, y un timestamp:

```
{
    "user_id": 123,
    "metric": "steps",
    "count": 12000 // La cantidad total de documentos con steps como métrica suma 12000
    "timestamp": "16:08:55-07/05/2023"
}
```

### Métricas de Sistema

Agregan la cantidad de ocurrencias por tipo:

```
{
    "name": "training_plan.created.running",
    "count": 12,
    "timestamp": "16:08:55-07/05/2023"
}
```


### Métricas de entrenamientos

Se guardan como vienen:

```
{
    "user_id": 123,
    "training_id": 456,
    "data": {
        "calories": 555,
        "time_in_minutes": 45,
        "distance_in_meters": 3000  
    },
    "timestamp": "16:08:55-07/05/2023"
}
```