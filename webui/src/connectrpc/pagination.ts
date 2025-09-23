import * as v from "valibot";

import { PersistedState } from "runed";

function createPaginationSchema(key: string, defaultSize: number) {
  const state = new PersistedState(key, defaultSize);
  return v.object({
    size: v.pipe(
      v.optional(v.string(), state.current.toString()),
      v.transform((val) => {
        const num = parseInt(val);
        return isNaN(num) ? state.current : num;
      }),
    ),
    nextToken: v.pipe(
      v.undefinedable(v.string()),
      v.transform((val) => {
        if (!val) return 0;
        const num = parseInt(val);
        if (!num) return 0;
        return num;
      }),
    ),
    prevToken: v.pipe(
      v.undefinedable(v.string()),
      v.transform((val) => {
        if (!val) return 0;
        const num = parseInt(val);
        if (!num) return 0;
        return num;
      }),
    ),
  });
}

const devicePaginationSchema = createPaginationSchema("device.list.size", 25);
export const DevicePaginationSchema = v.optional(devicePaginationSchema, {
  size: v.getDefault(devicePaginationSchema.entries.size),
  nextToken: v.getDefault(devicePaginationSchema.entries.nextToken),
  prevToken: v.getDefault(devicePaginationSchema.entries.size),
});

export type DevicePagination = v.InferOutput<typeof DevicePaginationSchema>;
