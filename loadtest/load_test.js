import http from "k6/http";
import { check, sleep } from "k6";

// configuração de cenários de carga crescente
export const options = {
  stages: [
    { duration: "30s", target: 10 },   // rampa: 0 → 10 VUs
    { duration: "1m",  target: 50 },   // sustentado: 50 VUs
    { duration: "30s", target: 100 },  // pico: 100 VUs
    { duration: "30s", target: 0 },    // rampa de descida
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"],  // 95% das requisições < 500ms
    http_req_failed:   ["rate<0.01"],  // menos de 1% de erro
  },
};

const SENSOR_TYPES = ["temperature", "humidity", "presence", "vibration", "luminosity", "reservoir"];
const READING_TYPES = ["analog", "discrete"];

function randomSensorType() {
  return SENSOR_TYPES[Math.floor(Math.random() * SENSOR_TYPES.length)];
}

function randomValue(sensorType) {
  if (sensorType === "presence") return Math.random() > 0.5 ? 1 : 0;
  return parseFloat((Math.random() * 100).toFixed(2));
}

export default function () {
  const sensorType = randomSensorType();
  const readingType = sensorType === "presence" ? "discrete" : "analog";

  const payload = JSON.stringify({
    device_id:    `device-${Math.floor(Math.random() * 100)}`,
    timestamp:    new Date().toISOString(),
    sensor_type:  sensorType,
    reading_type: readingType,
    value:        randomValue(sensorType),
  });

  const res = http.post("http://backend:8080/telemetry", payload, {
    headers: { "Content-Type": "application/json" },
  });

  check(res, {
    "status is 202": (r) => r.status === 202,
  });

  sleep(0.1);
}