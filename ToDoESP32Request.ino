#include <TFT_eSPI.h>
#include <WiFi.h>
#include <ArduinoJson.h>
#include <HTTPClient.h>

TFT_eSPI tft = TFT_eSPI();

/*
Insert your network credentials
Insert the source of your json data
*/
const char* ssid = "YOUR_SSID";
const char* password = "YOUR_PASSWORD";
const char* serverURL = "http://YOUR_IP:8080/top3";

const int fontSize = 2;

void setup() {

  Serial.begin(115200);

  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(1000);
    Serial.println("Connecting to Wifi...");
  }
  Serial.println("Connected to Wifi!");

  tft.init();

  tft.setRotation(3);
  tft.fillScreen(TFT_BLACK);

  int x = 320 / 2;
  int y = 16;

  tft.setTextColor(TFT_RED, TFT_BLACK);
  tft.drawCentreString("! URGENT TASKS !", x, y, fontSize);

  int rectTop = 45;
  int rectSide = 30;
  tft.drawRect(rectSide, rectTop, 260, 165, TFT_WHITE);
}

void loop() {

  int x = 320 / 2;
  int y = 60;

  HTTPClient http;
  http.begin(serverURL);
  int httpCode = http.GET();

  if (httpCode > 0 ) {
      if (httpCode == HTTP_CODE_OK) {
      String payload = http.getString();
      Serial.println(payload);

      // Parse JSON
      StaticJsonDocument<1024> doc;
      DeserializationError error = deserializeJson(doc, payload);
      if (error) {
        Serial.print("JSON parse failed: ");
        Serial.println(error.c_str());
      } 
      else {
        for (JsonObject task : doc.as<JsonArray>()) {
          const char* id = task["ID"];
          const char* name = task["Name"];
          const char* due = task["Due"];
          Serial.printf("Task %s: %s due %s\n", id, name, due);
          String data = String(task["ID"].as<const char *>()) + ": " + task["Name"].as<const char *>();
          tft.setTextColor(TFT_WHITE, TFT_BLACK);
          tft.drawCentreString(data, x, y, fontSize);
          y += 16;
          data = "Due: " + String(due);
          tft.drawCentreString(data, x, y, fontSize);
          y += 24;
        }
      }
    }
  } 
  else {
    Serial.printf("HTTP error: %s\n", http.errorToString(httpCode).c_str());
  }

  http.end();


}