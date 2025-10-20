declare namespace Otel {
  const OTEL_EXPORTER_OTLP_ENDPOINT: string | undefined;
  const OTEL_EXPORTER_OTLP_LOGS_ENDPOINT: string | undefined;
  const OTEL_EXPORTER_OTLP_TRACES_ENDPOINT: string | undefined;
  /** if this value is "grpc" or undefined, use "http/protobuf" */
  const OTEL_EXPORTER_OTLP_PROTOCOL: "http/protobuf" | "http/json" | undefined;
  const OTEL_EXPORTER_OTLP_LOGS_PROTOCOL:
    | "http/protobuf"
    | "http/json"
    | undefined;
  const OTEL_EXPORTER_OTLP_TRACES_PROTOCOL:
    | "http/protobuf"
    | "http/json"
    | undefined;
  const OTEL_RESOURCE_ATTRIBUTES: string | undefined;

  /** DEV_MODE when set to true should show logs in console */
  const DEV_MODE: boolean | undefined;
}
