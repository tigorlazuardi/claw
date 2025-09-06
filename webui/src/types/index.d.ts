/**
 * type M is a utitlity type to strip protobuf specific fields to satisfy typescript checks.
 */
export type M<T extends {}> = Omit<T, "$typeName" | "$unknown">;
