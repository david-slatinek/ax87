#include <DHT.h>
#include <MQ135.h>
#include "MQ7.h"

#define DHTPIN 7
#define DHTTYPE DHT22
DHT dht(DHTPIN, DHTTYPE);

#define PIN_MQ135 A0
MQ135 mq135(PIN_MQ135);

#define PIN_MQ7 1
#define VOLTAGE 5
MQ7 mq7(PIN_MQ7, VOLTAGE);

#define PIN_R A2
#define MIN_VALUE 0
#define MAX_VALUE 1024
#define MIN_CAT 0
#define MAX_CAT 3

#define PIN_S A3
#define AIR_VALUE 489
#define WATER_VALUE 238

#define ARRAY_SIZE 12

void setup() {
  Serial.begin(9600);

  while (!Serial) {
    ;
  }

  dht.begin();

  //mq7.calibrate();
}

float handleMq135() {
  float temperature = dht.readTemperature();
  float humidity = dht.readHumidity();

  if (isnan(temperature) || isnan(humidity)) {
    return mq135.getPPM();
  }
  return mq135.getCorrectedPPM(temperature, humidity);
}

void respond(float a) {
  char x[ARRAY_SIZE];
  dtostrf(a, 8, 2, x);
  x[ARRAY_SIZE - 1] = '\0';
  Serial.write(x, ARRAY_SIZE);
}

void loop() {
  while (Serial.available() == 0) {
  }

  switch (Serial.parseInt()) {
    case 1:
      respond(mq7.readPpm());
      //Serial.print(mq7.readPpm());
      break;
    case 2:
      respond(handleMq135());
      //Serial.println(handleMq135());
      break;
    case 3:
      respond(map(analogRead(PIN_R), MIN_VALUE, MAX_VALUE, MAX_CAT, MIN_CAT));
      //Serial.println(map(analogRead(PIN_R), MIN_VALUE, MAX_VALUE, MAX_CAT, MIN_CAT));
      break;
    case 4:
      respond(map(analogRead(PIN_S), AIR_VALUE, WATER_VALUE, 0, 100));
      //Serial.println(map(analogRead(PIN_S), AIR_VALUE, WATER_VALUE, 0, 100));
      break;
  }
}
