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

#define PIN_S A3
#define AIR_VALUE 489
#define WATER_VALUE 238

void setup() {
  Serial.begin(9600);

  while (!Serial) {
    ;
  }

  dht.begin();

  Serial.println("Calibrating MQ7");
  mq7.calibrate();
  Serial.println("Calibration done!");
}

float handleMq135() {
  float temperature = dht.readTemperature();
  float humidity = dht.readHumidity();

  if (isnan(temperature) || isnan(humidity)) {
    return mq135.getPPM();
  }
  return mq135.getCorrectedPPM(temperature, humidity);
}

void loop() {
  while (Serial.available() == 0) {
  }

  switch (Serial.parseInt()) {
    case 1:
      //Serial.println("mq7");
      Serial.println(mq7.readPpm());
      break;
    case 2:
      //Serial.println("mq135");
      Serial.println(handleMq135());
      break;
    case 3:
      //Serial.println("raindrop");
      Serial.println(analogRead(PIN_R));
      break;
    case 4:
      //Serial.println("moisture");
      int moistureValue = analogRead(PIN_S);
      int moisturePercent = map(moistureValue, AIR_VALUE, WATER_VALUE, 0, 100);
      //Serial.println(moistureValue);

      if (moisturePercent >= 100) {
        Serial.println(100);
      } else if (moisturePercent <= 0) {
        Serial.println(0);
      } else {
        Serial.println(moisturePercent);
      }
      break;
  }
}
