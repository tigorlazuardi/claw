import type { LogRecord } from "@opentelemetry/api-logs";
import type { LogBackend } from "./backend";

export const OtelBackend: LogBackend = {
  log(level, message, attrs, opts) {
    import("@opentelemetry/api-logs").then((mod) => {
      const log = mod.logs.getLogger("default");
      const record: LogRecord = opts ?? {};
      record.severityNumber = level;
      record.severityText = mod.SeverityNumber[level];
      record.body = message;
      if (attrs) {
        record.attributes = attrs;
      }
      log.emit(record);
    });
  },
};
