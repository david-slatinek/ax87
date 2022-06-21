#include <DHT.h>
#include <MQ135.h>

#define DHTPIN 7
#define DHTTYPE DHT22
#define PIN A2

DHT dht(DHTPIN, DHTTYPE);
MQ135 mq135(PIN);

void setup() {
  Serial.begin(9600);

  while (!Serial) {
    ;
  }

  dht.begin();
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
      Serial.println("mq7");
      break;
    case 2:
      Serial.println(handleMq135());
      //Serial.println("mq135");
      break;
    case 3:
      Serial.println("raindrop");
      break;
    case 4:
      Serial.println("moisture");
      break;
  }
}
