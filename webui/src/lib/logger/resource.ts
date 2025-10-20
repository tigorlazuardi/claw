import {
  type DetectedResourceAttributes,
  type Resource,
  detectResources,
  resourceFromAttributes,
} from "@opentelemetry/resources";
import { browserDetector } from "@opentelemetry/opentelemetry-browser-detector";

export function createResource(): Resource {
  let res = detectResources({
    detectors: [browserDetector],
  });
  res = res.merge(
    resourceFromAttributes({
      "service.namespace": "claw",
    }),
  );
  if (Otel.OTEL_RESOURCE_ATTRIBUTES) {
    const parsed = parseStringAttributes(Otel.OTEL_RESOURCE_ATTRIBUTES);
    res = res.merge(parsed);
  }
  res.merge(resourceFromAttributes({ "service.name": "webui" }));
  return res;
}

export function parseStringAttributes(s: string): Resource {
  const kvs = s.split(",");
  const attributes: DetectedResourceAttributes = {};
  for (const raw of kvs) {
    const [key, value] = raw.split("=", 2);
    if (!key || !value) {
      continue;
    }
    attributes[key.trim()] = value.trim();
  }
  return resourceFromAttributes(attributes);
}
