import * as v from "valibot";

import { deviceListSize } from "#/store/searchQuery";

function createPaginationSchema(defaultSize: number) {
  return v.object({
    size: v.fallback(
      v.pipe(v.string(), v.transform(parseInt), v.number()),
      defaultSize,
    ),
    nextToken: v.fallback(
      v.pipe(v.string(), v.transform(parseInt), v.number()),
      0,
    ),
    prevToken: v.fallback(
      v.pipe(v.string(), v.transform(parseInt), v.number()),
      0,
    ),
  });
}

/**
 * getDevicePaginationSchema gets device pagination schema that is synced
 * with the local storage of the browser
 */
export function getDevicePaginationSchema() {
  const devicePaginationSchema = createPaginationSchema(deviceListSize.current);
  return v.fallback(devicePaginationSchema, {
    size: deviceListSize.current,
    nextToken: 0,
    prevToken: 0,
  });
}

export type DevicePagination = v.InferOutput<
  ReturnType<typeof getDevicePaginationSchema>
>;
