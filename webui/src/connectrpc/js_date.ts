import type { Timestamp } from "@bufbuild/protobuf/wkt";
import type { M } from "../types";

/**
 * Converts a `Timestamp` to a JavaScript `Date`.
 *
 * @param timestamp - The `Timestamp` to convert.
 * @returns The corresponding JavaScript `Date`.
 */
export function toDate(timestamp: M<Timestamp>): Date {
  return new Date(
    Number(timestamp.seconds) * 1000 + timestamp.nanos / 1_000_000,
  );
}

/**
 * Converts a JavaScript `Date` to a `Timestamp`.
 *
 * Note that this conversion does create protobuf `Timestamp` objects, it
 * merely creates plain JavaScript objects that match the `Timestamp` shape.
 *
 * @param date - The JavaScript `Date` to convert.
 * @returns The corresponding `Timestamp`.
 */
export function fromDate(date: Date): M<Timestamp> {
  const seconds = Math.floor(date.getTime() / 1000);
  const nanos = (date.getTime() % 1000) * 1_000_000;
  return { seconds: BigInt(seconds), nanos };
}
