import type { LogBackend } from "./backend";

export const ConsoleLogBackend: LogBackend = {
  log(level, message, attrs, opts) {
    let engine = console.debug;
    if (level >= 9) {
      engine = console.info;
    }
    if (level >= 13) {
      engine = console.warn;
    }
    if (level >= 17) {
      engine = console.error;
    }
    if (opts?.eventName) {
      message = `[${opts.eventName}] ${message}`;
    }
    if (attrs) {
      engine(message, attrs);
      return;
    }
    engine(message);
  },
};
