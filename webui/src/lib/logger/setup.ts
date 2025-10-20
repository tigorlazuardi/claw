export async function setup() {
  if (Otel.DEV_MODE) {
    const { ConsoleLogBackend } = await import("./backend_console");
    const mod = await import("./backend");
    mod.backends.push(ConsoleLogBackend);
  }
  let endpoint = Otel.OTEL_EXPORTER_OTLP_LOGS_ENDPOINT;
  if (!endpoint) {
    endpoint = Otel.OTEL_EXPORTER_OTLP_ENDPOINT;
  }
  if (!endpoint) {
    return;
  }
  const { OtelBackend } = await import("./backend_otel");
  const mod = await import("./backend");
  mod.backends.push(OtelBackend);
}
