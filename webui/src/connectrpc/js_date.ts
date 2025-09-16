import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { timestampDate, timestampFromDate } from "@bufbuild/protobuf/wkt";

/**
 * Converts a `Timestamp` to a JavaScript `Date`.
 *
 * @param timestamp - The `Timestamp` to convert.
 * @returns The corresponding JavaScript `Date`.
 */
export function toDate(timestamp: Timestamp): Date {
  return timestampDate(timestamp);
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
export function fromDate(date: Date): Timestamp {
  return timestampFromDate(date);
}
