import {
  type SeverityNumber,
  type AnyValueMap,
  type LogRecord,
} from "@opentelemetry/api-logs";

export type LogOption = Omit<
  LogRecord,
  "body" | "severityNumber" | "severityText" | "attributes"
>;

export interface LogBackend {
  log: (
    level: SeverityNumber,
    message: string,
    attrs?: AnyValueMap,
    opts?: LogOption,
  ) => void;
}

export const backends: LogBackend[] = [];

export function info(message: string, attrs?: AnyValueMap, opts?: LogOption) {
  log(9, message, attrs, opts);
}

export function error(message: string, attrs?: AnyValueMap, opts?: LogOption) {
  log(17, message, attrs, opts);
}

export function debug(message: string, attrs?: AnyValueMap, opts?: LogOption) {
  log(5, message, attrs, opts);
}

export function warn(message: string, attrs?: AnyValueMap, opts?: LogOption) {
  log(13, message, attrs, opts);
}

export function log(
  level: SeverityNumber,
  message: string,
  attrs?: AnyValueMap,
  opts?: LogOption,
) {
  for (const backend of backends) {
    backend.log(level, message, attrs, opts);
  }
}
