import { NSFWMode } from "#/gen/claw/v1/nsfw_pb";
import { PersistedState } from "runed";

export const deviceListSize = new PersistedState("device.list.size", 25);
export const sourceListSize = new PersistedState("source.list.size", 25);
export const imageListSize = new PersistedState("image.list.size", 100);
export const nsfwState = new PersistedState(
  "nsfw",
  NSFWMode.NSFW_MODE_DISALLOW,
);

export function parseNSFWState(s: string) {
  const v = parseInt(s);
  return Object.values(NSFWMode).includes(v)
    ? (v as NSFWMode)
    : nsfwState.current;
}
